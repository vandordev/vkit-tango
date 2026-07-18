package main

import (
	"context"
	"testing"
)

func TestRunRejectsInvalidDatabaseURL(t *testing.T) {
	if err := run(context.Background(), "not-a-postgres-url", "../../database/migrations", "up"); err == nil {
		t.Fatal("run() error = nil, want invalid database URL error")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	if err := run(context.Background(), "postgresql://database", "../../database/migrations", "unknown"); err == nil {
		t.Fatal("run() error = nil, want unknown command error")
	}
}
