package chirouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Context struct {
	Ctx    context.Context
	Writer http.ResponseWriter
	Req    *http.Request
}

func (c *Context) JSON(statusCode int, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(statusCode)
	json.NewEncoder(c.Writer).Encode(data)
}

func (c *Context) String(statusCode int, format string, values ...interface{}) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(statusCode)

	if len(values) > 0 {
		fmt.Fprintf(c.Writer, format, values...)
	} else {
		c.Writer.Write([]byte(format))
	}
}

func (c *Context) BindJSON(v interface{}) error {
	if c.Req.Body == nil {
		return errors.New("request body is empty")
	}

	return json.NewDecoder(c.Req.Body).Decode(v)
}

func (c *Context) Body() ([]byte, error) {
	return io.ReadAll(c.Req.Body)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return chi.URLParam(c.Req, key)
}

func (c *Context) Header(key string) string {
	return c.Req.Header.Get(key)
}
