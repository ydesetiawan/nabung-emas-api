# EmasGo Backend - Implementation Status & Guide

## ✅ What Has Been Implemented

### 1. Project Structure ✓
- Complete folder structure following Go best practices
- Separation of concerns (handlers, services, repositories)
- Configuration management
- Database migrations

### 2. Configuration & Setup ✓
- ✅ `.env.example` - Environment variables template
- ✅ `.gitignore` - Git ignore rules
- ✅ `go.mod` - Go module dependencies
- ✅ `README.md` - Comprehensive documentation
- ✅ `migrations/001_initial_schema.sql` - Complete database schema

### 3. Core Infrastructure ✓
- ✅ `internal/config/config.go` - Configuration loader
- ✅ `internal/database/postgres.go` - Database connection
- ✅ `cmd/server/main.go` - Application entry point
- ✅ `internal/routes/routes.go` - Complete route definitions

### 4. Models ✓
All data models are complete:
- ✅ User models with auth requests
- ✅ TypePocket models
- ✅ Pocket models with CRUD requests
- ✅ Transaction models with CRUD requests
- ✅ Analytics models (portfolio, trends, distributions)
- ✅ Settings models
- ✅ GoldPrice models
- ✅ Response models (API responses, pagination)

### 5. Utilities ✓
- ✅ `internal/utils/jwt.go` - JWT token generation and validation
- ✅ `internal/utils/password.go` - Password hashing with bcrypt
- ✅ `internal/utils/validator.go` - Request validation
- ✅ `internal/utils/response.go` - Standardized API responses

### 6. Middleware ✓
- ✅ `internal/middleware/auth.go` - JWT authentication
- ✅ `internal/middleware/cors.go` - CORS configuration
- ✅ `internal/middleware/logger.go` - Request logging

### 7. Repositories ✓
All repositories with raw SQL:
- ✅ `user_repository.go` - User CRUD, stats, email checking
- ✅ `type_pocket_repository.go` - Type pocket retrieval
- ✅ `pocket_repository.go` - Pocket CRUD with pagination & filtering
- ✅ `transaction_repository.go` - Transaction CRUD with advanced filtering
- ✅ `analytics_repository.go` - Portfolio analytics & distributions
- ✅ `settings_repository.go` - User settings management

## ⚠️ What Needs to Be Implemented

### 1. Services Layer (CRITICAL)
You need to create these service files:

#### `internal/services/auth_service.go`
```go
package services

import (
	"errors"
	"nabung-emas-api/internal/config"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
	"nabung-emas-api/internal/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, *models.TokenResponse, error) {
	// 1. Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, errors.New("email already registered")
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	// 3. Create user
	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// 4. Generate tokens
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
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	// 2. Compare password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	// 3. Generate tokens
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

// Implement other methods: ForgotPassword, ResetPassword, RefreshToken, etc.
```

#### `internal/services/user_service.go`
#### `internal/services/pocket_service.go`
#### `internal/services/transaction_service.go`
#### `internal/services/analytics_service.go`
#### `internal/services/settings_service.go`

### 2. Handlers Layer (CRITICAL)
You need to create these handler files:

#### `internal/handlers/auth_handler.go`
```go
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	user, tokens, err := h.authService.Register(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Account created successfully", map[string]interface{}{
		"user":          user,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	user, tokens, err := h.authService.Login(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Login successful", map[string]interface{}{
		"user":          user,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
	})
}

// Implement other methods: ForgotPassword, ResetPassword, RefreshToken, Logout, GetCurrentUser
```

#### Other handlers to create:
- `internal/handlers/user_handler.go`
- `internal/handlers/type_pocket_handler.go`
- `internal/handlers/pocket_handler.go`
- `internal/handlers/transaction_handler.go`
- `internal/handlers/analytics_handler.go`
- `internal/handlers/settings_handler.go`

### 3. Optional Features
- Gold price tracking service (external API integration)
- File upload service (local storage or AWS S3)
- Email service for password reset
- Rate limiting middleware

## 🚀 How to Complete the Implementation

### Step 1: Install Go
First, you need to install Go on your system:
```bash
# On macOS
brew install go

# Verify installation
go version
```

### Step 2: Install Dependencies
```bash
cd /Users/185772.edy/GitHub/nabung-emas-api
go mod download
```

### Step 3: Setup Database
```bash
# Create PostgreSQL database
createdb nabung_emas

# Run migrations
psql -d nabung_emas -f migrations/001_initial_schema.sql
```

### Step 4: Configure Environment
```bash
cp .env.example .env
# Edit .env with your actual database credentials and JWT secret
```

### Step 5: Implement Missing Services
Create all service files in `internal/services/` following the pattern shown above.

### Step 6: Implement Missing Handlers
Create all handler files in `internal/handlers/` following the pattern shown above.

### Step 7: Test the Application
```bash
# Run the server
go run cmd/server/main.go

# Test health endpoint
curl http://localhost:8080/health

# Test registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "email": "test@example.com",
    "phone": "+62 812 3456 7890",
    "password": "Test1234",
    "confirm_password": "Test1234"
  }'
```

## 📝 Implementation Priority

1. **HIGH PRIORITY** (Core functionality):
   - ✅ AuthService & AuthHandler (partially shown above)
   - ✅ UserService & UserHandler
   - ✅ PocketService & PocketHandler
   - ✅ TransactionService & TransactionHandler

2. **MEDIUM PRIORITY** (Important features):
   - ✅ AnalyticsService & AnalyticsHandler
   - ✅ SettingsService & SettingsHandler
   - ✅ TypePocketService & TypePocketHandler

3. **LOW PRIORITY** (Nice to have):
   - ⚠️ GoldPriceService (external API integration)
   - ⚠️ StorageService (file uploads)
   - ⚠️ EmailService (password reset emails)

## 🔧 Troubleshooting

### Common Issues:

1. **"go: command not found"**
   - Install Go using Homebrew: `brew install go`

2. **Database connection errors**
   - Ensure PostgreSQL is running
   - Check DATABASE_URL in .env file
   - Verify database exists: `psql -l`

3. **Import errors**
   - Run `go mod tidy` to clean up dependencies
   - Ensure all files use correct package names

4. **Validation errors**
   - Check that request structs have proper validation tags
   - Ensure validator is initialized in main.go

## 📚 Next Steps

1. Complete all service implementations
2. Complete all handler implementations
3. Add unit tests for services
4. Add integration tests for handlers
5. Implement file upload functionality
6. Add rate limiting for auth endpoints
7. Set up CI/CD pipeline
8. Deploy to production

## 🎯 Quick Reference

### Service Pattern:
```go
type XService struct {
    repo *repositories.XRepository
}

func NewXService(repo *repositories.XRepository) *XService {
    return &XService{repo: repo}
}

func (s *XService) MethodName(params) (result, error) {
    // Business logic here
    return s.repo.MethodName(params)
}
```

### Handler Pattern:
```go
type XHandler struct {
    service *services.XService
}

func NewXHandler(service *services.XService) *XHandler {
    return &XHandler{service: service}
}

func (h *XHandler) MethodName(c echo.Context) error {
    // 1. Bind and validate request
    // 2. Call service method
    // 3. Return response
}
```

---

**Note**: The foundation is solid. You now need to implement the service and handler layers following the patterns shown above. All the infrastructure, models, repositories, and utilities are ready to use!
