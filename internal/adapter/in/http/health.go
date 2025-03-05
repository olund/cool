package http

import (
	"net/http"
)

func Health() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := encode(w, r, http.StatusOK, "ok"); err != nil {
			return
		}
	}

}
