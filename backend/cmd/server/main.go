package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leonza26/RTMServer_Monitoring/backend/internal/config"
	"github.com/leonza26/RTMServer_Monitoring/backend/internal/handler"
)

const shutdownTimeout = 10 * time.Second

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

	srv := &http.Server{
		Addr:              cfg.App.Addr(),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second, // lindungi dari client yang sengaja mengirim header sangat lambat
	}

	// ctx otomatis "selesai" saat program menerima CtrlC (SIGINT) atau SIGTERM (dari systemd).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Jalankan server di goroutine supaya main() bisa lanjut ke bawah dan menunggu sinyal.
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting server",
			slog.String("addr", srv.Addr),
			slog.String("env", cfg.App.Env),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Tunggu salah satu: server gagal start, ATAU ada sinyal berhenti.
	select {
	case err := <-serverErr:
		logger.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	// Kembalikan perilaku CtrlC normal: tekan sekali lagi untuk paksa keluar.
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown timed out, forcing close", slog.Any("error", err))
		_ = srv.Close()
	}

	// Phase 2: tutup koneksi database di sini.

	logger.Info("server stopped gracefully")
}

// newLogger membuat logger: format teks (mudah dibaca) untuk development,
// format JSON (mudah diproses mesin) untuk production.
func newLogger(app config.AppConfig) *slog.Logger {
	if app.IsProduction() {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
