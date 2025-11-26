#!/bin/bash

# EmasGo API Test Script
# This script tests all major endpoints to verify the API is working correctly

echo "🧪 Testing EmasGo API..."
echo ""

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
PASSED=0
FAILED=0

# Helper function to test endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5
    local token=$6
    
    echo -n "Testing: $name... "
    
    if [ -n "$token" ]; then
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $token" \
                -d "$data")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$endpoint" \
                -H "Authorization: Bearer $token")
        fi
    else
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$endpoint")
        fi
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ PASSED${NC} (HTTP $status_code)"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (Expected HTTP $expected_status, got $status_code)"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# 1. Test Health Check
echo -e "${BLUE}=== Health Check ===${NC}"
test_endpoint "Health Check" "GET" "/health" "" "200"
echo ""

# 2. Test Type Pockets (Public)
echo -e "${BLUE}=== Type Pockets (Public) ===${NC}"
test_endpoint "Get All Type Pockets" "GET" "/type-pockets" "" "200"
echo ""

# 3. Test Authentication
echo -e "${BLUE}=== Authentication ===${NC}"

# Register a new user
REGISTER_DATA='{
  "full_name": "API Test User",
  "email": "apitest@example.com",
  "phone": "+62 812 9999 8888",
  "password": "TestPass123",
  "confirm_password": "TestPass123"
}'

response=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "$REGISTER_DATA")

if echo "$response" | grep -q "success.*true"; then
    echo -e "Testing: Register New User... ${GREEN}✓ PASSED${NC}"
    PASSED=$((PASSED + 1))
    ACCESS_TOKEN=$(echo "$response" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['access_token'])" 2>/dev/null)
else
    echo -e "Testing: Register New User... ${RED}✗ FAILED${NC}"
    FAILED=$((FAILED + 1))
    # Try to login instead
    LOGIN_DATA='{"email": "test@example.com", "password": "SecurePass123"}'
    response=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "$LOGIN_DATA")
    ACCESS_TOKEN=$(echo "$response" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['access_token'])" 2>/dev/null)
fi

# Login
LOGIN_DATA='{"email": "apitest@example.com", "password": "TestPass123"}'
test_endpoint "Login" "POST" "/auth/login" "$LOGIN_DATA" "200"

# Get Current User
test_endpoint "Get Current User" "GET" "/auth/me" "" "200" "$ACCESS_TOKEN"
echo ""

# 4. Test Profile
echo -e "${BLUE}=== User Profile ===${NC}"
test_endpoint "Get Profile" "GET" "/profile" "" "200" "$ACCESS_TOKEN"
echo ""

# 5. Test Pockets
echo -e "${BLUE}=== Pockets ===${NC}"
test_endpoint "Get All Pockets" "GET" "/pockets" "" "200" "$ACCESS_TOKEN"

# Create a pocket
CREATE_POCKET_DATA='{
  "type_pocket_id": "fd60cd5a-363a-4a48-b491-f3519c4092d9",
  "name": "Test Emergency Fund",
  "description": "Testing pocket creation",
  "target_weight": 25.0
}'
response=$(curl -s -X POST "$API_URL/pockets" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$CREATE_POCKET_DATA")

if echo "$response" | grep -q "success.*true"; then
    echo -e "Testing: Create Pocket... ${GREEN}✓ PASSED${NC}"
    PASSED=$((PASSED + 1))
    POCKET_ID=$(echo "$response" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
else
    echo -e "Testing: Create Pocket... ${RED}✗ FAILED${NC}"
    FAILED=$((FAILED + 1))
fi
echo ""

# 6. Test Transactions
if [ -n "$POCKET_ID" ]; then
    echo -e "${BLUE}=== Transactions ===${NC}"
    test_endpoint "Get All Transactions" "GET" "/transactions" "" "200" "$ACCESS_TOKEN"
    
    # Create a transaction
    CREATE_TRANSACTION_DATA="{
      \"pocket_id\": \"$POCKET_ID\",
      \"transaction_date\": \"2025-11-25\",
      \"brand\": \"Antam\",
      \"weight\": 2.5,
      \"price_per_gram\": 1050000,
      \"total_price\": 2625000,
      \"description\": \"Test transaction\"
    }"
    test_endpoint "Create Transaction" "POST" "/transactions" "$CREATE_TRANSACTION_DATA" "201" "$ACCESS_TOKEN"
    echo ""
fi

# 7. Test Analytics
echo -e "${BLUE}=== Analytics ===${NC}"
test_endpoint "Get Dashboard" "GET" "/analytics/dashboard" "" "200" "$ACCESS_TOKEN"
test_endpoint "Get Portfolio" "GET" "/analytics/portfolio" "" "200" "$ACCESS_TOKEN"
test_endpoint "Get Brand Distribution" "GET" "/analytics/brand-distribution" "" "200" "$ACCESS_TOKEN"
echo ""

# 8. Test Settings
echo -e "${BLUE}=== Settings ===${NC}"
test_endpoint "Get Settings" "GET" "/settings" "" "200" "$ACCESS_TOKEN"
echo ""

# Summary
echo -e "${BLUE}=== Test Summary ===${NC}"
TOTAL=$((PASSED + FAILED))
echo "Total Tests: $TOTAL"
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi
