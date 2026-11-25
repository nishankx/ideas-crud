package http

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/nisha/gravehard/internal/auth"
	"github.com/nisha/gravehard/internal/config"
	"github.com/nisha/gravehard/internal/http/handler"
	"github.com/nisha/gravehard/internal/http/middleware"
)

// RouterDeps holds all dependencies required to build the HTTP router.
type RouterDeps struct {
	Config      *config.Config
	Logger      *slog.Logger
	JWTManager  *auth.JWTManager
	OpsHandler  *handler.OpsHandler
	AuthHandler *handler.AuthHandler
	IdeaHandler *handler.IdeaHandler
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
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	genLimiter := middleware.NewRateLimiter(deps.Config.RateLimitGenRPM)
	authLimiter := middleware.NewRateLimiter(deps.Config.RateLimitAuthPM)

	v1 := r.Group("/api/v1")
	v1.Use(genLimiter.Middleware())

	authGroup := v1.Group("/auth")
	authGroup.Use(authLimiter.Middleware())
	{
		authGroup.POST("/register", deps.AuthHandler.Register)
		authGroup.POST("/login", deps.AuthHandler.Login)
		authGroup.POST("/refresh", deps.AuthHandler.Refresh)
		authGroup.POST("/logout", deps.AuthHandler.Logout)
		authGroup.GET("/me", middleware.RequireAuth(deps.JWTManager), deps.AuthHandler.Me)
	}

	ideasGroup := v1.Group("/ideas")
	ideasGroup.Use(middleware.RequireAuth(deps.JWTManager))
	{
		ideasGroup.POST("", deps.IdeaHandler.Create)
		ideasGroup.GET("", deps.IdeaHandler.List)
		ideasGroup.GET("/:id", deps.IdeaHandler.Get)
		ideasGroup.PATCH("/:id", deps.IdeaHandler.Update)
		ideasGroup.DELETE("/:id", deps.IdeaHandler.Delete)
		ideasGroup.POST("/:id/resurrect", deps.IdeaHandler.Resurrect)
	}

	return r
}
