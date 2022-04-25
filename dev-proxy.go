package vueglue

import (
	"log"
	"net/http"
)

// Redirector for dev server

func (vg *VueGlue) DevServerRedirector() http.Handler {

	handler := func(w http.ResponseWriter, r *http.Request) {
		original := r.URL.Path
		prefix := "/dev/"
		if len(original) < len(prefix) || original[:len(prefix)] != prefix {
			http.NotFound(w, r)
			return
		}

		rest := original[len(prefix)-1:]
		log.Println("rest: ", rest)
		w.Header().Set("Content-Type", "application/javascript")
		http.Redirect(w, r, vg.DevServer+rest, http.StatusPermanentRedirect)
	}

	return http.HandlerFunc(handler)
}
