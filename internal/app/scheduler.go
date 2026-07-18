package app

import (
	"database/sql"

	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/config"
	generatedfx "github.com/vandordev/vkit-tango/internal/generated/fx"
	platformriver "github.com/vandordev/vkit-tango/internal/platform/river"
	schedulerriver "github.com/vandordev/vkit-tango/internal/scheduler/river"
	"go.uber.org/fx"
)

func NewSchedulerSettings() (config.Scheduler, error) {
	return config.LoadScheduler(config.Loader{Directory: configDirectory(), Environment: environment()})
}
func schedulerDatabase(settings config.Scheduler) config.Database { return settings.Database }
func NewPeriodicJobs(registrar *schedulerriver.BaselineRegistrar) []*riverqueue.PeriodicJob {
	return registrar.RegisterPeriodicJobs()
}
func NewSchedulerClient(database *sql.DB, settings config.Scheduler, jobs []*riverqueue.PeriodicJob) (*riverqueue.Client[*sql.Tx], error) {
	return platformriver.NewClient(database, riverqueue.NewWorkers(), settings.MaxWorkers, jobs)
}
func manageScheduler(lifecycle fx.Lifecycle, client *riverqueue.Client[*sql.Tx]) {
	lifecycle.Append(fx.Hook{OnStart: client.Start, OnStop: client.Stop})
}

var SchedulerModule = fx.Options(CommonModule, generatedfx.UsecaseModule, generatedfx.SchedulerModule, fx.Provide(NewSchedulerSettings, schedulerDatabase, NewPeriodicJobs, NewSchedulerClient), fx.Invoke(manageScheduler))
