package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labopase/flevance/chirouter"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"user,omitempty"`
}

func main() {
	cfg := chirouter.DefaultConfig()
	e, _ := chirouter.New(cfg)

	e.GET("/", func(ctx *chirouter.Context) error {
		ctx.String(http.StatusOK, "Welcome to Chirouter")
		return nil
	})

	e.GET("/user", func(ctx *chirouter.Context) error {
		ctx.JSON(http.StatusOK, Response{
			Error:   false,
			Message: "Success",
			Data: map[string]interface{}{
				"id":    1,
				"name":  "John Doe",
				"email": "[EMAIL_ADDRESS]",
				"age":   30,
			},
		})
		return nil
	})

	go func() {
		log.Print("Starting server on ", cfg.Addr())
		if err := e.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := e.Shutdown(context.Background()); err != nil {
		log.Fatal("Error shutting down HTTP server: ", err)
	}
}
