package imagestore

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You have reached: %s\n", r.URL)
}
