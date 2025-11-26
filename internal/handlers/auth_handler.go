package handlers

import (
	"net/http"
	"strings"

	"nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	user, tokens, err := h.authService.Register(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Account created successfully", map[string]interface{}{
		"user":          user,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	user, tokens, err := h.authService.Login(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Login successful", map[string]interface{}{
		"user":          user,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
	})
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req models.ForgotPasswordRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	if err := h.authService.ForgotPassword(req.Email); err != nil {
		// Always return success for security
		return utils.SuccessResponse(c, http.StatusOK, "Password reset link sent to your email", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Password reset link sent to your email", nil)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req models.ResetPasswordRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	if err := h.authService.ResetPassword(req.Token, req.Password); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req models.RefreshTokenRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", tokens)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// Extract the token from the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header")
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
	}

	accessToken := parts[1]

	// Blacklist the token
	if err := h.authService.Logout(accessToken, userID); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to logout")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}

func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", user)
}
