package http

import "net/http"

func NewServer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", Health())
	mux.HandleFunc("GET /hello", GetHelloWorld())

	return mux
}
