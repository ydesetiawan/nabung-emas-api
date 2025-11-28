package services

import (
	"context"
	"fmt"
	"log"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// GoldScraperService handles web scraping of gold prices from various sources
type GoldScraperService struct {
	repo *repositories.GoldPricingHistoryRepository
}

// NewGoldScraperService creates a new instance of GoldScraperService
func NewGoldScraperService(repo *repositories.GoldPricingHistoryRepository) *GoldScraperService {
	return &GoldScraperService{repo: repo}
}

// ScrapedGoldData represents the raw data scraped from the website
type ScrapedGoldData struct {
	GoldType    string
	BuyPrice    string
	SellPrice   string
	ProductName string // Full product name for category detection
}

// ScrapeResult represents the result of a scraping operation
type ScrapeResult struct {
	Success      bool                        `json:"success"`
	Message      string                      `json:"message"`
	PricingDate  time.Time                   `json:"pricing_date"`
	TotalScraped int                         `json:"total_scraped"`
	SavedCount   int                         `json:"saved_count"`
	UpdatedCount int                         `json:"updated_count"`
	FailedCount  int                         `json:"failed_count"`
	Errors       []string                    `json:"errors"`
	Duration     string                      `json:"duration"`
	Data         []models.GoldPricingHistory `json:"data,omitempty"`
}

// ScrapeLogamMulia scrapes gold prices from logammulia.com (Antam source)
func (s *GoldScraperService) ScrapeLogamMulia() (*ScrapeResult, error) {
	return s.scrapeWithRetry(context.Background(), models.GoldSourceAntam, "https://logammulia.com/id/harga-emas-hari-ini", 3)
}

// scrapeWithRetry implements retry logic with exponential backoff
func (s *GoldScraperService) scrapeWithRetry(ctx context.Context, source models.GoldSource, url string, maxRetries int) (*ScrapeResult, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("� Scraping attempt %d/%d for source: %s", attempt, maxRetries, source)

		result, err := s.scrape(ctx, source, url)
		if err == nil {
			return result, nil
		}

		lastErr = err
		if attempt < maxRetries {
			backoff := time.Duration(attempt*attempt) * time.Second
			log.Printf("⏳ Retry in %v due to error: %v", backoff, err)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// scrape performs the actual scraping operation
func (s *GoldScraperService) scrape(ctx context.Context, source models.GoldSource, url string) (*ScrapeResult, error) {
	startTime := time.Now()
	log.Printf("🕷️  Starting gold price scraping from %s...", url)

	// Create collector with context and configuration
	c := colly.NewCollector(
		colly.AllowedDomains("logammulia.com", "www.logammulia.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.Async(false),
	)

	// Set timeouts and delays
	c.SetRequestTimeout(30 * time.Second)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       2 * time.Second, // Respectful delay
	})

	scrapedData := []ScrapedGoldData{}
	errors := []string{}

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		errorMsg := fmt.Sprintf("Request failed: %v (Status: %d)", err, r.StatusCode)
		log.Printf("❌ %s", errorMsg)
		errors = append(errors, errorMsg)
	})

	// Before making a request
	c.OnRequest(func(r *colly.Request) {
		log.Printf("🌐 Visiting: %s", r.URL.String())
	})

	// Parse the HTML and extract gold prices
	c.OnHTML("table", func(e *colly.HTMLElement) {
		log.Println("📊 Found table, parsing gold prices...")

		// Look for table rows
		e.ForEach("tbody tr", func(i int, row *colly.HTMLElement) {
			// Extract data from each row
			goldType := strings.TrimSpace(row.ChildText("td:nth-child(1)"))
			buyPrice := strings.TrimSpace(row.ChildText("td:nth-child(2)"))
			sellPrice := strings.TrimSpace(row.ChildText("td:nth-child(3)"))

			// Skip empty rows or header rows
			if goldType == "" || goldType == "Jenis" || goldType == "Type" {
				return
			}

			// Clean up the data
			goldType = cleanText(goldType)
			buyPrice = cleanPrice(buyPrice)
			sellPrice = cleanPrice(sellPrice)

			// Only add if we have valid data
			if goldType != "" && (buyPrice != "" || sellPrice != "") {
				scrapedData = append(scrapedData, ScrapedGoldData{
					GoldType:    goldType,
					BuyPrice:    buyPrice,
					SellPrice:   sellPrice,
					ProductName: goldType, // Use gold type as product name for category detection
				})

				log.Printf("✅ Scraped: %s - Buy: %s, Sell: %s", goldType, buyPrice, sellPrice)
			}
		})
	})

	// Alternative parsing for different HTML structure
	c.OnHTML(".price-table, .gold-price-table, [class*='price'], [class*='gold']", func(e *colly.HTMLElement) {
		log.Println("📊 Found alternative price table structure...")

		e.ForEach("tr", func(i int, row *colly.HTMLElement) {
			cells := row.ChildTexts("td")
			if len(cells) >= 3 {
				goldType := cleanText(cells[0])
				buyPrice := cleanPrice(cells[1])
				sellPrice := cleanPrice(cells[2])

				// Skip header rows
				if goldType == "" || goldType == "Jenis" || goldType == "Type" {
					return
				}

				scrapedData = append(scrapedData, ScrapedGoldData{
					GoldType:    goldType,
					BuyPrice:    buyPrice,
					SellPrice:   sellPrice,
					ProductName: goldType,
				})

				log.Printf("✅ Scraped (alt): %s - Buy: %s, Sell: %s", goldType, buyPrice, sellPrice)
			}
		})
	})

	// Response handler
	c.OnResponse(func(r *colly.Response) {
		log.Printf("✅ Received response from %s (Status: %d)", r.Request.URL, r.StatusCode)
	})

	// Visit the URL with context
	err := c.Visit(url)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to visit website: %v", err)
		log.Printf("❌ %s", errorMsg)
		return &ScrapeResult{
			Success:      false,
			Message:      errorMsg,
			TotalScraped: 0,
			SavedCount:   0,
			Errors:       append(errors, errorMsg),
			Duration:     time.Since(startTime).String(),
		}, err
	}

	// Wait for async requests to complete
	c.Wait()

	log.Printf("📦 Total items scraped: %d", len(scrapedData))

	// If no data was scraped, return early
	if len(scrapedData) == 0 {
		message := "No gold price data found on the website"
		log.Printf("⚠️  %s", message)
		return &ScrapeResult{
			Success:      false,
			Message:      message,
			TotalScraped: 0,
			SavedCount:   0,
			Errors:       append(errors, "No data found in expected HTML structure"),
			Duration:     time.Since(startTime).String(),
		}, nil
	}

	// Save scraped data to database
	pricingDate := time.Now().Truncate(24 * time.Hour) // Today's date at midnight
	savedCount, updatedCount, failedCount, saveErrors := s.saveScrapedData(scrapedData, source, pricingDate)

	// Append save errors
	errors = append(errors, saveErrors...)

	// Fetch the saved data to return
	savedData, _ := s.repo.GetByDate(pricingDate)

	duration := time.Since(startTime)
	log.Printf("💾 Scraping complete: %d scraped, %d saved, %d updated, %d failed in %v",
		len(scrapedData), savedCount, updatedCount, failedCount, duration)

	success := failedCount == 0
	message := fmt.Sprintf("Successfully scraped %d items: %d new, %d updated",
		len(scrapedData), savedCount, updatedCount)

	if failedCount > 0 {
		message = fmt.Sprintf("Scraped %d items: %d new, %d updated, %d failed",
			len(scrapedData), savedCount, updatedCount, failedCount)
	}

	return &ScrapeResult{
		Success:      success,
		Message:      message,
		PricingDate:  pricingDate,
		TotalScraped: len(scrapedData),
		SavedCount:   savedCount,
		UpdatedCount: updatedCount,
		FailedCount:  failedCount,
		Errors:       errors,
		Duration:     duration.String(),
		Data:         savedData,
	}, nil
}

// saveScrapedData saves the scraped data to the database
func (s *GoldScraperService) saveScrapedData(data []ScrapedGoldData, source models.GoldSource, pricingDate time.Time) (int, int, int, []string) {
	if len(data) == 0 {
		return 0, 0, 0, nil
	}

	// Convert scraped data to create models
	createData := make([]models.GoldPricingHistoryCreate, 0, len(data))
	errors := []string{}

	for _, item := range data {
		// Parse sell price
		sellPrice, err := parsePrice(item.SellPrice)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to parse sell price for %s: %v", item.GoldType, err)
			log.Printf("⚠️  %s", errorMsg)
			errors = append(errors, errorMsg)
			continue
		}

		// Detect category from product name
		category, err := detectCategory(item.ProductName)
		if err != nil {
			log.Printf("⚠️  %v - using default category", err)
			category = models.GoldCategoryEmasBatangan // Default fallback
		}

		createData = append(createData, models.GoldPricingHistoryCreate{
			PricingDate: pricingDate,
			GoldType:    item.GoldType,
			SellPrice:   sellPrice,
			Source:      source,
			Category:    category,
		})
	}

	// Save to database using batch insert
	savedCount, updatedCount, err := s.repo.CreateBatch(createData)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to save batch: %v", err)
		errors = append(errors, errorMsg)
		return 0, 0, len(data), errors
	}

	failedCount := len(data) - len(createData)
	log.Printf("💾 Batch insert complete: %d new, %d updated, %d failed", savedCount, updatedCount, failedCount)

	return savedCount, updatedCount, failedCount, errors
}

// GetAllPrices retrieves all gold prices with optional filters
func (s *GoldScraperService) GetAllPrices(filter models.GoldPricingHistoryFilter) ([]models.GoldPricingHistory, error) {
	return s.repo.GetAll(filter)
}

// GetLatestPrices retrieves the latest price for each gold type
func (s *GoldScraperService) GetLatestPrices() ([]models.GoldPricingHistory, error) {
	return s.repo.GetLatest()
}

// GetPriceByID retrieves a gold price by ID
func (s *GoldScraperService) GetPriceByID(id int) (*models.GoldPricingHistory, error) {
	return s.repo.GetByID(id)
}

// Helper functions

// cleanText removes extra whitespace and newlines from text
func cleanText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	// Replace multiple spaces with single space
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}
	return text
}

// cleanPrice removes currency symbols and formats price string
func cleanPrice(price string) string {
	price = cleanText(price)
	// Remove common currency symbols and separators
	price = strings.ReplaceAll(price, "Rp", "")
	price = strings.ReplaceAll(price, "IDR", "")
	price = strings.ReplaceAll(price, ".", "")
	price = strings.ReplaceAll(price, ",", "")
	price = strings.TrimSpace(price)
	return price
}

// parsePrice converts a price string to int64
// Handles formats like "Rp1.234.567" or "1234567"
func parsePrice(priceStr string) (int64, error) {
	if priceStr == "" {
		return 0, fmt.Errorf("empty price string")
	}

	// Clean the price string
	cleaned := cleanPrice(priceStr)

	// Convert to int64
	price, err := strconv.ParseInt(cleaned, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price '%s': %w", priceStr, err)
	}

	if price < 0 {
		return 0, fmt.Errorf("negative price not allowed: %d", price)
	}

	return price, nil
}

// detectCategory determines the category based on product name
func detectCategory(productName string) (models.GoldCategory, error) {
	productLower := strings.ToLower(productName)

	// Check for most specific patterns first (order matters!)

	// Check for liontin/pendant with batik (most specific)
	if (strings.Contains(productLower, "liontin") && strings.Contains(productLower, "batik")) ||
		(strings.Contains(productLower, "pendant") && strings.Contains(productLower, "batik")) {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryLiontinBatikSeriIII, productName)
		return models.GoldCategoryLiontinBatikSeriIII, nil
	}

	// Check for specific batik series
	if strings.Contains(productLower, "batik seri iii") || strings.Contains(productLower, "batik seri 3") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryEmasBatanganBatikSeriIII, productName)
		return models.GoldCategoryEmasBatanganBatikSeriIII, nil
	}

	// Check for gift series
	if strings.Contains(productLower, "gift series") || strings.Contains(productLower, "gift") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryEmasBatanganGiftSeries, productName)
		return models.GoldCategoryEmasBatanganGiftSeries, nil
	}

	// Check for Idul Fitri/Lebaran
	if strings.Contains(productLower, "idul fitri") ||
		strings.Contains(productLower, "lebaran") ||
		strings.Contains(productLower, "selamat idul fitri") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryEmasBatanganSelamatIdulFitri, productName)
		return models.GoldCategoryEmasBatanganSelamatIdulFitri, nil
	}

	// Check for Imlek/Chinese New Year
	if strings.Contains(productLower, "imlek") ||
		strings.Contains(productLower, "chinese new year") ||
		strings.Contains(productLower, "tahun baru imlek") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryEmasBatanganImlek, productName)
		return models.GoldCategoryEmasBatanganImlek, nil
	}

	// Check for heritage silver
	if (strings.Contains(productLower, "perak") || strings.Contains(productLower, "silver")) &&
		strings.Contains(productLower, "heritage") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryPerakHeritage, productName)
		return models.GoldCategoryPerakHeritage, nil
	}

	// Check for silver/perak (general)
	if strings.Contains(productLower, "perak") || strings.Contains(productLower, "silver") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryPerakMurni, productName)
		return models.GoldCategoryPerakMurni, nil
	}

	// Check for liontin/pendant (general)
	if strings.Contains(productLower, "liontin") || strings.Contains(productLower, "pendant") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryLiontinBatikSeriIII, productName)
		return models.GoldCategoryLiontinBatikSeriIII, nil
	}

	// Check for batik (general - after more specific checks)
	if strings.Contains(productLower, "batik") {
		log.Printf("🏷️  Detected category '%s' from product: %s", models.GoldCategoryEmasBatanganBatikSeriIII, productName)
		return models.GoldCategoryEmasBatanganBatikSeriIII, nil
	}

	// Default to standard gold bars
	log.Printf("🏷️  Using default category '%s' for product: %s", models.GoldCategoryEmasBatangan, productName)
	return models.GoldCategoryEmasBatangan, nil
}
