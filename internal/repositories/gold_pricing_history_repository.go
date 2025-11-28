package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nabung-emas-api/internal/models"
)

type GoldPricingHistoryRepository struct {
	db *sql.DB
}

func NewGoldPricingHistoryRepository(db *sql.DB) *GoldPricingHistoryRepository {
	return &GoldPricingHistoryRepository{db: db}
}

func (r *GoldPricingHistoryRepository) Create(price *models.GoldPricingHistory) error {
	query := `
		INSERT INTO gold_pricing_histories (
			pricing_date, gold_type, base_price, buy_price, sell_price, 
			include_tax, source, category, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	now := time.Now()
	price.CreatedAt = now
	price.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		price.PricingDate,
		price.GoldType,
		price.BasePrice,
		price.BuyPrice,
		price.SellPrice,
		price.IncludeTax,
		price.Source,
		price.Category,
		price.CreatedAt,
		price.UpdatedAt,
	).Scan(&price.ID)

	return err
}

func (r *GoldPricingHistoryRepository) BulkCreate(prices []models.GoldPricingHistory) error {
	if len(prices) == 0 {
		return nil
	}

	query := `
		INSERT INTO gold_pricing_histories (
			pricing_date, gold_type, base_price, buy_price, sell_price, 
			include_tax, source, category, created_at, updated_at
		) VALUES 
	`

	values := []interface{}{}
	placeholders := []string{}
	now := time.Now()

	for i, price := range prices {
		n := i * 10
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10))

		values = append(values,
			price.PricingDate,
			price.GoldType,
			price.BasePrice,
			price.BuyPrice,
			price.SellPrice,
			price.IncludeTax,
			price.Source,
			price.Category,
			now,
			now,
		)
	}

	query += strings.Join(placeholders, ",")

	// Optional: Add ON CONFLICT clause if you want to update existing records
	// query += ` ON CONFLICT (pricing_date, source, gold_type, category) DO UPDATE SET ...`

	_, err := r.db.Exec(query, values...)
	return err
}

func (r *GoldPricingHistoryRepository) GetByDateAndSource(date time.Time, source models.Source) ([]models.GoldPricingHistory, error) {
	query := `
		SELECT id, pricing_date, gold_type, base_price, buy_price, sell_price, 
			   include_tax, source, category, created_at, updated_at
		FROM gold_pricing_histories
		WHERE pricing_date = $1 AND source = $2
	`

	rows, err := r.db.Query(query, date, source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []models.GoldPricingHistory
	for rows.Next() {
		var p models.GoldPricingHistory
		err := rows.Scan(
			&p.ID,
			&p.PricingDate,
			&p.GoldType,
			&p.BasePrice,
			&p.BuyPrice,
			&p.SellPrice,
			&p.IncludeTax,
			&p.Source,
			&p.Category,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}

	return prices, nil
}
