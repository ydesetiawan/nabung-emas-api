package middleware

import (
	"net/http"
	"strings"

	"nabung-emas-api/internal/config"
	"nabung-emas-api/internal/repositories"
	"nabung-emas-api/internal/utils"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	config             *config.Config
	tokenBlacklistRepo *repositories.TokenBlacklistRepository
}

func NewAuthMiddleware(cfg *config.Config, tokenBlacklistRepo *repositories.TokenBlacklistRepository) *AuthMiddleware {
	return &AuthMiddleware{
		config:             cfg,
		tokenBlacklistRepo: tokenBlacklistRepo,
	}
}

func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := parts[1]

		// Check if token is blacklisted
		isBlacklisted, err := m.tokenBlacklistRepo.IsBlacklisted(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify token")
		}
		if isBlacklisted {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token has been revoked")
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString, m.config.JWTSecret)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		return next(c)
	}
}

func GetUserID(c echo.Context) string {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return ""
	}
	return userID
}

func GetUserEmail(c echo.Context) string {
	email, ok := c.Get("user_email").(string)
	if !ok {
		return ""
	}
	return email
}
