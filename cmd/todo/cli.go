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

		index := -1

		for {

			prompt := promptui.Select{
				HideSelected: true,
				Label:        "Todos",
				Items:        todos,
				CursorPos:    index,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}",
					Active:   "\U0001F336 {{ .Name | cyan }} ({{if .Done }}\U00002714{{else}}\U00002716{{end}})",
					Inactive: "  {{ .Name | cyan }} ({{if .Done }}\U00002714{{else}}\U00002716{{end}})",
					Selected: "\U0001F336 {{ .Name | cyan }} {{ .Done | red }}",

					Details: `
--------- Todo ----------
{{ "Id:" | faint }}	{{ .Id }}
{{ "Name:" | faint }}	{{ .Name }}
{{ "Heat Unit:" | faint }}	{{ .Description }}
{{ "Peppers:" | faint }}	{{ .Done }}`,
				},
				Size: 10,
				Searcher: func(input string, index int) bool {
					td := todos[index]
					return strings.Contains(strings.ToLower(td.Name), strings.ToLower(input))
				},
			}

			index, _, err = prompt.Run()
			if err != nil {
				return fmt.Errorf("prompt select: %w", err)
			}

			todos[index].Done = !todos[index].Done
			// todo: update database as well.

			//_, _ = fmt.Fprintf(os.Stdout, "Selected %d: %s\n", index, td)
		}

		//for _, td := range todos {
		//	_, _ = fmt.Fprintf(os.Stdout, "Todo %d: %s\n", td.Id, td.Name)
		//}

		return nil
	}
}
