package repositories

import (
	"database/sql"
	"fmt"
	"nabung-emas-api/internal/models"
	"time"
)

// GoldPricingHistoryRepository handles database operations for gold pricing histories
type GoldPricingHistoryRepository struct {
	db *sql.DB
}

// NewGoldPricingHistoryRepository creates a new instance of GoldPricingHistoryRepository
func NewGoldPricingHistoryRepository(db *sql.DB) *GoldPricingHistoryRepository {
	return &GoldPricingHistoryRepository{db: db}
}

// calculateBuyPrice calculates the buy price as 94% of sell price (6% discount)
func calculateBuyPrice(sellPrice int64) int64 {
	// Calculate 94% (6% discount)
	return int64(float64(sellPrice) * 0.94)
}

// Create inserts a new gold pricing history record with UPSERT logic
func (r *GoldPricingHistoryRepository) Create(data *models.GoldPricingHistoryCreate) (*models.GoldPricingHistory, error) {
	buyPrice := calculateBuyPrice(data.SellPrice)

	query := `
		INSERT INTO gold_pricing_histories (pricing_date, gold_type, buy_price, sell_price, source)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (pricing_date, gold_type, source) 
		DO UPDATE SET 
			buy_price = EXCLUDED.buy_price,
			sell_price = EXCLUDED.sell_price,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, pricing_date, gold_type, buy_price, sell_price, source, created_at, updated_at
	`

	var history models.GoldPricingHistory
	err := r.db.QueryRow(
		query,
		data.PricingDate,
		data.GoldType,
		buyPrice,
		data.SellPrice,
		data.Source,
	).Scan(
		&history.ID,
		&history.PricingDate,
		&history.GoldType,
		&history.BuyPrice,
		&history.SellPrice,
		&history.Source,
		&history.CreatedAt,
		&history.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gold pricing history: %w", err)
	}

	return &history, nil
}

// CreateBatch inserts multiple gold pricing history records with UPSERT logic
func (r *GoldPricingHistoryRepository) CreateBatch(data []models.GoldPricingHistoryCreate) (int, int, error) {
	if len(data) == 0 {
		return 0, 0, nil
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO gold_pricing_histories (pricing_date, gold_type, buy_price, sell_price, source)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (pricing_date, gold_type, source) 
		DO UPDATE SET 
			buy_price = EXCLUDED.buy_price,
			sell_price = EXCLUDED.sell_price,
			updated_at = CURRENT_TIMESTAMP
		RETURNING (xmax = 0) AS inserted
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	savedCount := 0
	updatedCount := 0

	for _, item := range data {
		buyPrice := calculateBuyPrice(item.SellPrice)

		var inserted bool
		err := stmt.QueryRow(
			item.PricingDate,
			item.GoldType,
			buyPrice,
			item.SellPrice,
			item.Source,
		).Scan(&inserted)

		if err != nil {
			return 0, 0, fmt.Errorf("failed to insert/update record: %w", err)
		}

		if inserted {
			savedCount++
		} else {
			updatedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return savedCount, updatedCount, nil
}

// GetAll retrieves gold pricing histories with optional filters
func (r *GoldPricingHistoryRepository) GetAll(filter models.GoldPricingHistoryFilter) ([]models.GoldPricingHistory, error) {
	query := `
		SELECT id, pricing_date, gold_type, buy_price, sell_price, source, created_at, updated_at
		FROM gold_pricing_histories
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// Add filters
	if filter.GoldType != "" {
		query += fmt.Sprintf(" AND gold_type ILIKE $%d", argCount)
		args = append(args, "%"+filter.GoldType+"%")
		argCount++
	}

	if filter.Source != "" {
		query += fmt.Sprintf(" AND source = $%d", argCount)
		args = append(args, filter.Source)
		argCount++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND pricing_date >= $%d", argCount)
		args = append(args, *filter.StartDate)
		argCount++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND pricing_date <= $%d", argCount)
		args = append(args, *filter.EndDate)
		argCount++
	}

	// Order by pricing_date descending (latest first)
	query += " ORDER BY pricing_date DESC, gold_type ASC"

	// Add offset
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
		argCount++
	}

	// Add limit
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query gold pricing histories: %w", err)
	}
	defer rows.Close()

	histories := []models.GoldPricingHistory{}

	for rows.Next() {
		var history models.GoldPricingHistory
		err := rows.Scan(
			&history.ID,
			&history.PricingDate,
			&history.GoldType,
			&history.BuyPrice,
			&history.SellPrice,
			&history.Source,
			&history.CreatedAt,
			&history.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// GetLatest retrieves the latest price for each gold type and source
func (r *GoldPricingHistoryRepository) GetLatest() ([]models.GoldPricingHistory, error) {
	query := `
		SELECT DISTINCT ON (gold_type, source) 
			id, pricing_date, gold_type, buy_price, sell_price, source, created_at, updated_at
		FROM gold_pricing_histories
		ORDER BY gold_type, source, pricing_date DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest gold prices: %w", err)
	}
	defer rows.Close()

	histories := []models.GoldPricingHistory{}

	for rows.Next() {
		var history models.GoldPricingHistory
		err := rows.Scan(
			&history.ID,
			&history.PricingDate,
			&history.GoldType,
			&history.BuyPrice,
			&history.SellPrice,
			&history.Source,
			&history.CreatedAt,
			&history.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// GetByID retrieves a gold pricing history by ID
func (r *GoldPricingHistoryRepository) GetByID(id int) (*models.GoldPricingHistory, error) {
	query := `
		SELECT id, pricing_date, gold_type, buy_price, sell_price, source, created_at, updated_at
		FROM gold_pricing_histories
		WHERE id = $1
	`

	var history models.GoldPricingHistory
	err := r.db.QueryRow(query, id).Scan(
		&history.ID,
		&history.PricingDate,
		&history.GoldType,
		&history.BuyPrice,
		&history.SellPrice,
		&history.Source,
		&history.CreatedAt,
		&history.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get gold pricing history: %w", err)
	}

	return &history, nil
}

// GetByDate retrieves all gold pricing histories for a specific date
func (r *GoldPricingHistoryRepository) GetByDate(date time.Time) ([]models.GoldPricingHistory, error) {
	query := `
		SELECT id, pricing_date, gold_type, buy_price, sell_price, source, created_at, updated_at
		FROM gold_pricing_histories
		WHERE pricing_date = $1
		ORDER BY gold_type ASC, source ASC
	`

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query gold pricing histories by date: %w", err)
	}
	defer rows.Close()

	histories := []models.GoldPricingHistory{}

	for rows.Next() {
		var history models.GoldPricingHistory
		err := rows.Scan(
			&history.ID,
			&history.PricingDate,
			&history.GoldType,
			&history.BuyPrice,
			&history.SellPrice,
			&history.Source,
			&history.CreatedAt,
			&history.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// DeleteOldRecords deletes records older than the specified number of days
func (r *GoldPricingHistoryRepository) DeleteOldRecords(days int) (int64, error) {
	query := `
		DELETE FROM gold_pricing_histories
		WHERE pricing_date < CURRENT_DATE - INTERVAL '%d days'
	`

	result, err := r.db.Exec(fmt.Sprintf(query, days))
	if err != nil {
		return 0, fmt.Errorf("failed to delete old records: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetStats returns statistics about the gold pricing histories
func (r *GoldPricingHistoryRepository) GetStats() (*models.GoldPricingStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_records,
			COUNT(DISTINCT gold_type) as unique_gold_types,
			COUNT(DISTINCT source) as unique_sources,
			MIN(pricing_date) as oldest_date,
			MAX(pricing_date) as latest_date
		FROM gold_pricing_histories
	`

	var stats models.GoldPricingStats
	var oldestDate, latestDate sql.NullTime

	err := r.db.QueryRow(query).Scan(
		&stats.TotalRecords,
		&stats.UniqueGoldTypes,
		&stats.UniqueSources,
		&oldestDate,
		&latestDate,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	if oldestDate.Valid {
		stats.OldestDate = oldestDate.Time
	}

	if latestDate.Valid {
		stats.LatestDate = latestDate.Time
	}

	return &stats, nil
}

// GetVendorList returns a list of all unique vendors
func (r *GoldPricingHistoryRepository) GetVendorList() ([]string, error) {
	query := `
		SELECT DISTINCT source
		FROM gold_pricing_histories
		ORDER BY source
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vendor list: %w", err)
	}
	defer rows.Close()

	vendors := []string{}

	for rows.Next() {
		var vendor string
		if err := rows.Scan(&vendor); err != nil {
			return nil, fmt.Errorf("failed to scan vendor: %w", err)
		}
		vendors = append(vendors, vendor)
	}

	return vendors, nil
}

// CheckDuplicates checks if a record already exists for the given date, type, and source
func (r *GoldPricingHistoryRepository) CheckDuplicates(date time.Time, goldType string, source models.GoldSource) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM gold_pricing_histories
		WHERE pricing_date = $1 AND gold_type = $2 AND source = $3
	`

	var count int
	err := r.db.QueryRow(query, date, goldType, source).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check duplicates: %w", err)
	}

	return count > 0, nil
}
