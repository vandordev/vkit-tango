package main

import (
	"context"
	"testing"
)

func TestRunRejectsInvalidDatabaseURL(t *testing.T) {
	if err := run(context.Background(), "not-a-postgres-url", "../../database/migrations"); err == nil {
		t.Fatal("run() error = nil, want invalid database URL error")
	}
}
