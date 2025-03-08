package http

import (
	"github.com/olund/cool/internal/core/ports"
	"net/http"
)

func NewServer(authors ports.Authors) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", Health())

	mux.HandleFunc("GET /hello", GetHelloWorld())

	// Authors
	mux.HandleFunc("POST /author", LoggingMiddleware(CreateAuthor(authors)))
	mux.HandleFunc("GET /author/{id}", LoggingMiddleware(GetAuthorById(authors)))

	return mux
}
