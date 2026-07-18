package main

import (
	app "github.com/vandordev/vkit-tango/internal/app"
	"go.uber.org/fx"
)

func main() { fx.New(app.SchedulerModule).Run() }
