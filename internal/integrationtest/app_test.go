package integrationtest

import (
	"context"
	"fmt"
	"github.com/olund/cool/internal"
	"github.com/olund/cool/internal/config"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
	)

	if err != nil {
		panic(err)
	}

	getenv := func(key string) string {
		switch key {
		case "POSTGRES_CONNECTION_STRING":
			return postgresContainer.MustConnectionString(ctx)
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

func TestApp_Authors(t *testing.T) {

	t.Run("Create authors - no name - bad request", func(t *testing.T) {
		requestBodyWithoutName := struct {
			Bio string `json:"bio"`
		}{
			Bio: "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit",
		}

		apitest.New().
			EnableNetworking(http.DefaultClient).
			Postf("http://%s:%s/author", cfg.Host, cfg.Port).
			JSON(requestBodyWithoutName).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"error":"Bad Request"}`).
			End()
	})

	t.Run("Create authors - no bio - bad request", func(t *testing.T) {
		requestBodyWithoutName := struct {
			Name string `json:"name"`
		}{
			Name: "a name",
		}

		apitest.New().
			EnableNetworking(http.DefaultClient).
			Postf("http://%s:%s/author", cfg.Host, cfg.Port).
			JSON(requestBodyWithoutName).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"error":"Bad Request"}`).
			End()
	})

	t.Run("Create one author, and then get it by id", func(t *testing.T) {
		type authorRequest struct {
			Name string `json:"name"`
			Bio  string `json:"bio"`
		}
		type authorResponse struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
			Bio  string `json:"bio"`
		}

		body := authorRequest{
			Name: "test1",
			Bio:  "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit",
		}

		createResponseBody := authorResponse{}
		apitest.New().
			EnableNetworking(http.DefaultClient).
			Postf("http://%s:%s/author", cfg.Host, cfg.Port).
			JSON(body).
			Expect(t).
			Status(http.StatusCreated).
			End().JSON(&createResponseBody)

		assert.NotEmpty(t, createResponseBody.Id)
		assert.Equal(t, createResponseBody.Name, body.Name)
		assert.Equal(t, createResponseBody.Bio, body.Bio)

		getByIdResponseBody := authorResponse{}
		apitest.New().
			EnableNetworking(http.DefaultClient).
			Getf("http://%s:%s/author/%d", cfg.Host, cfg.Port, createResponseBody.Id).
			JSON(body).
			Expect(t).
			Status(http.StatusOK).
			End().JSON(&getByIdResponseBody)

		assert.Equal(t, createResponseBody, getByIdResponseBody)
	})

}
