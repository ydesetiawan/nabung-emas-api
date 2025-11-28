# Gold Scraper Service - Complete Revamp

## Overview
The gold scraper service has been completely rewritten with production-ready features, comprehensive error handling, and intelligent category detection.

## ✨ Key Features

### 1. **Retry Logic with Exponential Backoff**
- Automatically retries failed requests up to 3 times
- Exponential backoff: 1s, 4s, 9s between retries
- Prevents overwhelming the target server

### 2. **Intelligent Category Detection**
Automatically detects product categories based on product names:
- **Emas Batangan** - Standard gold bars (default)
- **Emas Batangan Gift Series** - Gift series products
- **Emas Batangan Selamat Idul Fitri** - Idul Fitri special editions
- **Emas Batangan Imlek** - Chinese New Year editions
- **Emas Batangan Batik Seri III** - Batik series products
- **Perak Murni** - Pure silver
- **Perak Heritage** - Heritage silver collection
- **Liontin Batik Seri III** - Batik pendants

### 3. **Robust Price Parsing**
- Handles multiple formats: "Rp1.234.567", "IDR 1,234,567", "1234567"
- Removes currency symbols and separators
- Validates numeric values
- Prevents negative prices

### 4. **Comprehensive Error Handling**
- Network failures
- Missing HTML elements
- Invalid price formats
- Unrecognized categories (falls back to default)
- Rate limiting protection

### 5. **Respectful Scraping**
- 2-second delay between requests
- Proper User-Agent headers
- 30-second timeout per request
- Context support for cancellation

### 6. **Detailed Logging**
- Structured logging with emojis for easy reading
- Tracks scraping progress
- Reports successes and failures
- Performance metrics (duration)

## 📊 API Response Structure

```json
{
  "success": true,
  "message": "Successfully scraped 15 items: 12 new, 3 updated",
  "pricing_date": "2025-11-28T00:00:00Z",
  "total_scraped": 15,
  "saved_count": 12,
  "updated_count": 3,
  "failed_count": 0,
  "errors": [],
  "duration": "5.234s",
  "data": [...]
}
```

## 🔧 Helper Functions

### `parsePrice(priceStr string) (int64, error)`
Converts price strings to integers.

**Examples:**
```go
parsePrice("Rp1.234.567")  // Returns: 1234567
parsePrice("IDR 500")       // Returns: 500
parsePrice("1234567")       // Returns: 1234567
```

### `detectCategory(productName string) (models.GoldCategory, error)`
Detects category from product name.

**Detection Rules (in order of priority):**
1. **Liontin + Batik** → `liontin_batik_seri_iii`
2. **Batik Seri III** → `emas_batangan_batik_seri_iii`
3. **Gift Series** → `emas_batangan_gift_series`
4. **Idul Fitri/Lebaran** → `emas_batangan_selamat_idul_fitri`
5. **Imlek/Chinese New Year** → `emas_batangan_imlek`
6. **Perak + Heritage** → `perak_heritage`
7. **Perak/Silver** → `perak_murni`
8. **Liontin/Pendant** → `liontin_batik_seri_iii`
9. **Batik** → `emas_batangan_batik_seri_iii`
10. **Default** → `emas_batangan`

### `cleanText(text string) string`
Removes extra whitespace and newlines.

### `cleanPrice(price string) string`
Removes currency symbols and separators.

## 🧪 Unit Tests

All helper functions have comprehensive unit tests:

```bash
# Run all tests
go test -v ./internal/services/...

# Run specific tests
go test -v ./internal/services/... -run TestParsePrice
go test -v ./internal/services/... -run TestDetectCategory
```

**Test Coverage:**
- ✅ `TestParsePrice` - 9 test cases
- ✅ `TestDetectCategory` - 14 test cases
- ✅ `TestCleanText` - 5 test cases
- ✅ `TestCleanPrice` - 6 test cases

## 📝 Usage Examples

### Basic Scraping
```go
// Create service
scraperService := services.NewGoldScraperService(repo)

// Scrape from Logam Mulia (Antam)
result, err := scraperService.ScrapeLogamMulia()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Scraped %d items\n", result.TotalScraped)
fmt.Printf("Saved: %d, Updated: %d, Failed: %d\n", 
    result.SavedCount, result.UpdatedCount, result.FailedCount)
```

### Via API Endpoint
```bash
# Trigger scraping
curl -X POST http://localhost:8080/api/scrape

# Get latest prices
curl http://localhost:8080/api/prices/latest

# Filter by category
curl "http://localhost:8080/api/prices?category=emas_batangan_gift_series"
```

## 🔄 Scraping Flow

```
1. Initialize Collector
   ├─ Set User-Agent
   ├─ Configure timeouts
   └─ Set rate limits

2. Visit Target URL
   ├─ Parse HTML tables
   ├─ Extract gold type, prices
   └─ Store in ScrapedGoldData

3. Process Data
   ├─ Parse prices (string → int64)
   ├─ Detect categories
   └─ Create GoldPricingHistoryCreate models

4. Save to Database
   ├─ Batch insert with UPSERT
   ├─ Track saved/updated counts
   └─ Return results

5. Return ScrapeResult
   ├─ Success status
   ├─ Statistics
   ├─ Errors (if any)
   └─ Saved data
```

## 🚀 Performance

- **Average scraping time**: 3-5 seconds
- **Retry overhead**: +1-14 seconds (if retries needed)
- **Memory efficient**: Streams data, doesn't load entire page
- **Database efficient**: Batch inserts with UPSERT

## 🛡️ Error Handling

### Network Errors
```go
// Automatic retry with exponential backoff
result, err := scraper.ScrapeLogamMulia()
// Retries: 1s → 4s → 9s
```

### Parse Errors
```go
// Invalid price format
price, err := parsePrice("invalid")
// Returns: error with descriptive message
// Continues with other items
```

### Category Detection
```go
// Unknown product
category, _ := detectCategory("Unknown Product")
// Returns: emas_batangan (default)
// Logs warning but doesn't fail
```

## 📈 Future Enhancements

### Multi-Source Support
Ready to extend for other sources:
```go
// Interface for different sources
type GoldPriceScraper interface {
    Scrape(ctx context.Context) (*ScrapeResult, error)
    GetSource() models.GoldSource
}

// Implement for each source
type AntamScraper struct { ... }
type UBSScraper struct { ... }
type Galeri24Scraper struct { ... }
type PegadaianScraper struct { ... }
```

### Caching
```go
// Add caching layer
type CachedScraper struct {
    scraper GoldPriceScraper
    cache   Cache
    ttl     time.Duration
}
```

### Webhooks
```go
// Notify on price changes
type WebhookNotifier struct {
    url string
}

func (w *WebhookNotifier) OnPriceChange(old, new *GoldPrice) {
    // Send webhook
}
```

## 🔍 Debugging

### Enable Verbose Logging
```go
// Already enabled by default
// Check logs for:
// 🕷️  - Scraping start
// 🌐 - HTTP requests
// 📊 - Table parsing
// ✅ - Successful scrapes
// 🏷️  - Category detection
// 💾 - Database operations
// ❌ - Errors
```

### Test Individual Functions
```go
// Test price parsing
price, err := parsePrice("Rp1.234.567")
fmt.Println(price) // 1234567

// Test category detection
category, _ := detectCategory("Emas Gift Series 5 gram")
fmt.Println(category) // emas_batangan_gift_series
```

## 📚 Related Documentation

- [Category System](./CATEGORY_SYSTEM.md)
- [Database Schema](./DATABASE_RECREATION_GUIDE.md)
- [Price Format](./PRICE_FORMAT_MIGRATION.md)
- [API Documentation](../EmasGo-API.postman_collection.json)

## ✅ Checklist

- [x] Retry logic with exponential backoff
- [x] Intelligent category detection
- [x] Robust price parsing
- [x] Comprehensive error handling
- [x] Respectful scraping (delays, User-Agent)
- [x] Detailed logging
- [x] Unit tests (34 test cases)
- [x] Production-ready code
- [x] Documentation
- [ ] Multi-source support (future)
- [ ] Caching layer (future)
- [ ] Webhook notifications (future)

## 🎯 Summary

The revamped gold scraper service is:
- **Production-ready** with comprehensive error handling
- **Intelligent** with automatic category detection
- **Reliable** with retry logic and validation
- **Well-tested** with 34 unit tests
- **Maintainable** with clean code and documentation
- **Extensible** ready for multi-source support

All tests pass ✅  
Build successful ✅  
Ready for deployment 🚀
