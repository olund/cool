package http

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "HTTP", "method", r.Method, "path", r.URL.Path)
		next(w, r)
	}
}
