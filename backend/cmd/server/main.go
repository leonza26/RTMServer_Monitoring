package main

import (
	"log"

	"github.com/leonza26/RTMServer_Monitoring/backend/internal/config"
	"github.com/leonza26/RTMServer_Monitoring/backend/internal/handler"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("starting server on %s (env=%s)", cfg.App.Addr(), cfg.App.Env)

	router, err := handler.NewRouter(cfg)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	if err := router.Run(cfg.App.Addr()); err != nil {
		log.Fatal(err)
	}

}
