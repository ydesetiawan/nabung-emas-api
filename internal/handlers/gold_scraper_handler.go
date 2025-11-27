package handlers

import (
	"log"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GoldScraperHandler handles HTTP requests for gold price scraping
type GoldScraperHandler struct {
	service *services.GoldScraperService
}

// NewGoldScraperHandler creates a new instance of GoldScraperHandler
func NewGoldScraperHandler(service *services.GoldScraperService) *GoldScraperHandler {
	return &GoldScraperHandler{service: service}
}

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Count   int         `json:"count,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// ScrapeGoldPrices handles POST /api/scrape
// @Summary Scrape gold prices from website
// @Description Scrapes gold prices from logammulia.com and saves to database
// @Tags Gold Scraper
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/scrape [post]
func (h *GoldScraperHandler) ScrapeGoldPrices(c echo.Context) error {
	log.Println("🚀 Starting gold price scraping...")

	// Perform scraping
	result, err := h.service.ScrapeLogamMulia()
	if err != nil {
		log.Printf("❌ Scraping failed: %v", err)
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to scrape gold prices",
			Errors:  []string{err.Error()},
		})
	}

	// Check if scraping was successful
	if !result.Success {
		log.Printf("⚠️  Scraping completed with warnings: %s", result.Message)
		return c.JSON(http.StatusOK, APIResponse{
			Success: false,
			Message: result.Message,
			Count:   result.SavedCount,
			Data:    result.Data,
			Errors:  result.Errors,
		})
	}

	log.Printf("✅ Successfully scraped and saved %d gold prices", result.SavedCount)

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: result.Message,
		Count:   result.SavedCount,
		Data:    result.Data,
	})
}

// GetAllPrices handles GET /api/prices
// @Summary Get all gold prices
// @Description Retrieves all gold prices with optional filters (type, source, limit)
// @Tags Gold Scraper
// @Accept json
// @Produce json
// @Param type query string false "Filter by gold type (partial match)"
// @Param source query string false "Filter by source (antam or usb)"
// @Param limit query int false "Limit number of results"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/prices [get]
func (h *GoldScraperHandler) GetAllPrices(c echo.Context) error {
	log.Println("📋 Fetching all gold prices...")

	// Parse query parameters
	filter := models.GoldPricingHistoryFilter{
		GoldType: c.QueryParam("type"),
		Source:   models.GoldSource(c.QueryParam("source")),
		Limit:    0,
	}

	// Parse limit
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("⚠️  Invalid limit parameter: %v", err)
			return c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Invalid limit parameter",
				Errors:  []string{err.Error()},
			})
		}
		filter.Limit = limit
	}

	// Validate source if provided
	if filter.Source != "" && filter.Source != models.GoldSourceAntam && filter.Source != models.GoldSourceUSB {
		log.Printf("⚠️  Invalid source parameter: %s", filter.Source)
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid source parameter. Must be 'antam' or 'usb'",
		})
	}

	// Fetch prices
	prices, err := h.service.GetAllPrices(filter)
	if err != nil {
		log.Printf("❌ Failed to fetch prices: %v", err)
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch gold prices",
			Errors:  []string{err.Error()},
		})
	}

	log.Printf("✅ Successfully fetched %d gold prices", len(prices))

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Successfully retrieved gold prices",
		Count:   len(prices),
		Data:    prices,
	})
}

// GetLatestPrices handles GET /api/prices/latest
// @Summary Get latest gold prices
// @Description Retrieves the latest price for each gold type
// @Tags Gold Scraper
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/prices/latest [get]
func (h *GoldScraperHandler) GetLatestPrices(c echo.Context) error {
	log.Println("🔍 Fetching latest gold prices...")

	// Fetch latest prices
	prices, err := h.service.GetLatestPrices()
	if err != nil {
		log.Printf("❌ Failed to fetch latest prices: %v", err)
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch latest gold prices",
			Errors:  []string{err.Error()},
		})
	}

	log.Printf("✅ Successfully fetched %d latest gold prices", len(prices))

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Successfully retrieved latest gold prices",
		Count:   len(prices),
		Data:    prices,
	})
}

// GetPriceByID handles GET /api/prices/:id
// @Summary Get gold price by ID
// @Description Retrieves a specific gold price by ID
// @Tags Gold Scraper
// @Accept json
// @Produce json
// @Param id path int true "Price ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/prices/{id} [get]
func (h *GoldScraperHandler) GetPriceByID(c echo.Context) error {
	// Parse ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("⚠️  Invalid ID parameter: %v", err)
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid ID parameter",
			Errors:  []string{err.Error()},
		})
	}

	log.Printf("🔍 Fetching gold price with ID: %d", id)

	// Fetch price by ID
	price, err := h.service.GetPriceByID(id)
	if err != nil {
		log.Printf("❌ Failed to fetch price: %v", err)
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch gold price",
			Errors:  []string{err.Error()},
		})
	}

	// Check if price exists
	if price == nil {
		log.Printf("⚠️  Price not found with ID: %d", id)
		return c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Gold price not found",
		})
	}

	log.Printf("✅ Successfully fetched gold price with ID: %d", id)

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Successfully retrieved gold price",
		Count:   1,
		Data:    price,
	})
}
