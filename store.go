package imagestore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/image"

	"github.com/mjibson/appstats"
)

type Image struct {
	ID      int64             `datastore:"-" json:"id"`
	BlobKey appengine.BlobKey `datastore:"bk" json:"-"`
	URL     string            `datastore:"u" json:"url,omitempty"`
}

// TODO(arunjit): Instead of panic-recover, writer a simple wrapper so that
// the handler can return an error.

func recoverHTTP(w http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(w, fmt.Sprintf("ERROR: %v\n", r), 500) // TODO(arunjit): status
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	defer recoverHTTP(w)
	mime := r.Header.Get("Content-Type")
	if mime == "" {
		panic("Missing image Content-Type [header]")
	}

	c := appstats.NewContext(r)
	defer c.Save()

	bw, err := blobstore.Create(c, mime)
	if err != nil {
		panic("Couldn't create blobstore.Writer: " + err.Error())
	}
	io.Copy(bw, r.Body)
	r.Body.Close()
	bw.Close()

	bkey, err := bw.Key()
	if err != nil {
		panic("Couldn't get key for blob: " + err.Error())
	}

	// opts := &image.ServingURLOptions{Secure:true}
	url, err := image.ServingURL(c, bkey, nil)
	if err != nil {
		panic("Couldn't create serving URL for blob: " + err.Error())
	}

	img := &Image{BlobKey: bkey, URL: url.String()}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "IMG", nil), img)
	if err != nil {
		panic("Couldn't save to datastore: " + err.Error())
	}

	img.ID = key.IntID()

	err = json.NewEncoder(w).Encode(img)
	if err != nil {
		panic("Couldn't create JSON: " + err.Error())
	}
}
