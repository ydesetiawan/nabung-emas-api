package handlers

import (
	"net/http"
	"strconv"

	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"

	"github.com/labstack/echo/v4"
)

type GoldScraperHandler struct {
	antamScraperService *services.AntamScraperService
}

func NewGoldScraperHandler(service *services.AntamScraperService) *GoldScraperHandler {
	return &GoldScraperHandler{antamScraperService: service}
}

func (h *GoldScraperHandler) Scrape(c echo.Context) error {

	scraper := h.antamScraperService
	prices, err := scraper.Scrape()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to scrape prices")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success:"+strconv.Itoa(len(prices)), prices)
}
