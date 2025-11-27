#!/bin/bash

# Galeri24 Gold Scraper Test Script
# Tests all API endpoints

BASE_URL="http://localhost:8080/api/v1/galeri24-scraper"

echo "🧪 Testing Galeri24 Gold Scraper API"
echo "===================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test 1: Scrape Gold Prices
echo -e "${BLUE}Test 1: Scrape Gold Prices${NC}"
echo "POST $BASE_URL/scrape"
echo "---"
curl -X POST "$BASE_URL/scrape" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'
echo ""
echo ""

# Wait a bit for scraping to complete
sleep 2

# Test 2: Get All Prices (Limited)
echo -e "${BLUE}Test 2: Get All Prices (Limited to 10)${NC}"
echo "GET $BASE_URL/prices?limit=10"
echo "---"
curl "$BASE_URL/prices?limit=10" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.success, .message, .count'
echo ""
echo ""

# Test 3: Get Latest Prices
echo -e "${BLUE}Test 3: Get Latest Prices${NC}"
echo "GET $BASE_URL/prices/latest"
echo "---"
curl "$BASE_URL/prices/latest" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.success, .message, .count'
echo ""
echo ""

# Test 4: Filter by Gold Type
echo -e "${BLUE}Test 4: Filter by Gold Type (1 gram)${NC}"
echo "GET $BASE_URL/prices?type=1&limit=5"
echo "---"
curl "$BASE_URL/prices?type=1&limit=5" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.success, .message, .count, .data[0]'
echo ""
echo ""

# Test 5: Filter by Source
echo -e "${BLUE}Test 5: Filter by Source (GALERI_24)${NC}"
echo "GET $BASE_URL/prices?source=GALERI_24&limit=5"
echo "---"
curl "$BASE_URL/prices?source=GALERI_24&limit=5" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.success, .message, .count, .data[0]'
echo ""
echo ""

# Test 6: Get Price by ID
echo -e "${BLUE}Test 6: Get Price by ID (ID: 1)${NC}"
echo "GET $BASE_URL/prices/1"
echo "---"
curl "$BASE_URL/prices/1" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'
echo ""
echo ""

# Test 7: Get Prices by Date
TODAY=$(date +%Y-%m-%d)
echo -e "${BLUE}Test 7: Get Prices by Date ($TODAY)${NC}"
echo "GET $BASE_URL/prices/date/$TODAY"
echo "---"
curl "$BASE_URL/prices/date/$TODAY" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.success, .message, .count, .date'
echo ""
echo ""

# Test 8: Get Statistics
echo -e "${BLUE}Test 8: Get Statistics${NC}"
echo "GET $BASE_URL/stats"
echo "---"
curl "$BASE_URL/stats" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'
echo ""
echo ""

# Test 9: Date Range Filter
echo -e "${BLUE}Test 9: Date Range Filter (November 2025)${NC}"
echo "GET $BASE_URL/prices?start_date=2025-11-01&end_date=2025-11-30&limit=10"
echo "---"
curl "$BASE_URL/prices?start_date=2025-11-01&end_date=2025-11-30&limit=10" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.success, .message, .count'
echo ""
echo ""

# Test 10: Test Duplicate Prevention (Scrape Again)
echo -e "${YELLOW}Test 10: Test Duplicate Prevention (Scrape Again)${NC}"
echo "POST $BASE_URL/scrape"
echo "---"
echo "This should UPDATE existing records, not create duplicates"
curl -X POST "$BASE_URL/scrape" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.success, .message, .total_scraped, .saved_count, .updated_count'
echo ""
echo ""

echo -e "${GREEN}✅ All tests completed!${NC}"
echo ""
echo "Summary:"
echo "- Tested scraping functionality"
echo "- Tested all query endpoints"
echo "- Tested filtering (type, source, date range)"
echo "- Tested duplicate prevention (UPSERT)"
echo ""
