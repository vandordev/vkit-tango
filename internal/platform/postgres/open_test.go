package postgres_test

import (
	"context"
	"testing"

	"github.com/vandordev/vkit-fast/internal/platform/postgres"
)

func TestOpenRejectsInvalidDatabaseURL(t *testing.T) {
	database, client, err := postgres.Open(context.Background(), "not-a-postgres-url")
	if err == nil {
		if client != nil {
			client.Close()
		}
		if database != nil {
			database.Close()
		}
		t.Fatal("Open() error = nil, want invalid database URL error")
	}
}
