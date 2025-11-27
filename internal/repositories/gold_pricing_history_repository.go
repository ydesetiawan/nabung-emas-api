package repositories

import (
	"database/sql"
	"fmt"
	"nabung-emas-api/internal/models"
)

// GoldPricingHistoryRepository handles database operations for gold pricing histories
type GoldPricingHistoryRepository struct {
	db *sql.DB
}

// NewGoldPricingHistoryRepository creates a new instance of GoldPricingHistoryRepository
func NewGoldPricingHistoryRepository(db *sql.DB) *GoldPricingHistoryRepository {
	return &GoldPricingHistoryRepository{db: db}
}

// Create inserts a new gold pricing history record into the database
func (r *GoldPricingHistoryRepository) Create(data *models.GoldPricingHistoryCreate) (*models.GoldPricingHistory, error) {
	query := `
		INSERT INTO gold_pricing_histories (gold_type, buy_price, sell_price, unit, source)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, gold_type, buy_price, sell_price, unit, source, scraped_at, created_at
	`

	var history models.GoldPricingHistory
	err := r.db.QueryRow(
		query,
		data.GoldType,
		data.BuyPrice,
		data.SellPrice,
		data.Unit,
		data.Source,
	).Scan(
		&history.ID,
		&history.GoldType,
		&history.BuyPrice,
		&history.SellPrice,
		&history.Unit,
		&history.Source,
		&history.ScrapedAt,
		&history.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gold pricing history: %w", err)
	}

	return &history, nil
}

// CreateBatch inserts multiple gold pricing history records in a single transaction
func (r *GoldPricingHistoryRepository) CreateBatch(data []models.GoldPricingHistoryCreate) ([]models.GoldPricingHistory, error) {
	if len(data) == 0 {
		return []models.GoldPricingHistory{}, nil
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO gold_pricing_histories (gold_type, buy_price, sell_price, unit, source)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, gold_type, buy_price, sell_price, unit, source, scraped_at, created_at
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	histories := make([]models.GoldPricingHistory, 0, len(data))

	for _, item := range data {
		var history models.GoldPricingHistory
		err := stmt.QueryRow(
			item.GoldType,
			item.BuyPrice,
			item.SellPrice,
			item.Unit,
			item.Source,
		).Scan(
			&history.ID,
			&history.GoldType,
			&history.BuyPrice,
			&history.SellPrice,
			&history.Unit,
			&history.Source,
			&history.ScrapedAt,
			&history.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to insert record: %w", err)
		}

		histories = append(histories, history)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return histories, nil
}

// GetAll retrieves gold pricing histories with optional filters
func (r *GoldPricingHistoryRepository) GetAll(filter models.GoldPricingHistoryFilter) ([]models.GoldPricingHistory, error) {
	query := `
		SELECT id, gold_type, buy_price, sell_price, unit, source, scraped_at, created_at
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

	// Order by scraped_at descending (latest first)
	query += " ORDER BY scraped_at DESC"

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
			&history.GoldType,
			&history.BuyPrice,
			&history.SellPrice,
			&history.Unit,
			&history.Source,
			&history.ScrapedAt,
			&history.CreatedAt,
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

// GetLatest retrieves the latest price for each gold type
func (r *GoldPricingHistoryRepository) GetLatest() ([]models.GoldPricingHistory, error) {
	query := `
		SELECT DISTINCT ON (gold_type, source) 
			id, gold_type, buy_price, sell_price, unit, source, scraped_at, created_at
		FROM gold_pricing_histories
		ORDER BY gold_type, source, scraped_at DESC
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
			&history.GoldType,
			&history.BuyPrice,
			&history.SellPrice,
			&history.Unit,
			&history.Source,
			&history.ScrapedAt,
			&history.CreatedAt,
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
		SELECT id, gold_type, buy_price, sell_price, unit, source, scraped_at, created_at
		FROM gold_pricing_histories
		WHERE id = $1
	`

	var history models.GoldPricingHistory
	err := r.db.QueryRow(query, id).Scan(
		&history.ID,
		&history.GoldType,
		&history.BuyPrice,
		&history.SellPrice,
		&history.Unit,
		&history.Source,
		&history.ScrapedAt,
		&history.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get gold pricing history: %w", err)
	}

	return &history, nil
}

// DeleteOldRecords deletes records older than the specified number of days
func (r *GoldPricingHistoryRepository) DeleteOldRecords(days int) (int64, error) {
	query := `
		DELETE FROM gold_pricing_histories
		WHERE scraped_at < NOW() - INTERVAL '%d days'
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
func (r *GoldPricingHistoryRepository) GetStats() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_records,
			COUNT(DISTINCT gold_type) as unique_gold_types,
			COUNT(DISTINCT source) as unique_sources,
			MIN(scraped_at) as oldest_record,
			MAX(scraped_at) as latest_record
		FROM gold_pricing_histories
	`

	var totalRecords, uniqueGoldTypes, uniqueSources int
	var oldestRecord, latestRecord sql.NullTime

	err := r.db.QueryRow(query).Scan(
		&totalRecords,
		&uniqueGoldTypes,
		&uniqueSources,
		&oldestRecord,
		&latestRecord,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	stats := map[string]interface{}{
		"total_records":     totalRecords,
		"unique_gold_types": uniqueGoldTypes,
		"unique_sources":    uniqueSources,
		"oldest_record":     nil,
		"latest_record":     nil,
	}

	if oldestRecord.Valid {
		stats["oldest_record"] = oldestRecord.Time
	}

	if latestRecord.Valid {
		stats["latest_record"] = latestRecord.Time
	}

	return stats, nil
}
