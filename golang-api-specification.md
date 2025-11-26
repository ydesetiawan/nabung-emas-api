# EmasGo (Gold Savings) - Complete API Specification for Golang Backend

**Version:** 1.0.0  
**Last Updated:** 2025-11-26  
**Base URL:** `/api/v1`  
**Framework:** Echo (Golang)  
**Database:** PostgreSQL (Raw SQL - No ORM)

---

## Table of Contents

1. [Overview](#overview)
2. [Technology Stack](#technology-stack)
3. [Authentication](#authentication)
4. [Common Response Formats](#common-response-formats)
5. [Error Handling](#error-handling)
6. [API Endpoints](#api-endpoints)
   - [Authentication](#authentication-endpoints)
   - [User Profile](#user-profile-endpoints)
   - [Type Pockets (Categories)](#type-pockets-endpoints)
   - [Pockets](#pockets-endpoints)
   - [Transactions](#transactions-endpoints)
   - [Analytics](#analytics-endpoints)
   - [Gold Price](#gold-price-endpoints)
   - [Settings](#settings-endpoints)
7. [Data Models](#data-models)
8. [Database Schema](#database-schema)
9. [Implementation Guide](#implementation-guide)

---

## Overview

EmasGo is a gold savings tracking application that allows users to manage multiple gold savings pockets, track transactions, and analyze their portfolio performance.

### Key Features
- User authentication (login, signup, forgot password)
- Multiple gold savings pockets with categories
- Transaction management with receipt uploads
- Real-time portfolio analytics
- Multi-language support (English/Indonesian)
- Dark mode support
- User profile and settings management

---

## Technology Stack

### Backend
- **Framework:** Echo v4 (Golang)
- **Database:** PostgreSQL 14+
- **Database Driver:** `github.com/lib/pq`
- **Authentication:** JWT tokens
- **Password Hashing:** bcrypt
- **Validation:** `github.com/go-playground/validator/v10`
- **File Upload:** AWS S3 or local storage
- **CORS:** `github.com/labstack/echo/v4/middleware`

### Recommended Packages
```go
github.com/labstack/echo/v4
github.com/lib/pq
github.com/golang-jwt/jwt/v5
github.com/go-playground/validator/v10
golang.org/x/crypto/bcrypt
github.com/google/uuid
```

---

## Authentication

All protected endpoints require JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

### JWT Token Structure
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "exp": 1700000000,
  "iat": 1699900000
}
```

### Token Expiration
- Access Token: 24 hours
- Refresh Token: 7 days

---

## Common Response Formats

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {}
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "errors": {
    "field_name": ["Error 1", "Error 2"]
  }
}
```

---

## Error Handling

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 204 | No Content - Request successful, no content to return |
| 400 | Bad Request - Invalid input data |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 422 | Unprocessable Entity - Validation failed |
| 500 | Internal Server Error - Server error |

---

## API Endpoints

### Authentication Endpoints

#### 1. Register (Sign Up)

```
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone": "+62 812 3456 7890",
  "password": "SecurePass123",
  "confirm_password": "SecurePass123"
}
```

**Validation:**
- `full_name`: Required, min 3 chars, max 100 chars
- `email`: Required, valid email format, unique
- `phone`: Required, valid phone format
- `password`: Required, min 8 chars, must contain uppercase, lowercase, and number
- `confirm_password`: Required, must match password

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Account created successfully",
  "data": {
    "user": {
      "id": "uuid",
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+62 812 3456 7890",
      "created_at": "2025-11-26T08:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

#### 2. Login

```
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123",
  "remember_me": false
}
```

**Validation:**
- `email`: Required, valid email format
- `password`: Required, min 6 chars
- `remember_me`: Optional, boolean (extends token expiry to 30 days)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid",
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+62 812 3456 7890",
      "avatar": "https://storage.example.com/avatars/uuid.jpg",
      "created_at": "2025-11-26T08:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Invalid email or password"
}
```

#### 3. Forgot Password

```
POST /api/v1/auth/forgot-password
```

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset link sent to your email"
}
```

**Note:** Always return success even if email doesn't exist (security best practice)

#### 4. Reset Password

```
POST /api/v1/auth/reset-password
```

**Request Body:**
```json
{
  "token": "reset_token_from_email",
  "password": "NewSecurePass123",
  "confirm_password": "NewSecurePass123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

#### 5. Refresh Token

```
POST /api/v1/auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

#### 6. Logout

```
POST /api/v1/auth/logout
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

#### 7. Get Current User

```
GET /api/v1/auth/me
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+62 812 3456 7890",
    "avatar": "https://storage.example.com/avatars/uuid.jpg",
    "created_at": "2025-11-26T08:00:00Z"
  }
}
```

---

### User Profile Endpoints

#### 1. Get User Profile

```
GET /api/v1/profile
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+62 812 3456 7890",
      "avatar": "https://storage.example.com/avatars/uuid.jpg",
      "created_at": "2025-11-26T08:00:00Z"
    },
    "stats": {
      "total_pockets": 5,
      "total_transactions": 42,
      "total_weight": 25.5,
      "total_value": 26775000
    }
  }
}
```

#### 2. Update Profile

```
PATCH /api/v1/profile
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "full_name": "John Doe Updated",
  "phone": "+62 812 9999 8888"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": "uuid",
    "full_name": "John Doe Updated",
    "email": "john@example.com",
    "phone": "+62 812 9999 8888",
    "avatar": "https://storage.example.com/avatars/uuid.jpg",
    "updated_at": "2025-11-26T09:00:00Z"
  }
}
```

#### 3. Upload Avatar

```
POST /api/v1/profile/avatar
```

**Headers:** `Authorization: Bearer <token>`

**Request:** (multipart/form-data)
```
avatar: <file>
```

**Validation:**
- File type: jpg, jpeg, png, webp
- Max size: 2MB
- Min dimensions: 200x200px

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Avatar uploaded successfully",
  "data": {
    "avatar_url": "https://storage.example.com/avatars/uuid.jpg"
  }
}
```

#### 4. Change Password

```
POST /api/v1/profile/change-password
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "current_password": "OldPass123",
  "new_password": "NewSecurePass123",
  "confirm_password": "NewSecurePass123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

---

### Type Pockets Endpoints

Type Pockets are predefined categories for organizing gold savings.

#### 1. Get All Type Pockets

```
GET /api/v1/type-pockets
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Emergency Fund",
      "description": "Savings for emergency situations",
      "icon": "heroicons:shield-check",
      "color": "blue",
      "created_at": "2025-11-25T10:00:00Z"
    },
    {
      "id": "uuid",
      "name": "Wedding",
      "description": "Savings for wedding expenses",
      "icon": "heroicons:heart",
      "color": "pink",
      "created_at": "2025-11-25T10:00:00Z"
    }
  ]
}
```

#### 2. Get Type Pocket by ID

```
GET /api/v1/type-pockets/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Emergency Fund",
    "description": "Savings for emergency situations",
    "icon": "heroicons:shield-check",
    "color": "blue",
    "created_at": "2025-11-25T10:00:00Z"
  }
}
```

---

### Pockets Endpoints

Pockets are individual gold savings accounts belonging to a user.

#### 1. Get All Pockets

```
GET /api/v1/pockets
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `type_pocket_id` (optional): Filter by type pocket ID
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Items per page
- `sort_by` (optional, default: 'created_at'): Sort field
- `sort_order` (optional, default: 'desc'): 'asc' or 'desc'

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "type_pocket_id": "uuid",
      "name": "Emergency Fund - 2025",
      "description": "My emergency gold savings",
      "aggregate_total_price": 15000000,
      "aggregate_total_weight": 14.5,
      "target_weight": 20.0,
      "type_pocket": {
        "id": "uuid",
        "name": "Emergency Fund",
        "icon": "heroicons:shield-check",
        "color": "blue"
      },
      "created_at": "2025-11-25T10:00:00Z",
      "updated_at": "2025-11-25T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

#### 2. Get Pocket by ID

```
GET /api/v1/pockets/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "type_pocket_id": "uuid",
    "name": "Emergency Fund - 2025",
    "description": "My emergency gold savings",
    "aggregate_total_price": 15000000,
    "aggregate_total_weight": 14.5,
    "target_weight": 20.0,
    "type_pocket": {
      "id": "uuid",
      "name": "Emergency Fund",
      "description": "Savings for emergency situations",
      "icon": "heroicons:shield-check",
      "color": "blue"
    },
    "transaction_count": 12,
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T10:00:00Z"
  }
}
```

#### 3. Create Pocket

```
POST /api/v1/pockets
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "type_pocket_id": "uuid",
  "name": "Wedding Fund",
  "description": "Saving for my wedding",
  "target_weight": 50.0
}
```

**Validation:**
- `type_pocket_id`: Required, must be valid UUID and exist
- `name`: Required, 3-100 characters
- `description`: Optional, max 500 characters
- `target_weight`: Optional, must be > 0 if provided

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Pocket created successfully",
  "data": {
    "id": "uuid",
    "type_pocket_id": "uuid",
    "name": "Wedding Fund",
    "description": "Saving for my wedding",
    "aggregate_total_price": 0,
    "aggregate_total_weight": 0,
    "target_weight": 50.0,
    "created_at": "2025-11-26T08:00:00Z"
  }
}
```

#### 4. Update Pocket

```
PATCH /api/v1/pockets/:id
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "name": "Wedding Fund 2026",
  "description": "Updated description",
  "target_weight": 60.0
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Pocket updated successfully",
  "data": {
    "id": "uuid",
    "type_pocket_id": "uuid",
    "name": "Wedding Fund 2026",
    "description": "Updated description",
    "aggregate_total_price": 5000000,
    "aggregate_total_weight": 4.8,
    "target_weight": 60.0,
    "updated_at": "2025-11-26T09:00:00Z"
  }
}
```

#### 5. Delete Pocket

```
DELETE /api/v1/pockets/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Pocket deleted successfully"
}
```

**Note:** Implement cascade delete for related transactions or prevent deletion if transactions exist.

#### 6. Get Pocket Statistics

```
GET /api/v1/pockets/:id/stats
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `current_gold_price` (optional): Current market price per gram for profit/loss calculation

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total_weight": 14.5,
    "total_value": 15000000,
    "average_price_per_gram": 1034482.76,
    "current_gold_price": 1055000,
    "current_value": 15297500,
    "profit_loss": 297500,
    "profit_loss_percentage": 1.98,
    "transaction_count": 12
  }
}
```

---

### Transactions Endpoints

Transactions represent gold purchases added to pockets.

#### 1. Get All Transactions

```
GET /api/v1/transactions
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `pocket_id` (optional): Filter by pocket ID
- `start_date` (optional): Filter from date (YYYY-MM-DD)
- `end_date` (optional): Filter to date (YYYY-MM-DD)
- `brand` (optional): Filter by gold brand
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Items per page
- `sort_by` (optional, default: 'transaction_date'): Sort field
- `sort_order` (optional, default: 'desc'): 'asc' or 'desc'

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "pocket_id": "uuid",
      "transaction_date": "2025-11-20",
      "brand": "Antam",
      "weight": 2.5,
      "price_per_gram": 1050000,
      "total_price": 2625000,
      "description": "Monthly gold purchase",
      "receipt_image": "https://storage.example.com/receipts/abc123.jpg",
      "pocket": {
        "id": "uuid",
        "name": "Emergency Fund - 2025",
        "type_pocket": {
          "name": "Emergency Fund",
          "color": "blue"
        }
      },
      "created_at": "2025-11-25T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 42,
    "total_pages": 3
  }
}
```

#### 2. Get Transaction by ID

```
GET /api/v1/transactions/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "pocket_id": "uuid",
    "transaction_date": "2025-11-20",
    "brand": "Antam",
    "weight": 2.5,
    "price_per_gram": 1050000,
    "total_price": 2625000,
    "description": "Monthly gold purchase",
    "receipt_image": "https://storage.example.com/receipts/abc123.jpg",
    "pocket": {
      "id": "uuid",
      "name": "Emergency Fund - 2025",
      "type_pocket_id": "uuid",
      "type_pocket": {
        "id": "uuid",
        "name": "Emergency Fund",
        "icon": "heroicons:shield-check",
        "color": "blue"
      }
    },
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T10:00:00Z"
  }
}
```

#### 3. Create Transaction

```
POST /api/v1/transactions
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "pocket_id": "uuid",
  "transaction_date": "2025-11-25",
  "brand": "Antam",
  "weight": 2.5,
  "price_per_gram": 1050000,
  "total_price": 2625000,
  "description": "Monthly gold purchase",
  "receipt_image": "base64_encoded_image_or_url"
}
```

**Validation:**
- `pocket_id`: Required, must be valid UUID and belong to user
- `transaction_date`: Required, valid date (not future)
- `brand`: Required, one of: Antam, UBS, Pegadaian, King Halim, Custom
- `weight`: Required, 0.1 - 1000 grams
- `price_per_gram`: Required, 1000 - 10000000 IDR
- `total_price`: Required, must equal weight * price_per_gram
- `description`: Optional, max 500 characters
- `receipt_image`: Optional, valid image URL or base64

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Transaction created successfully",
  "data": {
    "id": "uuid",
    "pocket_id": "uuid",
    "transaction_date": "2025-11-25",
    "brand": "Antam",
    "weight": 2.5,
    "price_per_gram": 1050000,
    "total_price": 2625000,
    "description": "Monthly gold purchase",
    "receipt_image": "https://storage.example.com/receipts/abc123.jpg",
    "created_at": "2025-11-26T08:00:00Z"
  }
}
```

**Side Effects:**
- Updates pocket's `aggregate_total_price` and `aggregate_total_weight`

#### 4. Update Transaction

```
PATCH /api/v1/transactions/:id
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "weight": 3.0,
  "price_per_gram": 1055000,
  "total_price": 3165000,
  "description": "Updated monthly gold purchase"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Transaction updated successfully",
  "data": {
    "id": "uuid",
    "pocket_id": "uuid",
    "transaction_date": "2025-11-25",
    "brand": "Antam",
    "weight": 3.0,
    "price_per_gram": 1055000,
    "total_price": 3165000,
    "description": "Updated monthly gold purchase",
    "updated_at": "2025-11-26T09:00:00Z"
  }
}
```

**Side Effects:**
- Recalculates pocket's `aggregate_total_price` and `aggregate_total_weight`

#### 5. Delete Transaction

```
DELETE /api/v1/transactions/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Transaction deleted successfully"
}
```

**Side Effects:**
- Recalculates pocket's `aggregate_total_price` and `aggregate_total_weight`

#### 6. Upload Receipt Image

```
POST /api/v1/transactions/:id/receipt
```

**Headers:** `Authorization: Bearer <token>`

**Request:** (multipart/form-data)
```
receipt: <file>
```

**Validation:**
- File type: jpg, jpeg, png, pdf
- Max size: 5MB

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Receipt uploaded successfully",
  "data": {
    "receipt_image": "https://storage.example.com/receipts/abc123.jpg"
  }
}
```

---

### Analytics Endpoints

Analytics endpoints provide aggregated data and insights.

#### 1. Get Dashboard Summary

```
GET /api/v1/analytics/dashboard
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `current_gold_price` (optional): Current market price for calculations

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "portfolio": {
      "total_value": 15000000,
      "total_weight": 14.5,
      "total_pockets": 3,
      "total_transactions": 12,
      "average_price_per_gram": 1034482.76,
      "current_gold_price": 1055000,
      "current_value": 15297500,
      "profit_loss": 297500,
      "profit_loss_percentage": 1.98
    },
    "recent_transactions": [
      {
        "id": "uuid",
        "pocket_id": "uuid",
        "transaction_date": "2025-11-20",
        "brand": "Antam",
        "weight": 2.5,
        "total_price": 2625000,
        "pocket": {
          "id": "uuid",
          "name": "Emergency Fund - 2025",
          "type_pocket": {
            "color": "blue"
          }
        }
      }
    ]
  }
}
```

#### 2. Get Portfolio Analytics

```
GET /api/v1/analytics/portfolio
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `current_gold_price` (optional): Current market price

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total_value": 15000000,
    "total_weight": 14.5,
    "average_price_per_gram": 1034482.76,
    "current_market_price": 1055000,
    "current_value": 15297500,
    "profit_loss": 297500,
    "profit_loss_percentage": 1.98,
    "distribution": [
      {
        "pocket_id": "uuid",
        "pocket_name": "Emergency Fund - 2025",
        "type_pocket_name": "Emergency Fund",
        "type_pocket_color": "blue",
        "weight": 8.5,
        "value": 8500000,
        "percentage": 56.67
      }
    ]
  }
}
```

#### 3. Get Monthly Purchase Analytics

```
GET /api/v1/analytics/monthly-purchases
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `months` (optional, default: 6): Number of months to include
- `pocket_id` (optional): Filter by specific pocket

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "monthly_data": [
      {
        "month": "2025-06",
        "weight": 3.5,
        "amount": 3675000,
        "count": 2,
        "average_price_per_gram": 1050000
      },
      {
        "month": "2025-07",
        "weight": 2.5,
        "amount": 2625000,
        "count": 1,
        "average_price_per_gram": 1050000
      }
    ],
    "average_monthly_purchase": 3150000,
    "total_period_weight": 6.0,
    "total_period_amount": 6300000
  }
}
```

#### 4. Get Brand Distribution

```
GET /api/v1/analytics/brand-distribution
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "brand": "Antam",
      "weight": 10.5,
      "value": 11025000,
      "transaction_count": 8,
      "percentage": 72.41
    },
    {
      "brand": "UBS",
      "weight": 4.0,
      "value": 3975000,
      "transaction_count": 4,
      "percentage": 27.59
    }
  ]
}
```

#### 5. Get Transaction Trends

```
GET /api/v1/analytics/trends
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `period` (optional, default: '6m'): '1m', '3m', '6m', '1y', 'all'
- `group_by` (optional, default: 'month'): 'day', 'week', 'month'

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "trends": [
      {
        "period": "2025-06",
        "total_weight": 3.5,
        "total_value": 3675000,
        "transaction_count": 2,
        "average_price_per_gram": 1050000
      }
    ],
    "summary": {
      "total_weight": 14.5,
      "total_value": 15000000,
      "transaction_count": 12,
      "average_price_per_gram": 1034482.76,
      "lowest_price_per_gram": 1020000,
      "highest_price_per_gram": 1055000
    }
  }
}
```

---

### Gold Price Endpoints

External gold price tracking and historical data.

#### 1. Get Current Gold Price

```
GET /api/v1/gold-price/current
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "price_per_gram": 1055000,
    "currency": "IDR",
    "source": "Antam",
    "last_updated": "2025-11-26T08:00:00Z",
    "change_24h": 5000,
    "change_percentage_24h": 0.48
  }
}
```

#### 2. Get Historical Gold Prices

```
GET /api/v1/gold-price/history
```

**Query Parameters:**
- `start_date` (optional): Start date (YYYY-MM-DD)
- `end_date` (optional): End date (YYYY-MM-DD)
- `period` (optional, default: '1m'): '1w', '1m', '3m', '6m', '1y'
- `source` (optional): Filter by source (Antam, UBS, etc.)

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "date": "2025-11-25",
      "price_per_gram": 1055000,
      "source": "Antam"
    },
    {
      "date": "2025-11-24",
      "price_per_gram": 1050000,
      "source": "Antam"
    }
  ]
}
```

---

### Settings Endpoints

#### 1. Get User Settings

```
GET /api/v1/settings
```

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "language": "en",
    "theme": "dark",
    "currency": "IDR",
    "notifications": {
      "email": true,
      "push": false,
      "price_alerts": true
    }
  }
}
```

#### 2. Update Settings

```
PATCH /api/v1/settings
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "language": "id",
  "theme": "light",
  "notifications": {
    "email": true,
    "push": true,
    "price_alerts": false
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Settings updated successfully",
  "data": {
    "language": "id",
    "theme": "light",
    "currency": "IDR",
    "notifications": {
      "email": true,
      "push": true,
      "price_alerts": false
    }
  }
}
```

---

## Data Models

### User
```go
type User struct {
    ID        string    `json:"id"`
    FullName  string    `json:"full_name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Password  string    `json:"-"` // Never return in JSON
    Avatar    *string   `json:"avatar"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### TypePocket
```go
type TypePocket struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Icon        string    `json:"icon"`
    Color       string    `json:"color"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Pocket
```go
type Pocket struct {
    ID                   string      `json:"id"`
    UserID               string      `json:"user_id"`
    TypePocketID         string      `json:"type_pocket_id"`
    Name                 string      `json:"name"`
    Description          *string     `json:"description"`
    AggregateTotalPrice  float64     `json:"aggregate_total_price"`
    AggregateTotalWeight float64     `json:"aggregate_total_weight"`
    TargetWeight         *float64    `json:"target_weight"`
    TypePocket           *TypePocket `json:"type_pocket,omitempty"`
    CreatedAt            time.Time   `json:"created_at"`
    UpdatedAt            time.Time   `json:"updated_at"`
}
```

### Transaction
```go
type Transaction struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`
    PocketID        string    `json:"pocket_id"`
    TransactionDate time.Time `json:"transaction_date"`
    Brand           string    `json:"brand"`
    Weight          float64   `json:"weight"`
    PricePerGram    float64   `json:"price_per_gram"`
    TotalPrice      float64   `json:"total_price"`
    Description     *string   `json:"description"`
    ReceiptImage    *string   `json:"receipt_image"`
    Pocket          *Pocket   `json:"pocket,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### GoldPrice
```go
type GoldPrice struct {
    ID           string    `json:"id"`
    Date         time.Time `json:"date"`
    PricePerGram float64   `json:"price_per_gram"`
    Source       string    `json:"source"`
    CreatedAt    time.Time `json:"created_at"`
}
```

---

## Database Schema

### PostgreSQL Schema

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Password reset tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Type Pockets (Categories)
CREATE TABLE type_pockets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    color VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pockets
CREATE TABLE pockets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type_pocket_id UUID NOT NULL REFERENCES type_pockets(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    aggregate_total_price DECIMAL(15, 2) DEFAULT 0,
    aggregate_total_weight DECIMAL(10, 3) DEFAULT 0,
    target_weight DECIMAL(10, 3),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_pocket_name_per_user UNIQUE(user_id, name)
);

-- Transactions
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pocket_id UUID NOT NULL REFERENCES pockets(id) ON DELETE CASCADE,
    transaction_date DATE NOT NULL,
    brand VARCHAR(50) NOT NULL,
    weight DECIMAL(10, 3) NOT NULL CHECK (weight >= 0.1 AND weight <= 1000),
    price_per_gram DECIMAL(15, 2) NOT NULL CHECK (price_per_gram >= 1000 AND price_per_gram <= 10000000),
    total_price DECIMAL(15, 2) NOT NULL,
    description TEXT,
    receipt_image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gold Prices (for historical tracking)
CREATE TABLE gold_prices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date DATE NOT NULL,
    price_per_gram DECIMAL(15, 2) NOT NULL,
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_price_per_date_source UNIQUE(date, source)
);

-- User Settings
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language VARCHAR(10) DEFAULT 'en',
    theme VARCHAR(20) DEFAULT 'light',
    currency VARCHAR(10) DEFAULT 'IDR',
    email_notifications BOOLEAN DEFAULT TRUE,
    push_notifications BOOLEAN DEFAULT FALSE,
    price_alerts BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_settings_per_user UNIQUE(user_id)
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_pockets_user_id ON pockets(user_id);
CREATE INDEX idx_pockets_type_pocket_id ON pockets(type_pocket_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_pocket_id ON transactions(pocket_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_brand ON transactions(brand);
CREATE INDEX idx_gold_prices_date ON gold_prices(date);
CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_type_pockets_updated_at BEFORE UPDATE ON type_pockets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pockets_updated_at BEFORE UPDATE ON pockets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update pocket aggregates
CREATE OR REPLACE FUNCTION update_pocket_aggregates()
RETURNS TRIGGER AS $$
DECLARE
    v_pocket_id UUID;
BEGIN
    -- Determine which pocket to update
    IF TG_OP = 'DELETE' THEN
        v_pocket_id := OLD.pocket_id;
    ELSE
        v_pocket_id := NEW.pocket_id;
    END IF;

    -- Update pocket aggregates
    UPDATE pockets
    SET 
        aggregate_total_weight = COALESCE((
            SELECT SUM(weight) FROM transactions WHERE pocket_id = v_pocket_id
        ), 0),
        aggregate_total_price = COALESCE((
            SELECT SUM(total_price) FROM transactions WHERE pocket_id = v_pocket_id
        ), 0),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_pocket_id;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update pocket aggregates
CREATE TRIGGER trigger_update_pocket_aggregates
AFTER INSERT OR UPDATE OR DELETE ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_pocket_aggregates();

-- Seed data for type_pockets
INSERT INTO type_pockets (name, description, icon, color) VALUES
    ('Emergency Fund', 'Savings for emergency situations', 'heroicons:shield-check', 'blue'),
    ('Wedding', 'Savings for wedding expenses', 'heroicons:heart', 'pink'),
    ('Investment', 'General investment savings', 'heroicons:chart-bar', 'gold'),
    ('Education', 'Savings for education expenses', 'heroicons:academic-cap', 'green'),
    ('Retirement', 'Savings for retirement', 'heroicons:home', 'purple'),
    ('Vacation', 'Savings for vacation and travel', 'heroicons:globe-alt', 'cyan'),
    ('Business', 'Savings for business capital', 'heroicons:briefcase', 'orange');
```

---

## Implementation Guide

### Project Structure

```
nabung-emas-backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── profile.go
│   │   ├── type_pockets.go
│   │   ├── pockets.go
│   │   ├── transactions.go
│   │   ├── analytics.go
│   │   ├── gold_price.go
│   │   └── settings.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── models/
│   │   ├── user.go
│   │   ├── type_pocket.go
│   │   ├── pocket.go
│   │   ├── transaction.go
│   │   ├── gold_price.go
│   │   └── response.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── type_pocket_repository.go
│   │   ├── pocket_repository.go
│   │   ├── transaction_repository.go
│   │   ├── gold_price_repository.go
│   │   └── settings_repository.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── pocket_service.go
│   │   ├── transaction_service.go
│   │   ├── analytics_service.go
│   │   ├── gold_price_service.go
│   │   └── storage_service.go
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── password.go
│   │   ├── validator.go
│   │   └── response.go
│   └── routes/
│       └── routes.go
├── migrations/
│   └── 001_initial_schema.sql
├── .env.example
├── go.mod
└── go.sum
```

### Example: main.go

```go
package main

import (
    "log"
    "os"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "nabung-emas-backend/internal/config"
    "nabung-emas-backend/internal/database"
    "nabung-emas-backend/internal/routes"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Initialize database
    db, err := database.NewPostgresDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Initialize Echo
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // Setup routes
    routes.Setup(e, db, cfg)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := e.Start(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### Example: User Repository (Raw SQL)

```go
package repositories

import (
    "database/sql"
    "errors"
    "time"

    "github.com/google/uuid"
    "nabung-emas-backend/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (id, full_name, email, phone, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

    user.ID = uuid.New().String()
    now := time.Now()

    err := r.db.QueryRow(
        query,
        user.ID,
        user.FullName,
        user.Email,
        user.Phone,
        user.Password,
        now,
        now,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

    return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    query := `
        SELECT id, full_name, email, phone, password_hash, avatar, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    user := &models.User{}
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.FullName,
        &user.Email,
        &user.Phone,
        &user.Password,
        &user.Avatar,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }

    return user, err
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
    query := `
        SELECT id, full_name, email, phone, avatar, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    user := &models.User{}
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.FullName,
        &user.Email,
        &user.Phone,
        &user.Avatar,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }

    return user, err
}

func (r *UserRepository) Update(user *models.User) error {
    query := `
        UPDATE users
        SET full_name = $1, phone = $2, avatar = $3, updated_at = $4
        WHERE id = $5
        RETURNING updated_at
    `

    err := r.db.QueryRow(
        query,
        user.FullName,
        user.Phone,
        user.Avatar,
        time.Now(),
        user.ID,
    ).Scan(&user.UpdatedAt)

    return err
}
```

### Example: Auth Handler

```go
package handlers

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "nabung-emas-backend/internal/models"
    "nabung-emas-backend/internal/services"
    "nabung-emas-backend/internal/utils"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
    var req models.RegisterRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
    }

    if err := c.Validate(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
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
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
    }

    if err := c.Validate(&req); err != nil {
        return utils.ValidationErrorResponse(c, err)
    }

    user, tokens, err := h.authService.Login(&req)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
    }

    return utils.SuccessResponse(c, http.StatusOK, "Login successful", map[string]interface{}{
        "user":          user,
        "access_token":  tokens.AccessToken,
        "refresh_token": tokens.RefreshToken,
        "expires_in":    tokens.ExpiresIn,
    })
}
```

### Environment Variables (.env.example)

```env
# Server
PORT=8080
ENV=development

# Database
DATABASE_URL=postgres://user:password@localhost:5432/nabung_emas?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

# Storage
STORAGE_TYPE=local
AWS_REGION=ap-southeast-1
AWS_BUCKET=nabung-emas-storage
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key

# CORS
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Email (for password reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@emasgo.com
```

---

## Best Practices

### 1. Security
- Always hash passwords with bcrypt (cost factor: 12)
- Use prepared statements to prevent SQL injection
- Validate all user inputs
- Implement rate limiting on auth endpoints
- Use HTTPS in production
- Set secure JWT secret (min 32 characters)

### 2. Database
- Use transactions for operations that modify multiple tables
- Create proper indexes on frequently queried columns
- Use connection pooling
- Implement database migrations

### 3. Error Handling
- Never expose internal errors to clients
- Log all errors with context
- Use consistent error response format
- Return appropriate HTTP status codes

### 4. Performance
- Implement pagination for list endpoints
- Use database indexes effectively
- Cache frequently accessed data (type pockets, gold prices)
- Optimize SQL queries (avoid N+1 queries)

### 5. Testing
- Write unit tests for services and repositories
- Write integration tests for handlers
- Test edge cases and error scenarios
- Aim for >80% code coverage

---

**End of API Specification**
