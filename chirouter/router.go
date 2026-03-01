package chirouter

import "github.com/go-chi/chi/v5"

type HandlerFunc func(ctx *Context) error

type router struct {
	chi chi.Router
}
