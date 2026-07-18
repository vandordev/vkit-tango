package app

import (
	"database/sql"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/vandordev/vkit-tango/internal/config"
	"go.uber.org/fx"
)

func TestAPIModuleBuilds(t *testing.T) {
	app := fx.New(APIModule, fx.Replace(&sql.DB{}), fx.Invoke(func(api huma.API) {
		if api.OpenAPI().Paths["/api/v1/system-metadata/{key}"] == nil {
			t.Fatal("metadata route was not registered")
		}
	}))
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestWorkerModuleBuilds(t *testing.T) {
	app := fx.New(WorkerModule, fx.Replace(config.Worker{MaxWorkers: 1}, &sql.DB{}))
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestSchedulerModuleBuilds(t *testing.T) {
	app := fx.New(SchedulerModule, fx.Replace(config.Scheduler{MaxWorkers: 1}, &sql.DB{}))
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
}
