package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vandordev/vkit-fast/internal/config"
	"github.com/vandordev/vkit-fast/internal/platform/postgres"
	transport "github.com/vandordev/vkit-fast/internal/transport/http"
	"github.com/vandordev/vkit-fast/internal/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	loaded, err := config.LoadAPI(config.Loader{Directory: "config", Environment: environment()})
	if err != nil {
		log.Fatal(err)
	}
	database, client, err := postgres.Open(ctx, loaded.Database.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	metadata := usecase.SystemMetadataService{Client: client}
	server := newServer(fmt.Sprintf("%s:%d", loaded.HTTPAPI.Host, loaded.HTTPAPI.Port), transport.NewHandler(func() error {
		return database.PingContext(context.Background())
	}, metadata))
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("api server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("api shutdown error: %v", err)
	}
}

func newServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func environment() map[string]string {
	values := make(map[string]string)
	for _, pair := range os.Environ() {
		name, value, found := strings.Cut(pair, "=")
		if found {
			values[name] = value
		}
	}
	return values
}
