package services

import (
	"errors"
	"log"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
)

type PocketService struct {
	pocketRepo     *repositories.PocketRepository
	typePocketRepo *repositories.TypePocketRepository
}

func NewPocketService(pocketRepo *repositories.PocketRepository, typePocketRepo *repositories.TypePocketRepository) *PocketService {
	return &PocketService{
		pocketRepo:     pocketRepo,
		typePocketRepo: typePocketRepo,
	}
}

func (s *PocketService) GetAll(userID string, typePocketID *string, page, limit int, sortBy, sortOrder string) ([]models.Pocket, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.pocketRepo.FindAll(userID, typePocketID, page, limit, sortBy, sortOrder)
}

func (s *PocketService) GetByID(id, userID string) (*models.Pocket, error) {
	return s.pocketRepo.FindByID(id, userID)
}

func (s *PocketService) Create(userID string, req *models.CreatePocketRequest) (*models.Pocket, error) {
	// Validate type pocket exists
	exists, err := s.typePocketRepo.Exists(req.TypePocketID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("invalid type pocket ID")
	}

	// Check if pocket name already exists for user
	nameExists, err := s.pocketRepo.NameExistsForUser(userID, req.Name, "")
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, errors.New("pocket name already exists")
	}

	pocket := &models.Pocket{
		UserID:       userID,
		TypePocketID: req.TypePocketID,
		Name:         req.Name,
		Description:  req.Description,
		TargetWeight: req.TargetWeight,
	}

	if err := s.pocketRepo.Create(pocket); err != nil {
		return nil, err
	}

	return pocket, nil
}

func (s *PocketService) Update(id, userID string, req *models.UpdatePocketRequest) (*models.Pocket, error) {
	pocket, err := s.pocketRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if new name already exists
		nameExists, err := s.pocketRepo.NameExistsForUser(userID, req.Name, id)
		if err != nil {
			return nil, err
		}
		if nameExists {
			return nil, errors.New("pocket name already exists")
		}
		pocket.Name = req.Name
	}
	if req.Description != nil {
		pocket.Description = req.Description
	}
	if req.TargetWeight != nil {
		pocket.TargetWeight = req.TargetWeight
	}

	if req.TypePocketID != "" {
		log.Println("Type pocket ID updated:" + req.TypePocketID)
		// Validate type pocket exists
		exists, err := s.typePocketRepo.Exists(req.TypePocketID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("invalid type pocket ID")
		}
		pocket.TypePocketID = req.TypePocketID
	}
	log.Println("Type pocket all:" + pocket.TypePocketID)
	if err := s.pocketRepo.Update(pocket); err != nil {
		return nil, err
	}

	return pocket, nil
}

func (s *PocketService) Delete(id, userID string) error {
	return s.pocketRepo.Delete(id, userID)
}

func (s *PocketService) GetStats(id, userID string, currentGoldPrice *float64) (*models.PocketStats, error) {
	stats, err := s.pocketRepo.GetStats(id, userID)
	if err != nil {
		return nil, err
	}

	// Calculate profit/loss if current price provided
	if currentGoldPrice != nil && *currentGoldPrice > 0 && stats.TotalWeight > 0 {
		currentValue := stats.TotalWeight * (*currentGoldPrice)
		profitLoss := currentValue - stats.TotalValue
		profitLossPercentage := (profitLoss / stats.TotalValue) * 100

		stats.CurrentGoldPrice = currentGoldPrice
		stats.CurrentValue = &currentValue
		stats.ProfitLoss = &profitLoss
		stats.ProfitLossPercentage = &profitLossPercentage
	}

	return stats, nil
}
