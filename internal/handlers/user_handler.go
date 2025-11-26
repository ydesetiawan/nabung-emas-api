package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	user, stats, err := h.userService.GetProfile(userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", map[string]interface{}{
		"user":  user,
		"stats": stats,
	})
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var req models.UpdateProfileRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	user, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", user)
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var req models.ChangePasswordRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

func (h *UserHandler) UploadAvatar(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// TODO: Implement file upload logic
	// For now, return not implemented
	return utils.ErrorResponse(c, http.StatusNotImplemented, "Avatar upload not yet implemented")
}
