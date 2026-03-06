package echohttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/labopase/flevance/packages/logger"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type engine struct {
	echo *echo.Echo
	cfg  *Config
	log  logger.Logger
}

type Engine interface {
	Start(ctx context.Context) error
}

func New(cfg *Config, log logger.Logger) (Engine, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	if log == nil {
		return nil, errors.New("logger is required")
	}

	cfg.applyDefaults()

	e := echo.New()

	return &engine{
		echo: e,
		cfg:  cfg,
		log:  log,
	}, nil
}

func (e *engine) Start(ctx context.Context) error {
	sc := &echo.StartConfig{
		Address:         e.cfg.Addr(),
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: e.cfg.ShutdownTimeout,
		BeforeServeFunc: func(s *http.Server) error {
			e.log.Infof("Starting server on %s", e.cfg.Addr())

			s.ReadTimeout = e.cfg.ReadTimeout
			s.WriteTimeout = e.cfg.WriteTimeout
			s.IdleTimeout = e.cfg.IdleTimeout
			s.MaxHeaderBytes = e.cfg.MaxHeaderBytes

			return nil
		},
		OnShutdownError: func(err error) {
			e.log.Errorf("Server shutdown error: %v", err)
		},
	}

	return sc.Start(ctx, e.echo)
}

func (e *engine) SetupDefaultMiddleware() {
	e.echo.Use(middleware.RequestLogger())
	e.echo.Use(middleware.Recover())
}
