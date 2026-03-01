package chirouter

import (
	"context"
	"errors"
	"net/http"
)

type engine struct {
	server *http.Server
	config *Config
}

type Engine interface {
	Start() error
	Shutdown(ctx context.Context) error
}

func New(cfg *Config) (Engine, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	cfg.applyDefaults()

	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      nil,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &engine{
		server: server,
		config: cfg,
	}, nil
}

func (e *engine) Start() error {
	return e.server.ListenAndServe()
}

func (e *engine) Shutdown(ctx context.Context) error {
	ctxShutdown, cancel := context.WithTimeout(ctx, e.config.ShutdownTimeout)
	defer cancel()

	return e.server.Shutdown(ctxShutdown)
}
