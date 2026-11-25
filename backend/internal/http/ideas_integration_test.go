package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/nisha/gravehard/internal/auth"
	"github.com/nisha/gravehard/internal/config"
	"github.com/nisha/gravehard/internal/domain"
	apphttp "github.com/nisha/gravehard/internal/http"
	"github.com/nisha/gravehard/internal/http/handler"
	"github.com/nisha/gravehard/internal/service"
)

func newFullTestRouter(pool *pgxpool.Pool) *gin.Engine {
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
	ideaSvc := service.NewIdeaService(pool, cfg.DBQueryTimeout)

	return apphttp.NewRouter(apphttp.RouterDeps{
		Config:      cfg,
		Logger:      logger,
		JWTManager:  jwtMgr,
		OpsHandler:  handler.NewOpsHandler(pool, cfg.DBQueryTimeout),
		AuthHandler: handler.NewAuthHandler(authSvc, cfg.JWTRefreshTTL, false, ""),
		IdeaHandler: handler.NewIdeaHandler(ideaSvc),
	})
}

func registerUserAndToken(t *testing.T, router *gin.Engine, email string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"email":        email,
		"password":     "password123",
		"display_name": "Founder",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp handler.AuthResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	return resp.AccessToken
}

func TestIdeasCRUDIntegration(t *testing.T) {
	pool := setupTestPostgres(t)
	router := newFullTestRouter(pool)

	tokenA := registerUserAndToken(t, router, "a@example.com")
	tokenB := registerUserAndToken(t, router, "b@example.com")

	// 1. Create Idea
	createBody := `{
		"name": "UberForCats",
		"tagline": "On-demand feline rideshare",
		"category": "marketplace",
		"born_at": "2024-01-10",
		"died_at": "2024-03-15",
		"primary_cause": "no_market_need",
		"secondary_causes": ["lost_motivation"],
		"epitaph": "Cats refused to get in the Uber.",
		"autopsy": "Talk to cat owners before building dispatch software.",
		"validation_level": "none",
		"money_spent_cents": 45000,
		"hours_invested": 120,
		"would_resurrect": false
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ideas", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var created domain.Idea
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
	require.Equal(t, domain.StatusBuried, created.Status)

	// 2. User B gets 404 (not 403) for User A's idea
	reqOther := httptest.NewRequest(http.MethodGet, "/api/v1/ideas/"+created.ID.String(), nil)
	reqOther.Header.Set("Authorization", "Bearer "+tokenB)
	recOther := httptest.NewRecorder()
	router.ServeHTTP(recOther, reqOther)
	require.Equal(t, http.StatusNotFound, recOther.Code)

	// 3. Full-text search + filter
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/ideas?q=feline&cause=no_market_need", nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenA)
	recList := httptest.NewRecorder()
	router.ServeHTTP(recList, reqList)
	require.Equal(t, http.StatusOK, recList.Code)

	var listResp domain.PaginatedIdeas
	require.NoError(t, json.Unmarshal(recList.Body.Bytes(), &listResp))
	require.Equal(t, int64(1), listResp.Total)

	// 4. Resurrect: buried -> on_life_support -> resurrected -> 409
	for _, expected := range []domain.IdeaStatus{domain.StatusOnLifeSupport, domain.StatusResurrected} {
		reqRes := httptest.NewRequest(http.MethodPost, "/api/v1/ideas/"+created.ID.String()+"/resurrect", nil)
		reqRes.Header.Set("Authorization", "Bearer "+tokenA)
		recRes := httptest.NewRecorder()
		router.ServeHTTP(recRes, reqRes)
		require.Equal(t, http.StatusOK, recRes.Code)

		var updated domain.Idea
		require.NoError(t, json.Unmarshal(recRes.Body.Bytes(), &updated))
		require.Equal(t, expected, updated.Status)
	}

	reqResAgain := httptest.NewRequest(http.MethodPost, "/api/v1/ideas/"+created.ID.String()+"/resurrect", nil)
	reqResAgain.Header.Set("Authorization", "Bearer "+tokenA)
	recResAgain := httptest.NewRecorder()
	router.ServeHTTP(recResAgain, reqResAgain)
	require.Equal(t, http.StatusConflict, recResAgain.Code)

	// 5. Soft delete -> 204, then GET -> 404
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/ideas/"+created.ID.String(), nil)
	reqDel.Header.Set("Authorization", "Bearer "+tokenA)
	recDel := httptest.NewRecorder()
	router.ServeHTTP(recDel, reqDel)
	require.Equal(t, http.StatusNoContent, recDel.Code)

	reqGetAfter := httptest.NewRequest(http.MethodGet, "/api/v1/ideas/"+created.ID.String(), nil)
	reqGetAfter.Header.Set("Authorization", "Bearer "+tokenA)
	recGetAfter := httptest.NewRecorder()
	router.ServeHTTP(recGetAfter, reqGetAfter)
	require.Equal(t, http.StatusNotFound, recGetAfter.Code)
}
