package chihttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HandlerFunc func(ctx *Context)

type router struct {
	chi chi.Router
}

type Router interface {
	http.Handler
	RegisterRoute(method string, path string, handler HandlerFunc)
}

func NewRouter() Router {
	return &router{
		chi: chi.NewRouter(),
	}
}

func (r *router) RegisterRoute(method string, path string, handler HandlerFunc) {
	r.chi.MethodFunc(method, path, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			Writer:  w,
			Request: r,
			Ctx:     r.Context(),
		}
		handler(ctx)
	})
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}
