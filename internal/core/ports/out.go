package ports

import (
	"context"
	"github.com/olund/cool/internal/core/domain"
)

type AuthorStore interface {
	ListAuthors(ctx context.Context) ([]domain.Author, error)
	Insert(ctx context.Context, req domain.CreateAuthorRequest) (domain.Author, error)
	GetById(ctx context.Context, id int64) (domain.Author, error)
}
