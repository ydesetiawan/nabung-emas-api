package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GoldCategory represents the category enum
type GoldCategory string

const (
	EmasBatangan                 GoldCategory = "emas_batangan"
	EmasBatanganGiftSeries       GoldCategory = "emas_batangan_gift_series"
	EmasBatanganSelamatIdulFitri GoldCategory = "emas_batangan_selamat_idul_fitri"
	EmasBatanganImlek            GoldCategory = "emas_batangan_imlek"
	EmasBatanganBatikSeriIII     GoldCategory = "emas_batangan_batik_seri_iii"
	PerakMurni                   GoldCategory = "perak_murni"
	PerakHeritage                GoldCategory = "perak_heritage"
	LontinBatikSeriIII           GoldCategory = "liontin_batik_seri_iii"
)

// Source represents the source enum
type Source string

const (
	SourceAntam     Source = "antam"
	SourceUBS       Source = "ubs"
	SourceGalery24  Source = "galery24"
	SourcePegadaian Source = "pegadaian"
)

// GoldPrice represents a gold price record
type GoldPrice struct {
	ID          int          `json:"id,omitempty"`
	PricingDate time.Time    `json:"pricing_date"`
	GoldType    string       `json:"gold_type"`
	BuyPrice    int          `json:"buy_price"`
	SellPrice   int          `json:"sell_price"`
	Source      Source       `json:"source"`
	Description string       `json:"description"`
	Category    GoldCategory `json:"category"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
}

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
func (s *AntamScraperService) Scrape() ([]GoldPrice, error) {
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
func (s *AntamScraperService) extractPrices(doc *goquery.Document) ([]GoldPrice, error) {
	var prices []GoldPrice
	pricingDate := s.getPricingDate()

	// Look for price tables or cards
	// The website likely uses a table or card layout for different gold types
	doc.Find(".product-card, .price-item, tr.price-row, .gold-price-item").Each(func(i int, sel *goquery.Selection) {
		price, err := s.parsePrice(sel, pricingDate)
		if err == nil {
			prices = append(prices, price)
		}
	})

	log.Println(prices)

	// Alternative: Look for JSON data in script tags
	if len(prices) == 0 {
		doc.Find("script[type='application/json'], script#__NEXT_DATA__").Each(func(i int, sel *goquery.Selection) {
			jsonData := sel.Text()
			extractedPrices := s.parseJSONData(jsonData, pricingDate)
			prices = append(prices, extractedPrices...)
		})
	}

	// Fallback: Parse table structure
	if len(prices) == 0 {
		prices = s.parseTableStructure(doc, pricingDate)
	}

	return prices, nil
}

// parsePrice parses a single price element
func (s *AntamScraperService) parsePrice(sel *goquery.Selection, pricingDate time.Time) (GoldPrice, error) {
	var price GoldPrice

	// Extract gold type (weight)
	goldType := strings.TrimSpace(sel.Find(".weight, .gram, .product-name").First().Text())
	if goldType == "" {
		goldType, _ = sel.Attr("data-weight")
	}

	// Extract category from product name
	productName := strings.TrimSpace(sel.Find(".product-title, .category-name").First().Text())
	category := s.detectCategory(productName)

	// Extract buy price (harga beli)
	buyPriceStr := strings.TrimSpace(sel.Find(".buy-price, .harga-beli").First().Text())
	buyPrice, err := s.parseRupiah(buyPriceStr)
	if err != nil {
		return price, fmt.Errorf("failed to parse buy price: %w", err)
	}

	// Extract sell price (harga jual)
	sellPriceStr := strings.TrimSpace(sel.Find(".sell-price, .harga-jual").First().Text())
	sellPrice, err := s.parseRupiah(sellPriceStr)
	if err != nil {
		return price, fmt.Errorf("failed to parse sell price: %w", err)
	}

	price = GoldPrice{
		PricingDate: pricingDate,
		GoldType:    goldType,
		BuyPrice:    buyPrice,
		SellPrice:   sellPrice,
		Source:      SourceAntam,
		Description: productName,
		Category:    category,
	}

	return price, nil
}

// parseTableStructure parses traditional table structure
func (s *AntamScraperService) parseTableStructure(doc *goquery.Document, pricingDate time.Time) []GoldPrice {
	var prices []GoldPrice

	log.Println("1111")
	log.Println(doc)

	doc.Find("table tbody tr, .price-table tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 3 {
			return
		}

		goldType := strings.TrimSpace(cells.Eq(0).Text())
		buyPriceStr := strings.TrimSpace(cells.Eq(1).Text())
		sellPriceStr := strings.TrimSpace(cells.Eq(2).Text())

		buyPrice, errBuy := s.parseRupiah(buyPriceStr)
		sellPrice, errSell := s.parseRupiah(sellPriceStr)

		if errBuy == nil && errSell == nil && goldType != "" {
			prices = append(prices, GoldPrice{
				PricingDate: pricingDate,
				GoldType:    goldType,
				BuyPrice:    buyPrice,
				SellPrice:   sellPrice,
				Source:      SourceAntam,
				Category:    EmasBatangan, // Default category
			})
		}
	})

	return prices
}

// parseJSONData extracts prices from JSON data
func (s *AntamScraperService) parseJSONData(jsonStr string, pricingDate time.Time) []GoldPrice {
	var prices []GoldPrice

	// Try to extract JSON object
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return prices
	}

	// Navigate through the JSON structure to find price data
	// This is placeholder logic - adjust based on actual JSON structure
	if products, ok := data["products"].([]interface{}); ok {
		for _, item := range products {
			if product, ok := item.(map[string]interface{}); ok {
				price := GoldPrice{
					PricingDate: pricingDate,
					Source:      SourceAntam,
					Category:    EmasBatangan,
				}

				if weight, ok := product["weight"].(string); ok {
					price.GoldType = weight
				}
				if buyPrice, ok := product["buyPrice"].(float64); ok {
					price.BuyPrice = int(buyPrice)
				}
				if sellPrice, ok := product["sellPrice"].(float64); ok {
					price.SellPrice = int(sellPrice)
				}

				if price.GoldType != "" && price.BuyPrice > 0 && price.SellPrice > 0 {
					prices = append(prices, price)
				}
			}
		}
	}

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
func (s *AntamScraperService) detectCategory(productName string) GoldCategory {
	productName = strings.ToLower(productName)

	switch {
	case strings.Contains(productName, "gift series"):
		return EmasBatanganGiftSeries
	case strings.Contains(productName, "idul fitri"), strings.Contains(productName, "selamat"):
		return EmasBatanganSelamatIdulFitri
	case strings.Contains(productName, "imlek"):
		return EmasBatanganImlek
	case strings.Contains(productName, "batik seri iii"), strings.Contains(productName, "batik seri 3"):
		return EmasBatanganBatikSeriIII
	case strings.Contains(productName, "perak murni"):
		return PerakMurni
	case strings.Contains(productName, "heritage"):
		return PerakHeritage
	case strings.Contains(productName, "liontin batik"):
		return LontinBatikSeriIII
	default:
		return EmasBatangan
	}
}

// getPricingDate returns today's date at 00:00:00 UTC
func (s *AntamScraperService) getPricingDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// Example usage function
func ExampleUsage() {
	scraper := NewAntamScraperService()

	prices, err := scraper.Scrape()
	if err != nil {
		fmt.Printf("Error scraping: %v\n", err)
		return
	}

	fmt.Printf("Successfully scraped %d prices\n", len(prices))
	for _, price := range prices {
		fmt.Printf("Gold Type: %s, Buy: %d, Sell: %d, Category: %s\n",
			price.GoldType, price.BuyPrice, price.SellPrice, price.Category)
	}
}
