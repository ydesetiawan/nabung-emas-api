package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetDashboard(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// Parse current gold price if provided
	var currentGoldPrice *float64
	if priceStr := c.QueryParam("current_gold_price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			currentGoldPrice = &price
		}
	}

	dashboard, err := h.service.GetDashboard(userID, currentGoldPrice)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dashboard data")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", dashboard)
}

func (h *AnalyticsHandler) GetPortfolio(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// Parse current gold price if provided
	var currentGoldPrice *float64
	if priceStr := c.QueryParam("current_gold_price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			currentGoldPrice = &price
		}
	}

	portfolio, err := h.service.GetPortfolio(userID, currentGoldPrice)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolio analytics")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", portfolio)
}

func (h *AnalyticsHandler) GetMonthlyPurchases(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	months, _ := strconv.Atoi(c.QueryParam("months"))
	pocketID := c.QueryParam("pocket_id")

	var pocketIDPtr *string
	if pocketID != "" {
		pocketIDPtr = &pocketID
	}

	analytics, err := h.service.GetMonthlyPurchases(userID, months, pocketIDPtr)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch monthly purchase analytics")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", analytics)
}

func (h *AnalyticsHandler) GetBrandDistribution(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	distribution, err := h.service.GetBrandDistribution(userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch brand distribution")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", distribution)
}

func (h *AnalyticsHandler) GetTrends(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	period := c.QueryParam("period")
	groupBy := c.QueryParam("group_by")

	trends, err := h.service.GetTrends(userID, period, groupBy)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch trends")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", trends)
}
