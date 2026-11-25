package handler

import "github.com/gin-gonic/gin"

// ErrorEnvelope is the consistent JSON error response shape across the API.
type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func RespondError(c *gin.Context, status int, code, message string, details map[string]string) {
	c.AbortWithStatusJSON(status, ErrorEnvelope{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
