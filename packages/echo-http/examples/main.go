package main

import (
	"context"

	echohttp "github.com/labopase/flevance/packages/echo-http"
	"github.com/labopase/flevance/packages/logger"
	"go.uber.org/fx"
)

var (
	e   echohttp.Engine
	log logger.Logger
)

func main() {
	// StartServer()
	StartServerFx()
}

func StartServer() {
	var err error

	log, err = logger.NewZapLogger(logger.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	e, err = echohttp.New(echohttp.DefaultConfig(), log)
	if err != nil {
		log.Fatal(err)
	}

	if err := e.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func StartServerFx() {
	fx.New(
		fx.Options(
			fx.Provide(func() *logger.Config {
				return logger.DefaultConfig()
			}),
			fx.Provide(func() *echohttp.Config {
				return echohttp.DefaultConfig()
			}),
			fx.Provide(logger.NewZapLogger),
			fx.Provide(echohttp.New),
		),
		fx.Invoke(func(lc fx.Lifecycle, lg logger.Logger) {
			lc.Append(
				fx.Hook{
					OnStop: func(ctx context.Context) error {
						return lg.Sync()
					},
				},
			)
		}),
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
	).Run()
}
