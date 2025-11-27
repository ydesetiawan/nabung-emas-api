# Galeri24 Gold Scraper - Quick Start Guide

## 🚀 Quick Start

### 1. Run Migrations
```bash
make migrate
```

### 2. Start Server
```bash
make run
# Or
./bin/nabung-emas-api
```

### 3. Test Scraper
```bash
./test-galeri24-scraper.sh
```

---

## 📋 API Endpoints

### Base URL
```
http://localhost:8080/api/v1/galeri24-scraper
```

### 1. Scrape Gold Prices
```bash
curl -X POST http://localhost:8080/api/v1/galeri24-scraper/scrape
```

**Response:**
```json
{
  "success": true,
  "message": "Successfully scraped and saved 170 records (85 new, 85 updated)",
  "pricing_date": "2025-11-27T00:00:00Z",
  "total_scraped": 170,
  "saved_count": 85,
  "updated_count": 85,
  "failed_count": 0,
  "duration": "5.234s"
}
```

### 2. Get All Prices
```bash
curl "http://localhost:8080/api/v1/galeri24-scraper/prices?limit=50"
```

### 3. Get Latest Prices
```bash
curl http://localhost:8080/api/v1/galeri24-scraper/prices/latest
```

### 4. Filter by Gold Type
```bash
curl "http://localhost:8080/api/v1/galeri24-scraper/prices?type=1&limit=20"
```

### 5. Filter by Source/Vendor
```bash
curl "http://localhost:8080/api/v1/galeri24-scraper/prices?source=GALERI_24&limit=20"
```

### 6. Date Range Filter
```bash
curl "http://localhost:8080/api/v1/galeri24-scraper/prices?start_date=2025-11-01&end_date=2025-11-30"
```

### 7. Get Prices by Date
```bash
curl http://localhost:8080/api/v1/galeri24-scraper/prices/date/2025-11-27
```

### 8. Get Statistics
```bash
curl http://localhost:8080/api/v1/galeri24-scraper/stats
```

---

## 🔒 Duplicate Prevention

The scraper uses **UPSERT logic** to prevent duplicates:

- **Database Constraint**: `UNIQUE(pricing_date, gold_type, source)`
- **Repository Logic**: `ON CONFLICT DO UPDATE`
- **Result**: If you scrape the same date multiple times, existing records are **UPDATED**, not duplicated

**Example:**
```bash
# First scrape
curl -X POST http://localhost:8080/api/v1/galeri24-scraper/scrape
# Result: 170 new records

# Second scrape (same day)
curl -X POST http://localhost:8080/api/v1/galeri24-scraper/scrape
# Result: 0 new, 170 updated (NO DUPLICATES!)
```

---

## 🗄️ Supported Vendors

17 gold vendors are supported:

1. GALERI_24
2. DINAR_G24
3. BABY_GALERI_24
4. ANTAM
5. UBS
6. ANTAM_MULIA_RETRO
7. ANTAM_NON_PEGADAIAN
8. LOTUS_ARCHI
9. UBS_DISNEY
10. UBS_ELSA
11. UBS_ANNA
12. UBS_MICKEY_FULLBODY
13. LOTUS_ARCHI_GIFT
14. UBS_HELLO_KITTY
15. BABY_SERIES_TUMBUHAN
16. BABY_SERIES_INVESTASI
17. BATIK_SERIES

---

## 📊 Query Parameters

### GET /prices

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `type` | string | Filter by gold type (partial match) | `type=1` |
| `source` | string | Filter by vendor | `source=GALERI_24` |
| `start_date` | date | Start date (YYYY-MM-DD) | `start_date=2025-11-01` |
| `end_date` | date | End date (YYYY-MM-DD) | `end_date=2025-11-30` |
| `limit` | int | Limit results | `limit=50` |
| `offset` | int | Offset for pagination | `offset=100` |

---

## 🧪 Testing

### Run All Tests
```bash
./test-galeri24-scraper.sh
```

### Manual Tests

1. **Scrape and verify no duplicates:**
```bash
# Scrape once
curl -X POST http://localhost:8080/api/v1/galeri24-scraper/scrape

# Scrape again - should update, not duplicate
curl -X POST http://localhost:8080/api/v1/galeri24-scraper/scrape
```

2. **Check database directly:**
```bash
psql -d nabung_emas -c "SELECT COUNT(*) FROM gold_pricing_histories;"
psql -d nabung_emas -c "SELECT pricing_date, gold_type, source, COUNT(*) FROM gold_pricing_histories GROUP BY pricing_date, gold_type, source HAVING COUNT(*) > 1;"
# Should return 0 rows (no duplicates)
```

---

## 📝 Response Format

### Success Response
```json
{
  "success": true,
  "message": "Successfully retrieved gold prices",
  "count": 170,
  "data": [...]
}
```

### Error Response
```json
{
  "success": false,
  "message": "Failed to scrape gold prices",
  "errors": ["timeout exceeded"]
}
```

---

## 🔧 Configuration

### Scraper Settings
- **Timeout**: 30 seconds
- **Rate Limit**: 1s delay + 500ms random
- **User-Agent**: Modern Chrome
- **Target**: https://galeri24.co.id/harga-emas

### Database Settings
- **Table**: `gold_pricing_histories`
- **Unique Constraint**: `(pricing_date, gold_type, source)`
- **Indexes**: 7 optimized indexes

---

## 📚 Documentation

- **Full Summary**: `GALERI24_SCRAPER_SUMMARY.md`
- **Test Script**: `test-galeri24-scraper.sh`
- **Postman Collection**: `EmasGo-API.postman_collection.json`

---

## ⚡ Performance

- **Scraping Duration**: ~5-10 seconds
- **Records per Scrape**: ~170 (10 weights × 17 vendors)
- **Database Operation**: Single transaction (batch insert/update)
- **Duplicate Check**: Automatic via UNIQUE constraint

---

## ✅ Key Features

- ✅ **No Duplicates**: UPSERT logic prevents duplicate records
- ✅ **Indonesian Date Support**: Parses "Diperbarui Kamis, 27 November 2025"
- ✅ **Multi-Vendor**: Supports 17 different gold vendors
- ✅ **Comprehensive Filtering**: Type, source, date range, pagination
- ✅ **Error Handling**: Timeout protection, retry mechanism
- ✅ **Logging**: Detailed console logging with emojis

---

## 🎯 Use Cases

### Daily Price Updates
```bash
# Run daily via cron
0 9 * * * curl -X POST http://localhost:8080/api/v1/galeri24-scraper/scrape
```

### Price Comparison
```bash
# Compare prices across vendors for 1 gram gold
curl "http://localhost:8080/api/v1/galeri24-scraper/prices?type=1&limit=100"
```

### Historical Analysis
```bash
# Get all prices for November 2025
curl "http://localhost:8080/api/v1/galeri24-scraper/prices?start_date=2025-11-01&end_date=2025-11-30"
```

---

**Status**: ✅ Production Ready  
**Version**: 1.0.0  
**Last Updated**: 2025-11-27
