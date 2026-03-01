package chirouter

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HandlerFunc func(ctx *Context) error

type router struct {
	chi chi.Router
}

type Router interface {
	http.Handler
	RegisterRoute(method string, pattern string, handler HandlerFunc)
}

func NewRouter() Router {
	return &router{
		chi: chi.NewRouter(),
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

func (r *router) RegisterRoute(method string, pattern string, handler HandlerFunc) {
	r.chi.MethodFunc(method, pattern, func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{
			Writer: w,
			Req:    req,
			Ctx:    req.Context(),
		}
		handler(ctx)
	})
}
