package services

import (
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
)

type TypePocketService struct {
	repo *repositories.TypePocketRepository
}

func NewTypePocketService(repo *repositories.TypePocketRepository) *TypePocketService {
	return &TypePocketService{repo: repo}
}

func (s *TypePocketService) GetAll() ([]models.TypePocket, error) {
	return s.repo.FindAll()
}

func (s *TypePocketService) GetByID(id string) (*models.TypePocket, error) {
	return s.repo.FindByID(id)
}
