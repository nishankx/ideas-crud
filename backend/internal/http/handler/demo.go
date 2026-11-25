package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nisha/gravehard/internal/repository"
	"github.com/nisha/gravehard/internal/service"
)

type DemoHandler struct {
	users        *repository.UserRepository
	ideaSvc      *service.IdeaService
	insightsSvc  *service.InsightsService
	queryTimeout time.Duration
}

func NewDemoHandler(
	pool *pgxpool.Pool,
	ideaSvc *service.IdeaService,
	insightsSvc *service.InsightsService,
	queryTimeout time.Duration,
) *DemoHandler {
	return &DemoHandler{
		users:        repository.NewUserRepository(pool),
		ideaSvc:      ideaSvc,
		insightsSvc:  insightsSvc,
		queryTimeout: queryTimeout,
	}
}

func (h *DemoHandler) demoUserID(c *gin.Context) (uuid.UUID, bool) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.queryTimeout)
	defer cancel()

	u, err := h.users.GetDemoUser(ctx)
	if err != nil {
		RespondError(c, http.StatusNotFound, "DEMO_NOT_SEEDED", "demo graveyard has not been seeded yet", nil)
		return uuid.Nil, false
	}
	return u.ID, true
}

// ListIdeas godoc
// @Summary List public demo graveyard ideas
// @Tags demo
// @Produce json
// @Success 200 {object} domain.PaginatedIdeas
// @Router /demo/ideas [get]
func (h *DemoHandler) ListIdeas(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	filter, err := ParseIdeaFilterFromQuery(c)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	res, err := h.ideaSvc.List(c.Request.Context(), uid, filter)
	if err != nil {
		handleIdeaServiceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetIdea godoc
// @Summary Get a single public demo idea
// @Tags demo
// @Produce json
// @Param id path string true "Idea UUID"
// @Success 200 {object} domain.Idea
// @Router /demo/ideas/{id} [get]
func (h *DemoHandler) GetIdea(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "idea not found", nil)
		return
	}
	idea, err := h.ideaSvc.Get(c.Request.Context(), uid, ideaID)
	if err != nil {
		handleIdeaServiceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, idea)
}

func (h *DemoHandler) Summary(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	res, err := h.insightsSvc.Summary(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute demo summary", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DemoHandler) Causes(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	res, err := h.insightsSvc.Causes(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute demo causes", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DemoHandler) Timeline(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	res, err := h.insightsSvc.Timeline(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute demo timeline", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DemoHandler) Lifespan(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	res, err := h.insightsSvc.Lifespan(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute demo lifespan", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DemoHandler) Validation(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	res, err := h.insightsSvc.Validation(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute demo validation", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DemoHandler) Patterns(c *gin.Context) {
	uid, ok := h.demoUserID(c)
	if !ok {
		return
	}
	res, err := h.insightsSvc.Patterns(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to compute demo patterns", nil)
		return
	}
	c.JSON(http.StatusOK, res)
}
