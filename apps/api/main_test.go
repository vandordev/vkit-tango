package main

import (
	"database/sql"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	app "github.com/vandordev/vkit-tango/internal/app"
	"go.uber.org/fx"
)

func TestAPIRootBuildsGeneratedRoutes(t *testing.T) {
	fxApp := fx.New(app.APIModule, fx.Replace(&sql.DB{}), fx.Invoke(func(api huma.API) {
		if api.OpenAPI().Paths["/api/v1/system-metadata/{key}"] == nil {
			t.Fatal("metadata route is missing")
		}
	}))
	if err := fxApp.Err(); err != nil {
		t.Fatal(err)
	}
}
