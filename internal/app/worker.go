package app

import (
	"database/sql"

	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/config"
	generatedfx "github.com/vandordev/vkit-tango/internal/generated/fx"
	platformrealtime "github.com/vandordev/vkit-tango/internal/platform/realtime"
	platformriver "github.com/vandordev/vkit-tango/internal/platform/river"
	"go.uber.org/fx"
)

func NewWorkerSettings() (config.Worker, error) {
	return config.LoadWorker(config.Loader{Directory: configDirectory(), Environment: environment()})
}
func workerDatabase(settings config.Worker) config.Database { return settings.Database }
func NewPublisher(settings config.Worker) platformrealtime.Publisher {
	return platformrealtime.HTTPPublisher{BaseURL: settings.Realtime.PublicURL, APIKey: settings.Realtime.InternalAPIKey}
}
func NewWorkers() *riverqueue.Workers { return riverqueue.NewWorkers() }
func NewWorkerClient(database *sql.DB, workers *riverqueue.Workers, settings config.Worker) (*riverqueue.Client[*sql.Tx], error) {
	return platformriver.NewClient(database, workers, settings.MaxWorkers, nil)
}
func manageWorker(lifecycle fx.Lifecycle, client *riverqueue.Client[*sql.Tx]) {
	lifecycle.Append(fx.Hook{OnStart: client.Start, OnStop: client.Stop})
}

var WorkerModule = fx.Options(CommonModule, generatedfx.UsecaseModule, generatedfx.WorkerModule, fx.Provide(NewWorkerSettings, workerDatabase, NewPublisher, NewWorkers, NewWorkerClient), fx.Invoke(manageWorker))
