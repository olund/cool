package integrationtest

import (
	"context"
	"fmt"
	"github.com/olund/cool/internal"
	"github.com/olund/cool/internal/config"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	cfg = config.Config{
		Host: "localhost",
		Port: "8080",
	}
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	app := internal.NewApp()

	getenv := func(key string) string {
		switch key {
		case "MIGRATIONS_DIR":
			return "../migrations"
		default:
			return ""
		}
	}

	<-time.After(5 * time.Second)

	go func() {
		err := app.Run(ctx, os.Stdout, getenv, cfg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "app run: %s\n", err)
			os.Exit(1)
		}
	}()

	if err := waitForReady(ctx, 5*time.Second, fmt.Sprintf("http://%s:%s/health", cfg.Host, cfg.Port)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "wait for ready failed: %s\n", err)
		os.Exit(1)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestApp_HelloWorld(t *testing.T) {

	t.Run("GET /hello", func(t *testing.T) {
		apitest.New().
			EnableNetworking(http.DefaultClient).
			Getf("http://%s:%s/hello", cfg.Host, cfg.Port).
			Expect(t).
			Body(`Hello World`).
			Status(http.StatusOK).
			End()
	})

}

func TestApp_Todo(t *testing.T) {

	t.Run("Create todo - no name - bad request", func(t *testing.T) {
		requestBodyWithoutName := struct {
			Description string `json:"description"`
		}{
			Description: "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit",
		}

		apitest.New().
			EnableNetworking(http.DefaultClient).
			Postf("http://%s:%s/todo", cfg.Host, cfg.Port).
			JSON(requestBodyWithoutName).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"error":"Bad Request"}`).
			End()
	})

	t.Run("Create todo - no description - bad request", func(t *testing.T) {
		requestBodyWithoutName := struct {
			Name string `json:"name"`
		}{
			Name: "a name",
		}

		apitest.New().
			EnableNetworking(http.DefaultClient).
			Postf("http://%s:%s/todo", cfg.Host, cfg.Port).
			JSON(requestBodyWithoutName).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"error":"Bad Request"}`).
			End()
	})

	t.Run("Create one todo, and then get it by id", func(t *testing.T) {
		type todoRequest struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		type todoResponse struct {
			Id          int64  `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Done        bool   `json:"done"`
		}

		body := todoRequest{
			Name:        "test1",
			Description: "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit",
		}

		createResponseBody := todoResponse{}
		apitest.New().
			EnableNetworking(http.DefaultClient).
			Postf("http://%s:%s/todo", cfg.Host, cfg.Port).
			JSON(body).
			Expect(t).
			Status(http.StatusCreated).
			End().JSON(&createResponseBody)

		assert.NotEmpty(t, createResponseBody.Id)
		assert.Equal(t, createResponseBody.Name, body.Name)
		assert.Equal(t, createResponseBody.Description, body.Description)
		assert.False(t, createResponseBody.Done)

		getByIdResponseBody := todoResponse{}
		apitest.New().
			EnableNetworking(http.DefaultClient).
			Getf("http://%s:%s/todo/%d", cfg.Host, cfg.Port, createResponseBody.Id).
			JSON(body).
			Expect(t).
			Status(http.StatusOK).
			End().JSON(&getByIdResponseBody)

		assert.Equal(t, createResponseBody, getByIdResponseBody)
	})

}
