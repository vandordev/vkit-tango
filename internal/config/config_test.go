package config_test

import (
	"path/filepath"
	"testing"

	"github.com/vandordev/vkit-fast/internal/config"
)

func TestLoadAPIParsesSelectedModules(t *testing.T) {
	loaded, err := config.LoadAPI(config.Loader{
		Directory: filepath.Join("..", "..", "config"),
		Environment: map[string]string{
			"DATABASE_URL":              "postgresql://database",
			"REALTIME_TICKET_SECRET":    "ticket-secret",
			"REALTIME_INTERNAL_API_KEY": "internal-key",
		},
	})
	if err != nil {
		t.Fatalf("LoadAPI() error = %v", err)
	}

	if loaded.Database.URL != "postgresql://database" {
		t.Fatalf("Database.URL = %q, want postgresql://database", loaded.Database.URL)
	}
	if loaded.HTTPAPI.Port != 4101 {
		t.Fatalf("HTTPAPI.Port = %d, want 4101", loaded.HTTPAPI.Port)
	}
}
