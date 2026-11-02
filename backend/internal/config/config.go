package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
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

// DatabaseConfig berisi konfigurasi koneksi PostgreSQL.
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// DSN membentuk connection string PostgreSQL (contoh: postgres://user:pass@host:5445/db?sslmode=disable).
// Menggunakan url.URL agar karakter spesial pada password otomatis di-escape dengan aman.
func (d DatabaseConfig) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   d.Name,
	}
	q := u.Query()
	q.Set("sslmode", d.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

// Load membaca file .env (jika ada), lalu menyusun Config dari environment variable.
func Load() (*Config, error) {

	// Di production tidak ada file .env (env diset oleh systemd), jadi "file tidak ada" bukan error.
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	appPort, err := getEnvInt("APP_PORT", 8080)
	if err != nil {
		return nil, err
	}

	dbPort, err := getEnvInt("DATABASE_PORT", 5432)
	if err != nil {
		return nil, err
	}

	dbPassword := getEnv("DATABASE_PASSWORD", "")
	if dbPassword == "" {
		return nil, errors.New("DATABASE_PASSWORD is required")
	}

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Host: getEnv("APP_HOST", "127.0.0.1"),
			Port: appPort,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "127.0.0.1"),
			Port:     dbPort,
			Name:     getEnv("DATABASE_NAME", "rtms_monitoring"),
			User:     getEnv("DATABASE_USER", "rtms"),
			Password: dbPassword,
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
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

func getEnvInt(key string, fallback int) (int, error) {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number: %w", key, err)
	}
	return v, nil
}
