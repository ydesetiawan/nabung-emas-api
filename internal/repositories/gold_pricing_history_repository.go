package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"nabung-emas-api/internal/models"
	"strings"
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

// Create inserts a single gold pricing record
// Uses UPSERT logic - if record exists (same date, type, source), it will UPDATE instead of creating duplicate
func (r *GoldPricingHistoryRepository) Create(data *models.GoldPricingHistoryCreate) (*models.GoldPricingHistory, error) {
	query := `
		INSERT INTO gold_pricing_histories (pricing_date, gold_type, buy_price, sell_price, source, scraped_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (pricing_date, gold_type, source) 
		DO UPDATE SET 
			buy_price = EXCLUDED.buy_price,
			sell_price = EXCLUDED.sell_price,
			scraped_at = EXCLUDED.scraped_at,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, pricing_date, gold_type, buy_price, sell_price, source, scraped_at, created_at, updated_at
	`

	var result models.GoldPricingHistory
	err := r.db.QueryRow(
		query,
		data.PricingDate,
		data.GoldType,
		data.BuyPrice,
		data.SellPrice,
		data.Source,
		time.Now(),
	).Scan(
		&result.ID,
		&result.PricingDate,
		&result.GoldType,
		&result.BuyPrice,
		&result.SellPrice,
		&result.Source,
		&result.ScrapedAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create/update gold pricing history: %w", err)
	}

	return &result, nil
}

// CreateBatch inserts multiple gold pricing records in a single transaction
// Uses UPSERT logic - replaces existing records with same date/type/source
func (r *GoldPricingHistoryRepository) CreateBatch(data []models.GoldPricingHistoryCreate) (savedCount int, updatedCount int, err error) {
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
		INSERT INTO gold_pricing_histories (pricing_date, gold_type, buy_price, sell_price, source, scraped_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (pricing_date, gold_type, source) 
		DO UPDATE SET 
			buy_price = EXCLUDED.buy_price,
			sell_price = EXCLUDED.sell_price,
			scraped_at = EXCLUDED.scraped_at,
			updated_at = CURRENT_TIMESTAMP
		RETURNING (xmax = 0) AS inserted
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	scrapedAt := time.Now()

	for _, item := range data {
		var wasInserted bool
		err := stmt.QueryRow(
			item.PricingDate,
			item.GoldType,
			item.BuyPrice,
			item.SellPrice,
			item.Source,
			scrapedAt,
		).Scan(&wasInserted)

		if err != nil {
			log.Printf("⚠️  Failed to insert/update record for %s %s %s: %v", item.PricingDate.Format("2006-01-02"), item.GoldType, item.Source, err)
			continue
		}

		if wasInserted {
			savedCount++
		} else {
			updatedCount++
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return savedCount, updatedCount, nil
}

// GetAll retrieves gold pricing histories with optional filters
func (r *GoldPricingHistoryRepository) GetAll(filter models.GoldPricingHistoryFilter) ([]models.GoldPricingHistory, error) {
	query := `
		SELECT id, pricing_date, gold_type, buy_price, sell_price, source, scraped_at, created_at, updated_at
		FROM gold_pricing_histories
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	// Apply filters
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

	// Order by pricing_date DESC, then by gold_type
	query += " ORDER BY pricing_date DESC, gold_type ASC"

	// Apply pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query gold pricing histories: %w", err)
	}
	defer rows.Close()

	var results []models.GoldPricingHistory
	for rows.Next() {
		var item models.GoldPricingHistory
		err := rows.Scan(
			&item.ID,
			&item.PricingDate,
			&item.GoldType,
			&item.BuyPrice,
			&item.SellPrice,
			&item.Source,
			&item.ScrapedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			log.Printf("⚠️  Failed to scan row: %v", err)
			continue
		}
		results = append(results, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// GetLatest retrieves the latest price for each gold type and source combination
func (r *GoldPricingHistoryRepository) GetLatest() ([]models.GoldPricingHistory, error) {
	query := `
		SELECT DISTINCT ON (gold_type, source) 
			id, pricing_date, gold_type, buy_price, sell_price, source, scraped_at, created_at, updated_at
		FROM gold_pricing_histories
		ORDER BY gold_type, source, pricing_date DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest gold prices: %w", err)
	}
	defer rows.Close()

	var results []models.GoldPricingHistory
	for rows.Next() {
		var item models.GoldPricingHistory
		err := rows.Scan(
			&item.ID,
			&item.PricingDate,
			&item.GoldType,
			&item.BuyPrice,
			&item.SellPrice,
			&item.Source,
			&item.ScrapedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			log.Printf("⚠️  Failed to scan row: %v", err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

// GetByID retrieves a gold pricing history by ID
func (r *GoldPricingHistoryRepository) GetByID(id int) (*models.GoldPricingHistory, error) {
	query := `
		SELECT id, pricing_date, gold_type, buy_price, sell_price, source, scraped_at, created_at, updated_at
		FROM gold_pricing_histories
		WHERE id = $1
	`

	var result models.GoldPricingHistory
	err := r.db.QueryRow(query, id).Scan(
		&result.ID,
		&result.PricingDate,
		&result.GoldType,
		&result.BuyPrice,
		&result.SellPrice,
		&result.Source,
		&result.ScrapedAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get gold pricing history: %w", err)
	}

	return &result, nil
}

// GetByDate retrieves all gold prices for a specific date
func (r *GoldPricingHistoryRepository) GetByDate(date time.Time) ([]models.GoldPricingHistory, error) {
	query := `
		SELECT id, pricing_date, gold_type, buy_price, sell_price, source, scraped_at, created_at, updated_at
		FROM gold_pricing_histories
		WHERE pricing_date = $1
		ORDER BY source, gold_type
	`

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query gold prices by date: %w", err)
	}
	defer rows.Close()

	var results []models.GoldPricingHistory
	for rows.Next() {
		var item models.GoldPricingHistory
		err := rows.Scan(
			&item.ID,
			&item.PricingDate,
			&item.GoldType,
			&item.BuyPrice,
			&item.SellPrice,
			&item.Source,
			&item.ScrapedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			log.Printf("⚠️  Failed to scan row: %v", err)
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

// DeleteOldRecords deletes records older than the specified number of days
func (r *GoldPricingHistoryRepository) DeleteOldRecords(daysToKeep int) (int64, error) {
	query := `
		DELETE FROM gold_pricing_histories
		WHERE pricing_date < CURRENT_DATE - INTERVAL '%d days'
	`

	result, err := r.db.Exec(fmt.Sprintf(query, daysToKeep))
	if err != nil {
		return 0, fmt.Errorf("failed to delete old records: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetStats retrieves statistics about the gold pricing data
func (r *GoldPricingHistoryRepository) GetStats() (*models.GoldPricingStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_records,
			COUNT(DISTINCT source) as unique_vendors,
			COUNT(DISTINCT gold_type) as unique_gold_types,
			MAX(scraped_at) as latest_scraped_at,
			MIN(pricing_date) as oldest_pricing_date,
			MAX(pricing_date) as latest_pricing_date
		FROM gold_pricing_histories
	`

	var stats models.GoldPricingStats
	err := r.db.QueryRow(query).Scan(
		&stats.TotalRecords,
		&stats.UniqueVendors,
		&stats.UniqueGoldTypes,
		&stats.LatestScrapedAt,
		&stats.OldestPricingDate,
		&stats.LatestPricingDate,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &stats, nil
}

// GetVendorList retrieves list of all vendors in the database
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

	var vendors []string
	for rows.Next() {
		var vendor string
		if err := rows.Scan(&vendor); err != nil {
			continue
		}
		vendors = append(vendors, vendor)
	}

	return vendors, nil
}

// CheckDuplicates checks if records already exist for a given date
func (r *GoldPricingHistoryRepository) CheckDuplicates(pricingDate time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM gold_pricing_histories
		WHERE pricing_date = $1
	`

	var count int
	err := r.db.QueryRow(query, pricingDate).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check duplicates: %w", err)
	}

	return count, nil
}

// Helper function to build WHERE clause dynamically
func buildWhereClause(conditions []string) string {
	if len(conditions) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(conditions, " AND ")
}
