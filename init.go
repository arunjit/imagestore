package imagestore

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	err    error
	status int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", e.status, e.err.Error())
}

func serve(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			msg := err.Error()
			code := http.StatusInternalServerError
			if httpErr, ok := err.(HTTPError); ok {
				code = httpErr.status
			}
			http.Error(w, msg, code)
		}
	}
}

func init() {
	http.HandleFunc("/upload", serve(upload))
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You have reached: %s\n", r.URL)
}
