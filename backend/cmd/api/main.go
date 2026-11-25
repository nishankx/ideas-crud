package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/nisha/gravehard/docs"
	"github.com/nisha/gravehard/internal/auth"
	"github.com/nisha/gravehard/internal/config"
	apphttp "github.com/nisha/gravehard/internal/http"
	"github.com/nisha/gravehard/internal/http/handler"
	"github.com/nisha/gravehard/internal/service"
)

// @title Startup Idea Graveyard API
// @version 1.0
// @description API for burying abandoned startup ideas and analyzing failure patterns.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed_to_load_config", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed_to_parse_database_url", slog.Any("error", err))
		os.Exit(1)
	}
	poolCfg.MaxConns = cfg.DBMaxConns
	poolCfg.MinConns = cfg.DBMinConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		logger.Error("failed_to_create_db_pool", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessTTL)
	authSvc := service.NewAuthService(pool, jwtMgr, cfg.JWTRefreshTTL, cfg.DBQueryTimeout)

	opsHandler := handler.NewOpsHandler(pool, cfg.DBQueryTimeout)
	authHandler := handler.NewAuthHandler(authSvc, cfg.JWTRefreshTTL, cfg.CookieSecure, cfg.CookieDomain)

	router := apphttp.NewRouter(apphttp.RouterDeps{
		Config:      cfg,
		Logger:      logger,
		JWTManager:  jwtMgr,
		OpsHandler:  opsHandler,
		AuthHandler: authHandler,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: cfg.DBQueryTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server_starting", slog.String("port", cfg.Port), slog.String("env", cfg.AppEnv))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		logger.Error("server_failed", slog.Any("error", err))
		os.Exit(1)
	case <-ctx.Done():
		logger.Info("shutdown_signal_received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server_forced_shutdown", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("server_stopped_gracefully")
}
