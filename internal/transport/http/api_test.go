package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	transport "github.com/vandordev/vkit-fast/internal/transport/http"
	"github.com/vandordev/vkit-fast/internal/usecase"
)

type metadataSetter func(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error)

func (setter metadataSetter) Execute(ctx context.Context, input usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
	return setter(ctx, input)
}

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

func TestHandlerInvokesUsecaseForVersionedMutation(t *testing.T) {
	called := false
	handler := transport.NewHandler(func() error { return nil }, metadataSetter(func(_ context.Context, input usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
		called = input.Key == "feature" && input.Value["enabled"] == true
		return usecase.SetSystemMetadataResult{Key: input.Key}, nil
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPut, "/api/v1/system-metadata/feature", strings.NewReader(`{"value":{"enabled":true}}`)))
	if response.Code != http.StatusOK || !called {
		t.Fatalf("status=%d called=%v", response.Code, called)
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
	handler := transport.NewHandler(func() error { return nil }, metadataSetter(func(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
		return usecase.SetSystemMetadataResult{}, nil
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if body := response.Body.String(); !strings.Contains(body, `"/api/v1/status"`) || !strings.Contains(body, `"operationId":"v1_get_system_status"`) || !strings.Contains(body, `"operationId":"v1_set_system_metadata"`) {
		t.Fatalf("OpenAPI body does not contain v1 status operation: %s", body)
	}
}
