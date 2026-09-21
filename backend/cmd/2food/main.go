package main

import (
	"context"
	"log/slog"

	"gitlab.com/2foood/backend/config"
	"gitlab.com/2foood/backend/internal/app"
	"gitlab.com/2foood/backend/pkg/logger"
)

// @title           2food API
// @version         1.0
// @description     API documentation for the 2food backend service.

// @BasePath        /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1Ni..."
func main() {
	// config
	cfg, err := config.Init() // reads .env (if there is one) + env vars
	if err != nil {
		panic(err) // can't run without config, so just crash
	}

	// logger
	log := logger.New(cfg.Logger)

	// app
	app := app.New(context.Background(), log, cfg)
	// Run() blocks until ctrl+c / SIGTERM, so we only get past it on shutdown
	if err := app.Run(); err != nil {
		log.Error("app terminated with error", slog.Any("error", err))
	}
}