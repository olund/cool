package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/olund/cool/internal/core/ports"
	"log"
	"log/slog"
	"os"

	"github.com/olund/cool/internal/adapter/out/sqlite/todo"
	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/service"
	"github.com/olund/cool/internal/migrations"
	"github.com/urfave/cli/v2"
	_ "modernc.org/sqlite"
)

func main() {
	// DB
	db, err := sql.Open("sqlite", "todo.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := migrations.Run(context.Background(), db, "../../internal/migrations"); err != nil {
		log.Fatal(err)
	}

	todoDb := todo.New(db)

	todoStore := todo.NewTodoStore(todoDb)
	todoService := service.NewTodoService(todoStore)

	cliApp := &cli.App{
		Name:           "Todo",
		Usage:          "A CLI todo list application",
		DefaultCommand: "list",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a todo",
				Action:  AddTodo(todoService),
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List all todos",
				Action:  ListTodos(todoService),
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func AddTodo(todos ports.Todos) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		// Todo validate input

		var name, description string

		if c.NArg() == 1 {
			name = c.Args().First()
		} else if c.NArg() == 2 {
			name = c.Args().First()
			description = c.Args().Get(1)
		} else {
			return fmt.Errorf("plase provide the name of the todo")
		}

		create, err := todos.Create(c.Context, domain.CreateTodoRequest{
			Name:        name,
			Description: description,
		})

		if err != nil {
			slog.ErrorContext(c.Context, "add", slog.String("err", err.Error()))
			return fmt.Errorf("failed to create todo: %w", err)
		}

		_, _ = fmt.Fprintf(os.Stdout, "Added Todo %d: %s, %s\n", create.Id, create.Name, create.Description)
		return nil
	}
}

func ListTodos(todoService ports.Todos) func(c *cli.Context) error {
	return func(c *cli.Context) error {

		todos, err := todoService.ListAll(c.Context)

		if err != nil {
			slog.ErrorContext(c.Context, "List", slog.String("err", err.Error()))
			return fmt.Errorf("failed to list todos: %w", err)
		}

		for _, td := range todos {
			_, _ = fmt.Fprintf(os.Stdout, "Todo %d: %s\n", td.Id, td.Name)
		}

		return nil
	}
}
