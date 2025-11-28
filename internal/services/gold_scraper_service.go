package services

import (
	"fmt"
	"log"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// GoldScraperService handles web scraping of gold prices
type GoldScraperService struct {
	repo *repositories.GoldPricingHistoryRepository
}

// NewGoldScraperService creates a new instance of GoldScraperService
func NewGoldScraperService(repo *repositories.GoldPricingHistoryRepository) *GoldScraperService {
	return &GoldScraperService{repo: repo}
}

// ScrapedGoldData represents the data scraped from the website
type ScrapedGoldData struct {
	GoldType  string
	BuyPrice  string
	SellPrice string
}

// ScrapeResult represents the result of a scraping operation
type ScrapeResult struct {
	Success      bool                        `json:"success"`
	Message      string                      `json:"message"`
	ScrapedCount int                         `json:"scraped_count"`
	SavedCount   int                         `json:"saved_count"`
	Data         []models.GoldPricingHistory `json:"data,omitempty"`
	Errors       []string                    `json:"errors,omitempty"`
}

// ScrapeLogamMulia scrapes gold prices from logammulia.com
func (s *GoldScraperService) ScrapeLogamMulia() (*ScrapeResult, error) {
	log.Println("🕷️  Starting gold price scraping from logammulia.com...")

	// Create a new collector with configuration
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
		Delay:       1 * time.Second,
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

			// Filter: Only scrape "Emas Batangan" (gold bars)
			// Skip jewelry, coins, and other products
			// goldTypeLower := strings.ToLower(goldType)
			// if !isGoldBar(goldTypeLower) {
			// 	log.Printf("⏭️  Skipping non-gold-bar item: %s", goldType)
			// 	return
			// }

			// Only add if we have valid data
			if goldType != "" && (buyPrice != "" || sellPrice != "") {
				scrapedData = append(scrapedData, ScrapedGoldData{
					GoldType:  goldType,
					BuyPrice:  buyPrice,
					SellPrice: sellPrice,
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

				// Filter: Only scrape gold bars
				goldTypeLower := strings.ToLower(goldType)
				if !isGoldBar(goldTypeLower) {
					log.Printf("⏭️  Skipping non-gold-bar item (alt): %s", goldType)
					return
				}

				scrapedData = append(scrapedData, ScrapedGoldData{
					GoldType:  goldType,
					BuyPrice:  buyPrice,
					SellPrice: sellPrice,
				})

				log.Printf("✅ Scraped (alt): %s - Buy: %s, Sell: %s", goldType, buyPrice, sellPrice)
			}
		})
	})

	// Response handler
	c.OnResponse(func(r *colly.Response) {
		log.Printf("✅ Received response from %s (Status: %d)", r.Request.URL, r.StatusCode)
	})

	// Visit the URL
	err := c.Visit("https://logammulia.com/id/harga-emas-hari-ini")
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to visit website: %v", err)
		log.Printf("❌ %s", errorMsg)
		return &ScrapeResult{
			Success:      false,
			Message:      errorMsg,
			ScrapedCount: 0,
			SavedCount:   0,
			Errors:       append(errors, errorMsg),
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
			ScrapedCount: 0,
			SavedCount:   0,
			Errors:       append(errors, "No data found in expected HTML structure"),
		}, nil
	}

	// Save scraped data to database
	savedData, saveErr := s.SaveScrapedData(scrapedData, models.GoldSourceAntam)
	if saveErr != nil {
		errorMsg := fmt.Sprintf("Failed to save data: %v", saveErr)
		log.Printf("❌ %s", errorMsg)
		return &ScrapeResult{
			Success:      false,
			Message:      errorMsg,
			ScrapedCount: len(scrapedData),
			SavedCount:   0,
			Errors:       append(errors, errorMsg),
		}, saveErr
	}

	log.Printf("💾 Successfully saved %d records to database", len(savedData))

	return &ScrapeResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully scraped and saved %d gold prices", len(savedData)),
		ScrapedCount: len(scrapedData),
		SavedCount:   len(savedData),
		Data:         savedData,
		Errors:       errors,
	}, nil
}

// SaveScrapedData saves the scraped data to the database
func (s *GoldScraperService) SaveScrapedData(data []ScrapedGoldData, source models.GoldSource) ([]models.GoldPricingHistory, error) {
	if len(data) == 0 {
		return []models.GoldPricingHistory{}, nil
	}

	// Convert scraped data to create models
	createData := make([]models.GoldPricingHistoryCreate, 0, len(data))
	pricingDate := time.Now().Truncate(24 * time.Hour) // Today's date at midnight

	for _, item := range data {
		// Convert string price to int64
		sellPrice, err := parsePrice(item.SellPrice)
		if err != nil {
			log.Printf("⚠️  Failed to parse sell price for %s: %v", item.GoldType, err)
			continue
		}

		createData = append(createData, models.GoldPricingHistoryCreate{
			PricingDate: pricingDate,
			GoldType:    item.GoldType,
			SellPrice:   sellPrice, // Buy price will be calculated as 94% of sell price
			Source:      source,
		})
	}

	// Save to database using batch insert (returns savedCount, updatedCount, error)
	savedCount, updatedCount, err := s.repo.CreateBatch(createData)
	if err != nil {
		return nil, err
	}

	log.Printf("💾 Batch insert complete: %d new, %d updated", savedCount, updatedCount)

	// Fetch the saved data to return
	return s.repo.GetByDate(pricingDate)
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
	// Clean the price string
	cleaned := cleanPrice(priceStr)

	// Convert to int64
	price, err := fmt.Sscanf(cleaned, "%d", new(int64))
	if err != nil || price != 1 {
		// Try alternative parsing
		var result int64
		_, err := fmt.Sscanf(cleaned, "%d", &result)
		if err != nil {
			return 0, fmt.Errorf("failed to parse price '%s': %w", priceStr, err)
		}
		return result, nil
	}

	var result int64
	fmt.Sscanf(cleaned, "%d", &result)
	return result, nil
}

// isGoldBar checks if the item is a gold bar (Emas Batangan)
// Returns true for gold bars, false for jewelry, coins, and other products
func isGoldBar(goldType string) bool {
	goldType = strings.ToLower(goldType)

	// Keywords that indicate it's NOT a gold bar (jewelry, coins, etc.)
	excludeKeywords := []string{
		"anting",      // earrings
		"cincin",      // ring
		"gelang",      // bracelet
		"kalung",      // necklace
		"liontin",     // pendant
		"perhiasan",   // jewelry
		"koin",        // coin
		"dinar",       // dinar coin
		"dirham",      // dirham coin
		"medali",      // medal
		"gift",        // gift items
		"souvenir",    // souvenir
		"hello kitty", // character items
		"disney",      // character items
	}

	// Check if it contains any exclude keywords
	for _, keyword := range excludeKeywords {
		if strings.Contains(goldType, keyword) {
			return false
		}
	}

	// Keywords that indicate it IS a gold bar
	includeKeywords := []string{
		"batang",      // bar
		"logam mulia", // logam mulia brand
		"antam",       // antam brand
		"ubs",         // ubs brand
		"gram",        // weight indicator
		"gr",          // weight indicator
		"g ",          // weight indicator
	}

	// Check if it contains any include keywords
	for _, keyword := range includeKeywords {
		if strings.Contains(goldType, keyword) {
			return true
		}
	}

	// If it's just a number followed by gram/gr/g, it's likely a gold bar
	// Examples: "1 gram", "5 gr", "10 g", "100 gram"
	if strings.Contains(goldType, "gram") ||
		strings.Contains(goldType, " gr") ||
		strings.Contains(goldType, " g") {
		return true
	}

	// Default: exclude if we're not sure
	return false
}
