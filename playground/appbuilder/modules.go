package appbuilder

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	chihttp "github.com/labopase/flevance/packages/chi-http"
	echohttp "github.com/labopase/flevance/packages/echo-http"
	"github.com/labopase/flevance/packages/logger"
	"go.uber.org/fx"
)

var (

	// AppModule is the module for the application configuration.
	// It provides the AppConfig struct to the application.
	// It also validates the configuration to ensure that the environment and name are set.
	AppModule = fx.Module(
		"app",
		fx.Provide(func() (*AppConfig, error) {
			cfg, err := BindConfig[*AppConfig]("app")
			if err != nil {
				return nil, err
			}

			if cfg.Environment == "" {
				return nil, errors.New("environment is required")
			}

			if cfg.Name == "" {
				return nil, errors.New("name is required")
			}

			return cfg, nil
		}),
	)

	// LoggerModule is the module for the logger configuration.
	// It depends on the AppModule to get the environment configuration.
	// It also provides the logger.Logger interface to the application.
	LoggerModule = fx.Module(
		"logger",
		fx.Provide(
			func(appConfig *AppConfig) (*logger.Config, error) {
				cfg, err := BindConfig[*logger.Config]("logger")
				if err != nil {
					return nil, err
				}

				cfg.Environment = logger.Mode(appConfig.Environment)

				return cfg, nil
			},
			logger.NewZapLogger,
		),
		fx.Invoke(func(lc fx.Lifecycle, log logger.Logger) {
			lc.Append(
				fx.Hook{
					OnStop: func(_ context.Context) error {
						return log.Sync()
					},
				},
			)
		}),
	)

	ChiModule = fx.Module("chi",
		fx.Provide(
			func() (*chihttp.Config, error) {
				cfg, err := BindConfig[*chihttp.Config]("server")
				if err != nil {
					return nil, err
				}

				return cfg, nil
			},
			chihttp.New,
		),
		fx.Invoke(func(lc fx.Lifecycle, eng chihttp.Engine) {
			lc.Append(
				fx.Hook{
					OnStart: func(context.Context) error {
						go func() {
							if err := eng.Start(); err != nil {
								return
							}
						}()

						return nil
					},
					OnStop: func(context.Context) error {
						quit := make(chan os.Signal, 1)
						signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

						<-quit

						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()

						return eng.Shutdown(ctx)
					},
				},
			)
		}),
	)

	// EchoModule is the module for the echo configuration.
	// It depends on the AppModule to get the environment configuration.
	// It also provides the echohttp.Engine interface to the application.
	EchoModule = fx.Module(
		"echo",
		fx.Provide(
			func() (*echohttp.Config, error) {
				cfg, err := BindConfig[*echohttp.Config]("server")
				if err != nil {
					return nil, err
				}

				return cfg, nil
			},
			echohttp.New,
		),
		fx.Invoke(func(lc fx.Lifecycle, eng echohttp.Engine) {
			var (
				serverCtx    context.Context
				serverCancel context.CancelFunc
				serverDone   = make(chan error, 1)
			)

			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						serverCtx, serverCancel = context.WithCancel(context.Background())
						go func() {
							if err := eng.Start(serverCtx); err != nil {
								serverDone <- err
								return
							}
							serverDone <- nil
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						serverCancel()

						select {
						case err := <-serverDone:
							return err
						case <-ctx.Done():
							return ctx.Err()
						}
					},
				},
			)
		}),
	)
)
