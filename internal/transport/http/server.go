package http

import (
	"context"
	"errors"
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vandordev/vkit-tango/internal/config"
	"go.uber.org/fx"
)

func NewServer(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, settings config.HTTPAPI, router chi.Router) *stdhttp.Server {
	server := &stdhttp.Server{Addr: fmt.Sprintf("%s:%d", settings.Host, settings.Port), Handler: router, ReadHeaderTimeout: 5 * time.Second}
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
					_ = shutdowner.Shutdown(fx.ExitCode(1))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdown, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			return server.Shutdown(shutdown)
		},
	})
	return server
}
