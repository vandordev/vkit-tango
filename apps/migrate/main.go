package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pressly/goose/v3"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/vandordev/vkit-fast/internal/config"
	"github.com/vandordev/vkit-fast/internal/platform/postgres"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	loaded, err := config.LoadMigrate(config.Loader{Directory: "config", Environment: environment()})
	if err != nil {
		log.Fatal(err)
	}
	if err := run(ctx, loaded.Database.URL, "database/migrations"); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, databaseURL, migrationsDirectory string) error {
	database, client, err := postgres.Open(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.UpContext(ctx, database, migrationsDirectory); err != nil {
		return fmt.Errorf("apply goose migrations: %w", err)
	}

	migrator, err := rivermigrate.New(riverdatabasesql.New(database), nil)
	if err != nil {
		return fmt.Errorf("create river migrator: %w", err)
	}
	if _, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil); err != nil {
		return fmt.Errorf("apply river migrations: %w", err)
	}

	return nil
}

func environment() map[string]string {
	values := make(map[string]string)
	for _, pair := range os.Environ() {
		name, value, found := strings.Cut(pair, "=")
		if found {
			values[name] = value
		}
	}
	return values
}
