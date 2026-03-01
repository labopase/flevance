package chirouter

import (
	"context"
	"net/http"
)

type Context struct {
	ctx    context.Context
	writer http.ResponseWriter
	req    *http.Request
}
