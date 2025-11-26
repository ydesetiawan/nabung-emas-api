package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"nabung-emas-api/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, user_id, pocket_id, transaction_date, brand, weight, 
			price_per_gram, total_price, description, receipt_image, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	transaction.ID = uuid.New().String()
	now := time.Now()

	err := r.db.QueryRow(
		query,
		transaction.ID,
		transaction.UserID,
		transaction.PocketID,
		transaction.TransactionDate,
		transaction.Brand,
		transaction.Weight,
		transaction.PricePerGram,
		transaction.TotalPrice,
		transaction.Description,
		transaction.ReceiptImage,
		now,
		now,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)

	return err
}

func (r *TransactionRepository) FindAll(userID string, pocketID, brand, startDate, endDate *string, page, limit int, sortBy, sortOrder string) ([]models.Transaction, int, error) {
	// Build count query
	countQuery := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`
	args := []interface{}{userID}
	argCount := 1

	whereClause := ""
	if pocketID != nil && *pocketID != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND pocket_id = $%d", argCount)
		args = append(args, *pocketID)
	}
	if brand != nil && *brand != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND brand = $%d", argCount)
		args = append(args, *brand)
	}
	if startDate != nil && *startDate != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND transaction_date >= $%d", argCount)
		args = append(args, *startDate)
	}
	if endDate != nil && *endDate != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND transaction_date <= $%d", argCount)
		args = append(args, *endDate)
	}

	countQuery += whereClause

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build main query
	query := `
		SELECT 
			t.id, t.user_id, t.pocket_id, t.transaction_date, t.brand,
			t.weight, t.price_per_gram, t.total_price, t.description,
			t.receipt_image, t.created_at, t.updated_at,
			p.id, p.name,
			tp.name, tp.color
		FROM transactions t
		LEFT JOIN pockets p ON p.id = t.pocket_id
		LEFT JOIN type_pockets tp ON tp.id = p.type_pocket_id
		WHERE t.user_id = $1
	`

	query += whereClause

	// Add sorting
	if sortBy == "" {
		sortBy = "transaction_date"
	}
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY t.%s %s", sortBy, sortOrder)

	// Add pagination
	offset := (page - 1) * limit
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, limit)
	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var p models.Pocket
		var tp models.TypePocket

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.PocketID,
			&t.TransactionDate,
			&t.Brand,
			&t.Weight,
			&t.PricePerGram,
			&t.TotalPrice,
			&t.Description,
			&t.ReceiptImage,
			&t.CreatedAt,
			&t.UpdatedAt,
			&p.ID,
			&p.Name,
			&tp.Name,
			&tp.Color,
		)
		if err != nil {
			return nil, 0, err
		}

		p.TypePocket = &tp
		t.Pocket = &p
		transactions = append(transactions, t)
	}

	return transactions, total, rows.Err()
}

func (r *TransactionRepository) FindByID(id, userID string) (*models.Transaction, error) {
	query := `
		SELECT 
			t.id, t.user_id, t.pocket_id, t.transaction_date, t.brand,
			t.weight, t.price_per_gram, t.total_price, t.description,
			t.receipt_image, t.created_at, t.updated_at,
			p.id, p.name, p.type_pocket_id,
			tp.id, tp.name, tp.icon, tp.color
		FROM transactions t
		LEFT JOIN pockets p ON p.id = t.pocket_id
		LEFT JOIN type_pockets tp ON tp.id = p.type_pocket_id
		WHERE t.id = $1 AND t.user_id = $2
	`

	t := &models.Transaction{}
	p := &models.Pocket{}
	tp := &models.TypePocket{}

	err := r.db.QueryRow(query, id, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.PocketID,
		&t.TransactionDate,
		&t.Brand,
		&t.Weight,
		&t.PricePerGram,
		&t.TotalPrice,
		&t.Description,
		&t.ReceiptImage,
		&t.CreatedAt,
		&t.UpdatedAt,
		&p.ID,
		&p.Name,
		&p.TypePocketID,
		&tp.ID,
		&tp.Name,
		&tp.Icon,
		&tp.Color,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}

	if err != nil {
		return nil, err
	}

	p.TypePocket = tp
	t.Pocket = p

	return t, nil
}

func (r *TransactionRepository) Update(transaction *models.Transaction) error {
	query := `
		UPDATE transactions
		SET transaction_date = $1, brand = $2, weight = $3, price_per_gram = $4,
		    total_price = $5, description = $6, updated_at = $7
		WHERE id = $8 AND user_id = $9
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		transaction.TransactionDate,
		transaction.Brand,
		transaction.Weight,
		transaction.PricePerGram,
		transaction.TotalPrice,
		transaction.Description,
		time.Now(),
		transaction.ID,
		transaction.UserID,
	).Scan(&transaction.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("transaction not found")
	}

	return err
}

func (r *TransactionRepository) Delete(id, userID string) error {
	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

func (r *TransactionRepository) UpdateReceiptImage(id, userID, receiptURL string) error {
	query := `
		UPDATE transactions
		SET receipt_image = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4
	`

	result, err := r.db.Exec(query, receiptURL, time.Now(), id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

func (r *TransactionRepository) GetRecentTransactions(userID string, limit int) ([]models.Transaction, error) {
	query := `
		SELECT 
			t.id, t.pocket_id, t.transaction_date, t.brand, t.weight, t.total_price,
			p.id, p.name,
			tp.color
		FROM transactions t
		LEFT JOIN pockets p ON p.id = t.pocket_id
		LEFT JOIN type_pockets tp ON tp.id = p.type_pocket_id
		WHERE t.user_id = $1
		ORDER BY t.transaction_date DESC, t.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var p models.Pocket
		var tp models.TypePocket

		err := rows.Scan(
			&t.ID,
			&t.PocketID,
			&t.TransactionDate,
			&t.Brand,
			&t.Weight,
			&t.TotalPrice,
			&p.ID,
			&p.Name,
			&tp.Color,
		)
		if err != nil {
			return nil, err
		}

		p.TypePocket = &tp
		t.Pocket = &p
		transactions = append(transactions, t)
	}

	return transactions, rows.Err()
}
