package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menampung seluruh konfigurasi aplikasi.
type Config struct {
	App AppConfig
}

// AppConfig berisi konfigurasi dasar HTTP server.
type AppConfig struct {
	Env  string
	Host string
	Port int
}

// Addr mengembalikan alamat listen, contoh "127.0.0.1:8080".
func (a AppConfig) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

// IsProduction bernilai true jika APP_ENV=production.
func (a AppConfig) IsProduction() bool {
	return a.Env == "production"
}

// Load membaca file .env (jika ada), lalu menyusun Config dari environment variable.
func Load() (*Config, error) {
	// Di production tidak ada file .env (env diset oleh systemd), jadi "file tidak ada" bukan error.
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}
	port, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("APP_PORT must be a number: %w", err)
	}
	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Host: getEnv("APP_HOST", "127.0.0.1"),
			Port: port,
		},
	}
	return cfg, nil
}

// getEnv mengambil environment variable, atau fallback jika kosong / tidak ada.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
