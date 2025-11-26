package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"nabung-emas-api/internal/config"
	"nabung-emas-api/internal/database"
	custommiddleware "nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/routes"
	"nabung-emas-api/internal/utils"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Echo
	e := echo.New()

	// Hide banner
	e.HideBanner = true

	// Custom validator
	e.Validator = utils.NewValidator()

	// Middleware
	e.Use(custommiddleware.SetupLogger())
	e.Use(middleware.Recover())
	e.Use(custommiddleware.SetupCORS(cfg))

	// Setup routes
	routes.Setup(e, db, cfg)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"service": "nabung-emas-api",
		})
	})

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📝 Environment: %s", cfg.Env)
	log.Printf("🔗 API Base URL: http://localhost:%s/api/v1", port)
	
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
