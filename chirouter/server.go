package chirouter

import (
	"context"
	"errors"
	"net/http"
)

type engine struct {
	server *http.Server
	config *Config
	router Router
}

type Engine interface {
	Start() error
	Shutdown(ctx context.Context) error

	GET(pattern string, handler HandlerFunc)
	POST(pattern string, handler HandlerFunc)
	PUT(pattern string, handler HandlerFunc)
	DELETE(pattern string, handler HandlerFunc)
	PATCH(pattern string, handler HandlerFunc)
	HEAD(pattern string, handler HandlerFunc)
	OPTIONS(pattern string, handler HandlerFunc)
	TRACE(pattern string, handler HandlerFunc)
	CUSTOM(method string, pattern string, handler HandlerFunc)
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
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &engine{
		server: server,
		config: cfg,
		router: NewRouter(),
	}, nil
}

func (e *engine) Start() error {
	e.server.Handler = e.router
	return e.server.ListenAndServe()
}

func (e *engine) Shutdown(ctx context.Context) error {
	ctxShutdown, cancel := context.WithTimeout(ctx, e.config.ShutdownTimeout)
	defer cancel()

	return e.server.Shutdown(ctxShutdown)
}

func (e *engine) GET(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodGet, pattern, handler)
}

func (e *engine) POST(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodPost, pattern, handler)
}

func (e *engine) PUT(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodPut, pattern, handler)
}

func (e *engine) DELETE(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodDelete, pattern, handler)
}

func (e *engine) PATCH(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodPatch, pattern, handler)
}

func (e *engine) HEAD(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodHead, pattern, handler)
}

func (e *engine) OPTIONS(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodOptions, pattern, handler)
}

func (e *engine) TRACE(pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodTrace, pattern, handler)
}

func (e *engine) CUSTOM(method string, pattern string, handler HandlerFunc) {
	e.router.RegisterRoute(method, pattern, handler)
}
