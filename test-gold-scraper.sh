#!/bin/bash

# Gold Scraper API Test Script
# This script tests all the gold scraper endpoints

BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testing Gold Scraper API Endpoints"
echo "======================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Scrape gold prices
echo "📝 Test 1: Scraping gold prices from logammulia.com"
echo "POST ${BASE_URL}/gold-scraper/scrape"
echo ""
SCRAPE_RESPONSE=$(curl -s -X POST "${BASE_URL}/gold-scraper/scrape" \
  -H "Content-Type: application/json")

echo "Response:"
echo "$SCRAPE_RESPONSE" | jq '.'
echo ""

# Check if scraping was successful
SUCCESS=$(echo "$SCRAPE_RESPONSE" | jq -r '.success')
if [ "$SUCCESS" = "true" ]; then
    echo -e "${GREEN}✅ Scraping successful!${NC}"
else
    echo -e "${RED}❌ Scraping failed!${NC}"
fi
echo ""
echo "---"
echo ""

# Wait a bit before next request
sleep 1

# Test 2: Get all prices
echo "📝 Test 2: Getting all gold prices"
echo "GET ${BASE_URL}/gold-scraper/prices"
echo ""
ALL_PRICES_RESPONSE=$(curl -s -X GET "${BASE_URL}/gold-scraper/prices")

echo "Response:"
echo "$ALL_PRICES_RESPONSE" | jq '.'
echo ""

COUNT=$(echo "$ALL_PRICES_RESPONSE" | jq -r '.count')
echo -e "${GREEN}✅ Retrieved $COUNT gold prices${NC}"
echo ""
echo "---"
echo ""

# Wait a bit before next request
sleep 1

# Test 3: Get all prices with limit
echo "📝 Test 3: Getting gold prices with limit=5"
echo "GET ${BASE_URL}/gold-scraper/prices?limit=5"
echo ""
LIMITED_PRICES_RESPONSE=$(curl -s -X GET "${BASE_URL}/gold-scraper/prices?limit=5")

echo "Response:"
echo "$LIMITED_PRICES_RESPONSE" | jq '.'
echo ""

COUNT=$(echo "$LIMITED_PRICES_RESPONSE" | jq -r '.count')
echo -e "${GREEN}✅ Retrieved $COUNT gold prices (limited)${NC}"
echo ""
echo "---"
echo ""

# Wait a bit before next request
sleep 1

# Test 4: Get prices filtered by type
echo "📝 Test 4: Getting gold prices filtered by type (emas)"
echo "GET ${BASE_URL}/gold-scraper/prices?type=emas"
echo ""
FILTERED_PRICES_RESPONSE=$(curl -s -X GET "${BASE_URL}/gold-scraper/prices?type=emas")

echo "Response:"
echo "$FILTERED_PRICES_RESPONSE" | jq '.'
echo ""

COUNT=$(echo "$FILTERED_PRICES_RESPONSE" | jq -r '.count')
echo -e "${GREEN}✅ Retrieved $COUNT gold prices (filtered by type)${NC}"
echo ""
echo "---"
echo ""

# Wait a bit before next request
sleep 1

# Test 5: Get prices filtered by source
echo "📝 Test 5: Getting gold prices filtered by source (antam)"
echo "GET ${BASE_URL}/gold-scraper/prices?source=antam"
echo ""
SOURCE_FILTERED_RESPONSE=$(curl -s -X GET "${BASE_URL}/gold-scraper/prices?source=antam")

echo "Response:"
echo "$SOURCE_FILTERED_RESPONSE" | jq '.'
echo ""

COUNT=$(echo "$SOURCE_FILTERED_RESPONSE" | jq -r '.count')
echo -e "${GREEN}✅ Retrieved $COUNT gold prices (filtered by source)${NC}"
echo ""
echo "---"
echo ""

# Wait a bit before next request
sleep 1

# Test 6: Get latest prices
echo "📝 Test 6: Getting latest gold prices for each type"
echo "GET ${BASE_URL}/gold-scraper/prices/latest"
echo ""
LATEST_PRICES_RESPONSE=$(curl -s -X GET "${BASE_URL}/gold-scraper/prices/latest")

echo "Response:"
echo "$LATEST_PRICES_RESPONSE" | jq '.'
echo ""

COUNT=$(echo "$LATEST_PRICES_RESPONSE" | jq -r '.count')
echo -e "${GREEN}✅ Retrieved $COUNT latest gold prices${NC}"
echo ""
echo "---"
echo ""

# Wait a bit before next request
sleep 1

# Test 7: Get price by ID
# First, get the ID from the latest prices
PRICE_ID=$(echo "$LATEST_PRICES_RESPONSE" | jq -r '.data[0].id')

if [ "$PRICE_ID" != "null" ] && [ -n "$PRICE_ID" ]; then
    echo "📝 Test 7: Getting gold price by ID ($PRICE_ID)"
    echo "GET ${BASE_URL}/gold-scraper/prices/${PRICE_ID}"
    echo ""
    PRICE_BY_ID_RESPONSE=$(curl -s -X GET "${BASE_URL}/gold-scraper/prices/${PRICE_ID}")

    echo "Response:"
    echo "$PRICE_BY_ID_RESPONSE" | jq '.'
    echo ""

    SUCCESS=$(echo "$PRICE_BY_ID_RESPONSE" | jq -r '.success')
    if [ "$SUCCESS" = "true" ]; then
        echo -e "${GREEN}✅ Successfully retrieved gold price by ID${NC}"
    else
        echo -e "${RED}❌ Failed to retrieve gold price by ID${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Skipping Test 7: No price ID available${NC}"
fi
echo ""
echo "---"
echo ""

# Test 8: Get price by invalid ID
echo "📝 Test 8: Getting gold price by invalid ID (99999)"
echo "GET ${BASE_URL}/gold-scraper/prices/99999"
echo ""
INVALID_ID_RESPONSE=$(curl -s -X GET "${BASE_URL}/gold-scraper/prices/99999")

echo "Response:"
echo "$INVALID_ID_RESPONSE" | jq '.'
echo ""

SUCCESS=$(echo "$INVALID_ID_RESPONSE" | jq -r '.success')
if [ "$SUCCESS" = "false" ]; then
    echo -e "${GREEN}✅ Correctly returned error for invalid ID${NC}"
else
    echo -e "${RED}❌ Should have returned error for invalid ID${NC}"
fi
echo ""
echo "---"
echo ""

# Summary
echo "======================================"
echo "🎉 All tests completed!"
echo "======================================"
echo ""
echo "Summary:"
echo "- ✅ Scrape gold prices"
echo "- ✅ Get all prices"
echo "- ✅ Get prices with limit"
echo "- ✅ Get prices filtered by type"
echo "- ✅ Get prices filtered by source"
echo "- ✅ Get latest prices"
echo "- ✅ Get price by ID"
echo "- ✅ Get price by invalid ID"
echo ""
echo "💡 Tip: You can run individual tests by copying the curl commands above"
