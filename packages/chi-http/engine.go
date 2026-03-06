package chihttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/labopase/flevance/packages/logger"
)

type engine struct {
	config *Config
	router Router
	server *http.Server
	log    logger.Logger
}

type Engine interface {
	Start() error
	Shutdown(ctx context.Context) error

	// HTTP-method routing along `pattern`

	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	PATCH(path string, handler HandlerFunc)
	HEAD(path string, handler HandlerFunc)
	OPTIONS(path string, handler HandlerFunc)
	TRACE(path string, handler HandlerFunc)
	CONNECT(path string, handler HandlerFunc)
}

func New(cfg *Config, log logger.Logger) (Engine, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	if log == nil {
		return nil, errors.New("logger is required")
	}

	cfg.applyDefaults()

	handler := chi.NewRouter()

	server := &http.Server{
		Addr:           cfg.Addr(),
		Handler:        handler,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	return &engine{
		config: cfg,
		server: server,
		log:    log,
	}, nil

}

func (e *engine) Start() error {
	e.log.Infof("Starting server on %s", e.config.Addr())
	return e.server.ListenAndServe()
}

func (e *engine) Shutdown(ctx context.Context) error {
	return e.server.Shutdown(ctx)
}

func (e *engine) GET(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodGet, path, handler)
}

func (e *engine) POST(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodPost, path, handler)
}

func (e *engine) PUT(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodPut, path, handler)
}

func (e *engine) DELETE(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodDelete, path, handler)
}

func (e *engine) PATCH(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodPatch, path, handler)
}

func (e *engine) HEAD(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodHead, path, handler)
}

func (e *engine) OPTIONS(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodOptions, path, handler)
}

func (e *engine) TRACE(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodTrace, path, handler)
}

func (e *engine) CONNECT(path string, handler HandlerFunc) {
	e.router.RegisterRoute(http.MethodConnect, path, handler)
}
