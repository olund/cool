package http

import (
	"github.com/olund/cool/internal/core/ports"
	"net/http"
)

func NewServer(todos ports.Todos) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", Health())

	mux.HandleFunc("GET /hello", GetHelloWorld())

	// Todos
	mux.HandleFunc("GET /todo", LoggingMiddleware(ListTodos(todos)))
	mux.HandleFunc("POST /todo", LoggingMiddleware(CreateTodo(todos)))
	mux.HandleFunc("GET /todo/{id}", LoggingMiddleware(GetTodoById(todos)))

	return mux
}
