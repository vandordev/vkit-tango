package http_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/vandordev/vkit-tango/internal/contract"
	transport "github.com/vandordev/vkit-tango/internal/transport/http"
	"github.com/vandordev/vkit-tango/internal/transport/http/handler/system_metadata"
	"github.com/vandordev/vkit-tango/internal/usecase"
	"go.uber.org/fx"
)

type metadataSetter func(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error)

func (setter metadataSetter) Execute(ctx context.Context, input usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
	return setter(ctx, input)
}

func newHandler(t *testing.T, command contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult]) chi.Router {
	t.Helper()
	var router chi.Router
	app := fx.New(
		transport.Module,
		fx.Provide(func() *sql.DB { return &sql.DB{} }),
		fx.Provide(func() contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult] {
			return command
		}),
		fx.Provide(system_metadata.NewSetSystemMetadataHandler),
		fx.Invoke(func(handler *system_metadata.SetSystemMetadataHandler) { handler.RegisterRoutes() }),
		fx.Populate(&router),
	)
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
	return router
}

func TestHandlerServesVersionedStatus(t *testing.T) {
	handler := newHandler(t, metadataSetter(func(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
		return usecase.SetSystemMetadataResult{}, nil
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/status", nil))
	if response.Code != http.StatusOK || response.Body.String() != "{\"success\":true,\"data\":{\"status\":\"ok\"}}\n" {
		t.Fatalf("response = (%d, %s)", response.Code, response.Body.String())
	}
}

func TestHandlerInvokesUsecaseForVersionedMutation(t *testing.T) {
	called := false
	handler := newHandler(t, metadataSetter(func(_ context.Context, input usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
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
	handler := newHandler(t, metadataSetter(func(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
		return usecase.SetSystemMetadataResult{}, nil
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/status", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
}

func TestHandlerPublishesVersionedOpenAPI(t *testing.T) {
	handler := newHandler(t, metadataSetter(func(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
		return usecase.SetSystemMetadataResult{}, nil
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"/api/v1/system-metadata/{key}"`) {
		t.Fatalf("response = (%d, %s)", response.Code, response.Body.String())
	}
}
