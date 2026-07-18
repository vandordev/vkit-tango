package app

import (
	"os"
	"strings"

	"github.com/vandordev/vkit-tango/internal/config"
	generatedfx "github.com/vandordev/vkit-tango/internal/generated/fx"
	transport "github.com/vandordev/vkit-tango/internal/transport/http"
	"go.uber.org/fx"
)

func environment() map[string]string {
	values := map[string]string{}
	for _, pair := range os.Environ() {
		if key, value, ok := strings.Cut(pair, "="); ok {
			values[key] = value
		}
	}
	return values
}
func configDirectory() string {
	if _, err := os.Stat("config"); err == nil {
		return "config"
	}
	return "../../config"
}
func NewAPISettings() (config.API, error) {
	return config.LoadAPI(config.Loader{Directory: configDirectory(), Environment: environment()})
}
func apiDatabase(settings config.API) config.Database { return settings.Database }
func apiHTTP(settings config.API) config.HTTPAPI      { return settings.HTTPAPI }

var APIModule = fx.Options(CommonModule, generatedfx.UsecaseModule, transport.Module, generatedfx.HttpModule, fx.Provide(NewAPISettings, apiDatabase, apiHTTP))
