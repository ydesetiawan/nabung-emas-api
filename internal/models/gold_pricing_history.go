package models

import "time"

// GoldSource represents the source of gold pricing data
type GoldSource string

const (
	GoldSourceAntam GoldSource = "antam"
	GoldSourceUSB   GoldSource = "usb"
)

// GoldPricingHistory represents a record of scraped gold prices
type GoldPricingHistory struct {
	ID        int        `json:"id" db:"id"`
	GoldType  string     `json:"gold_type" db:"gold_type"`
	BuyPrice  string     `json:"buy_price" db:"buy_price"`
	SellPrice string     `json:"sell_price" db:"sell_price"`
	Unit      string     `json:"unit" db:"unit"`
	Source    GoldSource `json:"source" db:"source"`
	ScrapedAt time.Time  `json:"scraped_at" db:"scraped_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// GoldPricingHistoryCreate represents the data needed to create a new gold pricing history record
type GoldPricingHistoryCreate struct {
	GoldType  string     `json:"gold_type" validate:"required"`
	BuyPrice  string     `json:"buy_price" validate:"required"`
	SellPrice string     `json:"sell_price" validate:"required"`
	Unit      string     `json:"unit" validate:"required"`
	Source    GoldSource `json:"source" validate:"required,oneof=antam usb"`
}

// GoldPricingHistoryFilter represents filters for querying gold pricing histories
type GoldPricingHistoryFilter struct {
	GoldType string     `query:"type"`
	Source   GoldSource `query:"source"`
	Limit    int        `query:"limit"`
}
