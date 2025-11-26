package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nabung-emas-api/internal/models"

	"github.com/google/uuid"
)

type PocketRepository struct {
	db *sql.DB
}

func NewPocketRepository(db *sql.DB) *PocketRepository {
	return &PocketRepository{db: db}
}

func (r *PocketRepository) Create(pocket *models.Pocket) error {
	query := `
		INSERT INTO pockets (id, user_id, type_pocket_id, name, description, target_weight, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, aggregate_total_price, aggregate_total_weight, created_at, updated_at
	`

	pocket.ID = uuid.New().String()
	now := time.Now()

	err := r.db.QueryRow(
		query,
		pocket.ID,
		pocket.UserID,
		pocket.TypePocketID,
		pocket.Name,
		pocket.Description,
		pocket.TargetWeight,
		now,
		now,
	).Scan(
		&pocket.ID,
		&pocket.AggregateTotalPrice,
		&pocket.AggregateTotalWeight,
		&pocket.CreatedAt,
		&pocket.UpdatedAt,
	)

	return err
}

func (r *PocketRepository) GetTopPockets(userID string) ([]models.Pocket, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.type_pocket_id, p.name, p.description,
			p.aggregate_total_price, p.aggregate_total_weight, p.target_weight,
			p.created_at, p.updated_at,
			tp.id, tp.name, tp.icon, tp.color
		FROM pockets p
		LEFT JOIN type_pockets tp ON tp.id = p.type_pocket_id
		WHERE p.user_id = $1
		ORDER BY p.aggregate_total_weight DESC
		LIMIT 3
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pockets []models.Pocket
	for rows.Next() {
		var p models.Pocket
		var tp models.TypePocket

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.TypePocketID,
			&p.Name,
			&p.Description,
			&p.AggregateTotalPrice,
			&p.AggregateTotalWeight,
			&p.TargetWeight,
			&p.CreatedAt,
			&p.UpdatedAt,
			&tp.ID,
			&tp.Name,
			&tp.Icon,
			&tp.Color,
		)
		if err != nil {
			return nil, err
		}

		p.TypePocket = &tp
		pockets = append(pockets, p)
	}

	return pockets, nil
}

func (r *PocketRepository) FindAll(userID string, typePocketID *string, page, limit int, sortBy, sortOrder string) ([]models.Pocket, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM pockets WHERE user_id = $1`
	args := []interface{}{userID}
	argCount := 1

	if typePocketID != nil && *typePocketID != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND type_pocket_id = $%d", argCount)
		args = append(args, *typePocketID)
	}

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get pockets with type_pocket info
	query := `
		SELECT 
			p.id, p.user_id, p.type_pocket_id, p.name, p.description,
			p.aggregate_total_price, p.aggregate_total_weight, p.target_weight,
			p.created_at, p.updated_at,
			tp.id, tp.name, tp.icon, tp.color
		FROM pockets p
		LEFT JOIN type_pockets tp ON tp.id = p.type_pocket_id
		WHERE p.user_id = $1
	`

	queryArgs := []interface{}{userID}
	argIdx := 1

	if typePocketID != nil && *typePocketID != "" {
		argIdx++
		query += fmt.Sprintf(" AND p.type_pocket_id = $%d", argIdx)
		queryArgs = append(queryArgs, *typePocketID)
	}

	// Add sorting
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY p.%s %s", sortBy, sortOrder)

	// Add pagination
	offset := (page - 1) * limit
	argIdx++
	query += fmt.Sprintf(" LIMIT $%d", argIdx)
	queryArgs = append(queryArgs, limit)
	argIdx++
	query += fmt.Sprintf(" OFFSET $%d", argIdx)
	queryArgs = append(queryArgs, offset)

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pockets []models.Pocket
	for rows.Next() {
		var p models.Pocket
		var tp models.TypePocket

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.TypePocketID,
			&p.Name,
			&p.Description,
			&p.AggregateTotalPrice,
			&p.AggregateTotalWeight,
			&p.TargetWeight,
			&p.CreatedAt,
			&p.UpdatedAt,
			&tp.ID,
			&tp.Name,
			&tp.Icon,
			&tp.Color,
		)
		if err != nil {
			return nil, 0, err
		}

		p.TypePocket = &tp
		pockets = append(pockets, p)
	}

	return pockets, total, rows.Err()
}

func (r *PocketRepository) FindByID(id, userID string) (*models.Pocket, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.type_pocket_id, p.name, p.description,
			p.aggregate_total_price, p.aggregate_total_weight, p.target_weight,
			p.created_at, p.updated_at,
			tp.id, tp.name, tp.description, tp.icon, tp.color,
			(SELECT COUNT(*) FROM transactions WHERE pocket_id = p.id) as transaction_count
		FROM pockets p
		LEFT JOIN type_pockets tp ON tp.id = p.type_pocket_id
		WHERE p.id = $1 AND p.user_id = $2
	`

	p := &models.Pocket{}
	tp := &models.TypePocket{}
	var transactionCount int

	err := r.db.QueryRow(query, id, userID).Scan(
		&p.ID,
		&p.UserID,
		&p.TypePocketID,
		&p.Name,
		&p.Description,
		&p.AggregateTotalPrice,
		&p.AggregateTotalWeight,
		&p.TargetWeight,
		&p.CreatedAt,
		&p.UpdatedAt,
		&tp.ID,
		&tp.Name,
		&tp.Description,
		&tp.Icon,
		&tp.Color,
		&transactionCount,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("pocket not found")
	}

	if err != nil {
		return nil, err
	}

	p.TypePocket = tp
	p.TransactionCount = &transactionCount

	return p, nil
}

func (r *PocketRepository) Update(pocket *models.Pocket) error {
	query := `
		UPDATE pockets
		SET name = $1, description = $2, target_weight = $3, updated_at = $4, type_pocket_id = $5
		WHERE id = $6 AND user_id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		pocket.Name,
		pocket.Description,
		pocket.TargetWeight,
		time.Now(),
		pocket.TypePocketID,
		pocket.ID,
		pocket.UserID,
	).Scan(&pocket.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("pocket not found")
	}

	return err
}

func (r *PocketRepository) Delete(id, userID string) error {
	query := `DELETE FROM pockets WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("pocket not found")
	}

	return nil
}

func (r *PocketRepository) GetStats(pocketID, userID string) (*models.PocketStats, error) {
	query := `
		SELECT 
			p.aggregate_total_weight,
			p.aggregate_total_price,
			CASE 
				WHEN p.aggregate_total_weight > 0 
				THEN p.aggregate_total_price / p.aggregate_total_weight 
				ELSE 0 
			END as average_price_per_gram,
			(SELECT COUNT(*) FROM transactions WHERE pocket_id = p.id) as transaction_count
		FROM pockets p
		WHERE p.id = $1 AND p.user_id = $2
	`

	stats := &models.PocketStats{}
	err := r.db.QueryRow(query, pocketID, userID).Scan(
		&stats.TotalWeight,
		&stats.TotalValue,
		&stats.AveragePricePerGram,
		&stats.TransactionCount,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("pocket not found")
	}

	return stats, err
}

func (r *PocketRepository) NameExistsForUser(userID, name, excludeID string) (bool, error) {
	var query string
	var args []interface{}

	if excludeID != "" {
		query = `SELECT EXISTS(SELECT 1 FROM pockets WHERE user_id = $1 AND name = $2 AND id != $3)`
		args = []interface{}{userID, name, excludeID}
	} else {
		query = `SELECT EXISTS(SELECT 1 FROM pockets WHERE user_id = $1 AND name = $2)`
		args = []interface{}{userID, name}
	}

	var exists bool
	err := r.db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}
