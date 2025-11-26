package services

import (
	"errors"
	"time"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
	pocketRepo      *repositories.PocketRepository
}

func NewTransactionService(transactionRepo *repositories.TransactionRepository, pocketRepo *repositories.PocketRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		pocketRepo:      pocketRepo,
	}
}

func (s *TransactionService) GetAll(userID string, pocketID, brand, startDate, endDate *string, page, limit int, sortBy, sortOrder string) ([]models.Transaction, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.transactionRepo.FindAll(userID, pocketID, brand, startDate, endDate, page, limit, sortBy, sortOrder)
}

func (s *TransactionService) GetByID(id, userID string) (*models.Transaction, error) {
	return s.transactionRepo.FindByID(id, userID)
}

func (s *TransactionService) Create(userID string, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	// Validate pocket belongs to user
	_, err := s.pocketRepo.FindByID(req.PocketID, userID)
	if err != nil {
		return nil, errors.New("pocket not found")
	}

	// Validate total price matches weight * price_per_gram
	expectedTotal := req.Weight * req.PricePerGram
	if req.TotalPrice != expectedTotal {
		return nil, errors.New("total price must equal weight * price_per_gram")
	}

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return nil, errors.New("invalid transaction date format, use YYYY-MM-DD")
	}

	// Validate date is not in future
	if transactionDate.After(time.Now()) {
		return nil, errors.New("transaction date cannot be in the future")
	}

	transaction := &models.Transaction{
		UserID:          userID,
		PocketID:        req.PocketID,
		TransactionDate: transactionDate,
		Brand:           req.Brand,
		Weight:          req.Weight,
		PricePerGram:    req.PricePerGram,
		TotalPrice:      req.TotalPrice,
		Description:     req.Description,
		ReceiptImage:    req.ReceiptImage,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) Update(id, userID string, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.TransactionDate != "" {
		transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
		if err != nil {
			return nil, errors.New("invalid transaction date format, use YYYY-MM-DD")
		}
		if transactionDate.After(time.Now()) {
			return nil, errors.New("transaction date cannot be in the future")
		}
		transaction.TransactionDate = transactionDate
	}
	if req.Brand != "" {
		transaction.Brand = req.Brand
	}
	if req.Weight > 0 {
		transaction.Weight = req.Weight
	}
	if req.PricePerGram > 0 {
		transaction.PricePerGram = req.PricePerGram
	}
	if req.TotalPrice > 0 {
		// Validate total price
		expectedTotal := transaction.Weight * transaction.PricePerGram
		if req.TotalPrice != expectedTotal {
			return nil, errors.New("total price must equal weight * price_per_gram")
		}
		transaction.TotalPrice = req.TotalPrice
	}
	if req.Description != nil {
		transaction.Description = req.Description
	}

	if err := s.transactionRepo.Update(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) Delete(id, userID string) error {
	return s.transactionRepo.Delete(id, userID)
}

func (s *TransactionService) UpdateReceipt(id, userID, receiptURL string) error {
	return s.transactionRepo.UpdateReceiptImage(id, userID, receiptURL)
}
