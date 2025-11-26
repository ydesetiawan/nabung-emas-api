package services

import (
	"errors"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
	"nabung-emas-api/internal/utils"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID string) (*models.User, *models.UserStats, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}

	stats, err := s.userRepo.GetStats(userID)
	if err != nil {
		return nil, nil, err
	}

	return user, stats, nil
}

func (s *UserService) UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ChangePassword(userID string, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByEmail("")
	if err != nil {
		// Get user with password
		user, err = s.userRepo.FindByID(userID)
		if err != nil {
			return err
		}
	}

	// Verify current password
	if err := utils.ComparePassword(user.Password, req.CurrentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

func (s *UserService) UpdateAvatar(userID, avatarURL string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Avatar = &avatarURL
	return s.userRepo.Update(user)
}
