package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	transport "github.com/vandordev/vkit-fast/internal/transport/http"
)

func TestHandlerServesVersionedStatus(t *testing.T) {
	handler := transport.NewHandler(func() error { return nil })
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/status", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if got := response.Body.String(); got != "{\"success\":true,\"data\":{\"status\":\"ok\"}}\n" {
		t.Fatalf("body = %s", got)
	}
}

func TestHandlerDoesNotServeUnversionedStatus(t *testing.T) {
	handler := transport.NewHandler(func() error { return nil })
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/status", nil))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
}

func TestHandlerPublishesVersionedOpenAPI(t *testing.T) {
	handler := transport.NewHandler(func() error { return nil })
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if body := response.Body.String(); !strings.Contains(body, `"/api/v1/status"`) || !strings.Contains(body, `"operationId":"v1_get_system_status"`) {
		t.Fatalf("OpenAPI body does not contain v1 status operation: %s", body)
	}
}
