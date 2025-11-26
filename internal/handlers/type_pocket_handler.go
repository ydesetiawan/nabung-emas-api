package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"
)

type TypePocketHandler struct {
	service *services.TypePocketService
}

func NewTypePocketHandler(service *services.TypePocketService) *TypePocketHandler {
	return &TypePocketHandler{service: service}
}

func (h *TypePocketHandler) GetAll(c echo.Context) error {
	typePockets, err := h.service.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch type pockets")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", typePockets)
}

func (h *TypePocketHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	typePocket, err := h.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Type pocket not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", typePocket)
}
