package handlers

import (
	"log"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// Galeri24ScraperHandler handles HTTP requests for Galeri24 gold price scraping
type Galeri24ScraperHandler struct {
	service *services.Galeri24ScraperService
}

// NewGaleri24ScraperHandler creates a new instance of Galeri24ScraperHandler
func NewGaleri24ScraperHandler(service *services.Galeri24ScraperService) *Galeri24ScraperHandler {
	return &Galeri24ScraperHandler{service: service}
}

// ScrapeGaleri24Prices handles POST /api/v1/galeri24-scraper/scrape
// @Summary Scrape gold prices from Galeri24
// @Description Scrapes gold prices from galeri24.co.id and saves to database (replaces existing data for same date)
// @Tags Galeri24 Scraper
// @Accept json
// @Produce json
// @Success 200 {object} models.ScrapeResult
// @Failure 500 {object} models.ScrapeResult
// @Router /api/v1/galeri24-scraper/scrape [post]
func (h *Galeri24ScraperHandler) ScrapeGaleri24Prices(c echo.Context) error {
	log.Println("🚀 Starting Galeri24 gold price scraping...")

	// Perform scraping
	result, err := h.service.ScrapeGaleri24()
	if err != nil {
		log.Printf("❌ Scraping failed: %v", err)
		return c.JSON(http.StatusInternalServerError, result)
	}

	// Check if scraping was successful
	if !result.Success {
		log.Printf("⚠️  Scraping completed with warnings: %s", result.Message)
		return c.JSON(http.StatusOK, result)
	}

	log.Printf("✅ Successfully scraped: %d new, %d updated", result.SavedCount, result.UpdatedCount)

	return c.JSON(http.StatusOK, result)
}

// GetAllPrices handles GET /api/v1/galeri24-scraper/prices
// @Summary Get all gold prices
// @Description Retrieves all gold prices with optional filters (type, source, date range, limit)
// @Tags Galeri24 Scraper
// @Accept json
// @Produce json
// @Param type query string false "Filter by gold type (partial match)"
// @Param source query string false "Filter by source/vendor"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/galeri24-scraper/prices [get]
func (h *Galeri24ScraperHandler) GetAllPrices(c echo.Context) error {
	log.Println("📋 Fetching all gold prices...")

	// Parse query parameters
	filter := models.GoldPricingHistoryFilter{
		GoldType: c.QueryParam("type"),
		Source:   models.GoldSource(c.QueryParam("source")),
	}

	// Parse limit
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("⚠️  Invalid limit parameter: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid limit parameter",
				"errors":  []string{err.Error()},
			})
		}
		filter.Limit = limit
	}

	// Parse offset
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("⚠️  Invalid offset parameter: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid offset parameter",
				"errors":  []string{err.Error()},
			})
		}
		filter.Offset = offset
	}

	// Parse start_date
	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			log.Printf("⚠️  Invalid start_date parameter: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid start_date parameter. Use format: YYYY-MM-DD",
				"errors":  []string{err.Error()},
			})
		}
		filter.StartDate = &startDate
	}

	// Parse end_date
	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			log.Printf("⚠️  Invalid end_date parameter: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid end_date parameter. Use format: YYYY-MM-DD",
				"errors":  []string{err.Error()},
			})
		}
		filter.EndDate = &endDate
	}

	// Validate source if provided
	if filter.Source != "" && !filter.Source.IsValid() {
		log.Printf("⚠️  Invalid source parameter: %s", filter.Source)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid source parameter",
		})
	}

	// Fetch prices
	prices, err := h.service.GetAllPrices(filter)
	if err != nil {
		log.Printf("❌ Failed to fetch prices: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch gold prices",
			"errors":  []string{err.Error()},
		})
	}

	log.Printf("✅ Successfully fetched %d gold prices", len(prices))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully retrieved gold prices",
		"count":   len(prices),
		"data":    prices,
	})
}

// GetLatestPrices handles GET /api/v1/galeri24-scraper/prices/latest
// @Summary Get latest gold prices
// @Description Retrieves the latest price for each gold type and source combination
// @Tags Galeri24 Scraper
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/galeri24-scraper/prices/latest [get]
func (h *Galeri24ScraperHandler) GetLatestPrices(c echo.Context) error {
	log.Println("🔍 Fetching latest gold prices...")

	// Fetch latest prices
	prices, err := h.service.GetLatestPrices()
	if err != nil {
		log.Printf("❌ Failed to fetch latest prices: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch latest gold prices",
			"errors":  []string{err.Error()},
		})
	}

	log.Printf("✅ Successfully fetched %d latest gold prices", len(prices))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully retrieved latest gold prices",
		"count":   len(prices),
		"data":    prices,
	})
}

// GetPriceByID handles GET /api/v1/galeri24-scraper/prices/:id
// @Summary Get gold price by ID
// @Description Retrieves a specific gold price by ID
// @Tags Galeri24 Scraper
// @Accept json
// @Produce json
// @Param id path int true "Price ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/galeri24-scraper/prices/{id} [get]
func (h *Galeri24ScraperHandler) GetPriceByID(c echo.Context) error {
	// Parse ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("⚠️  Invalid ID parameter: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid ID parameter",
			"errors":  []string{err.Error()},
		})
	}

	log.Printf("🔍 Fetching gold price with ID: %d", id)

	// Fetch price by ID
	price, err := h.service.GetPriceByID(id)
	if err != nil {
		log.Printf("❌ Failed to fetch price: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch gold price",
			"errors":  []string{err.Error()},
		})
	}

	// Check if price exists
	if price == nil {
		log.Printf("⚠️  Price not found with ID: %d", id)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "Gold price not found",
		})
	}

	log.Printf("✅ Successfully fetched gold price with ID: %d", id)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully retrieved gold price",
		"data":    price,
	})
}

// GetPricesByDate handles GET /api/v1/galeri24-scraper/prices/date/:date
// @Summary Get gold prices by date
// @Description Retrieves all gold prices for a specific date
// @Tags Galeri24 Scraper
// @Accept json
// @Produce json
// @Param date path string true "Date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/galeri24-scraper/prices/date/{date} [get]
func (h *Galeri24ScraperHandler) GetPricesByDate(c echo.Context) error {
	// Parse date from path parameter
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("⚠️  Invalid date parameter: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid date parameter. Use format: YYYY-MM-DD",
			"errors":  []string{err.Error()},
		})
	}

	log.Printf("🔍 Fetching gold prices for date: %s", date.Format("2006-01-02"))

	// Fetch prices by date
	prices, err := h.service.GetPricesByDate(date)
	if err != nil {
		log.Printf("❌ Failed to fetch prices: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch gold prices",
			"errors":  []string{err.Error()},
		})
	}

	log.Printf("✅ Successfully fetched %d gold prices for date %s", len(prices), date.Format("2006-01-02"))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully retrieved gold prices",
		"count":   len(prices),
		"date":    date.Format("2006-01-02"),
		"data":    prices,
	})
}

// GetStats handles GET /api/v1/galeri24-scraper/stats
// @Summary Get statistics
// @Description Retrieves statistics about gold pricing data
// @Tags Galeri24 Scraper
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/galeri24-scraper/stats [get]
func (h *Galeri24ScraperHandler) GetStats(c echo.Context) error {
	log.Println("📊 Fetching gold pricing statistics...")

	// Fetch stats
	stats, err := h.service.GetStats()
	if err != nil {
		log.Printf("❌ Failed to fetch stats: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch statistics",
			"errors":  []string{err.Error()},
		})
	}

	log.Printf("✅ Successfully fetched statistics")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully retrieved statistics",
		"data":    stats,
	})
}
