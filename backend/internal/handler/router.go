package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/leonza26/RTMServer_Monitoring/backend/internal/config"
	"github.com/leonza26/RTMServer_Monitoring/backend/internal/middleware"
	"github.com/leonza26/RTMServer_Monitoring/backend/internal/response"
)

// NewRouter menyusun Gin engine beserta seluruh route aplikasi.
func NewRouter(cfg *config.Config, logger *slog.Logger) (*gin.Engine, error) {
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// gin.New() = router kosong. gin.Default() otomatis memasang Logger & Recovery milik Gin,
	// sekarang kita pasang versi sendiri.
	router := gin.New()

	// Urutan penting: RequestID dulu agar Logger bisa membaca ID-nya,
	// Recovery paling dalam agar panic tetap tercatat oleh Logger sebagai status 500.
	router.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		gin.CustomRecovery(func(c *gin.Context, _ any) {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}),
	)

	// Belum ada reverse proxy (Nginx) di depan API, jadi jangan percaya header X-Forwarded-For.
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("set trusted proxies: %w", err)
	}

	router.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "ROUTE_NOT_FOUND", "Route not found")
	})

	health := NewHealthHandler()

	api := router.Group("/api")

	api.GET("/health", health.Check)

	return router, nil
}
