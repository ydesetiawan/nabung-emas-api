#!/bin/bash

# Test script for logout functionality
BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testing Logout Functionality"
echo "================================"
echo ""

# Step 1: Login
echo "1️⃣  Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123",
    "remember_me": false
  }')

echo "Login Response: $LOGIN_RESPONSE"
echo ""

# Extract access token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
  echo "❌ Failed to get access token"
  exit 1
fi

echo "✅ Access Token: ${ACCESS_TOKEN:0:20}..."
echo ""

# Step 2: Test accessing protected route with token
echo "2️⃣  Testing protected route with valid token..."
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "Profile Response: $PROFILE_RESPONSE"
echo ""

# Step 3: Logout
echo "3️⃣  Logging out..."
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "Logout Response: $LOGOUT_RESPONSE"
echo ""

# Step 4: Try to access protected route with blacklisted token
echo "4️⃣  Testing protected route with blacklisted token..."
PROFILE_AFTER_LOGOUT=$(curl -s -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "Profile After Logout Response: $PROFILE_AFTER_LOGOUT"
echo ""

# Check if the response contains "revoked" or "unauthorized"
if echo "$PROFILE_AFTER_LOGOUT" | grep -qi "revoked\|unauthorized"; then
  echo "✅ Logout test PASSED! Token was successfully blacklisted."
else
  echo "❌ Logout test FAILED! Token is still valid after logout."
fi

echo ""
echo "================================"
echo "Test Complete!"
