package repositories

import (
	"database/sql"
	"errors"

	"nabung-emas-api/internal/models"
)

type TypePocketRepository struct {
	db *sql.DB
}

func NewTypePocketRepository(db *sql.DB) *TypePocketRepository {
	return &TypePocketRepository{db: db}
}

func (r *TypePocketRepository) FindAll() ([]models.TypePocket, error) {
	query := `
		SELECT id, name, description, icon, color, created_at, updated_at
		FROM type_pockets
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var typePockets []models.TypePocket
	for rows.Next() {
		var tp models.TypePocket
		err := rows.Scan(
			&tp.ID,
			&tp.Name,
			&tp.Description,
			&tp.Icon,
			&tp.Color,
			&tp.CreatedAt,
			&tp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		typePockets = append(typePockets, tp)
	}

	return typePockets, rows.Err()
}

func (r *TypePocketRepository) FindByID(id string) (*models.TypePocket, error) {
	query := `
		SELECT id, name, description, icon, color, created_at, updated_at
		FROM type_pockets
		WHERE id = $1
	`

	tp := &models.TypePocket{}
	err := r.db.QueryRow(query, id).Scan(
		&tp.ID,
		&tp.Name,
		&tp.Description,
		&tp.Icon,
		&tp.Color,
		&tp.CreatedAt,
		&tp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("type pocket not found")
	}

	return tp, err
}

func (r *TypePocketRepository) Exists(id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM type_pockets WHERE id = $1)`
	
	var exists bool
	err := r.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}
