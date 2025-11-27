# Gold Scraper API Documentation

## Overview

The Gold Scraper API provides endpoints to scrape gold prices from [logammulia.com](https://logammulia.com/id/harga-emas-hari-ini) and store them in a PostgreSQL database. This API is built using the Echo framework and Colly web scraping library.

## Features

✅ **Web Scraping**: Scrapes gold prices from logammulia.com using Colly  
✅ **Database Storage**: Stores scraped data in PostgreSQL with proper indexing  
✅ **RESTful API**: Clean REST endpoints for accessing gold price data  
✅ **Error Handling**: Graceful error handling with timeouts and retries  
✅ **Filtering**: Query prices by type, source, and limit  
✅ **Latest Prices**: Get the most recent price for each gold type  
✅ **Emoji Logging**: Console logging with emojis for better visibility  

## Database Schema

### Table: `gold_pricing_histories`

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL PRIMARY KEY | Auto-incrementing ID |
| `gold_type` | VARCHAR(255) | Type of gold (e.g., "Emas Antam 1 gram") |
| `buy_price` | VARCHAR(50) | Buy price in IDR |
| `sell_price` | VARCHAR(50) | Sell price in IDR |
| `unit` | VARCHAR(50) | Unit of measurement |
| `source` | ENUM | Source of data ('antam' or 'usb') |
| `scraped_at` | TIMESTAMP | When the data was scraped |
| `created_at` | TIMESTAMP | When the record was created |

### Indexes

- `idx_gold_pricing_histories_scraped_at` - Index on scraped_at
- `idx_gold_pricing_histories_gold_type` - Index on gold_type
- `idx_gold_pricing_histories_source` - Index on source
- `idx_gold_pricing_histories_gold_type_source` - Composite index
- `idx_gold_pricing_histories_latest` - Composite index for latest queries

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1/gold-scraper
```

---

### 1. Scrape Gold Prices

**Endpoint**: `POST /scrape`

**Description**: Scrapes gold prices from logammulia.com and saves them to the database.

**Request**:
```bash
curl -X POST http://localhost:8080/api/v1/gold-scraper/scrape
```

**Response** (Success):
```json
{
  "success": true,
  "message": "Successfully scraped and saved 10 gold prices",
  "count": 10,
  "data": [
    {
      "id": 1,
      "gold_type": "Emas Antam 1 gram",
      "buy_price": "1150000",
      "sell_price": "1050000",
      "unit": "gram",
      "source": "antam",
      "scraped_at": "2025-11-27T14:21:50+07:00",
      "created_at": "2025-11-27T14:21:50+07:00"
    }
    // ... more items
  ]
}
```

**Response** (Error):
```json
{
  "success": false,
  "message": "Failed to scrape gold prices",
  "errors": ["Request failed: timeout"]
}
```

---

### 2. Get All Prices

**Endpoint**: `GET /prices`

**Description**: Retrieves all gold prices with optional filters.

**Query Parameters**:
- `type` (optional) - Filter by gold type (partial match, case-insensitive)
- `source` (optional) - Filter by source ('antam' or 'usb')
- `limit` (optional) - Limit number of results

**Examples**:

Get all prices:
```bash
curl http://localhost:8080/api/v1/gold-scraper/prices
```

Get prices with limit:
```bash
curl http://localhost:8080/api/v1/gold-scraper/prices?limit=10
```

Filter by type:
```bash
curl http://localhost:8080/api/v1/gold-scraper/prices?type=emas
```

Filter by source:
```bash
curl http://localhost:8080/api/v1/gold-scraper/prices?source=antam
```

Combine filters:
```bash
curl http://localhost:8080/api/v1/gold-scraper/prices?type=emas&source=antam&limit=5
```

**Response**:
```json
{
  "success": true,
  "message": "Successfully retrieved gold prices",
  "count": 10,
  "data": [
    {
      "id": 1,
      "gold_type": "Emas Antam 1 gram",
      "buy_price": "1150000",
      "sell_price": "1050000",
      "unit": "gram",
      "source": "antam",
      "scraped_at": "2025-11-27T14:21:50+07:00",
      "created_at": "2025-11-27T14:21:50+07:00"
    }
    // ... more items
  ]
}
```

---

### 3. Get Latest Prices

**Endpoint**: `GET /prices/latest`

**Description**: Retrieves the latest price for each gold type.

**Request**:
```bash
curl http://localhost:8080/api/v1/gold-scraper/prices/latest
```

**Response**:
```json
{
  "success": true,
  "message": "Successfully retrieved latest gold prices",
  "count": 5,
  "data": [
    {
      "id": 15,
      "gold_type": "Emas Antam 1 gram",
      "buy_price": "1150000",
      "sell_price": "1050000",
      "unit": "gram",
      "source": "antam",
      "scraped_at": "2025-11-27T14:21:50+07:00",
      "created_at": "2025-11-27T14:21:50+07:00"
    }
    // ... more items (one per gold type)
  ]
}
```

---

### 4. Get Price by ID

**Endpoint**: `GET /prices/:id`

**Description**: Retrieves a specific gold price by ID.

**Request**:
```bash
curl http://localhost:8080/api/v1/gold-scraper/prices/1
```

**Response** (Success):
```json
{
  "success": true,
  "message": "Successfully retrieved gold price",
  "count": 1,
  "data": {
    "id": 1,
    "gold_type": "Emas Antam 1 gram",
    "buy_price": "1150000",
    "sell_price": "1050000",
    "unit": "gram",
    "source": "antam",
    "scraped_at": "2025-11-27T14:21:50+07:00",
    "created_at": "2025-11-27T14:21:50+07:00"
  }
}
```

**Response** (Not Found):
```json
{
  "success": false,
  "message": "Gold price not found"
}
```

---

## Response Structure

All API responses follow this standard structure:

```json
{
  "success": boolean,      // Indicates if the request was successful
  "message": string,       // Human-readable message
  "count": integer,        // Number of items returned (optional)
  "data": object|array,    // Response data (optional)
  "errors": [string]       // Array of error messages (optional)
}
```

## Error Handling

The API handles various error scenarios:

- **Timeout Errors**: 30-second timeout for scraping requests
- **Network Errors**: Graceful handling of connection failures
- **Invalid Parameters**: Validation of query parameters
- **Database Errors**: Proper error messages for database issues
- **Not Found**: 404 responses for non-existent resources

## Scraping Details

### User Agent
The scraper uses a modern browser user agent to avoid being blocked:
```
Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36
```

### Rate Limiting
- 1-second delay between requests
- Single parallel request at a time
- 30-second timeout per request

### Data Cleaning
The scraper automatically:
- Removes extra whitespace and newlines
- Cleans currency symbols (Rp, IDR)
- Removes thousand separators (dots and commas)
- Normalizes text formatting

## Installation & Setup

### 1. Install Dependencies
```bash
go get github.com/gocolly/colly/v2
go mod tidy
```

### 2. Run Migrations
```bash
make migrate
```

### 3. Start the Server
```bash
make run
```

## Testing

### Run Test Script
```bash
./test-gold-scraper.sh
```

### Manual Testing

1. **Scrape prices**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/gold-scraper/scrape
   ```

2. **View all prices**:
   ```bash
   curl http://localhost:8080/api/v1/gold-scraper/prices
   ```

3. **View latest prices**:
   ```bash
   curl http://localhost:8080/api/v1/gold-scraper/prices/latest
   ```

## Logging

The API uses emoji logging for better visibility:

- 🕷️ Starting scraping
- 🌐 Visiting URL
- 📊 Parsing data
- ✅ Success
- ❌ Error
- ⚠️ Warning
- 💾 Saving to database
- 🚀 Starting operation
- 📋 Fetching data
- 🔍 Searching

## Architecture

```
┌─────────────────┐
│   Echo Router   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Handler Layer  │  (gold_scraper_handler.go)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Service Layer  │  (gold_scraper_service.go)
└────────┬────────┘
         │
         ├──────────────────┐
         ▼                  ▼
┌─────────────────┐  ┌─────────────────┐
│  Colly Scraper  │  │  Repository     │
│  (Web Scraping) │  │  (Database)     │
└─────────────────┘  └─────────────────┘
```

## Code Structure

```
internal/
├── models/
│   └── gold_pricing_history.go      # Data models
├── repositories/
│   └── gold_pricing_history_repository.go  # Database operations
├── services/
│   └── gold_scraper_service.go      # Business logic & scraping
├── handlers/
│   └── gold_scraper_handler.go      # HTTP handlers
└── routes/
    └── routes.go                     # Route configuration
```

## Future Enhancements

- [ ] Add support for multiple gold price sources
- [ ] Implement scheduled scraping (cron jobs)
- [ ] Add price change notifications
- [ ] Implement caching for frequently accessed data
- [ ] Add GraphQL support
- [ ] Create price comparison charts
- [ ] Add historical price analysis
- [ ] Implement WebSocket for real-time updates

## License

This API is part of the Nabung Emas (Gold Savings) project.

## Support

For issues or questions, please refer to the main project documentation.
