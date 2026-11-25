package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/http/middleware"
	"github.com/nisha/gravehard/internal/service"
)

type IdeaHandler struct {
	ideaSvc *service.IdeaService
}

func NewIdeaHandler(ideaSvc *service.IdeaService) *IdeaHandler {
	return &IdeaHandler{ideaSvc: ideaSvc}
}

type CreateIdeaRequest struct {
	Name            string                 `json:"name" binding:"required,min=1,max=120"`
	Tagline         string                 `json:"tagline" binding:"max=200"`
	Description     string                 `json:"description" binding:"max=5000"`
	Category        domain.IdeaCategory    `json:"category" binding:"required"`
	BornAt          string                 `json:"born_at" binding:"required"`
	DiedAt          string                 `json:"died_at" binding:"required"`
	PrimaryCause    domain.DeathCause      `json:"primary_cause" binding:"required"`
	SecondaryCauses []domain.DeathCause    `json:"secondary_causes"`
	Epitaph         string                 `json:"epitaph" binding:"required,min=1,max=140"`
	Autopsy         string                 `json:"autopsy" binding:"max=5000"`
	ValidationLevel domain.ValidationLevel `json:"validation_level"`
	MoneySpentCents int64                  `json:"money_spent_cents" binding:"min=0"`
	HoursInvested   int32                  `json:"hours_invested" binding:"min=0"`
	WouldResurrect  bool                   `json:"would_resurrect"`
}

type UpdateIdeaRequest struct {
	Name            *string                 `json:"name" binding:"omitempty,min=1,max=120"`
	Tagline         *string                 `json:"tagline" binding:"omitempty,max=200"`
	Description     *string                 `json:"description" binding:"omitempty,max=5000"`
	Category        *domain.IdeaCategory    `json:"category"`
	BornAt          *string                 `json:"born_at"`
	DiedAt          *string                 `json:"died_at"`
	PrimaryCause    *domain.DeathCause      `json:"primary_cause"`
	SecondaryCauses *[]domain.DeathCause    `json:"secondary_causes"`
	Epitaph         *string                 `json:"epitaph" binding:"omitempty,min=1,max=140"`
	Autopsy         *string                 `json:"autopsy" binding:"omitempty,max=5000"`
	ValidationLevel *domain.ValidationLevel `json:"validation_level"`
	MoneySpentCents *int64                  `json:"money_spent_cents" binding:"omitempty,min=0"`
	HoursInvested   *int32                  `json:"hours_invested" binding:"omitempty,min=0"`
	Status          *domain.IdeaStatus      `json:"status"`
	WouldResurrect  *bool                   `json:"would_resurrect"`
}

// Create godoc
// @Summary Bury a new startup idea
// @Tags ideas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateIdeaRequest true "Idea payload"
// @Success 201 {object} domain.Idea
// @Failure 400 {object} ErrorEnvelope
// @Failure 401 {object} ErrorEnvelope
// @Router /ideas [post]
func (h *IdeaHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}

	var req CreateIdeaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, err)
		return
	}

	bornAt, err := parseDate(req.BornAt)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid born_at date", map[string]string{"born_at": "expected YYYY-MM-DD"})
		return
	}
	diedAt, err := parseDate(req.DiedAt)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid died_at date", map[string]string{"died_at": "expected YYYY-MM-DD"})
		return
	}

	idea, err := h.ideaSvc.Create(c.Request.Context(), userID, service.CreateIdeaInput{
		Name:            req.Name,
		Tagline:         req.Tagline,
		Description:     req.Description,
		Category:        req.Category,
		BornAt:          bornAt,
		DiedAt:          diedAt,
		PrimaryCause:    req.PrimaryCause,
		SecondaryCauses: req.SecondaryCauses,
		Epitaph:         req.Epitaph,
		Autopsy:         req.Autopsy,
		ValidationLevel: req.ValidationLevel,
		MoneySpentCents: req.MoneySpentCents,
		HoursInvested:   req.HoursInvested,
		WouldResurrect:  req.WouldResurrect,
	})
	if err != nil {
		handleIdeaServiceErr(c, err)
		return
	}

	c.JSON(http.StatusCreated, idea)
}

// List godoc
// @Summary List buried ideas with filters, full-text search, and pagination
// @Tags ideas
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param cause query string false "Primary cause filter"
// @Param category query string false "Category filter"
// @Param status query string false "Status filter"
// @Param validation query string false "Validation level filter"
// @Param from query string false "Died after YYYY-MM-DD"
// @Param to query string false "Died before YYYY-MM-DD"
// @Param q query string false "Full-text search query"
// @Param sort query string false "Sort field (died_at, created_at, money_spent, lifespan)"
// @Param order query string false "Sort direction (asc, desc)"
// @Success 200 {object} domain.PaginatedIdeas
// @Failure 401 {object} ErrorEnvelope
// @Router /ideas [get]
func (h *IdeaHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}

	filter, err := ParseIdeaFilterFromQuery(c)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	res, err := h.ideaSvc.List(c.Request.Context(), userID, filter)
	if err != nil {
		handleIdeaServiceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Get godoc
// @Summary Get a single buried idea by ID
// @Tags ideas
// @Security BearerAuth
// @Produce json
// @Param id path string true "Idea UUID"
// @Success 200 {object} domain.Idea
// @Failure 404 {object} ErrorEnvelope
// @Router /ideas/{id} [get]
func (h *IdeaHandler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}

	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "idea not found", nil)
		return
	}

	idea, err := h.ideaSvc.Get(c.Request.Context(), userID, ideaID)
	if err != nil {
		handleIdeaServiceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, idea)
}

// Update godoc
// @Summary Update a buried idea
// @Tags ideas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Idea UUID"
// @Param body body UpdateIdeaRequest true "Fields to update"
// @Success 200 {object} domain.Idea
// @Failure 400 {object} ErrorEnvelope
// @Failure 404 {object} ErrorEnvelope
// @Failure 409 {object} ErrorEnvelope
// @Router /ideas/{id} [patch]
func (h *IdeaHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}

	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "idea not found", nil)
		return
	}

	var req UpdateIdeaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, err)
		return
	}

	in := service.UpdateIdeaInput{
		Name:            req.Name,
		Tagline:         req.Tagline,
		Description:     req.Description,
		Category:        req.Category,
		PrimaryCause:    req.PrimaryCause,
		SecondaryCauses: req.SecondaryCauses,
		Epitaph:         req.Epitaph,
		Autopsy:         req.Autopsy,
		ValidationLevel: req.ValidationLevel,
		MoneySpentCents: req.MoneySpentCents,
		HoursInvested:   req.HoursInvested,
		Status:          req.Status,
		WouldResurrect:  req.WouldResurrect,
	}

	if req.BornAt != nil {
		d, err := parseDate(*req.BornAt)
		if err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid born_at date", map[string]string{"born_at": "expected YYYY-MM-DD"})
			return
		}
		in.BornAt = &d
	}
	if req.DiedAt != nil {
		d, err := parseDate(*req.DiedAt)
		if err != nil {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid died_at date", map[string]string{"died_at": "expected YYYY-MM-DD"})
			return
		}
		in.DiedAt = &d
	}

	idea, err := h.ideaSvc.Update(c.Request.Context(), userID, ideaID, in)
	if err != nil {
		handleIdeaServiceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, idea)
}

// Delete godoc
// @Summary Soft-delete a buried idea
// @Tags ideas
// @Security BearerAuth
// @Param id path string true "Idea UUID"
// @Success 204
// @Failure 404 {object} ErrorEnvelope
// @Router /ideas/{id} [delete]
func (h *IdeaHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}

	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "idea not found", nil)
		return
	}

	if err := h.ideaSvc.Delete(c.Request.Context(), userID, ideaID); err != nil {
		handleIdeaServiceErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Resurrect godoc
// @Summary Advance idea status (buried -> on_life_support -> resurrected)
// @Tags ideas
// @Security BearerAuth
// @Produce json
// @Param id path string true "Idea UUID"
// @Success 200 {object} domain.Idea
// @Failure 404 {object} ErrorEnvelope
// @Failure 409 {object} ErrorEnvelope
// @Router /ideas/{id}/resurrect [post]
func (h *IdeaHandler) Resurrect(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}

	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "idea not found", nil)
		return
	}

	idea, err := h.ideaSvc.Resurrect(c.Request.Context(), userID, ideaID)
	if err != nil {
		handleIdeaServiceErr(c, err)
		return
	}
	c.JSON(http.StatusOK, idea)
}

func ParseIdeaFilterFromQuery(c *gin.Context) (domain.IdeaFilter, error) {
	var f domain.IdeaFilter
	if p := c.Query("page"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return f, errors.New("invalid page parameter")
		}
		f.Page = n
	}
	if l := c.Query("limit"); l != "" {
		n, err := strconv.Atoi(l)
		if err != nil {
			return f, errors.New("invalid limit parameter")
		}
		f.Limit = n
	}
	if cause := c.Query("cause"); cause != "" {
		v := domain.DeathCause(cause)
		if !v.Valid() {
			return f, errors.New("invalid cause filter")
		}
		f.Cause = &v
	}
	if cat := c.Query("category"); cat != "" {
		v := domain.IdeaCategory(cat)
		if !v.Valid() {
			return f, errors.New("invalid category filter")
		}
		f.Category = &v
	}
	if st := c.Query("status"); st != "" {
		v := domain.IdeaStatus(st)
		if !v.Valid() {
			return f, errors.New("invalid status filter")
		}
		f.Status = &v
	}
	if val := c.Query("validation"); val != "" {
		v := domain.ValidationLevel(val)
		if !v.Valid() {
			return f, errors.New("invalid validation filter")
		}
		f.Validation = &v
	}
	if from := c.Query("from"); from != "" {
		d, err := parseDate(from)
		if err != nil {
			return f, errors.New("invalid from date (expected YYYY-MM-DD)")
		}
		f.From = &d
	}
	if to := c.Query("to"); to != "" {
		d, err := parseDate(to)
		if err != nil {
			return f, errors.New("invalid to date (expected YYYY-MM-DD)")
		}
		f.To = &d
	}
	f.Query = c.Query("q")
	f.Sort = c.Query("sort")
	f.Order = c.Query("order")
	return f, nil
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func handleIdeaServiceErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "idea not found", nil)
	case errors.Is(err, domain.ErrInvalidTransition):
		RespondError(c, http.StatusConflict, "INVALID_TRANSITION", "idea cannot be transitioned from its current status", nil)
	case errors.Is(err, domain.ErrValidation):
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
	default:
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to process idea request", nil)
	}
}
