package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/vandordev/vkit-tango/internal/transport/http/method"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

type statusResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Status string `json:"status"`
	} `json:"data"`
}

type statusOutput struct {
	Body statusResponse
}

type setSystemMetadataInput struct {
	Key  string `path:"key" minLength:"1"`
	Body struct {
		Value map[string]any `json:"value"`
	}
}

type setSystemMetadataOutput struct {
	Body struct {
		Success bool `json:"success"`
		Data    struct {
			Key string `json:"key"`
		} `json:"data"`
	}
}

func NewHandler(ready func() error, mutations ...usecase.SetSystemMetadata) http.Handler {
	mux := http.NewServeMux()
	config := huma.DefaultConfig("vkit-tango API", "1.0.0")
	config.CreateHooks = nil
	config.OpenAPIPath = "/api/openapi"
	config.DocsPath = "/api/docs"
	api := humago.New(mux, config)

	v1 := method.New("v1")
	statusPath, err := v1.Path("/status")
	if err != nil {
		panic(err)
	}
	huma.Register(api, huma.Operation{
		OperationID: "v1_get_system_status",
		Method:      http.MethodGet,
		Path:        statusPath,
		Summary:     "Get system status",
	}, func(context.Context, *struct{}) (*statusOutput, error) {
		response := statusResponse{Success: true}
		response.Data.Status = "ok"
		return &statusOutput{Body: response}, nil
	})
	if len(mutations) == 1 && mutations[0] != nil {
		metadataPath, err := v1.Path("/system-metadata/{key}")
		if err != nil {
			panic(err)
		}
		huma.Register(api, huma.Operation{OperationID: "v1_set_system_metadata", Method: http.MethodPut, Path: metadataPath, Summary: "Set system metadata"}, func(ctx context.Context, input *setSystemMetadataInput) (*setSystemMetadataOutput, error) {
			result, err := mutations[0].Execute(ctx, usecase.SetSystemMetadataInput{Key: input.Key, Value: input.Body.Value})
			if err != nil {
				return nil, err
			}
			output := &setSystemMetadataOutput{}
			output.Body.Success = true
			output.Body.Data.Key = result.Key
			return output, nil
		})
	}

	mux.HandleFunc("GET /health", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"status":    "healthy",
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			},
		})
	})
	mux.HandleFunc("GET /health/ready", func(writer http.ResponseWriter, request *http.Request) {
		if err := ready(); err != nil {
			writeJSON(writer, http.StatusServiceUnavailable, map[string]any{
				"success": false,
				"error":   "NOT_READY",
				"message": "Database is not ready",
			})
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{
			"success": true,
			"data":    map[string]string{"status": "ready"},
		})
	})

	return mux
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
