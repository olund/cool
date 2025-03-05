package integrationtest

import (
	"context"
	"fmt"
	"github.com/olund/cool/internal"
	"github.com/olund/cool/internal/config"
	"github.com/steinfletcher/apitest"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log/slog"
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

	slog.InfoContext(ctx, "Postgres container started", "connection string", postgresContainer.MustConnectionString(ctx))

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

// waitForReady calls the specified endpoint until it gets a 200
// response or until the context is cancelled or the timeout is
// reached.
func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			slog.DebugContext(ctx, "Endpoint is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait a little while between checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}
