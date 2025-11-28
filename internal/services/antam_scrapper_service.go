package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"nabung-emas-api/internal/models"

	"github.com/PuerkitoBio/goquery"
)

// AntamScraper handles scraping from Logam Mulia website
type AntamScraperService struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// NewAntamScraper creates a new Antam scraper instance
func NewAntamScraperService() *AntamScraperService {
	return &AntamScraperService{
		BaseURL: "https://www.logammulia.com/id/harga-emas-hari-ini",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
}

// Scrape fetches and parses gold prices from Logam Mulia website
func (s *AntamScraperService) Scrape() ([]models.GoldPricingHistory, error) {
	// Create request
	req, err := http.NewRequest("GET", s.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7")

	// Execute request with retry logic
	var resp *http.Response
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = s.HTTPClient.Do(req)
		if err == nil && resp.StatusCode == 200 {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(time.Duration(attempt) * 2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch page after retries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract prices
	prices, err := s.extractPrices(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract prices: %w", err)
	}

	return prices, nil
}

// extractPrices extracts gold prices from the HTML document
func (s *AntamScraperService) extractPrices(doc *goquery.Document) ([]models.GoldPricingHistory, error) {
	var prices []models.GoldPricingHistory
	pricingDate := s.getPricingDate()

	// Fallback: Parse table structure
	if len(prices) == 0 {
		prices = s.parseTableStructure(doc, pricingDate)
	}

	return prices, nil
}

// parseTableStructure parses traditional table structure
func (s *AntamScraperService) parseTableStructure(doc *goquery.Document, pricingDate time.Time) []models.GoldPricingHistory {
	var prices []models.GoldPricingHistory
	var currentCategory models.GoldCategory = models.EmasBatangan // Default

	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		// Check if it's a category header (no tds, likely ths)
		if row.Find("td").Length() == 0 {
			headerText := strings.TrimSpace(row.Text())
			// Filter out the main header row "Berat Harga Dasar ..."
			if strings.Contains(strings.ToLower(headerText), "berat") {
				return
			}

			if headerText != "" {
				currentCategory = s.detectCategory(headerText)
			}
			return
		}

		cells := row.Find("td")
		if cells.Length() < 3 {
			return
		}

		goldType := strings.TrimSpace(cells.Eq(0).Text())
		basePriceStr := strings.TrimSpace(cells.Eq(1).Text())
		sellPriceStr := strings.TrimSpace(cells.Eq(2).Text())

		basePrice, errBase := s.parseRupiah(basePriceStr)
		sellPrice, errSell := s.parseRupiah(sellPriceStr)

		if errBase == nil && errSell == nil && goldType != "" {
			prices = append(prices, models.GoldPricingHistory{
				PricingDate: pricingDate,
				GoldType:    goldType,
				BasePrice:   basePrice,
				BuyPrice:    0,
				SellPrice:   sellPrice,
				IncludeTax:  true,
				Source:      models.SourceAntam,
				Category:    currentCategory,
			})
		}
	})

	return prices
}

// parseRupiah converts Indonesian rupiah string to integer
// Example: "Rp 1.235.892" -> 1235892
func (s *AntamScraperService) parseRupiah(priceStr string) (int, error) {
	if priceStr == "" {
		return 0, fmt.Errorf("empty price string")
	}

	// Remove common prefixes and whitespace
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.ReplaceAll(priceStr, "Rp", "")
	priceStr = strings.ReplaceAll(priceStr, "IDR", "")
	priceStr = strings.TrimSpace(priceStr)

	// Remove dots (thousand separators)
	priceStr = strings.ReplaceAll(priceStr, ".", "")

	// Remove commas
	priceStr = strings.ReplaceAll(priceStr, ",", "")

	// Extract only numbers
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(priceStr, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("no numbers found in price string: %s", priceStr)
	}

	priceStr = strings.Join(matches, "")

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert to integer: %w", err)
	}

	return price, nil
}

// detectCategory determines the category based on product name
func (s *AntamScraperService) detectCategory(productName string) models.GoldCategory {
	productName = strings.ToLower(productName)

	switch {
	case strings.Contains(productName, "gift series"):
		return models.EmasBatanganGiftSeries
	case strings.Contains(productName, "idul fitri"), strings.Contains(productName, "selamat"):
		return models.EmasBatanganSelamatIdulFitri
	case strings.Contains(productName, "imlek"):
		return models.EmasBatanganImlek
	case strings.Contains(productName, "batik seri iii"), strings.Contains(productName, "batik seri 3"):
		return models.EmasBatanganBatikSeriIII
	case strings.Contains(productName, "perak murni"):
		return models.PerakMurni
	case strings.Contains(productName, "heritage"):
		return models.PerakHeritage
	case strings.Contains(productName, "liontin batik"):
		return models.LontinBatikSeriIII
	default:
		return models.EmasBatangan
	}
}

// getPricingDate returns today's date at 00:00:00 UTC
func (s *AntamScraperService) getPricingDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
