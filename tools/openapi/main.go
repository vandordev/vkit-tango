package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/vandordev/vkit-tango/internal/contract"
	transport "github.com/vandordev/vkit-tango/internal/transport/http"
	"github.com/vandordev/vkit-tango/internal/transport/http/handler/system_metadata"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

type openAPIMetadataSetter struct{}

func (openAPIMetadataSetter) Execute(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
	return usecase.SetSystemMetadataResult{}, nil
}

func main() {
	router := transport.NewRouter(&sql.DB{})
	api := transport.NewAPI(router)
	command := contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult](openAPIMetadataSetter{})
	system_metadata.NewSetSystemMetadataHandler(api, command).RegisterRoutes()
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))
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
