package ports

import (
	"context"
	"github.com/olund/cool/internal/core/domain"
)

type TodoStore interface {
	ListTodos(ctx context.Context) ([]domain.Todo, error)
	Insert(ctx context.Context, req domain.CreateTodoRequest) (domain.Todo, error)
	GetById(ctx context.Context, id int64) (domain.Todo, error)
}
