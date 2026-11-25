package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nisha/gravehard/internal/auth"
)

const UserIDKey = "user_id"

// RequireAuth validates the Bearer JWT access token and attaches the authenticated user ID to the context.
func RequireAuth(jwtMgr *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "missing or invalid authorization header",
				},
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.VerifyAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid or expired access token",
				},
			})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// GetUserID retrieves the authenticated user's UUID from the Gin context.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(UserIDKey)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
