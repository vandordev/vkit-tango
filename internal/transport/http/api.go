package http

import (
	"context"
	"database/sql"
	"encoding/json"
	stdhttp "net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vandordev/vkit-tango/internal/transport/http/method"
)

type statusResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Status string `json:"status"`
	} `json:"data"`
}

type statusOutput struct{ Body statusResponse }

func NewRouter(database *sql.DB) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, middleware.Timeout(30*time.Second))
	router.Get("/health", func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		writeJSON(writer, stdhttp.StatusOK, map[string]any{"success": true, "data": map[string]any{"status": "healthy", "timestamp": time.Now().UTC().Format(time.RFC3339Nano)}})
	})
	router.Get("/health/ready", func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		if err := database.PingContext(request.Context()); err != nil {
			writeJSON(writer, stdhttp.StatusServiceUnavailable, map[string]any{"success": false, "error": "NOT_READY", "message": "Database is not ready"})
			return
		}
		writeJSON(writer, stdhttp.StatusOK, map[string]any{"success": true, "data": map[string]string{"status": "ready"}})
	})
	return router
}

func NewAPI(router chi.Router) huma.API {
	config := huma.DefaultConfig("vkit-tango API", "1.0.0")
	config.CreateHooks = nil
	config.OpenAPIPath = "/api/openapi"
	config.DocsPath = ""
	api := humachi.New(router, config)
	method.GET(api, "/api/v1/status", "Get system status", []string{"system"}, func(context.Context, *struct{}) (*statusOutput, error) {
		response := statusResponse{Success: true}
		response.Data.Status = "ok"
		return &statusOutput{Body: response}, nil
	})
	return api
}

func writeJSON(writer stdhttp.ResponseWriter, status int, value any) {
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
