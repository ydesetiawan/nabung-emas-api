package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"
)

type SettingsHandler struct {
	service *services.SettingsService
}

func NewSettingsHandler(service *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

func (h *SettingsHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	settings, err := h.service.Get(userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch settings")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", settings)
}

func (h *SettingsHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var req models.UpdateSettingsRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	settings, err := h.service.Update(userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Settings updated successfully", settings)
}
