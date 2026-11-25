package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/nisha/gravehard/internal/auth"
	"github.com/nisha/gravehard/internal/config"
	apphttp "github.com/nisha/gravehard/internal/http"
	"github.com/nisha/gravehard/internal/http/handler"
	"github.com/nisha/gravehard/internal/service"
)

func setupTestPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if os.Getenv("SKIP_INTEGRATION") == "1" {
		t.Skip("SKIP_INTEGRATION=1")
	}

	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("gravehard_test"),
		postgres.WithUsername("gravehard"),
		postgres.WithPassword("gravehard"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Skipf("postgres testcontainer unavailable: %v", err)
	}
	t.Cleanup(func() {
		_ = pgContainer.Terminate(context.Background())
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	migrationPath := filepath.Join("..", "..", "migrations", "000001_init.up.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, string(sqlBytes))
	require.NoError(t, err)

	return pool
}

func newTestRouter(pool *pgxpool.Pool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		AppEnv:          "test",
		CORSOrigins:     []string{"http://localhost:5173"},
		DBQueryTimeout:  5 * time.Second,
		JWTSecret:       "integration-test-secret",
		JWTAccessTTL:    15 * time.Minute,
		JWTRefreshTTL:   7 * 24 * time.Hour,
		RateLimitGenRPM: 600,
		RateLimitAuthPM: 100,
		RateLimitDemoPM: 300,
	}

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessTTL)
	authSvc := service.NewAuthService(pool, jwtMgr, cfg.JWTRefreshTTL, cfg.DBQueryTimeout)

	return apphttp.NewRouter(apphttp.RouterDeps{
		Config:      cfg,
		Logger:      logger,
		JWTManager:  jwtMgr,
		OpsHandler:  handler.NewOpsHandler(pool, cfg.DBQueryTimeout),
		AuthHandler: handler.NewAuthHandler(authSvc, cfg.JWTRefreshTTL, false, ""),
	})
}

func TestAuthFlowIntegration(t *testing.T) {
	pool := setupTestPostgres(t)
	router := newTestRouter(pool)

	// 1. Register
	regBody := `{"email":"founder@example.com","password":"password123","display_name":"Ada"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var regResp handler.AuthResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &regResp))
	require.NotEmpty(t, regResp.AccessToken)
	require.Equal(t, "founder@example.com", regResp.User.Email)

	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	refreshCookie1 := cookies[0]
	require.True(t, refreshCookie1.HttpOnly)

	// 2. Duplicate register -> 409
	reqDup := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
	reqDup.Header.Set("Content-Type", "application/json")
	recDup := httptest.NewRecorder()
	router.ServeHTTP(recDup, reqDup)
	require.Equal(t, http.StatusConflict, recDup.Code)

	// 3. /auth/me with Bearer token
	reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+regResp.AccessToken)
	recMe := httptest.NewRecorder()
	router.ServeHTTP(recMe, reqMe)
	require.Equal(t, http.StatusOK, recMe.Code)

	// 4. Refresh rotation
	reqRef := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	reqRef.AddCookie(refreshCookie1)
	recRef := httptest.NewRecorder()
	router.ServeHTTP(recRef, reqRef)
	require.Equal(t, http.StatusOK, recRef.Code)
	refreshCookie2 := recRef.Result().Cookies()[0]
	require.NotEqual(t, refreshCookie1.Value, refreshCookie2.Value)

	// 5. Reuse of old refreshCookie1 -> 401 and revokes family
	reqReuse := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	reqReuse.AddCookie(refreshCookie1)
	recReuse := httptest.NewRecorder()
	router.ServeHTTP(recReuse, reqReuse)
	require.Equal(t, http.StatusUnauthorized, recReuse.Code)

	// 6. refreshCookie2 is now also revoked due to family revocation
	reqAfterReuse := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	reqAfterReuse.AddCookie(refreshCookie2)
	recAfterReuse := httptest.NewRecorder()
	router.ServeHTTP(recAfterReuse, reqAfterReuse)
	require.Equal(t, http.StatusUnauthorized, recAfterReuse.Code)

	// 7. Login works again
	loginBody := `{"email":"founder@example.com","password":"password123"}`
	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	recLogin := httptest.NewRecorder()
	router.ServeHTTP(recLogin, reqLogin)
	require.Equal(t, http.StatusOK, recLogin.Code)
}
