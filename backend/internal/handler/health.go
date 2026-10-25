package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler menangani endpoint health check.
type HealthHandler struct {
	startedAt time.Time
}

// NewHealthHandler membuat HealthHandler dan mencatat waktu server mulai berjalan.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{startedAt: time.Now()}
}

type healthData struct {
	Status string `json:"status"`
	App    string `json:"app"`
	Uptime string `json:"uptime"`
}

// Check menangani GET /api/health.
func (h *HealthHandler) Check(c *gin.Context) {
	Success(c, http.StatusOK, healthData{
		Status: "ok",
		App:    "RTMS_Monitoring",
		Uptime: time.Since(h.startedAt).Round(time.Second).String(),
	})
}
