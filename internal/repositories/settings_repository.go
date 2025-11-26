package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"nabung-emas-api/internal/models"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) FindByUserID(userID string) (*models.UserSettings, error) {
	query := `
		SELECT id, user_id, language, theme, currency, 
		       email_notifications, push_notifications, price_alerts,
		       created_at, updated_at
		FROM user_settings
		WHERE user_id = $1
	`

	settings := &models.UserSettings{}
	err := r.db.QueryRow(query, userID).Scan(
		&settings.ID,
		&settings.UserID,
		&settings.Language,
		&settings.Theme,
		&settings.Currency,
		&settings.EmailNotifications,
		&settings.PushNotifications,
		&settings.PriceAlerts,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create default settings if not found
		return r.CreateDefault(userID)
	}

	return settings, err
}

func (r *SettingsRepository) CreateDefault(userID string) (*models.UserSettings, error) {
	query := `
		INSERT INTO user_settings (id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, language, theme, currency,
		          email_notifications, push_notifications, price_alerts,
		          created_at, updated_at
	`

	settings := &models.UserSettings{}
	now := time.Now()
	id := uuid.New().String()

	err := r.db.QueryRow(query, id, userID, now, now).Scan(
		&settings.ID,
		&settings.UserID,
		&settings.Language,
		&settings.Theme,
		&settings.Currency,
		&settings.EmailNotifications,
		&settings.PushNotifications,
		&settings.PriceAlerts,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	return settings, err
}

func (r *SettingsRepository) Update(settings *models.UserSettings) error {
	query := `
		UPDATE user_settings
		SET language = $1, theme = $2, 
		    email_notifications = $3, push_notifications = $4, price_alerts = $5,
		    updated_at = $6
		WHERE user_id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		settings.Language,
		settings.Theme,
		settings.EmailNotifications,
		settings.PushNotifications,
		settings.PriceAlerts,
		time.Now(),
		settings.UserID,
	).Scan(&settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("settings not found")
	}

	return err
}
