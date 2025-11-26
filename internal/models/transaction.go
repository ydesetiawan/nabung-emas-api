package models

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	PocketID        string    `json:"pocket_id"`
	TransactionDate time.Time `json:"transaction_date"`
	Brand           string    `json:"brand"`
	Weight          float64   `json:"weight"`
	PricePerGram    float64   `json:"price_per_gram"`
	TotalPrice      float64   `json:"total_price"`
	Description     *string   `json:"description"`
	ReceiptImage    *string   `json:"receipt_image"`
	Pocket          *Pocket   `json:"pocket,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateTransactionRequest struct {
	PocketID        string  `json:"pocket_id" validate:"required,uuid"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
	Brand           string  `json:"brand" validate:"required,oneof=Antam UBS Pegadaian 'King Halim' Custom"`
	Weight          float64 `json:"weight" validate:"required,gte=0.1,lte=1000"`
	PricePerGram    float64 `json:"price_per_gram" validate:"required,gte=1000,lte=10000000"`
	TotalPrice      float64 `json:"total_price" validate:"required"`
	Description     *string `json:"description" validate:"omitempty,max=500"`
	ReceiptImage    *string `json:"receipt_image" validate:"omitempty"`
}

type UpdateTransactionRequest struct {
	TransactionDate string  `json:"transaction_date" validate:"omitempty"`
	Brand           string  `json:"brand" validate:"omitempty,oneof=Antam UBS Pegadaian 'King Halim' Custom"`
	Weight          float64 `json:"weight" validate:"omitempty,gte=0.1,lte=1000"`
	PricePerGram    float64 `json:"price_per_gram" validate:"omitempty,gte=1000,lte=10000000"`
	TotalPrice      float64 `json:"total_price" validate:"omitempty"`
	Description     *string `json:"description" validate:"omitempty,max=500"`
}
