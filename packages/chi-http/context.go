package chihttp

import (
	"context"
	"net/http"
)

type Context struct {
	Ctx     context.Context
	Writer  http.ResponseWriter
	Request *http.Request
}
