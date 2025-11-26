package models

import "time"

type GoldPrice struct {
	ID           string    `json:"id"`
	Date         time.Time `json:"date"`
	PricePerGram float64   `json:"price_per_gram"`
	Source       string    `json:"source"`
	CreatedAt    time.Time `json:"created_at"`
}

type CurrentGoldPrice struct {
	PricePerGram        float64   `json:"price_per_gram"`
	Currency            string    `json:"currency"`
	Source              string    `json:"source"`
	LastUpdated         time.Time `json:"last_updated"`
	Change24h           *float64  `json:"change_24h,omitempty"`
	ChangePercentage24h *float64  `json:"change_percentage_24h,omitempty"`
}
