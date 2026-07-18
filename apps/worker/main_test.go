package main

import (
	"database/sql"
	"testing"

	app "github.com/vandordev/vkit-tango/internal/app"
	"github.com/vandordev/vkit-tango/internal/config"
	"go.uber.org/fx"
)

func TestWorkerRootBuildsGeneratedRegistrars(t *testing.T) {
	fxApp := fx.New(app.WorkerModule, fx.Replace(config.Worker{MaxWorkers: 1}, &sql.DB{}))
	if err := fxApp.Err(); err != nil {
		t.Fatal(err)
	}
}
