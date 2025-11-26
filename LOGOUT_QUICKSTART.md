# Quick Start: Testing Logout Functionality

## 🚀 Start the Server

```bash
# Make sure database is running and migrated
make migrate

# Start the server
make run
```

## 🧪 Test with Script

```bash
# Run automated test
./test-logout.sh
```

Expected output:
```
🧪 Testing Logout Functionality
================================

1️⃣  Logging in...
✅ Access Token: eyJhbGciOiJIUzI1NiIs...

2️⃣  Testing protected route with valid token...
✅ Profile data retrieved

3️⃣  Logging out...
✅ Logout successful

4️⃣  Testing protected route with blacklisted token...
❌ 401 Unauthorized - Token has been revoked

✅ Logout test PASSED! Token was successfully blacklisted.
```

## 📮 Test with Postman

### Step 1: Login
```
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "SecurePass123",
  "remember_me": false
}
```

Response will include `access_token` - copy this!

### Step 2: Test Protected Route (Before Logout)
```
GET http://localhost:8080/api/v1/profile
Authorization: Bearer <your_access_token>
```

Should return: `200 OK` with user profile

### Step 3: Logout
```
POST http://localhost:8080/api/v1/auth/logout
Authorization: Bearer <your_access_token>
```

Should return:
```json
{
  "success": true,
  "message": "Logged out successfully",
  "data": null
}
```

### Step 4: Test Protected Route (After Logout)
```
GET http://localhost:8080/api/v1/profile
Authorization: Bearer <same_access_token>
```

Should return: `401 Unauthorized` with message "Token has been revoked"

## 🔍 Verify in Database

```bash
# Connect to database
psql -d nabung_emas

# Check blacklisted tokens
SELECT 
  id, 
  user_id, 
  LEFT(token_hash, 20) as token_hash_preview,
  expires_at,
  created_at
FROM token_blacklist
ORDER BY created_at DESC
LIMIT 5;
```

## 🐛 Troubleshooting

### Issue: "Token has been revoked" immediately after login
**Solution:** Clear the token blacklist table
```sql
TRUNCATE token_blacklist;
```

### Issue: Token still works after logout
**Check:**
1. Verify migration was applied: `\d token_blacklist` in psql
2. Check server logs for errors
3. Verify middleware is initialized with token blacklist repo

### Issue: 500 error on logout
**Check:**
1. Database connection is active
2. token_blacklist table exists
3. Server logs for specific error message

## 📊 Monitor Cleanup Service

The cleanup service runs every 24 hours. To test it immediately:

```go
// In internal/routes/routes.go, temporarily change:
cleanupService.StartTokenCleanup(1 * time.Minute) // Test with 1 minute
```

Check logs for:
```
Successfully cleaned up expired tokens
```

## 🔐 Security Best Practices

1. **Always use HTTPS in production** - Tokens in headers can be intercepted
2. **Set appropriate token expiration** - Balance security vs. user experience
3. **Implement refresh token rotation** - Additional security layer
4. **Monitor blacklist size** - Set up alerts if it grows too large
5. **Consider Redis for high traffic** - Faster lookups than PostgreSQL

## 📈 Performance Tips

1. **Database Indexes:** Already created on token_hash, user_id, expires_at
2. **Cleanup Frequency:** Adjust based on your token expiration times
3. **Connection Pooling:** Ensure database pool is properly configured
4. **Monitoring:** Track blacklist table size and query performance

## 🎯 Common Use Cases

### Logout from All Devices
Currently, logout only blacklists the current token. To implement "logout from all devices":

1. Store all active tokens per user
2. Add `RevokeAllUserTokens(userID)` method
3. Call on password change or security events

### Session Management
To view active sessions:

```sql
SELECT 
  user_id,
  COUNT(*) as active_sessions,
  MAX(expires_at) as latest_expiry
FROM token_blacklist
WHERE expires_at > NOW()
GROUP BY user_id;
```

### Force Logout User (Admin)
```go
// Add to auth service
func (s *AuthService) ForceLogoutUser(userID string) error {
    // Implementation to blacklist all user tokens
}
```

## ✅ Checklist

- [ ] Database migration applied
- [ ] Server compiles without errors
- [ ] Login works and returns token
- [ ] Protected routes work with valid token
- [ ] Logout successfully blacklists token
- [ ] Protected routes reject blacklisted token
- [ ] Cleanup service is running
- [ ] Postman collection updated

---

**Need Help?** Check `LOGOUT_DOCUMENTATION.md` for detailed information.
