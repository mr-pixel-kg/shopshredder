package main

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/api"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/config"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/database"
	"log"
	"net/http"
	"strconv"
)

// @title mpXsandbox API
// @version 1.0.0
// @description Management API for Docker Sandbox Environment
// @license.name MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.basic BasicAuth
// @schemes http https
func main() {
	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Sentry
	if cfg.Sentry.SentryDSN == "" {
		log.Println("Sentry DSN not provided, skipping Sentry initialization")
	} else {
		log.Println("Sentry DSN provided, initializing Sentry")

		if err := sentry.Init(sentry.ClientOptions{
			Dsn: cfg.Sentry.SentryDSN,
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
	}

	// Database
	database.ConnectDB(cfg.Database)

	// Echo framework
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())  // Loggt Anfragen
	e.Use(middleware.Recover()) // Fängt Panics ab und gibt 500 zurück
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Server.AllowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Access-Control-Allow-Origin"},
	}))
	if cfg.Sentry.SentryDSN != "" {
		// Add sentry middleware if enabled
		e.Use(sentryecho.New(sentryecho.Options{}))
	}

	// Register routes
	api.RegisterRoutes(e, cfg)

	// Start server
	port := cfg.Server.Port
	log.Printf("Starting server on http://localhost:%d", port)
	if err := e.Start(":" + strconv.Itoa(port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not start server: %v", err)
	}
}
