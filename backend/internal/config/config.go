package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppEnv          string
	Port            string
	ShutdownTimeout time.Duration
	CORSOrigins     []string

	DatabaseURL    string
	DBMaxConns     int32
	DBMinConns     int32
	DBQueryTimeout time.Duration

	JWTSecret       string
	JWTAccessTTL    time.Duration
	JWTRefreshTTL   time.Duration
	CookieSecure    bool
	CookieDomain    string
	RateLimitGenRPM int
	RateLimitAuthPM int
	RateLimitDemoPM int
}

// Load reads configuration from environment variables with safe defaults for local development.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		ShutdownTimeout: time.Duration(getEnvInt("SHUTDOWN_TIMEOUT_SECONDS", 10)) * time.Second,
		CORSOrigins:     splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),

		DatabaseURL:    os.Getenv("DATABASE_URL"),
		DBMaxConns:     int32(getEnvInt("DB_MAX_CONNS", 15)),
		DBMinConns:     int32(getEnvInt("DB_MIN_CONNS", 2)),
		DBQueryTimeout: time.Duration(getEnvInt("DB_QUERY_TIMEOUT_SECONDS", 5)) * time.Second,

		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTAccessTTL:    time.Duration(getEnvInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
		JWTRefreshTTL:   time.Duration(getEnvInt("JWT_REFRESH_TTL_DAYS", 7)) * 24 * time.Hour,
		CookieSecure:    getEnvBool("COOKIE_SECURE", false),
		CookieDomain:    os.Getenv("COOKIE_DOMAIN"),
		RateLimitGenRPM: getEnvInt("RATE_LIMIT_GENERAL_RPM", 120),
		RateLimitAuthPM: getEnvInt("RATE_LIMIT_AUTH_RPM", 20),
		RateLimitDemoPM: getEnvInt("RATE_LIMIT_DEMO_RPM", 60),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		if cfg.AppEnv == "production" {
			return nil, fmt.Errorf("JWT_SECRET is required in production")
		}
		cfg.JWTSecret = "dev-only-insecure-secret-change-me"
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
