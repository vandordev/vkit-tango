package app

import (
	"context"
	"database/sql"

	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/config"
	"github.com/vandordev/vkit-tango/internal/platform/db"
	"github.com/vandordev/vkit-tango/internal/platform/postgres"
	platformriver "github.com/vandordev/vkit-tango/internal/platform/river"
	"github.com/vandordev/vkit-tango/internal/usecase"
	"go.uber.org/fx"
)

func NewDatabase(lifecycle fx.Lifecycle, settings config.Database) (*sql.DB, *db.Client, error) {
	database, client, err := postgres.Open(context.Background(), settings.URL)
	if err != nil {
		return nil, nil, err
	}
	lifecycle.Append(fx.StopHook(func(context.Context) error { client.Close(); return database.Close() }))
	return database, client, nil
}

func NewProducer(database *sql.DB) (*riverqueue.Client[*sql.Tx], error) {
	return platformriver.NewProducer(database)
}

var CommonModule = fx.Options(fx.Provide(NewDatabase, fx.Annotate(NewProducer, fx.ResultTags(`name:"producer"`))), usecase.Module)
