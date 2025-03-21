package todo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/ports"
)

var _ ports.TodoStore = &TodoStore{}

type TodoStore struct {
	queries *Queries
}

func NewTodoStore(queries *Queries) *TodoStore {
	return &TodoStore{
		queries: queries,
	}
}

func (s *TodoStore) Insert(ctx context.Context, request domain.CreateTodoRequest) (domain.Todo, error) {
	created, err := s.queries.CreateTodo(ctx, CreateTodoParams{
		Name: request.Name,
		Description: sql.NullString{
			String: request.Description,
			Valid:  true,
		},
	})
	if err != nil {
		return domain.Todo{}, err
	}
	return ToTodo(created), err
}

func (s *TodoStore) GetById(ctx context.Context, id int64) (domain.Todo, error) {
	todo, err := s.queries.GetTodo(ctx, id)
	if err != nil {
		return domain.Todo{}, err
	}
	return ToTodo(todo), err
}

func (s *TodoStore) ListTodos(ctx context.Context) ([]domain.Todo, error) {
	todos, err := s.queries.ListTodos(ctx)
	if err != nil {
		return nil, err
	}

	var ret []domain.Todo
	for _, todo := range todos {
		ret = append(ret, ToTodo(todo))
	}
	return ret, nil
}

func (s *TodoStore) UpdateDone(ctx context.Context, request domain.UpdateDoneRequest) error {
	err := s.queries.UpdateTodoDoneState(ctx, UpdateTodoDoneStateParams{
		Done: sql.NullBool{
			Bool:  request.Done,
			Valid: true,
		},
		ID: request.Id,
	})

	if err != nil {
		return fmt.Errorf("failed persist done state for todo: %d: %w", request.Id, err)
	}

	return nil
}

func ToTodo(todo Todo) domain.Todo {
	return domain.Todo{
		Id:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description.String,
		Done:        todo.Done.Bool,
	}
}
