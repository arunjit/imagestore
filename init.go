package imagestore

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You have reached: %s\n", r.URL)
}
