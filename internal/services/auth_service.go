package services

import (
	"errors"
	"nabung-emas-api/internal/config"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
	"nabung-emas-api/internal/utils"
	"time"
)

type AuthService struct {
	userRepo           *repositories.UserRepository
	tokenBlacklistRepo *repositories.TokenBlacklistRepository
	config             *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, tokenBlacklistRepo *repositories.TokenBlacklistRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		tokenBlacklistRepo: tokenBlacklistRepo,
		config:             cfg,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, *models.TokenResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	// Create user
	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiry)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.RefreshTokenExpiry)
	if err != nil {
		return nil, nil, err
	}

	tokens := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWTExpiry.Seconds()),
	}

	// Clear password before returning
	user.Password = ""

	return user, tokens, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.User, *models.TokenResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	// Compare password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	// Generate tokens
	expiry := s.config.JWTExpiry
	if req.RememberMe {
		expiry = s.config.RefreshTokenExpiry
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, expiry)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.RefreshTokenExpiry)
	if err != nil {
		return nil, nil, err
	}

	tokens := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expiry.Seconds()),
	}

	// Clear password before returning
	user.Password = ""

	return user, tokens, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken, s.config.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Generate new tokens
	accessToken, err := utils.GenerateToken(claims.UserID, claims.Email, s.config.JWTSecret, s.config.JWTExpiry)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateToken(claims.UserID, claims.Email, s.config.JWTSecret, s.config.RefreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.JWTExpiry.Seconds()),
	}, nil
}

func (s *AuthService) GetCurrentUser(userID string) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) ForgotPassword(email string) error {
	// Check if user exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Return success even if user doesn't exist (security best practice)
		return nil
	}

	// TODO: Generate reset token and send email
	// For now, just return success
	_ = user
	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// TODO: Validate reset token and update password
	// For now, return not implemented error
	return errors.New("password reset not yet implemented")
}

func (s *AuthService) Logout(accessToken string, userID string) error {
	// Validate the token to get its expiration time
	claims, err := utils.ValidateToken(accessToken, s.config.JWTSecret)
	if err != nil {
		// Even if token is invalid/expired, we still want to blacklist it
		// Use a default expiration time
		expiresAt := time.Now().Add(s.config.JWTExpiry)
		return s.tokenBlacklistRepo.Add(accessToken, userID, expiresAt)
	}

	// Add token to blacklist with its actual expiration time
	return s.tokenBlacklistRepo.Add(accessToken, claims.UserID, claims.ExpiresAt.Time)
}
