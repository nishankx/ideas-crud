package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nisha/gravehard/internal/config"
)

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	_, err := config.Load()
	require.Error(t, err)
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db?sslmode=disable")
	t.Setenv("APP_ENV", "development")

	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "8080", cfg.Port)
	require.Equal(t, 15*time.Minute, cfg.JWTAccessTTL)
	require.Equal(t, 7*24*time.Hour, cfg.JWTRefreshTTL)
	require.Equal(t, []string{"http://localhost:5173"}, cfg.CORSOrigins)
}
