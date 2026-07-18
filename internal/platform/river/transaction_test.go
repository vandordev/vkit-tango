package river_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/vandordev/vkit-fast/internal/platform/db"
	"github.com/vandordev/vkit-fast/internal/platform/db/systemmetadata"
	"github.com/vandordev/vkit-fast/internal/platform/postgres"
)

type transactionProbeArgs struct{}

func (transactionProbeArgs) Kind() string { return "platform.transaction_probe.v1" }

func TestEntAndRiverShareTransaction(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not configured")
	}

	ctx := context.Background()
	database, client, err := postgres.Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer client.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("SetDialect() error = %v", err)
	}
	if err := goose.UpContext(ctx, database, filepath.Join("..", "..", "..", "database", "migrations")); err != nil {
		t.Fatalf("UpContext() error = %v", err)
	}
	migrator, err := rivermigrate.New(riverdatabasesql.New(database), nil)
	if err != nil {
		t.Fatalf("rivermigrate.New() error = %v", err)
	}
	if _, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	riverClient, err := river.NewClient(riverdatabasesql.New(database), &river.Config{})
	if err != nil {
		t.Fatalf("river.NewClient() error = %v", err)
	}

	assertTransactionOutcome(t, ctx, database, client, riverClient, false)
	assertTransactionOutcome(t, ctx, database, client, riverClient, true)
}

func assertTransactionOutcome(t *testing.T, ctx context.Context, database *sql.DB, client *db.Client, riverClient *river.Client[*sql.Tx], commit bool) {
	t.Helper()

	key := fmt.Sprintf("transaction-probe-%s", uuid.NewString())
	jobCountBefore := riverJobCount(t, ctx, database)
	transaction, err := database.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}

	transactionClient := db.NewClient(db.Driver(entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: transaction})))
	if _, err := transactionClient.SystemMetadata.Create().SetID(uuid.New()).SetKey(key).SetValue(map[string]any{"commit": commit}).Save(ctx); err != nil {
		transaction.Rollback()
		t.Fatalf("SystemMetadata.Create() error = %v", err)
	}
	if _, err := riverClient.InsertTx(ctx, transaction, transactionProbeArgs{}, nil); err != nil {
		transaction.Rollback()
		t.Fatalf("InsertTx() error = %v", err)
	}

	if commit {
		if err := transaction.Commit(); err != nil {
			t.Fatalf("Commit() error = %v", err)
		}
	} else if err := transaction.Rollback(); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	metadataCount, err := client.SystemMetadata.Query().Where(systemmetadata.KeyEQ(key)).Count(ctx)
	if err != nil {
		t.Fatalf("SystemMetadata.Count() error = %v", err)
	}
	jobCount := riverJobCount(t, ctx, database)

	if commit && (metadataCount != 1 || jobCount != jobCountBefore+1) {
		t.Fatalf("committed transaction metadata=%d jobs=%d, want metadata=1 and jobs=%d", metadataCount, jobCount, jobCountBefore+1)
	}
	if !commit && (metadataCount != 0 || jobCount != jobCountBefore) {
		t.Fatalf("rolled back transaction metadata=%d jobs=%d, want metadata=0 and jobs=%d", metadataCount, jobCount, jobCountBefore)
	}
}

func riverJobCount(t *testing.T, ctx context.Context, database *sql.DB) int {
	t.Helper()

	var count int
	if err := database.QueryRowContext(ctx, `SELECT COUNT(*) FROM river_job WHERE kind = $1`, transactionProbeArgs{}.Kind()).Scan(&count); err != nil {
		t.Fatalf("river job count error = %v", err)
	}
	return count
}
