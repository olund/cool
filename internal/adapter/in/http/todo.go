package http

import (
	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/ports"
	"log/slog"
	"net/http"
	"strconv"
)

func ListTodos(todoService ports.Todos) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := todoService.ListAll(r.Context())
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to list todos", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if todos == nil {
			todos = []domain.Todo{}
		}

		if err := encode(w, r, http.StatusOK, todos); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func GetTodoById(todos ports.Todos) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		todoIdStr := r.PathValue("id")
		slog.InfoContext(r.Context(), "GetTodoById", "path", r.URL.Path)
		todoId, err := strconv.ParseInt(todoIdStr, 10, 64)
		if err != nil {
			slog.ErrorContext(r.Context(), "GetTodoById", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		found, err := todos.GetById(r.Context(), todoId)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to list todos", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := encode(w, r, http.StatusOK, found); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func CreateTodo(todos ports.Todos) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		createTodoRequest, problems, err := decodeValid[domain.CreateTodoRequest](r)

		if len(problems) > 0 {
			slog.WarnContext(r.Context(), "CreateTodo", "problems", problems)
			if err := encode(w, r, http.StatusBadRequest, HttpError{Error: "Bad Request"}); err != nil {
				slog.ErrorContext(r.Context(), "CreateTodo", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		todo, err := todos.Create(r.Context(), createTodoRequest)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to create todo", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		if err := encode(w, r, http.StatusCreated, todo); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

}
