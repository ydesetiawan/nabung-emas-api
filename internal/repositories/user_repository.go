package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"nabung-emas-api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, full_name, email, phone, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	user.ID = uuid.New().String()
	now := time.Now()

	err := r.db.QueryRow(
		query,
		user.ID,
		user.FullName,
		user.Email,
		user.Phone,
		user.Password,
		now,
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, phone, password_hash, avatar, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, phone, avatar, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET full_name = $1, phone = $2, avatar = $3, updated_at = $4
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		user.FullName,
		user.Phone,
		user.Avatar,
		time.Now(),
		user.ID,
	).Scan(&user.UpdatedAt)

	return err
}

func (r *UserRepository) UpdatePassword(userID, hashedPassword string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, hashedPassword, time.Now(), userID)
	return err
}

func (r *UserRepository) GetStats(userID string) (*models.UserStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT p.id) as total_pockets,
			COUNT(t.id) as total_transactions,
			COALESCE(SUM(t.weight), 0) as total_weight,
			COALESCE(SUM(t.total_price), 0) as total_value
		FROM users u
		LEFT JOIN pockets p ON p.user_id = u.id
		LEFT JOIN transactions t ON t.user_id = u.id
		WHERE u.id = $1
		GROUP BY u.id
	`

	stats := &models.UserStats{}
	err := r.db.QueryRow(query, userID).Scan(
		&stats.TotalPockets,
		&stats.TotalTransactions,
		&stats.TotalWeight,
		&stats.TotalValue,
	)

	if err == sql.ErrNoRows {
		return &models.UserStats{}, nil
	}

	return stats, err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}
