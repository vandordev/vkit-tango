package method_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/vandordev/vkit-tango/internal/transport/http/method"
)

type putInput struct {
	Body struct {
		Value string `json:"value"`
	}
}

type putOutput struct {
	Body struct {
		Value string `json:"value"`
	}
}

func TestPUTRegistersTypedRoute(t *testing.T) {
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("test", "1.0.0"))
	method.PUT(api, "/api/v1/examples/{id}", "Set example", []string{"examples"}, func(_ context.Context, input *putInput) (*putOutput, error) {
		output := &putOutput{}
		output.Body.Value = input.Body.Value
		return output, nil
	})

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPut, "/api/v1/examples/one", strings.NewReader(`{"value":"ok"}`)))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"value":"ok"`) {
		t.Fatalf("response = (%d, %s)", response.Code, response.Body.String())
	}
}

func TestPUTSetsDeterministicOperationID(t *testing.T) {
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("test", "1.0.0"))
	method.PUT(api, "/api/v1/examples/{id}", "Set example", nil, func(context.Context, *putInput) (*putOutput, error) {
		return &putOutput{}, nil
	})

	openAPI := api.OpenAPI().Paths["/api/v1/examples/{id}"]
	if openAPI == nil || openAPI.Put == nil || openAPI.Put.OperationID != "put_api_v1_examples_id" {
		t.Fatalf("operation = %#v", openAPI)
	}
}
