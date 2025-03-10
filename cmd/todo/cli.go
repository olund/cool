package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/olund/cool/internal/adapter/out/sqlite/todo"
	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/ports"
	"github.com/olund/cool/internal/core/service"
	"github.com/olund/cool/internal/migrations"
	"github.com/urfave/cli/v2"
	"log"
	"log/slog"
	_ "modernc.org/sqlite"
	"os"
	"strings"
)

var cliApp *cli.App

var listTemplates = &promptui.SelectTemplates{
	Label:    "{{ . }}",
	Active:   "\U0001F336 {{ .Name | cyan }} ({{if .Done }}\U00002714{{else}}\U00002716{{end}})",
	Inactive: "  {{ .Name | cyan }} ({{if .Done }}\U00002714{{else}}\U00002716{{end}})",
	Selected: "\U0001F336 {{ .Name | cyan }} {{ .Done | red }}",

	Details: `
--------- Todo ----------
{{ "Id:" | faint }}	{{ .Id }}
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}
{{ "Done:" | faint }}	{{ .Done }}`,
}

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

	cliApp = &cli.App{
		Name:           "Todo",
		Usage:          "A CLI todo list application",
		DefaultCommand: "edit",
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
			{
				Name:    "edit",
				Aliases: []string{""},
				Usage:   "Edit todos",
				Action:  EditTodos(todoService),
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func AddTodo(todos ports.Todos) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return innerAddTodo(c.Context, todos)
	}
}

func innerAddTodo(ctx context.Context, todos ports.Todos) error {
	namePrompt := promptui.Prompt{
		Label: "Name",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("name cannot be empty")
			}
			if len(input) < 3 {
				return errors.New("name needs to be at least 3 characters")
			}
			return nil
		},
	}

	name, err := namePrompt.Run()
	if err != nil {
		return err
	}

	descriptionPrompt := promptui.Prompt{
		Label: "Description",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("description cannot be empty")
			}
			return nil
		},
	}

	description, err := descriptionPrompt.Run()
	if err != nil {
		return err
	}

	create, err := todos.Create(ctx, domain.CreateTodoRequest{
		Name:        name,
		Description: description,
	})

	if err != nil {
		slog.ErrorContext(ctx, "add", slog.String("err", err.Error()))
		return fmt.Errorf("failed to create todo: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Added Todo %d: %s, %s\n", create.Id, create.Name, create.Description)
	return nil
}

func ListTodos(todoService ports.Todos) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		todos, err := todoService.ListAll(c.Context)
		if err != nil {
			slog.ErrorContext(c.Context, "List", slog.String("err", err.Error()))
			return fmt.Errorf("failed to list todos: %w", err)
		}

		index := -1

		for {
			prompt := listTodoPrompt(todos, index)
			index, _, err = prompt.Run()
			if err != nil {
				return err
			}
			if err := toggleTodoDoneState(c, todos, index, todoService); err != nil {
				return err
			}
		}
	}
}

func EditTodos(todoService ports.Todos) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		todos, err := loadTodos(c.Context, todoService)
		if err != nil {
			return err
		}

		index := -1

		for {
			prompt := listTodoPrompt(todos, index)
			index, _, err = prompt.Run()
			if err != nil {
				if err.Error() == "^C" {
					choicesPrompt := promptui.Select{
						Label: "What do you want to do?",
						Items: []string{"Add", "List"},
					}

					_, result, err := choicesPrompt.Run()
					switch result {
					case "List":
						continue
					case "Add":
						if err := innerAddTodo(c.Context, todoService); err != nil {
							slog.ErrorContext(c.Context, "innerAddTodo", slog.String("err", err.Error()))
						}

						// Reload todos
						todos, err = loadTodos(c.Context, todoService)
						if err != nil {
							return err
						}

						continue
					}

					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("prompt failed: %w", err)
				}
			}

			if err := toggleTodoDoneState(c, todos, index, todoService); err != nil {
				return err
			}
		}
	}
}

func loadTodos(ctx context.Context, todoService ports.Todos) ([]domain.Todo, error) {
	todos, err := todoService.ListAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "List", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	return todos, nil
}

func toggleTodoDoneState(c *cli.Context, todos []domain.Todo, index int, todoService ports.Todos) error {
	currentTodo := todos[index]
	newDoneState := !currentTodo.Done
	if err := todoService.UpdateDone(c.Context, domain.UpdateDoneRequest{
		Id:   currentTodo.Id,
		Done: newDoneState,
	}); err != nil {
		return fmt.Errorf("failed to update done: %w", err)
	}

	currentTodo.Done = newDoneState
	todos[index] = currentTodo
	return nil
}

func listTodoPrompt(todos []domain.Todo, index int) promptui.Select {
	prompt := promptui.Select{
		HideSelected: true,
		Label:        "Todos",
		Items:        todos,
		CursorPos:    index,
		Templates:    listTemplates,
		Size:         10,
		Searcher: func(input string, index int) bool {
			td := todos[index]
			return strings.Contains(strings.ToLower(td.Name), strings.ToLower(input))
		},
	}
	return prompt
}
