package services

import (
	"log"
	"nabung-emas-api/internal/repositories"
	"time"
)

type CleanupService struct {
	tokenBlacklistRepo *repositories.TokenBlacklistRepository
}

func NewCleanupService(tokenBlacklistRepo *repositories.TokenBlacklistRepository) *CleanupService {
	return &CleanupService{
		tokenBlacklistRepo: tokenBlacklistRepo,
	}
}

// StartTokenCleanup starts a background goroutine that periodically cleans up expired tokens
func (s *CleanupService) StartTokenCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := s.tokenBlacklistRepo.CleanupExpired(); err != nil {
				log.Printf("Error cleaning up expired tokens: %v", err)
			} else {
				log.Println("Successfully cleaned up expired tokens")
			}
		}
	}()
}
