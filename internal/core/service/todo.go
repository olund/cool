package service

import (
	"context"
	"fmt"
	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/ports"
)

var _ ports.Todos = &TodoService{}

type TodoService struct {
	todoStore ports.TodoStore
}

func NewTodoService(todoStore ports.TodoStore) *TodoService {
	return &TodoService{todoStore: todoStore}
}

func (s *TodoService) Create(ctx context.Context, req domain.CreateTodoRequest) (domain.Todo, error) {
	// todo validation.
	return s.todoStore.Insert(ctx, req)
}

func (s *TodoService) GetById(ctx context.Context, id int64) (domain.Todo, error) {
	return s.todoStore.GetById(ctx, id)
}

func (s *TodoService) ListAll(ctx context.Context) ([]domain.Todo, error) {
	todos, err := s.todoStore.ListTodos(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	return todos, nil
}
