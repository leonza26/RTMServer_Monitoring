package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/leonza26/RTMServer_Monitoring/backend/internal/config"
)

// NewRouter menyusun Gin engine beserta seluruh route aplikasi.
func NewRouter(cfg *config.Config) (*gin.Engine, error) {
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Belum ada reverse proxy (Nginx) di depan API, jadi jangan percaya header X-Forwarded-For.
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("set trusted proxies: %w", err)
	}

	router.NoRoute(func(c *gin.Context) {
		Error(c, http.StatusNotFound, "ROUTE_NOT_FOUND", "Route not found")
	})

	health := NewHealthHandler()

	api := router.Group("/api")
	api.GET("/health", health.Check)

	return router, nil
}
