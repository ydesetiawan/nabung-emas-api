# EmasGo - Gold Savings API

A RESTful API backend for gold savings tracking application built with Go and PostgreSQL.

## Features

- 🔐 JWT-based authentication
- 👤 User profile management
- 💰 Multiple gold savings pockets
- 📊 Transaction tracking with receipt uploads
- 📈 Portfolio analytics and insights
- 💵 Gold price tracking
- ⚙️ User settings management

## Tech Stack

- **Framework:** Echo v4
- **Database:** PostgreSQL 14+
- **Authentication:** JWT
- **Validation:** go-playground/validator
- **Password Hashing:** bcrypt

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- (Optional) AWS S3 for file storage

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/nabung-emas-api.git
cd nabung-emas-api
```

2. Install dependencies:
```bash
go mod download
```

3. Copy environment file and configure:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Create PostgreSQL database:
```bash
createdb nabung_emas
```

5. Run database migrations:
```bash
psql -d nabung_emas -f migrations/001_initial_schema.sql
```

6. Run the application:
```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

## API Documentation

See [golang-api-specification.md](./golang-api-specification.md) for complete API documentation.

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Most endpoints require JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Project Structure

```
nabung-emas-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── database/                # Database connection
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # Custom middleware
│   ├── models/                  # Data models
│   ├── repositories/            # Database operations
│   ├── services/                # Business logic
│   ├── utils/                   # Utility functions
│   └── routes/                  # Route definitions
├── migrations/                  # Database migrations
├── .env.example                 # Environment variables template
├── go.mod                       # Go dependencies
└── README.md                    # This file
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| DATABASE_URL | PostgreSQL connection string | - |
| JWT_SECRET | Secret key for JWT signing | - |
| JWT_EXPIRY | Access token expiry | 24h |
| REFRESH_TOKEN_EXPIRY | Refresh token expiry | 168h |
| STORAGE_TYPE | Storage type (local/s3) | local |
| ALLOWED_ORIGINS | CORS allowed origins | - |

## Development

### Running Tests
```bash
go test ./...
```

### Running with Hot Reload
Install air:
```bash
go install github.com/cosmtrek/air@latest
```

Run with air:
```bash
air
```

### Database Migrations

To create a new migration:
```bash
# Create migration file in migrations/ directory
touch migrations/002_your_migration_name.sql
```

To run migrations:
```bash
psql -d nabung_emas -f migrations/002_your_migration_name.sql
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user

### User Profile
- `GET /api/v1/profile` - Get user profile
- `PATCH /api/v1/profile` - Update profile
- `POST /api/v1/profile/avatar` - Upload avatar
- `POST /api/v1/profile/change-password` - Change password

### Type Pockets (Categories)
- `GET /api/v1/type-pockets` - Get all type pockets
- `GET /api/v1/type-pockets/:id` - Get type pocket by ID

### Pockets
- `GET /api/v1/pockets` - Get all pockets
- `GET /api/v1/pockets/:id` - Get pocket by ID
- `POST /api/v1/pockets` - Create pocket
- `PATCH /api/v1/pockets/:id` - Update pocket
- `DELETE /api/v1/pockets/:id` - Delete pocket
- `GET /api/v1/pockets/:id/stats` - Get pocket statistics

### Transactions
- `GET /api/v1/transactions` - Get all transactions
- `GET /api/v1/transactions/:id` - Get transaction by ID
- `POST /api/v1/transactions` - Create transaction
- `PATCH /api/v1/transactions/:id` - Update transaction
- `DELETE /api/v1/transactions/:id` - Delete transaction
- `POST /api/v1/transactions/:id/receipt` - Upload receipt

### Analytics
- `GET /api/v1/analytics/dashboard` - Get dashboard summary
- `GET /api/v1/analytics/portfolio` - Get portfolio analytics
- `GET /api/v1/analytics/monthly-purchases` - Get monthly purchase analytics
- `GET /api/v1/analytics/brand-distribution` - Get brand distribution
- `GET /api/v1/analytics/trends` - Get transaction trends

### Gold Price
- `GET /api/v1/gold-price/current` - Get current gold price
- `GET /api/v1/gold-price/history` - Get historical gold prices

### Settings
- `GET /api/v1/settings` - Get user settings
- `PATCH /api/v1/settings` - Update settings

## Security Best Practices

- All passwords are hashed using bcrypt (cost factor: 12)
- SQL injection prevention using prepared statements
- Input validation on all endpoints
- JWT tokens with expiration
- CORS configuration
- Rate limiting on authentication endpoints (recommended)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For support, email support@emasgo.com or open an issue in the repository.
