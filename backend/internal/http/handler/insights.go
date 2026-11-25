package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nisha/gravehard/internal/http/middleware"
	"github.com/nisha/gravehard/internal/service"
)

type InsightsHandler struct {
	insightsSvc *service.InsightsService
}

func NewInsightsHandler(insightsSvc *service.InsightsService) *InsightsHandler {
	return &InsightsHandler{insightsSvc: insightsSvc}
}

// Summary godoc
// @Summary Get aggregate graveyard summary stats
// @Tags insights
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.InsightsSummary
// @Router /insights/summary [get]
func (h *InsightsHandler) Summary(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}
	res, err := h.insightsSvc.Summary(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute summary", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Causes godoc
// @Summary Get cause-of-death breakdown and co-occurring secondary causes
// @Tags insights
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.CausesResponse
// @Router /insights/causes [get]
func (h *InsightsHandler) Causes(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}
	res, err := h.insightsSvc.Causes(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute causes", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Timeline godoc
// @Summary Get monthly burial counts
// @Tags insights
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.TimelineResponse
// @Router /insights/timeline [get]
func (h *InsightsHandler) Timeline(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}
	res, err := h.insightsSvc.Timeline(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute timeline", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Lifespan godoc
// @Summary Get lifespan distribution buckets
// @Tags insights
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.LifespanResponse
// @Router /insights/lifespan [get]
func (h *InsightsHandler) Lifespan(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}
	res, err := h.insightsSvc.Lifespan(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute lifespan", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Validation godoc
// @Summary Get validation level vs cause-of-death matrix
// @Tags insights
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.ValidationResponse
// @Router /insights/validation [get]
func (h *InsightsHandler) Validation(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}
	res, err := h.insightsSvc.Validation(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute validation matrix", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Patterns godoc
// @Summary Get server-side plain-English failure patterns
// @Tags insights
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.PatternsResponse
// @Router /insights/patterns [get]
func (h *InsightsHandler) Patterns(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}
	res, err := h.insightsSvc.Patterns(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute patterns", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}
