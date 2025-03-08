package ports

import (
	"context"
	"github.com/olund/cool/internal/core/domain"
)

type Todos interface {
	ListAll(ctx context.Context) ([]domain.Todo, error)
	Create(ctx context.Context, req domain.CreateTodoRequest) (domain.Todo, error)
	GetById(ctx context.Context, id int64) (domain.Todo, error)
}
