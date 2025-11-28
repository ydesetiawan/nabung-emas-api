package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// GoldSource represents the source of gold pricing data
type GoldSource string

const (
	GoldSourceAntam GoldSource = "antam"
	GoldSourceUSB   GoldSource = "usb"
)

// IsValid checks if the GoldSource is valid
func (gs GoldSource) IsValid() bool {
	switch gs {
	case GoldSourceAntam, GoldSourceUSB:
		return true
	}
	return false
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
func (gs GoldSource) Value() (driver.Value, error) {
	return string(gs), nil
}

// GoldCategory represents the category of gold/silver product
type GoldCategory string

const (
	GoldCategoryEmasBatangan                 GoldCategory = "emas_batangan"
	GoldCategoryEmasBatanganGiftSeries       GoldCategory = "emas_batangan_gift_series"
	GoldCategoryEmasBatanganSelamatIdulFitri GoldCategory = "emas_batangan_selamat_idul_fitri"
	GoldCategoryEmasBatanganImlek            GoldCategory = "emas_batangan_imlek"
	GoldCategoryEmasBatanganBatikSeriIII     GoldCategory = "emas_batangan_batik_seri_iii"
	GoldCategoryPerakMurni                   GoldCategory = "perak_murni"
	GoldCategoryPerakHeritage                GoldCategory = "perak_heritage"
	GoldCategoryLiontinBatikSeriIII          GoldCategory = "liontin_batik_seri_iii"
)

// IsValid checks if the GoldCategory is valid
func (gc GoldCategory) IsValid() bool {
	switch gc {
	case GoldCategoryEmasBatangan,
		GoldCategoryEmasBatanganGiftSeries,
		GoldCategoryEmasBatanganSelamatIdulFitri,
		GoldCategoryEmasBatanganImlek,
		GoldCategoryEmasBatanganBatikSeriIII,
		GoldCategoryPerakMurni,
		GoldCategoryPerakHeritage,
		GoldCategoryLiontinBatikSeriIII:
		return true
	}
	return false
}

// Scan implements sql.Scanner for GoldCategory
func (gc *GoldCategory) Scan(value interface{}) error {
	if value == nil {
		*gc = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*gc = GoldCategory(v)
	case string:
		*gc = GoldCategory(v)
	default:
		return fmt.Errorf("cannot scan type %T into GoldCategory", value)
	}

	return nil
}

// Value implements driver.Valuer for GoldCategory
func (gc GoldCategory) Value() (driver.Value, error) {
	return string(gc), nil
}

// GetDisplayName returns a human-readable display name for the category
func (gc GoldCategory) GetDisplayName() string {
	switch gc {
	case GoldCategoryEmasBatangan:
		return "Emas Batangan"
	case GoldCategoryEmasBatanganGiftSeries:
		return "Emas Batangan Gift Series"
	case GoldCategoryEmasBatanganSelamatIdulFitri:
		return "Emas Batangan Selamat Idul Fitri"
	case GoldCategoryEmasBatanganImlek:
		return "Emas Batangan Imlek"
	case GoldCategoryEmasBatanganBatikSeriIII:
		return "Emas Batangan Batik Seri III"
	case GoldCategoryPerakMurni:
		return "Perak Murni"
	case GoldCategoryPerakHeritage:
		return "Perak Heritage"
	case GoldCategoryLiontinBatikSeriIII:
		return "Liontin Batik Seri III"
	default:
		return string(gc)
	}
}

// GoldPricingHistory represents a record of scraped gold prices
type GoldPricingHistory struct {
	ID          int          `json:"id" db:"id"`
	PricingDate time.Time    `json:"pricing_date" db:"pricing_date"`
	GoldType    string       `json:"gold_type" db:"gold_type"`
	BuyPrice    int64        `json:"buy_price" db:"buy_price"`
	SellPrice   int64        `json:"sell_price" db:"sell_price"`
	Source      GoldSource   `json:"source" db:"source"`
	Category    GoldCategory `json:"category" db:"category"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// GoldPricingHistoryCreate represents the data needed to create a new gold pricing history record
type GoldPricingHistoryCreate struct {
	PricingDate time.Time    `json:"pricing_date" validate:"required"`
	GoldType    string       `json:"gold_type" validate:"required"`
	SellPrice   int64        `json:"sell_price" validate:"required"`
	Source      GoldSource   `json:"source" validate:"required,oneof=antam usb"`
	Category    GoldCategory `json:"category"`
}

// CalculateBuyPrice calculates the buy price as 94% of sell price (6% discount)
func (g *GoldPricingHistoryCreate) CalculateBuyPrice() int64 {
	// Calculate 94% (6% discount)
	return int64(float64(g.SellPrice) * 0.94)
}

// GoldPricingHistoryFilter represents filters for querying gold pricing histories
type GoldPricingHistoryFilter struct {
	GoldType  string       `query:"type"`
	Source    GoldSource   `query:"source"`
	Category  GoldCategory `query:"category"`
	StartDate *time.Time   `query:"start_date"`
	EndDate   *time.Time   `query:"end_date"`
	Limit     int          `query:"limit"`
	Offset    int          `query:"offset"`
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
	Errors       []string             `json:"errors"`
	Duration     string               `json:"duration"`
	Data         []GoldPricingHistory `json:"data,omitempty"`
}

// GoldPricingStats represents statistics about gold pricing data
type GoldPricingStats struct {
	TotalRecords    int       `json:"total_records"`
	UniqueGoldTypes int       `json:"unique_gold_types"`
	UniqueSources   int       `json:"unique_sources"`
	LatestDate      time.Time `json:"latest_date"`
	OldestDate      time.Time `json:"oldest_date"`
}
