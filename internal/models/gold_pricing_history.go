package models

import (
	"fmt"
	"time"
)

// GoldSource represents the vendor/brand of gold
type GoldSource string

// Gold source enum values matching database enum
const (
	GoldSourceGaleri24            GoldSource = "GALERI_24"
	GoldSourceDinarG24            GoldSource = "DINAR_G24"
	GoldSourceBabyGaleri24        GoldSource = "BABY_GALERI_24"
	GoldSourceAntam               GoldSource = "ANTAM"
	GoldSourceUBS                 GoldSource = "UBS"
	GoldSourceAntamMuliaRetro     GoldSource = "ANTAM_MULIA_RETRO"
	GoldSourceAntamNonPegadaian   GoldSource = "ANTAM_NON_PEGADAIAN"
	GoldSourceLotusArchi          GoldSource = "LOTUS_ARCHI"
	GoldSourceUBSDisney           GoldSource = "UBS_DISNEY"
	GoldSourceUBSElsa             GoldSource = "UBS_ELSA"
	GoldSourceUBSAnna             GoldSource = "UBS_ANNA"
	GoldSourceUBSMickeyFullbody   GoldSource = "UBS_MICKEY_FULLBODY"
	GoldSourceLotusArchiGift      GoldSource = "LOTUS_ARCHI_GIFT"
	GoldSourceUBSHelloKitty       GoldSource = "UBS_HELLO_KITTY"
	GoldSourceBabySeriesTumbuhan  GoldSource = "BABY_SERIES_TUMBUHAN"
	GoldSourceBabySeriesInvestasi GoldSource = "BABY_SERIES_INVESTASI"
	GoldSourceBatikSeries         GoldSource = "BATIK_SERIES"
)

// GoldPricingHistory represents a gold price record from vendors
type GoldPricingHistory struct {
	ID          int        `json:"id" db:"id"`
	PricingDate time.Time  `json:"pricing_date" db:"pricing_date"`
	GoldType    string     `json:"gold_type" db:"gold_type"`   // Weight in grams (e.g., "0.5", "1", "2", "5", etc.)
	BuyPrice    string     `json:"buy_price" db:"buy_price"`   // Harga Buyback (e.g., "Rp1.132.000")
	SellPrice   string     `json:"sell_price" db:"sell_price"` // Harga Jual (e.g., "Rp1.271.000")
	Source      GoldSource `json:"source" db:"source"`         // Vendor name
	ScrapedAt   time.Time  `json:"scraped_at" db:"scraped_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// GoldPricingHistoryCreate is the DTO for creating a new gold pricing record
type GoldPricingHistoryCreate struct {
	PricingDate time.Time  `json:"pricing_date" validate:"required"`
	GoldType    string     `json:"gold_type" validate:"required"`
	BuyPrice    string     `json:"buy_price" validate:"required"`
	SellPrice   string     `json:"sell_price" validate:"required"`
	Source      GoldSource `json:"source" validate:"required"`
}

// GoldPricingHistoryFilter represents filter options for querying gold prices
type GoldPricingHistoryFilter struct {
	GoldType  string     `json:"gold_type,omitempty"`
	Source    GoldSource `json:"source,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}

// GoldPricingStats represents statistics about gold pricing data
type GoldPricingStats struct {
	TotalRecords      int       `json:"total_records"`
	UniqueVendors     int       `json:"unique_vendors"`
	UniqueGoldTypes   int       `json:"unique_gold_types"`
	LatestScrapedAt   time.Time `json:"latest_scraped_at"`
	OldestPricingDate time.Time `json:"oldest_pricing_date"`
	LatestPricingDate time.Time `json:"latest_pricing_date"`
}

// ScrapeResult represents the result of a scraping operation
type ScrapeResult struct {
	Success      bool                 `json:"success"`
	Message      string               `json:"message"`
	PricingDate  time.Time            `json:"pricing_date"`
	TotalScraped int                  `json:"total_scraped"`
	SavedCount   int                  `json:"saved_count"`
	UpdatedCount int                  `json:"updated_count"`
	FailedCount  int                  `json:"failed_count"`
	Data         []GoldPricingHistory `json:"data,omitempty"`
	Errors       []string             `json:"errors,omitempty"`
	Duration     string               `json:"duration"`
}

// VendorNameMapping maps website vendor names to enum values
var VendorNameMapping = map[string]GoldSource{
	"GALERI 24":             GoldSourceGaleri24,
	"DINAR G24":             GoldSourceDinarG24,
	"BABY GALERI 24":        GoldSourceBabyGaleri24,
	"ANTAM":                 GoldSourceAntam,
	"UBS":                   GoldSourceUBS,
	"ANTAM MULIA RETRO":     GoldSourceAntamMuliaRetro,
	"ANTAM NON PEGADAIAN":   GoldSourceAntamNonPegadaian,
	"LOTUS ARCHI":           GoldSourceLotusArchi,
	"UBS DISNEY":            GoldSourceUBSDisney,
	"UBS ELSA":              GoldSourceUBSElsa,
	"UBS ANNA":              GoldSourceUBSAnna,
	"UBS MICKEY FULLBODY":   GoldSourceUBSMickeyFullbody,
	"LOTUS ARCHI GIFT":      GoldSourceLotusArchiGift,
	"UBS HELLO KITTY":       GoldSourceUBSHelloKitty,
	"BABY SERIES TUMBUHAN":  GoldSourceBabySeriesTumbuhan,
	"BABY SERIES INVESTASI": GoldSourceBabySeriesInvestasi,
	"BATIK SERIES":          GoldSourceBatikSeries,
}

// IndonesianMonthMapping maps Indonesian month names to numbers
var IndonesianMonthMapping = map[string]time.Month{
	"Januari":   time.January,
	"Februari":  time.February,
	"Maret":     time.March,
	"April":     time.April,
	"Mei":       time.May,
	"Juni":      time.June,
	"Juli":      time.July,
	"Agustus":   time.August,
	"September": time.September,
	"Oktober":   time.October,
	"November":  time.November,
	"Desember":  time.December,
}

// Scan implements sql.Scanner for GoldSource
func (gs *GoldSource) Scan(value interface{}) error {
	if value == nil {
		*gs = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*gs = GoldSource(v)
	case string:
		*gs = GoldSource(v)
	default:
		return fmt.Errorf("cannot scan type %T into GoldSource", value)
	}

	return nil
}

// Value implements driver.Valuer for GoldSource
func (gs GoldSource) Value() (interface{}, error) {
	return string(gs), nil
}

// String returns the string representation of GoldSource
func (gs GoldSource) String() string {
	return string(gs)
}

// IsValid checks if the GoldSource is a valid enum value
func (gs GoldSource) IsValid() bool {
	switch gs {
	case GoldSourceGaleri24, GoldSourceDinarG24, GoldSourceBabyGaleri24,
		GoldSourceAntam, GoldSourceUBS, GoldSourceAntamMuliaRetro,
		GoldSourceAntamNonPegadaian, GoldSourceLotusArchi, GoldSourceUBSDisney,
		GoldSourceUBSElsa, GoldSourceUBSAnna, GoldSourceUBSMickeyFullbody,
		GoldSourceLotusArchiGift, GoldSourceUBSHelloKitty, GoldSourceBabySeriesTumbuhan,
		GoldSourceBabySeriesInvestasi, GoldSourceBatikSeries:
		return true
	}
	return false
}
