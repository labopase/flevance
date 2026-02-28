package main

import (
	"log"

	"github.com/labopase/flevance/logger"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	log, err := logger.NewZapLogger(&logger.Config{
		Level:        logger.DebugLevel,
		EnableCaller: false,
		EnableTrace:  false,
		Mode:         logger.ModeDevelopment,
	})

	if err != nil {
		panic(err)
	}

	defer log.Sync()

	log.Info("Info")
	log.Debug("Debug")
	log.Warn("Warn")
	log.Error("Error")
	// log.Fatal("Fatal")

	log.Infow("Info with fields", logger.String("key", "value"))
	log.Debugw("Debug with fields", logger.String("key", "value"))
	log.Warnw("Warn with fields", logger.String("key", "value"))
	log.Errorw("Error with fields", logger.String("key", "value"))
	// log.Fatalw("Fatal with fields", logger.String("key", "value"))

	log.Infof("Info with fields %s", "value")
	log.Debugf("Debug with fields %s", "value")
	log.Warnf("Warn with fields %s", "value")
	log.Errorf("Error with fields %s", "value")
	// log.Fatalf("Fatal with fields %s", "value")

}
