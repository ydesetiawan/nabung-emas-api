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
	Unit      string
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
			unit := strings.TrimSpace(row.ChildText("td:nth-child(4)"))

			// Skip empty rows or header rows
			if goldType == "" || goldType == "Jenis" || goldType == "Type" {
				return
			}

			// Clean up the data
			goldType = cleanText(goldType)
			buyPrice = cleanPrice(buyPrice)
			sellPrice = cleanPrice(sellPrice)
			unit = cleanText(unit)

			// Only add if we have valid data
			if goldType != "" && (buyPrice != "" || sellPrice != "") {
				scrapedData = append(scrapedData, ScrapedGoldData{
					GoldType:  goldType,
					BuyPrice:  buyPrice,
					SellPrice: sellPrice,
					Unit:      unit,
				})

				log.Printf("✅ Scraped: %s - Buy: %s, Sell: %s, Unit: %s", goldType, buyPrice, sellPrice, unit)
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
				unit := ""
				if len(cells) >= 4 {
					unit = cleanText(cells[3])
				}

				if goldType != "" && goldType != "Jenis" && goldType != "Type" {
					scrapedData = append(scrapedData, ScrapedGoldData{
						GoldType:  goldType,
						BuyPrice:  buyPrice,
						SellPrice: sellPrice,
						Unit:      unit,
					})

					log.Printf("✅ Scraped (alt): %s - Buy: %s, Sell: %s, Unit: %s", goldType, buyPrice, sellPrice, unit)
				}
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
	for _, item := range data {
		createData = append(createData, models.GoldPricingHistoryCreate{
			GoldType:  item.GoldType,
			BuyPrice:  item.BuyPrice,
			SellPrice: item.SellPrice,
			Unit:      item.Unit,
			Source:    source,
		})
	}

	// Save to database using batch insert
	return s.repo.CreateBatch(createData)
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
