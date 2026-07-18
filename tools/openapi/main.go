package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	transport "github.com/vandordev/vkit-tango/internal/transport/http"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

type openAPIMetadataSetter struct{}

func (openAPIMetadataSetter) Execute(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
	return usecase.SetSystemMetadataResult{}, nil
}

func main() {
	handler := transport.NewHandler(func() error { return nil }, openAPIMetadataSetter{})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))
	if response.Code != http.StatusOK {
		panic(fmt.Sprintf("OpenAPI export returned HTTP %d", response.Code))
	}

	path := filepath.Join("contracts", "openapi", "openapi.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, response.Body.Bytes(), 0o644); err != nil {
		panic(err)
	}
}
