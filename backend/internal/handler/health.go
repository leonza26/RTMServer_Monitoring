package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/leonza26/RTMServer_Monitoring/backend/internal/response"
)

const healthCheckTimeout = 2 * time.Second

// Pinger adalah kontrak untuk komponen yang bisa dicek koneksinya (mis. *pgxpool.Pool).
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthHandler menangani endpoint health check.
type HealthHandler struct {
	startedAt time.Time
	db        Pinger
}

// NewHealthHandler membuat HealthHandler dan mencatat waktu server mulai berjalan.
func NewHealthHandler(db Pinger) *HealthHandler {
	return &HealthHandler{
		startedAt: time.Now(),
		db:        db,
	}
}

type healthData struct {
	Status   string `json:"status"`
	App      string `json:"app"`
	Database string `json:"database"`
	Uptime   string `json:"uptime"`
}

// Check menangani GET /api/health.
func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), healthCheckTimeout)
	defer cancel()

	status := "ok"
	dbStatus := "up"
	httpStatus := http.StatusOK

	if err := h.db.Ping(ctx); err != nil {
		status = "degraded"
		dbStatus = "down"
		httpStatus = http.StatusServiceUnavailable
	}

	response.Success(c, httpStatus, healthData{
		Status:   status,
		App:      "RTMS_Monitoring",
		Database: dbStatus,
		Uptime:   time.Since(h.startedAt).Round(time.Second).String(),
	})
}
