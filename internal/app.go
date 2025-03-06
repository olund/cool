package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	ownhttp "github.com/olund/cool/internal/adapter/in/http"
	"github.com/olund/cool/internal/adapter/out/postgres/author"
	"github.com/olund/cool/internal/config"
	"github.com/olund/cool/internal/core/service"
	"github.com/olund/cool/internal/migrations"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type App struct {
}

func NewApp() *App {

	return &App{}
}

func (a *App) Run(ctx context.Context, w io.Writer, getenv func(string) string, config config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// DB
	connectionString := getenv("POSTGRES_CONNECTION_STRING")
	if err := migrations.Run(ctx, connectionString, getenv("MIGRATIONS_DIR")); err != nil {
		return err
	}

	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	authorDb := author.New(conn)

	authorStore := author.NewAuthorStore(authorDb)
	authorService := service.NewAuthorService(authorStore)
	
	server := ownhttp.NewServer(authorService)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: server,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			_, _ = fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}
