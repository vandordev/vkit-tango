package config_test

import (
	"path/filepath"
	"testing"

	"github.com/vandordev/vkit-tango/internal/config"
)

func TestLoaderResolvesRequiredAndDefaultValues(t *testing.T) {
	loader := config.Loader{
		Directory: filepath.Join("..", "..", "config", "testdata"),
		Environment: map[string]string{
			"DATABASE_URL": "postgresql://database",
		},
	}

	loaded, err := loader.Load("required", "defaults")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	database := loaded["database"].(map[string]any)
	if database["url"] != "postgresql://database" {
		t.Fatalf("database.url = %v, want postgresql://database", database["url"])
	}

	httpAPI := loaded["http_api"].(map[string]any)
	if httpAPI["port"] != "4101" {
		t.Fatalf("http_api.port = %v, want 4101", httpAPI["port"])
	}
}

func TestLoaderRejectsMissingRequiredValue(t *testing.T) {
	loader := config.Loader{Directory: filepath.Join("..", "..", "config", "testdata")}

	_, err := loader.Load("required")
	if err == nil || err.Error() != `missing configuration environment variable "DATABASE_URL" in module "required"` {
		t.Fatalf("Load() error = %v, want missing DATABASE_URL error", err)
	}
}
