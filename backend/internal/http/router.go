package http

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nisha/gravehard/internal/config"
	"github.com/nisha/gravehard/internal/http/handler"
	"github.com/nisha/gravehard/internal/http/middleware"
)

// RouterDeps holds all dependencies required to build the HTTP router.
type RouterDeps struct {
	Config     *config.Config
	Logger     *slog.Logger
	OpsHandler *handler.OpsHandler
}

// NewRouter constructs the Gin engine with global middleware and registered routes.
func NewRouter(deps RouterDeps) *gin.Engine {
	if deps.Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Recovery(deps.Logger),
		middleware.AccessLog(deps.Logger),
		middleware.SecurityHeaders(),
		middleware.BodySizeLimit(1<<20), // 1 MiB
		middleware.CORS(deps.Config.CORSOrigins),
	)

	r.GET("/healthz", deps.OpsHandler.Healthz)
	r.GET("/readyz", deps.OpsHandler.Readyz)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
