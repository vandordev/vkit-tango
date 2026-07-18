package usecase

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/platform/db"
)

type Transaction struct {
	Ent *db.Client
	SQL *sql.Tx
}
type Runner struct {
	Database *sql.DB
	River    *riverqueue.Client[*sql.Tx]
}

func NewRunner(database *sql.DB, river *riverqueue.Client[*sql.Tx]) Runner {
	return Runner{Database: database, River: river}
}

func (runner Runner) WithinTransaction(ctx context.Context, fn func(context.Context, Transaction) error) error {
	tx, err := runner.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ent := db.NewClient(db.Driver(entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: tx})))
	if err := fn(ctx, Transaction{Ent: ent, SQL: tx}); err != nil {
		return err
	}
	return tx.Commit()
}
