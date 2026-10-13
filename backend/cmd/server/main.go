package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/leonza26/RTMServer_Monitoring/backend/internal/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Belum ada reverse proxy (Nginx) di depan API, jadi jangan percaya header X-Forwarded-For.
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    "RTMS_Monitoring",
		})
	})

	log.Printf("starting server on %s (env=%s)", cfg.App.Addr(), cfg.App.Env)

	if err := router.Run(cfg.App.Addr()); err != nil {
		log.Fatal(err)
	}

}
