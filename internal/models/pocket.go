package models

import "time"

type Pocket struct {
	ID                   string      `json:"id"`
	UserID               string      `json:"user_id"`
	TypePocketID         string      `json:"type_pocket_id"`
	Name                 string      `json:"name"`
	Description          *string     `json:"description"`
	AggregateTotalPrice  float64     `json:"aggregate_total_price"`
	AggregateTotalWeight float64     `json:"aggregate_total_weight"`
	TargetWeight         *float64    `json:"target_weight"`
	TypePocket           *TypePocket `json:"type_pocket,omitempty"`
	TransactionCount     *int        `json:"transaction_count,omitempty"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

type CreatePocketRequest struct {
	TypePocketID string   `json:"type_pocket_id" validate:"required,uuid"`
	Name         string   `json:"name" validate:"required,min=3,max=100"`
	Description  *string  `json:"description" validate:"omitempty,max=500"`
	TargetWeight *float64 `json:"target_weight" validate:"omitempty,gt=0"`
}

type UpdatePocketRequest struct {
	Name         string   `json:"name" validate:"omitempty,min=3,max=100"`
	Description  *string  `json:"description" validate:"omitempty,max=500"`
	TargetWeight *float64 `json:"target_weight" validate:"omitempty,gt=0"`
}

type PocketStats struct {
	TotalWeight          float64  `json:"total_weight"`
	TotalValue           float64  `json:"total_value"`
	AveragePricePerGram  float64  `json:"average_price_per_gram"`
	CurrentGoldPrice     *float64 `json:"current_gold_price,omitempty"`
	CurrentValue         *float64 `json:"current_value,omitempty"`
	ProfitLoss           *float64 `json:"profit_loss,omitempty"`
	ProfitLossPercentage *float64 `json:"profit_loss_percentage,omitempty"`
	TransactionCount     int      `json:"transaction_count"`
}
