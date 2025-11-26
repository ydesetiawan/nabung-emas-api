package models

import "time"

type UserSettings struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Language           string    `json:"language"`
	Theme              string    `json:"theme"`
	Currency           string    `json:"currency"`
	EmailNotifications bool      `json:"email_notifications"`
	PushNotifications  bool      `json:"push_notifications"`
	PriceAlerts        bool      `json:"price_alerts"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type UpdateSettingsRequest struct {
	Language      *string                `json:"language" validate:"omitempty,oneof=en id"`
	Theme         *string                `json:"theme" validate:"omitempty,oneof=light dark"`
	Notifications *NotificationSettings  `json:"notifications"`
}

type NotificationSettings struct {
	Email       *bool `json:"email"`
	Push        *bool `json:"push"`
	PriceAlerts *bool `json:"price_alerts"`
}

type SettingsResponse struct {
	Language      string                       `json:"language"`
	Theme         string                       `json:"theme"`
	Currency      string                       `json:"currency"`
	Notifications NotificationSettingsResponse `json:"notifications"`
}

type NotificationSettingsResponse struct {
	Email       bool `json:"email"`
	Push        bool `json:"push"`
	PriceAlerts bool `json:"price_alerts"`
}
