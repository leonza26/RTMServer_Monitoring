package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/leonza26/RTMServer_Monitoring/backend/internal/config"
	"github.com/leonza26/RTMServer_Monitoring/backend/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		// Logger belum dibuat (butuh config), jadi pakai package log standar.
		log.Fatalf("failed to load config: %v", err)
	}

	logger := newLogger(cfg.App)
	slog.SetDefault(logger)

	router, err := handler.NewRouter(cfg, logger)
	if err != nil {
		logger.Error("failed to create router", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("starting server",
		slog.String("addr", cfg.App.Addr()),
		slog.String("env", cfg.App.Env),
	)

	if err := router.Run(cfg.App.Addr()); err != nil {
		logger.Error("server stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

// newLogger membuat logger: format teks (mudah dibaca) untuk development,
// format JSON (mudah diproses mesin) untuk production.
func newLogger(app config.AppConfig) *slog.Logger {
	if app.IsProduction() {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
