package chihttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Ctx     context.Context
	Writer  http.ResponseWriter
	Request *http.Request
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

func (c *Context) JSON(statusCode int, i interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(statusCode)

	if err := json.NewEncoder(c.Writer).Encode(i); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}
