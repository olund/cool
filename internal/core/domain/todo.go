package domain

import (
	"context"
	"fmt"
)

type Todo struct {
	Id          int64
	Name        string
	Description string
	Done        bool
}

func (t Todo) String() string {
	return fmt.Sprintf("%d: %s - %s - %t", t.Id, t.Name, t.Description, t.Done)
}

type CreateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c CreateTodoRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if len(c.Name) == 0 {
		problems["Name"] = "Name cannot be empty"
	}
	if len(c.Description) == 0 {
		problems["Description"] = "Description cannot be empty"
	}
	return problems
}

type UpdateDoneRequest struct {
	Id   int64 `json:"id"`
	Done bool  `json:"done"`
}

func (c UpdateDoneRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if c.Id < 0 {
		problems["Id"] = "Id must be a positive number"
	}
	return problems
}
