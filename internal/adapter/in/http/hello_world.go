package http

import (
	"log/slog"
	"net/http"
)

func GetHelloWorld() http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// use thing to handle request
			slog.InfoContext(r.Context(), "msg", "handleSomething", "val")
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`Hello World`))
		},
	)

}
