package http

import (
	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/ports"
	"log/slog"
	"net/http"
	"strconv"
)

func GetAuthorById(authors ports.Authors) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		authorIdStr := r.PathValue("id")
		slog.InfoContext(r.Context(), "GetAuthorsById", "path", r.URL.Path)
		authorId, err := strconv.ParseInt(authorIdStr, 10, 64)
		if err != nil {
			slog.ErrorContext(r.Context(), "GetAuthorById", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		found, err := authors.GetById(r.Context(), authorId)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to list authors", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := encode(w, r, http.StatusOK, found); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func CreateAuthor(authors ports.Authors) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		createAuthorRequest, err := decode[domain.CreateAuthorRequest](r)

		author, err := authors.Create(r.Context(), createAuthorRequest)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to create author", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		if err := encode(w, r, http.StatusCreated, author); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

}
