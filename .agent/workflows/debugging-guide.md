---
description: Step-by-step debugging guide for the nabung-emas-api
---

# Debugging Guide for nabung-emas-api

## Understanding the Issue

The error `pq: invalid input syntax for type uuid: ""` occurs when PostgreSQL receives an empty string where it expects a UUID value. This typically happens in:
- INSERT/UPDATE operations with UUID columns
- WHERE clauses comparing UUID columns
- Function parameters expecting UUID types

## Step-by-Step Debugging Process

### 1. **Identify the Error Location**

When you receive an error response:
```json
{
  "success": false,
  "message": "pq: invalid input syntax for type uuid: \"\""
}
```

**Steps:**
1. Look at the API endpoint that failed (e.g., `POST /api/v1/pockets`)
2. Find the corresponding handler in `internal/handlers/`
3. Trace the flow: Handler → Service → Repository

### 2. **Add Debug Logging**

Add strategic log statements to track data flow:

```go
// In handler (e.g., pocket_handler.go)
func (h *PocketHandler) Create(c echo.Context) error {
    userID := middleware.GetUserID(c)
    fmt.Printf("DEBUG: UserID from middleware: '%s'\n", userID)
    
    var req models.CreatePocketRequest
    if err := utils.BindAndValidate(c, &req); err != nil {
        return utils.HandleError(c, err)
    }
    fmt.Printf("DEBUG: Request data: %+v\n", req)
    
    pocket, err := h.service.Create(userID, &req)
    if err != nil {
        fmt.Printf("DEBUG: Service error: %v\n", err)
        return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
    }
    
    return utils.SuccessResponse(c, http.StatusCreated, "Pocket created successfully", pocket)
}
```

```go
// In service (e.g., pocket_service.go)
func (s *PocketService) Create(userID string, req *models.CreatePocketRequest) (*models.Pocket, error) {
    fmt.Printf("DEBUG Service: userID='%s', typePocketID='%s'\n", userID, req.TypePocketID)
    
    // ... rest of the code
}
```

```go
// In repository (e.g., pocket_repository.go)
func (r *PocketRepository) Create(pocket *models.Pocket) error {
    fmt.Printf("DEBUG Repo: Creating pocket - UserID='%s', TypePocketID='%s'\n", 
        pocket.UserID, pocket.TypePocketID)
    
    // ... rest of the code
}
```

### 3. **Check Authentication Middleware**

Verify that the JWT token is being parsed correctly:

```go
// In internal/middleware/auth.go
func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        authHeader := c.Request().Header.Get("Authorization")
        fmt.Printf("DEBUG Auth: Authorization header: %s\n", authHeader)
        
        // ... token parsing ...
        
        claims, err := utils.ValidateToken(tokenString, m.config.JWTSecret)
        if err != nil {
            fmt.Printf("DEBUG Auth: Token validation error: %v\n", err)
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
        }
        
        fmt.Printf("DEBUG Auth: Claims - UserID='%s', Email='%s'\n", claims.UserID, claims.Email)
        
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        
        return next(c)
    }
}
```

### 4. **Inspect Database Queries**

Enable PostgreSQL query logging to see actual SQL being executed:

**Option A: In your code (add to repository methods)**
```go
func (r *PocketRepository) Create(pocket *models.Pocket) error {
    query := `
        INSERT INTO pockets (id, user_id, type_pocket_id, name, description, target_weight, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, aggregate_total_price, aggregate_total_weight, created_at, updated_at
    `
    
    pocket.ID = uuid.New().String()
    now := time.Now()
    
    // Debug: Print query and parameters
    fmt.Printf("DEBUG SQL Query: %s\n", query)
    fmt.Printf("DEBUG SQL Params: id=%s, user_id=%s, type_pocket_id=%s, name=%s\n", 
        pocket.ID, pocket.UserID, pocket.TypePocketID, pocket.Name)
    
    err := r.db.QueryRow(
        query,
        pocket.ID,
        pocket.UserID,
        pocket.TypePocketID,
        pocket.Name,
        pocket.Description,
        pocket.TargetWeight,
        now,
        now,
    ).Scan(
        &pocket.ID,
        &pocket.AggregateTotalPrice,
        &pocket.AggregateTotalWeight,
        &pocket.CreatedAt,
        &pocket.UpdatedAt,
    )
    
    if err != nil {
        fmt.Printf("DEBUG SQL Error: %v\n", err)
    }
    
    return err
}
```

**Option B: Enable PostgreSQL logging**
```sql
-- Connect to your database and run:
ALTER DATABASE your_database_name SET log_statement = 'all';
ALTER DATABASE your_database_name SET log_duration = on;

-- Then check logs at:
-- macOS (Homebrew): /usr/local/var/log/postgresql@14.log
-- Linux: /var/log/postgresql/postgresql-14-main.log
```

### 5. **Use Postman/cURL with Verbose Output**

Test the endpoint with detailed request/response:

```bash
curl -v -X POST http://localhost:8080/api/v1/pockets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "type_pocket_id": "abe2731b-1e5e-4f38-8825-9f389b7273e7",
    "name": "My Emergency Fund",
    "description": "Savings for emergency situations",
    "target_weight": 50.0
  }'
```

### 6. **Validate Input Data**

Check that all UUID fields in the request are valid:

```go
// Add validation helper
func isValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}

// Use in handler or service
if !isValidUUID(req.TypePocketID) {
    return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid type_pocket_id format")
}
```

### 7. **Check Database Schema**

Verify that your database columns are correctly defined:

```sql
-- Check the pockets table structure
\d pockets

-- Expected output should show:
-- id: uuid
-- user_id: uuid
-- type_pocket_id: uuid
-- etc.

-- Verify data types
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'pockets';
```

### 8. **Test with psql**

Manually test the query in PostgreSQL:

```sql
-- Test the problematic query
SELECT EXISTS(
    SELECT 1 FROM pockets 
    WHERE user_id = '49f7a78a-b94b-4a08-958a-6cf162506816' 
    AND name = 'My Emergency Fund' 
    AND id != ''  -- This will fail!
);

-- Correct version (without empty UUID)
SELECT EXISTS(
    SELECT 1 FROM pockets 
    WHERE user_id = '49f7a78a-b94b-4a08-958a-6cf162506816' 
    AND name = 'My Emergency Fund'
);
```

## Common UUID-Related Issues

### Issue 1: Empty String Passed as UUID
**Symptom:** `pq: invalid input syntax for type uuid: ""`
**Solution:** Check for empty strings before using in UUID comparisons

```go
// Bad
query := `SELECT * FROM table WHERE id != $1`
db.Query(query, "")  // Error!

// Good
if excludeID != "" {
    query := `SELECT * FROM table WHERE id != $1`
    db.Query(query, excludeID)
} else {
    query := `SELECT * FROM table`
    db.Query(query)
}
```

### Issue 2: NULL vs Empty String
**Symptom:** Unexpected behavior with optional UUID fields
**Solution:** Use proper NULL handling

```go
// Use sql.NullString for optional UUID fields
var userID sql.NullString
if id != "" {
    userID = sql.NullString{String: id, Valid: true}
} else {
    userID = sql.NullString{Valid: false}
}
```

### Issue 3: Invalid UUID Format
**Symptom:** `pq: invalid input syntax for type uuid: "invalid-uuid"`
**Solution:** Validate UUID format before database operations

```go
import "github.com/google/uuid"

func validateUUID(id string) error {
    if _, err := uuid.Parse(id); err != nil {
        return fmt.Errorf("invalid UUID format: %s", id)
    }
    return nil
}
```

## IDE Debugging Setup

### Using VS Code

1. **Install Go extension** (if not already installed)

2. **Create launch.json** in `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch API",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/api",
            "env": {
                "ENV": "development"
            },
            "args": []
        },
        {
            "name": "Attach to Process",
            "type": "go",
            "request": "attach",
            "mode": "local",
            "processId": "${command:pickProcess}"
        }
    ]
}
```

3. **Set breakpoints:**
   - Click on the left margin of any line to set a breakpoint
   - Red dot will appear

4. **Start debugging:**
   - Press F5 or click "Run and Debug"
   - Make API request
   - Execution will pause at breakpoints

5. **Inspect variables:**
   - Hover over variables to see values
   - Use Debug Console to evaluate expressions
   - Check Call Stack panel

### Using Delve (CLI Debugger)

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debugging
dlv debug ./cmd/api

# Set breakpoint
(dlv) break pocket_handler.go:72

# Continue execution
(dlv) continue

# Inspect variables
(dlv) print userID
(dlv) print req

# Step through code
(dlv) next    # Next line
(dlv) step    # Step into function
(dlv) stepout # Step out of function

# List source code
(dlv) list

# Show goroutines
(dlv) goroutines

# Exit
(dlv) quit
```

## Quick Debugging Checklist

- [ ] Check error message for specific details
- [ ] Verify JWT token is valid and not expired
- [ ] Confirm userID is being extracted from token
- [ ] Validate all UUID fields in request body
- [ ] Check database schema matches code expectations
- [ ] Add debug logging at each layer (handler → service → repository)
- [ ] Test SQL queries directly in psql
- [ ] Verify middleware is properly configured in routes
- [ ] Check for empty strings being passed to UUID fields
- [ ] Review recent code changes that might affect the flow

## Testing the Fix

After applying fixes, test with:

```bash
# 1. Restart the server
make run

# 2. Get a fresh JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "your_password"
  }'

# 3. Test the create pocket endpoint
curl -X POST http://localhost:8080/api/v1/pockets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_NEW_TOKEN" \
  -d '{
    "type_pocket_id": "abe2731b-1e5e-4f38-8825-9f389b7273e7",
    "name": "My Emergency Fund",
    "description": "Savings for emergency situations",
    "target_weight": 50.0
  }'

# Expected response:
# {
#   "success": true,
#   "message": "Pocket created successfully",
#   "data": { ... }
# }
```

## Additional Resources

- [PostgreSQL UUID Documentation](https://www.postgresql.org/docs/current/datatype-uuid.html)
- [Go UUID Package](https://github.com/google/uuid)
- [Delve Debugger](https://github.com/go-delve/delve)
- [Echo Framework Debugging](https://echo.labstack.com/guide/debugging/)
