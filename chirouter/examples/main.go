package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labopase/flevance/chirouter"
)

func main() {
	cfg := chirouter.DefaultConfig()
	engine, _ := chirouter.New(cfg)

	go func() {
		log.Print("Starting server on ", cfg.Addr())
		if err := engine.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := engine.Shutdown(context.Background()); err != nil {
		log.Fatal("Error shutting down HTTP server: ", err)
	}
}
