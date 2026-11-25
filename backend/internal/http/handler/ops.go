package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OpsHandler struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

func NewOpsHandler(pool *pgxpool.Pool, timeout time.Duration) *OpsHandler {
	return &OpsHandler{pool: pool, timeout: timeout}
}

// Healthz godoc
// @Summary Liveness probe
// @Tags ops
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (h *OpsHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Readyz godoc
// @Summary Readiness probe
// @Tags ops
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} ErrorEnvelope
// @Router /readyz [get]
func (h *OpsHandler) Readyz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		RespondError(c, http.StatusServiceUnavailable, "DB_UNAVAILABLE", "database is not ready", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
