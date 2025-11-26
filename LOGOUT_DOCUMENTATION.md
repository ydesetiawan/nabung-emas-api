# Logout Functionality Documentation

## Overview

The logout functionality has been implemented with **token blacklisting** to ensure that JWT tokens are properly invalidated when users log out. This provides a secure logout mechanism even though JWTs are stateless.

## Implementation Details

### 1. Database Schema

A new table `token_blacklist` has been created to store blacklisted tokens:

```sql
CREATE TABLE token_blacklist (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Key Features:**
- Tokens are stored as SHA-256 hashes for security
- Includes expiration time to allow automatic cleanup
- Indexed for fast lookups

### 2. Components

#### TokenBlacklistRepository
**Location:** `internal/repositories/token_blacklist_repository.go`

**Methods:**
- `Add(token, userID, expiresAt)` - Adds a token to the blacklist
- `IsBlacklisted(token)` - Checks if a token is blacklisted
- `CleanupExpired()` - Removes expired tokens from the database
- `HashToken(token)` - Creates SHA-256 hash of token for storage

#### AuthService.Logout
**Location:** `internal/services/auth_service.go`

**Functionality:**
- Extracts token expiration time from JWT claims
- Adds token to blacklist with proper expiration
- Handles both valid and invalid tokens gracefully

#### AuthMiddleware.RequireAuth
**Location:** `internal/middleware/auth.go`

**Updated to:**
1. Check if token is blacklisted before validating
2. Reject requests with blacklisted tokens
3. Return appropriate error message

#### CleanupService
**Location:** `internal/services/cleanup_service.go`

**Purpose:**
- Runs background job to clean up expired tokens
- Configurable interval (default: 24 hours)
- Prevents database bloat from old tokens

### 3. API Endpoint

**Endpoint:** `POST /api/v1/auth/logout`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Logged out successfully",
  "data": null
}
```

**Error Responses:**

- **401 Unauthorized** - Missing or invalid token
- **500 Internal Server Error** - Failed to blacklist token

### 4. Flow Diagram

```
User Request (POST /auth/logout)
    ↓
AuthMiddleware.RequireAuth
    ↓
Check if token is blacklisted → Yes → Return 401 "Token has been revoked"
    ↓ No
Validate JWT token
    ↓
AuthHandler.Logout
    ↓
Extract token from Authorization header
    ↓
AuthService.Logout
    ↓
Validate token to get expiration time
    ↓
Add token hash to blacklist table
    ↓
Return success response
```

### 5. Security Considerations

1. **Token Hashing:** Tokens are hashed using SHA-256 before storage to prevent token leakage if database is compromised

2. **Expiration Tracking:** Blacklisted tokens include expiration time, allowing automatic cleanup

3. **Performance:** 
   - Indexed token_hash column for fast lookups
   - Automatic cleanup prevents database bloat
   - Minimal overhead on protected routes

4. **Edge Cases Handled:**
   - Invalid/expired tokens can still be blacklisted
   - Duplicate blacklist attempts are ignored (ON CONFLICT DO NOTHING)
   - Database errors return appropriate HTTP status codes

### 6. Testing

Use the provided test script:

```bash
./test-logout.sh
```

**Test Steps:**
1. Login to get access token
2. Access protected route with valid token (should succeed)
3. Logout to blacklist token
4. Try to access protected route with blacklisted token (should fail)

**Expected Result:**
```
✅ Logout test PASSED! Token was successfully blacklisted.
```

### 7. Postman Collection

The Postman collection has been updated with a test script for the logout endpoint that automatically clears environment variables after successful logout.

### 8. Maintenance

**Automatic Cleanup:**
The cleanup service runs every 24 hours to remove expired tokens from the blacklist. This can be configured in `internal/routes/routes.go`:

```go
cleanupService.StartTokenCleanup(24 * time.Hour) // Adjust interval as needed
```

**Manual Cleanup:**
You can also manually clean up expired tokens:

```sql
DELETE FROM token_blacklist WHERE expires_at <= NOW();
```

### 9. Future Enhancements

Potential improvements:
1. **Redis Integration:** Use Redis for faster blacklist lookups
2. **Revoke All User Tokens:** Implement functionality to revoke all tokens for a user (useful for password changes)
3. **Admin Dashboard:** Add admin interface to view and manage blacklisted tokens
4. **Metrics:** Track logout frequency and blacklist size

## Migration

To apply the token blacklist table:

```bash
make migrate
```

Or manually:

```bash
psql -d nabung_emas -f migrations/002_add_token_blacklist.sql
```

## Troubleshooting

**Issue:** Logout returns 500 error
- Check database connection
- Verify token_blacklist table exists
- Check application logs for specific error

**Issue:** Token still works after logout
- Verify middleware is checking blacklist before validation
- Check if token was successfully added to blacklist table
- Ensure cleanup service hasn't removed the token prematurely

**Issue:** Performance degradation
- Check blacklist table size
- Verify cleanup service is running
- Consider adding more indexes if needed
