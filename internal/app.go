package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	ownhttp "github.com/olund/cool/internal/adapter/in/http"
	"github.com/olund/cool/internal/adapter/out/sqlite/todo"
	"github.com/olund/cool/internal/config"
	"github.com/olund/cool/internal/core/service"
	"github.com/olund/cool/internal/migrations"
	_ "modernc.org/sqlite"
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
	dbName := getenv("DB_NAME")
	if dbName == "" {
		dbName = ":memory:"
	}

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return err
	}

	if err := migrations.Run(ctx, db, getenv("MIGRATIONS_DIR")); err != nil {
		return err
	}

	todoDb := todo.New(db)

	todoStore := todo.NewTodoStore(todoDb)
	todoService := service.NewTodoService(todoStore)

	server := ownhttp.NewServer(todoService)

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
