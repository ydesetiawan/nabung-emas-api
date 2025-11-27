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

// GoldPricingHistory represents a record of scraped gold prices
type GoldPricingHistory struct {
	ID          int        `json:"id" db:"id"`
	PricingDate time.Time  `json:"pricing_date" db:"pricing_date"`
	GoldType    string     `json:"gold_type" db:"gold_type"`
	BuyPrice    string     `json:"buy_price" db:"buy_price"`
	SellPrice   string     `json:"sell_price" db:"sell_price"`
	Source      GoldSource `json:"source" db:"source"`
	ScrapedAt   time.Time  `json:"scraped_at" db:"scraped_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// GoldPricingHistoryCreate represents the data needed to create a new gold pricing history record
type GoldPricingHistoryCreate struct {
	PricingDate time.Time  `json:"pricing_date" validate:"required"`
	GoldType    string     `json:"gold_type" validate:"required"`
	SellPrice   string     `json:"sell_price" validate:"required"`
	Source      GoldSource `json:"source" validate:"required,oneof=antam usb"`
}

// CalculateBuyPrice calculates the buy price as 94% of sell price (6% discount)
func (g *GoldPricingHistoryCreate) CalculateBuyPrice() string {
	// This will be calculated in the repository layer
	// For now, return empty string
	return ""
}

// GoldPricingHistoryFilter represents filters for querying gold pricing histories
type GoldPricingHistoryFilter struct {
	GoldType  string     `query:"type"`
	Source    GoldSource `query:"source"`
	StartDate *time.Time `query:"start_date"`
	EndDate   *time.Time `query:"end_date"`
	Limit     int        `query:"limit"`
	Offset    int        `query:"offset"`
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
