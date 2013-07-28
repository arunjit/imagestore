package imagestore

import (
	"encoding/json"
	"errors"
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

func upload(w http.ResponseWriter, r *http.Request) error {
	mime := r.Header.Get("Content-Type")
	if mime == "" {
		return HTTPError{
			errors.New("Missing image Content-Type [header]"),
			http.StatusBadRequest,
		}
	}

	c := appstats.NewContext(r)
	defer c.Save()

	bw, err := blobstore.Create(c, mime)
	if err != nil {
		return errors.New("Couldn't create blobstore.Writer: " + err.Error())
	}
	io.Copy(bw, r.Body)
	r.Body.Close()
	bw.Close()

	bkey, err := bw.Key()
	if err != nil {
		return errors.New("Couldn't get key for blob: " + err.Error())
	}

	// opts := &image.ServingURLOptions{Secure:true}
	url, err := image.ServingURL(c, bkey, nil)
	if err != nil {
		return errors.New("Couldn't create serving URL for blob: " + err.Error())
	}

	img := &Image{BlobKey: bkey, URL: url.String()}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "IMG", nil), img)
	if err != nil {
		return errors.New("Couldn't save to datastore: " + err.Error())
	}

	img.ID = key.IntID()

	err = json.NewEncoder(w).Encode(img)
	if err != nil {
		return errors.New("Couldn't create JSON: " + err.Error())
	}
	return nil
}
