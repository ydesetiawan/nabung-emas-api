package services

import (
	"encoding/json"
	"fmt"
	"log"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Galeri24ScraperService handles scraping gold prices from Galeri24
type Galeri24ScraperService struct {
	repo *repositories.GoldPricingHistoryRepository
}

// NewGaleri24ScraperService creates a new instance of Galeri24ScraperService
func NewGaleri24ScraperService(repo *repositories.GoldPricingHistoryRepository) *Galeri24ScraperService {
	return &Galeri24ScraperService{repo: repo}
}

// ScrapeGaleri24 scrapes gold prices from https://galeri24.co.id/harga-emas
func (s *Galeri24ScraperService) ScrapeGaleri24() (*models.ScrapeResult, error) {
	startTime := time.Now()
	log.Println("🕷️  Starting Galeri24 gold price scraping...")

	result := &models.ScrapeResult{
		Success: false,
		Errors:  []string{},
	}

	// Launch browser
	log.Println("🌐 Launching headless browser...")
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// Create page
	page := browser.MustPage("")
	defer page.MustClose()

	// Enable network tracking to capture API calls
	router := page.HijackRequests()
	defer router.Stop()

	// Capture API responses
	apiResponses := make(map[string]string)

	router.MustAdd("*", func(ctx *rod.Hijack) {
		// Log all requests
		if strings.Contains(ctx.Request.URL().String(), "api") ||
			strings.Contains(ctx.Request.URL().String(), "price") ||
			strings.Contains(ctx.Request.URL().String(), "product") {
			log.Printf("📡 API Request: %s", ctx.Request.URL().String())
		}

		// Continue the request
		ctx.MustLoadResponse()

		// Capture response if it looks like price data
		if strings.Contains(ctx.Request.URL().String(), "api") {
			body := ctx.Response.Body()
			apiResponses[ctx.Request.URL().String()] = body
			log.Printf("💾 Captured API response: %d bytes from %s", len(body), ctx.Request.URL().String())
		}
	})

	go router.Run()

	// Navigate to the page
	log.Println("🌐 Navigating to https://galeri24.co.id/harga-emas...")
	page.MustNavigate("https://galeri24.co.id/harga-emas")

	// Wait for page to load
	log.Println("⏳ Waiting for page to load...")
	page.MustWaitLoad()

	// Wait for JavaScript and API calls
	log.Println("⏳ Waiting for API calls to complete...")
	time.Sleep(8 * time.Second)

	// Get page HTML as fallback
	html := page.MustHTML()

	// Debug: Save HTML to file for inspection
	log.Printf("📄 Page HTML length: %d characters", len(html))
	log.Printf("📊 Captured %d API responses", len(apiResponses))

	// Try to extract data from API responses first
	var scrapedData []models.GoldPricingHistoryCreate
	pricingDate := time.Now()

	// Log all captured API URLs
	for url := range apiResponses {
		log.Printf("🔍 API URL: %s", url)
	}

	// Try to parse API responses for price data
	for url, body := range apiResponses {
		log.Printf("🔍 Analyzing API response from: %s", url)

		// Try to parse as JSON
		var jsonData interface{}
		if err := json.Unmarshal([]byte(body), &jsonData); err == nil {
			log.Printf("✅ Successfully parsed JSON from %s", url)
			// TODO: Extract price data from JSON
			// For now, just log that we got JSON
		}
	}

	// If no API data, try HTML parsing
	if len(scrapedData) == 0 {
		log.Println("⚠️  No data from API responses, trying HTML parsing...")

		// Extract pricing date
		var err error
		pricingDate, err = extractDateFromHTML(html)
		if err != nil {
			log.Printf("⚠️  Could not extract date, using today's date: %v", err)
			pricingDate = time.Now()
		}

		log.Printf("📅 Pricing date: %s", pricingDate.Format("2006-01-02"))
		result.PricingDate = pricingDate

		// Try JavaScript extraction
		scrapedData, err = extractGoldPrices(page, pricingDate) // Renamed from extractGoldPricesWithJS to match existing function
		if err != nil || len(scrapedData) == 0 {
			log.Printf("⚠️  JavaScript extraction failed or returned no data: %v", err)

			// Save HTML for debugging
			if err := saveHTMLToFile(html); err == nil {
				log.Println("💾 Saved HTML to /tmp/galeri24_debug.html for inspection")
			}

			result.Message = "Failed to extract gold prices"
			result.Errors = append(result.Errors, "No pricing data found in API responses or HTML")
			result.Duration = time.Since(startTime).String()
			return result, fmt.Errorf("no data extracted")
		}
	}

	if len(scrapedData) == 0 {
		result.Message = "No data scraped from website"
		result.Errors = append(result.Errors, "No pricing data found")
		result.Duration = time.Since(startTime).String()
		return result, fmt.Errorf("no data scraped")
	}

	log.Printf("📊 Scraped %d gold price records", len(scrapedData))
	result.TotalScraped = len(scrapedData)
	result.PricingDate = pricingDate

	// Save to database (with UPSERT logic - will replace duplicates)
	log.Println("💾 Saving scraped data to database...")
	savedCount, updatedCount, err := s.repo.CreateBatch(scrapedData)
	if err != nil {
		result.Message = "Failed to save data to database"
		result.Errors = append(result.Errors, err.Error())
		result.Duration = time.Since(startTime).String()
		return result, err
	}

	result.SavedCount = savedCount
	result.UpdatedCount = updatedCount
	result.FailedCount = result.TotalScraped - savedCount - updatedCount

	// Fetch saved data
	savedData, err := s.repo.GetByDate(pricingDate)
	if err == nil {
		result.Data = savedData
	}

	result.Success = true
	result.Message = fmt.Sprintf("Successfully scraped and saved %d records (%d new, %d updated)",
		savedCount+updatedCount, savedCount, updatedCount)
	result.Duration = time.Since(startTime).String()

	log.Printf("✅ Scraping completed: %d new, %d updated, %d failed in %s",
		savedCount, updatedCount, result.FailedCount, result.Duration)

	return result, nil
}

// extractDateFromHTML extracts the pricing date from HTML
func extractDateFromHTML(html string) (time.Time, error) {
	// Look for "Diperbarui" text
	re := regexp.MustCompile(`Diperbarui\s+\w+,\s+(\d+)\s+(\w+)\s+(\d{4})`)
	matches := re.FindStringSubmatch(html)

	if len(matches) != 4 {
		return time.Time{}, fmt.Errorf("date pattern not found in HTML")
	}

	day, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %s", matches[1])
	}

	monthName := matches[2]
	month, ok := models.IndonesianMonthMapping[monthName]
	if !ok {
		return time.Time{}, fmt.Errorf("unknown month: %s", monthName)
	}

	year, err := strconv.Atoi(matches[3])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year: %s", matches[3])
	}

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}

// extractGoldPrices extracts gold prices from the rendered page
func extractGoldPrices(page *rod.Page, pricingDate time.Time) ([]models.GoldPricingHistoryCreate, error) {
	var scrapedData []models.GoldPricingHistoryCreate

	// Try to find price elements using various selectors
	// The website likely uses cards or divs for each vendor

	// Get all text content and look for vendor names and prices
	html := page.MustHTML()

	// Look for vendor sections - they typically have vendor names followed by price data
	vendorPatterns := []string{
		"GALERI 24", "DINAR G24", "BABY GALERI 24", "ANTAM", "UBS",
		"ANTAM MULIA RETRO", "ANTAM NON PEGADAIAN", "LOTUS ARCHI",
		"UBS DISNEY", "UBS ELSA", "UBS ANNA", "UBS MICKEY FULLBODY",
		"LOTUS ARCHI GIFT", "UBS HELLO KITTY", "BABY SERIES TUMBUHAN",
		"BABY SERIES INVESTASI", "BATIK SERIES",
	}

	// Extract prices using JavaScript evaluation
	pricesJS := `
		() => {
			const results = [];
			// Try to find all price cards or sections
			const elements = document.querySelectorAll('[class*="card"], [class*="price"], [class*="vendor"], [class*="product"]');
			
			elements.forEach(el => {
				const text = el.innerText || el.textContent;
				if (text && (text.includes('Rp') || text.includes('gram'))) {
					results.push(text);
				}
			});
			
			return results;
		}
	`

	priceTexts, err := page.Eval(pricesJS)
	if err != nil {
		log.Printf("⚠️  JavaScript evaluation failed: %v", err)
		// Fall back to HTML parsing
		return parseHTMLForPrices(html, pricingDate, vendorPatterns)
	}

	// Process the extracted texts
	if priceTexts != nil && priceTexts.Value.Arr() != nil {
		for _, item := range priceTexts.Value.Arr() {
			text := item.String()
			// Parse each text block for vendor, weight, and prices
			prices := parseTextBlock(text, pricingDate)
			scrapedData = append(scrapedData, prices...)
		}
	}

	// If we didn't get data from JS, try HTML parsing
	if len(scrapedData) == 0 {
		return parseHTMLForPrices(html, pricingDate, vendorPatterns)
	}

	return scrapedData, nil
}

// parseHTMLForPrices parses HTML content for price data
func parseHTMLForPrices(html string, pricingDate time.Time, vendorPatterns []string) ([]models.GoldPricingHistoryCreate, error) {
	var scrapedData []models.GoldPricingHistoryCreate

	// Look for price patterns in HTML
	// Format: Rp1.234.000 or similar
	priceRe := regexp.MustCompile(`Rp[\d.,]+`)
	weightRe := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(?:gram|gr|g)\b`)

	// For each vendor, try to find their prices
	for _, vendorName := range vendorPatterns {
		source, ok := models.VendorNameMapping[vendorName]
		if !ok {
			continue
		}

		// Find vendor section in HTML
		vendorIndex := strings.Index(html, vendorName)
		if vendorIndex == -1 {
			continue
		}

		// Extract a section of HTML after the vendor name (next 5000 chars)
		section := html[vendorIndex:min(vendorIndex+5000, len(html))]

		// Find all prices in this section
		prices := priceRe.FindAllString(section, -1)
		weights := weightRe.FindAllStringSubmatch(section, -1)

		// Try to match weights with prices
		// Typically: weight, sell price, buy price
		for i := 0; i < len(weights) && i*2+1 < len(prices); i++ {
			weight := weights[i][1]
			sellPrice := prices[i*2]
			buyPrice := ""
			if i*2+1 < len(prices) {
				buyPrice = prices[i*2+1]
			}

			if buyPrice == "" {
				buyPrice = sellPrice // Use same price if buy price not found
			}

			scrapedData = append(scrapedData, models.GoldPricingHistoryCreate{
				PricingDate: pricingDate,
				GoldType:    weight,
				BuyPrice:    buyPrice,
				SellPrice:   sellPrice,
				Source:      source,
			})
		}
	}

	if len(scrapedData) == 0 {
		return nil, fmt.Errorf("no price data found in HTML")
	}

	return scrapedData, nil
}

// parseTextBlock parses a text block for price information
func parseTextBlock(text string, pricingDate time.Time) []models.GoldPricingHistoryCreate {
	var results []models.GoldPricingHistoryCreate

	// Determine vendor from text
	var source models.GoldSource
	vendorFound := false

	for vendorName, vendorSource := range models.VendorNameMapping {
		if strings.Contains(text, vendorName) {
			source = vendorSource
			vendorFound = true
			break
		}
	}

	if !vendorFound {
		return results
	}

	// Extract weights and prices
	priceRe := regexp.MustCompile(`Rp[\d.,]+`)
	weightRe := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(?:gram|gr|g)\b`)

	prices := priceRe.FindAllString(text, -1)
	weights := weightRe.FindAllStringSubmatch(text, -1)

	// Match weights with prices
	for i := 0; i < len(weights) && i*2+1 < len(prices); i++ {
		weight := weights[i][1]
		sellPrice := prices[i*2]
		buyPrice := prices[i*2+1]

		results = append(results, models.GoldPricingHistoryCreate{
			PricingDate: pricingDate,
			GoldType:    weight,
			BuyPrice:    buyPrice,
			SellPrice:   sellPrice,
			Source:      source,
		})
	}

	return results
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetAllPrices retrieves all gold prices with filters
func (s *Galeri24ScraperService) GetAllPrices(filter models.GoldPricingHistoryFilter) ([]models.GoldPricingHistory, error) {
	return s.repo.GetAll(filter)
}

// GetLatestPrices retrieves the latest prices for each gold type and source
func (s *Galeri24ScraperService) GetLatestPrices() ([]models.GoldPricingHistory, error) {
	return s.repo.GetLatest()
}

// GetPriceByID retrieves a gold price by ID
func (s *Galeri24ScraperService) GetPriceByID(id int) (*models.GoldPricingHistory, error) {
	return s.repo.GetByID(id)
}

// GetPricesByDate retrieves all prices for a specific date
func (s *Galeri24ScraperService) GetPricesByDate(date time.Time) ([]models.GoldPricingHistory, error) {
	return s.repo.GetByDate(date)
}

// GetStats retrieves statistics about gold pricing data
func (s *Galeri24ScraperService) GetStats() (*models.GoldPricingStats, error) {
	return s.repo.GetStats()
}

// saveHTMLToFile saves HTML content to a file for debugging
func saveHTMLToFile(html string) error {
	return os.WriteFile("/tmp/galeri24_debug.html", []byte(html), 0644)
}
