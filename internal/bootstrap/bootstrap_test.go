package bootstrap_test

import (
	"errors"
	"testing"

	"github.com/vandordev/vkit-fast/internal/bootstrap"
)

func TestNewRejectsNilDatabase(t *testing.T) {
	_, err := bootstrap.New(bootstrap.Dependencies{})
	if !errors.Is(err, bootstrap.ErrDatabaseRequired) {
		t.Fatalf("New() error = %v, want ErrDatabaseRequired", err)
	}
}
