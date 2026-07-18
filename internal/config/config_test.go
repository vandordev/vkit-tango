package config_test

import (
	"path/filepath"
	"testing"

	"github.com/vandordev/vkit-tango/internal/config"
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

func TestLoadMigrateParsesDatabaseModule(t *testing.T) {
	loaded, err := config.LoadMigrate(config.Loader{
		Directory: filepath.Join("..", "..", "config"),
		Environment: map[string]string{
			"DATABASE_URL":              "postgresql://database",
			"REALTIME_TICKET_SECRET":    "ticket-secret",
			"REALTIME_INTERNAL_API_KEY": "internal-key",
		},
	})
	if err != nil {
		t.Fatalf("LoadMigrate() error = %v", err)
	}

	if loaded.Database.URL != "postgresql://database" {
		t.Fatalf("Database.URL = %q, want postgresql://database", loaded.Database.URL)
	}
}

func TestLoadWorkerParsesConcurrency(t *testing.T) {
	loaded, err := config.LoadWorker(config.Loader{
		Directory: filepath.Join("..", "..", "config"),
		Environment: map[string]string{
			"DATABASE_URL":              "postgresql://database",
			"REALTIME_TICKET_SECRET":    "ticket-secret",
			"REALTIME_INTERNAL_API_KEY": "internal-key",
		},
	})
	if err != nil {
		t.Fatalf("LoadWorker() error = %v", err)
	}

	if loaded.MaxWorkers != 10 {
		t.Fatalf("MaxWorkers = %d, want 10", loaded.MaxWorkers)
	}
	if loaded.Realtime.InternalAPIKey != "internal-key" {
		t.Fatalf("Realtime.InternalAPIKey = %q, want internal-key", loaded.Realtime.InternalAPIKey)
	}
}

func TestLoadSchedulerParsesWorkerConcurrency(t *testing.T) {
	loaded, err := config.LoadScheduler(config.Loader{
		Directory: filepath.Join("..", "..", "config"),
		Environment: map[string]string{
			"DATABASE_URL":              "postgresql://database",
			"REALTIME_TICKET_SECRET":    "ticket-secret",
			"REALTIME_INTERNAL_API_KEY": "internal-key",
		},
	})
	if err != nil {
		t.Fatalf("LoadScheduler() error = %v", err)
	}
	if loaded.MaxWorkers != 10 {
		t.Fatalf("MaxWorkers = %d, want 10", loaded.MaxWorkers)
	}
}
