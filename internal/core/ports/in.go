package ports

import (
	"context"
	"github.com/olund/cool/internal/core/domain"
)

type Authors interface {
	ListAll(ctx context.Context) ([]domain.Author, error)
	Create(ctx context.Context, req domain.CreateAuthorRequest) (domain.Author, error)
	GetById(ctx context.Context, id int64) (domain.Author, error)
}
