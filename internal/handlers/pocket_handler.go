package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"
)

type PocketHandler struct {
	service *services.PocketService
}

func NewPocketHandler(service *services.PocketService) *PocketHandler {
	return &PocketHandler{service: service}
}

func (h *PocketHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// Parse query parameters
	typePocketID := c.QueryParam("type_pocket_id")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	var typePocketIDPtr *string
	if typePocketID != "" {
		typePocketIDPtr = &typePocketID
	}

	pockets, total, err := h.service.GetAll(userID, typePocketIDPtr, page, limit, sortBy, sortOrder)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch pockets")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	return utils.PaginatedResponse(c, pockets, page, limit, total)
}

func (h *PocketHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	pocket, err := h.service.GetByID(id, userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Pocket not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", pocket)
}

func (h *PocketHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var req models.CreatePocketRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	pocket, err := h.service.Create(userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Pocket created successfully", pocket)
}

func (h *PocketHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	var req models.UpdatePocketRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	pocket, err := h.service.Update(id, userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Pocket updated successfully", pocket)
}

func (h *PocketHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	if err := h.service.Delete(id, userID); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Pocket deleted successfully", nil)
}

func (h *PocketHandler) GetStats(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	// Parse current gold price if provided
	var currentGoldPrice *float64
	if priceStr := c.QueryParam("current_gold_price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			currentGoldPrice = &price
		}
	}

	stats, err := h.service.GetStats(id, userID, currentGoldPrice)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Pocket not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", stats)
}
