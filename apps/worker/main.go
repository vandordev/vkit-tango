package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vandordev/vkit-tango/internal/config"
	"github.com/vandordev/vkit-tango/internal/platform/postgres"
	platformrealtime "github.com/vandordev/vkit-tango/internal/platform/realtime"
	platformriver "github.com/vandordev/vkit-tango/internal/platform/river"
	"github.com/vandordev/vkit-tango/internal/usecase"
	workerriver "github.com/vandordev/vkit-tango/internal/worker/river"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	loaded, err := config.LoadWorker(config.Loader{Directory: "config", Environment: environment()})
	if err != nil {
		log.Fatal(err)
	}
	database, client, err := postgres.Open(ctx, loaded.Database.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	publisher := platformrealtime.HTTPPublisher{BaseURL: loaded.Realtime.PublicURL, APIKey: loaded.Realtime.InternalAPIKey}
	producer, err := platformriver.NewProducer(database)
	if err != nil {
		log.Fatal(err)
	}
	metadata := usecase.SystemMetadataService{Runner: usecase.Runner{Database: database, River: producer}}
	workers, err := workerriver.RegisterWorkers(publisher, metadata)
	if err != nil {
		log.Fatal(err)
	}
	riverClient, err := platformriver.NewClient(database, workers, loaded.MaxWorkers, workerriver.RegisterPeriodicJobs())
	if err != nil {
		log.Fatal(err)
	}
	if err := riverClient.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := riverClient.Stop(shutdownContext); err != nil {
		log.Printf("worker shutdown error: %v", err)
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
