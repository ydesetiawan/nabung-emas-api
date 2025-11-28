package models

import "time"

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

// GoldPricingHistory represents a gold price record
type GoldPricingHistory struct {
	ID          int          `json:"id,omitempty"`
	PricingDate time.Time    `json:"pricing_date"`
	GoldType    string       `json:"gold_type"`
	BasePrice   int          `json:"base_price"`
	BuyPrice    int          `json:"buy_price"`
	SellPrice   int          `json:"sell_price"`
	IncludeTax  bool         `json:"include_tax"`
	Source      Source       `json:"source"`
	Category    GoldCategory `json:"category"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
}

type CurrentGoldPrice struct {
	PricePerGram        float64   `json:"price_per_gram"`
	Currency            string    `json:"currency"`
	Source              string    `json:"source"`
	LastUpdated         time.Time `json:"last_updated"`
	Change24h           *float64  `json:"change_24h,omitempty"`
	ChangePercentage24h *float64  `json:"change_percentage_24h,omitempty"`
}
