package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/http/middleware"
	"github.com/nisha/gravehard/internal/service"
)

const RefreshCookieName = "gravehard_refresh"

type AuthHandler struct {
	authSvc      *service.AuthService
	refreshTTL   time.Duration
	cookieSecure bool
	cookieDomain string
}

func NewAuthHandler(
	authSvc *service.AuthService,
	refreshTTL time.Duration,
	cookieSecure bool,
	cookieDomain string,
) *AuthHandler {
	return &AuthHandler{
		authSvc:      authSvc,
		refreshTTL:   refreshTTL,
		cookieSecure: cookieSecure,
		cookieDomain: cookieDomain,
	}
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email,max=254"`
	Password    string `json:"password" binding:"required,min=8,max=128"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=64"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        *domain.User `json:"user"`
}

// Register godoc
// @Summary Register a new account
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Registration payload"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorEnvelope
// @Failure 409 {object} ErrorEnvelope
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, err)
		return
	}

	res, err := h.authSvc.Register(c.Request.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			RespondError(c, http.StatusConflict, "EMAIL_EXISTS", "an account with this email already exists", nil)
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to register user", nil)
		return
	}

	h.setRefreshCookie(c, res.RefreshToken)
	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken: res.AccessToken,
		User:        res.User,
	})
}

// Login godoc
// @Summary Log in with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login payload"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorEnvelope
// @Failure 401 {object} ErrorEnvelope
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, err)
		return
	}

	res, err := h.authSvc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password", nil)
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to log in", nil)
		return
	}

	h.setRefreshCookie(c, res.RefreshToken)
	c.JSON(http.StatusOK, AuthResponse{
		AccessToken: res.AccessToken,
		User:        res.User,
	})
}

// Refresh godoc
// @Summary Rotate refresh cookie and issue a new access token
// @Tags auth
// @Produce json
// @Success 200 {object} AuthResponse
// @Failure 401 {object} ErrorEnvelope
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	raw, err := c.Cookie(RefreshCookieName)
	if err != nil || raw == "" {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing refresh token", nil)
		return
	}

	res, err := h.authSvc.Refresh(c.Request.Context(), raw)
	if err != nil {
		h.clearRefreshCookie(c)
		if errors.Is(err, domain.ErrUnauthorized) {
			RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired refresh token", nil)
			return
		}
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to refresh session", nil)
		return
	}

	h.setRefreshCookie(c, res.RefreshToken)
	c.JSON(http.StatusOK, AuthResponse{
		AccessToken: res.AccessToken,
		User:        res.User,
	})
}

// Logout godoc
// @Summary Revoke refresh session and clear cookie
// @Tags auth
// @Produce json
// @Success 204
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	raw, _ := c.Cookie(RefreshCookieName)
	_ = h.authSvc.Logout(c.Request.Context(), raw)
	h.clearRefreshCookie(c)
	c.Status(http.StatusNoContent)
}

// Me godoc
// @Summary Get current authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} ErrorEnvelope
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
		return
	}

	user, err := h.authSvc.Me(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "user not found", nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		RefreshCookieName,
		token,
		int(h.refreshTTL.Seconds()),
		"/api/v1/auth",
		h.cookieDomain,
		h.cookieSecure,
		true,
	)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		RefreshCookieName,
		"",
		-1,
		"/api/v1/auth",
		h.cookieDomain,
		h.cookieSecure,
		true,
	)
}

// RespondValidationError converts Gin binding / validator errors into consistent field details.
func RespondValidationError(c *gin.Context, err error) {
	details := make(map[string]string)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			details[field] = fe.Tag()
		}
	}
	if len(details) == 0 {
		details["body"] = "invalid request payload"
	}
	RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "request validation failed", details)
}
