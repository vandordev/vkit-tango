package method_test

import (
	"testing"

	"github.com/vandordev/vkit-tango/internal/transport/http/method"
)

func TestV1PathBuildsVersionedRoute(t *testing.T) {
	path, err := method.New("v1").Path("/status")
	if err != nil {
		t.Fatalf("Path() error = %v", err)
	}
	if path != "/api/v1/status" {
		t.Fatalf("Path() = %q, want /api/v1/status", path)
	}
}

func TestPathRejectsEmbeddedAPIPrefix(t *testing.T) {
	_, err := method.New("v1").Path("/api/status")
	if err == nil {
		t.Fatal("Path() error = nil, want API prefix rejection")
	}
}

func TestPathRejectsEmbeddedVersionPrefix(t *testing.T) {
	_, err := method.New("v1").Path("/v1/status")
	if err == nil {
		t.Fatal("Path() error = nil, want version prefix rejection")
	}
}
