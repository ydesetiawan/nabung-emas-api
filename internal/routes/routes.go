package routes

import (
	"database/sql"
	"time"

	"nabung-emas-api/internal/config"
	"nabung-emas-api/internal/handlers"
	"nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/repositories"
	"nabung-emas-api/internal/services"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, db *sql.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	tokenBlacklistRepo := repositories.NewTokenBlacklistRepository(db)
	typePocketRepo := repositories.NewTypePocketRepository(db)
	pocketRepo := repositories.NewPocketRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	analyticsRepo := repositories.NewAnalyticsRepository(db)
	settingsRepo := repositories.NewSettingsRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenBlacklistRepo, cfg)
	userService := services.NewUserService(userRepo)
	typePocketService := services.NewTypePocketService(typePocketRepo)
	pocketService := services.NewPocketService(pocketRepo, typePocketRepo)
	transactionService := services.NewTransactionService(transactionRepo, pocketRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, transactionRepo, pocketRepo)
	settingsService := services.NewSettingsService(settingsRepo)
	antamScraperService := services.NewAntamScraperService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	typePocketHandler := handlers.NewTypePocketHandler(typePocketService)
	pocketHandler := handlers.NewPocketHandler(pocketService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	goldScraperHandler := handlers.NewGoldScraperHandler(antamScraperService)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg, tokenBlacklistRepo)

	// Initialize and start cleanup service for token blacklist
	cleanupService := services.NewCleanupService(tokenBlacklistRepo)
	cleanupService.StartTokenCleanup(24 * time.Hour) // Run cleanup once per day

	// API v1 group
	api := e.Group("/api/v1")

	// Public routes - Authentication
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Protected auth routes
		auth.POST("/logout", authHandler.Logout, authMiddleware.RequireAuth)
		auth.GET("/me", authHandler.GetCurrentUser, authMiddleware.RequireAuth)
	}

	// Public routes - Type Pockets
	typePockets := api.Group("/type-pockets")
	{
		typePockets.GET("", typePocketHandler.GetAll)
		typePockets.GET("/:id", typePocketHandler.GetByID)
	}

	// Protected routes - User Profile
	profile := api.Group("/profile", authMiddleware.RequireAuth)
	{
		profile.GET("", userHandler.GetProfile)
		profile.PATCH("", userHandler.UpdateProfile)
		profile.POST("/avatar", userHandler.UploadAvatar)
		profile.POST("/change-password", userHandler.ChangePassword)
	}

	// Protected routes - Pockets
	pockets := api.Group("/pockets", authMiddleware.RequireAuth)
	{
		pockets.GET("", pocketHandler.GetAll)
		pockets.GET("/:id", pocketHandler.GetByID)
		pockets.POST("", pocketHandler.Create)
		pockets.PATCH("/:id", pocketHandler.Update)
		pockets.DELETE("/:id", pocketHandler.Delete)
		pockets.GET("/:id/stats", pocketHandler.GetStats)
	}

	// Protected routes - Transactions
	transactions := api.Group("/transactions", authMiddleware.RequireAuth)
	{
		transactions.GET("", transactionHandler.GetAll)
		transactions.GET("/:id", transactionHandler.GetByID)
		transactions.POST("", transactionHandler.Create)
		transactions.PATCH("/:id", transactionHandler.Update)
		transactions.DELETE("/:id", transactionHandler.Delete)
		transactions.POST("/:id/receipt", transactionHandler.UploadReceipt)
	}

	// Protected routes - Analytics
	analytics := api.Group("/analytics", authMiddleware.RequireAuth)
	{
		analytics.GET("/dashboard", analyticsHandler.GetDashboard)
		analytics.GET("/portfolio", analyticsHandler.GetPortfolio)
		analytics.GET("/monthly-purchases", analyticsHandler.GetMonthlyPurchases)
		analytics.GET("/brand-distribution", analyticsHandler.GetBrandDistribution)
		analytics.GET("/trends", analyticsHandler.GetTrends)
	}

	// Protected routes - Settings
	settings := api.Group("/settings", authMiddleware.RequireAuth)
	{
		settings.GET("", settingsHandler.Get)
		settings.PATCH("", settingsHandler.Update)
	}

	// Gold Scraper routes
	api.POST("/gold-scraper/scrape", goldScraperHandler.Scrape)

	// Gold Price routes (can be public or protected based on requirements)
	// goldPrice := api.Group("/gold-price")
	// {
	// 	goldPrice.GET("/current", goldPriceHandler.GetCurrent)
	// 	goldPrice.GET("/history", goldPriceHandler.GetHistory)
	// }
}
