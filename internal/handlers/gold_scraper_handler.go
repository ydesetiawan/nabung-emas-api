package handlers

import (
	"net/http"
	"strconv"

	"nabung-emas-api/internal/repositories"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"

	"github.com/labstack/echo/v4"
)

type GoldScraperHandler struct {
	antamScraperService *services.AntamScraperService
	repo                *repositories.GoldPricingHistoryRepository
}

func NewGoldScraperHandler(service *services.AntamScraperService, repo *repositories.GoldPricingHistoryRepository) *GoldScraperHandler {
	return &GoldScraperHandler{
		antamScraperService: service,
		repo:                repo,
	}
}

func (h *GoldScraperHandler) Scrape(c echo.Context) error {
	scraper := h.antamScraperService
	prices, err := scraper.Scrape()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to scrape prices")
	}

	// Save to database
	if len(prices) > 0 {
		err = h.repo.BulkCreate(prices)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save prices to database: "+err.Error())
		}
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success: "+strconv.Itoa(len(prices))+" prices scraped and saved", prices)
}
