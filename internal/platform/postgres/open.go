package postgres

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vandordev/vkit-fast/internal/platform/db"
)

func Open(ctx context.Context, databaseURL string) (*sql.DB, *db.Client, error) {
	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}
	if err := database.PingContext(ctx); err != nil {
		database.Close()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	return database, db.NewClient(db.Driver(entsql.OpenDB("postgres", database))), nil
}
