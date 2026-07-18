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

	"github.com/vandordev/vkit-tango/internal/config"
	"github.com/vandordev/vkit-tango/internal/platform/postgres"
	platformriver "github.com/vandordev/vkit-tango/internal/platform/river"
	transport "github.com/vandordev/vkit-tango/internal/transport/http"
	"github.com/vandordev/vkit-tango/internal/usecase"
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

	producer, err := platformriver.NewProducer(database)
	if err != nil {
		log.Fatal(err)
	}
	metadata := usecase.SystemMetadataService{Runner: usecase.Runner{Database: database, River: producer}}
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
