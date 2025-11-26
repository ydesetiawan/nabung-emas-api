package services

import (
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
)

type SettingsService struct {
	repo *repositories.SettingsRepository
}

func NewSettingsService(repo *repositories.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) Get(userID string) (*models.SettingsResponse, error) {
	settings, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &models.SettingsResponse{
		Language: settings.Language,
		Theme:    settings.Theme,
		Currency: settings.Currency,
		Notifications: models.NotificationSettingsResponse{
			Email:       settings.EmailNotifications,
			Push:        settings.PushNotifications,
			PriceAlerts: settings.PriceAlerts,
		},
	}, nil
}

func (s *SettingsService) Update(userID string, req *models.UpdateSettingsRequest) (*models.SettingsResponse, error) {
	settings, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Language != nil {
		settings.Language = *req.Language
	}
	if req.Theme != nil {
		settings.Theme = *req.Theme
	}
	if req.Notifications != nil {
		if req.Notifications.Email != nil {
			settings.EmailNotifications = *req.Notifications.Email
		}
		if req.Notifications.Push != nil {
			settings.PushNotifications = *req.Notifications.Push
		}
		if req.Notifications.PriceAlerts != nil {
			settings.PriceAlerts = *req.Notifications.PriceAlerts
		}
	}

	if err := s.repo.Update(settings); err != nil {
		return nil, err
	}

	return &models.SettingsResponse{
		Language: settings.Language,
		Theme:    settings.Theme,
		Currency: settings.Currency,
		Notifications: models.NotificationSettingsResponse{
			Email:       settings.EmailNotifications,
			Push:        settings.PushNotifications,
			PriceAlerts: settings.PriceAlerts,
		},
	}, nil
}
