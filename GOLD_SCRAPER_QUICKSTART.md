# Gold Scraper API - Quick Start Guide

This guide will help you quickly get started with the Gold Scraper API.

## 🚀 Quick Start

### 1. Prerequisites

Make sure you have:
- ✅ Go 1.21 or higher installed
- ✅ PostgreSQL database running
- ✅ Database created (`nabung_emas`)
- ✅ Environment variables configured (`.env` file)

### 2. Install Dependencies

```bash
go get github.com/gocolly/colly/v2
go mod tidy
```

### 3. Run Migrations

```bash
make migrate
```

This will create the `gold_pricing_histories` table with proper indexes.

### 4. Start the Server

```bash
make run
```

The server will start on `http://localhost:8080`

## 📝 Basic Usage

### Scrape Gold Prices

```bash
curl -X POST http://localhost:8080/api/v1/gold-scraper/scrape
```

**Expected Output:**
```json
{
  "success": true,
  "message": "Successfully scraped and saved 10 gold prices",
  "count": 10,
  "data": [...]
}
```

### View All Prices

```bash
curl http://localhost:8080/api/v1/gold-scraper/prices
```

### View Latest Prices

```bash
curl http://localhost:8080/api/v1/gold-scraper/prices/latest
```

### View Specific Price

```bash
curl http://localhost:8080/api/v1/gold-scraper/prices/1
```

## 🧪 Testing

Run the comprehensive test script:

```bash
./test-gold-scraper.sh
```

This will test all endpoints and display results with colored output.

## 📊 Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/gold-scraper/scrape` | Scrape and save gold prices |
| GET | `/api/v1/gold-scraper/prices` | Get all prices (with filters) |
| GET | `/api/v1/gold-scraper/prices/latest` | Get latest prices |
| GET | `/api/v1/gold-scraper/prices/:id` | Get price by ID |

## 🔍 Query Parameters

### `/prices` endpoint supports:

- `type` - Filter by gold type (e.g., `?type=emas`)
- `source` - Filter by source (e.g., `?source=antam`)
- `limit` - Limit results (e.g., `?limit=10`)

**Examples:**

```bash
# Get first 5 prices
curl "http://localhost:8080/api/v1/gold-scraper/prices?limit=5"

# Filter by type
curl "http://localhost:8080/api/v1/gold-scraper/prices?type=emas"

# Filter by source
curl "http://localhost:8080/api/v1/gold-scraper/prices?source=antam"

# Combine filters
curl "http://localhost:8080/api/v1/gold-scraper/prices?type=emas&source=antam&limit=10"
```

## 📦 Database Schema

The `gold_pricing_histories` table structure:

```sql
CREATE TABLE gold_pricing_histories (
    id SERIAL PRIMARY KEY,
    gold_type VARCHAR(255) NOT NULL,
    buy_price VARCHAR(50) NOT NULL,
    sell_price VARCHAR(50) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    source gold_source NOT NULL,  -- ENUM: 'antam' or 'usb'
    scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🎯 Common Use Cases

### 1. Daily Price Scraping

Set up a cron job to scrape prices daily:

```bash
# Add to crontab (scrape at 9 AM daily)
0 9 * * * curl -X POST http://localhost:8080/api/v1/gold-scraper/scrape
```

### 2. Get Current Market Prices

```bash
curl http://localhost:8080/api/v1/gold-scraper/prices/latest | jq '.data'
```

### 3. Price History Analysis

```bash
# Get all prices for a specific gold type
curl "http://localhost:8080/api/v1/gold-scraper/prices?type=1%20gram" | jq '.data'
```

### 4. Compare Prices from Different Sources

```bash
# Get Antam prices
curl "http://localhost:8080/api/v1/gold-scraper/prices?source=antam" | jq '.data'

# Get USB prices
curl "http://localhost:8080/api/v1/gold-scraper/prices?source=usb" | jq '.data'
```

## 🐛 Troubleshooting

### Issue: Scraping returns no data

**Solution:**
- Check if the website structure has changed
- Verify network connectivity
- Check server logs for detailed error messages

### Issue: Database connection error

**Solution:**
- Verify PostgreSQL is running
- Check `.env` file for correct `DATABASE_URL`
- Ensure migrations have been run

### Issue: Timeout errors

**Solution:**
- The scraper has a 30-second timeout
- Check your internet connection
- The website might be temporarily unavailable

## 📚 Additional Resources

- [Full API Documentation](./GOLD_SCRAPER_API.md)
- [Main Project README](./README.md)
- [Implementation Guide](./IMPLEMENTATION_GUIDE.md)

## 🎨 Console Output

The API uses emoji logging for better visibility:

```
🕷️  Starting gold price scraping from logammulia.com...
🌐 Visiting: https://logammulia.com/id/harga-emas-hari-ini
📊 Found table, parsing gold prices...
✅ Scraped: Emas Antam 1 gram - Buy: 1150000, Sell: 1050000, Unit: gram
💾 Successfully saved 10 records to database
```

## 💡 Tips

1. **Use jq for JSON formatting**: Install `jq` to prettify JSON responses
   ```bash
   curl http://localhost:8080/api/v1/gold-scraper/prices/latest | jq '.'
   ```

2. **Save responses to file**:
   ```bash
   curl http://localhost:8080/api/v1/gold-scraper/prices > prices.json
   ```

3. **Check response status**:
   ```bash
   curl -i http://localhost:8080/api/v1/gold-scraper/prices/latest
   ```

4. **Verbose output for debugging**:
   ```bash
   curl -v http://localhost:8080/api/v1/gold-scraper/scrape
   ```

## 🔐 Security Considerations

- The scraper endpoints are currently **public** (no authentication required)
- Consider adding authentication middleware for production use
- Implement rate limiting to prevent abuse
- Use HTTPS in production

## 🚦 Next Steps

1. ✅ Test all endpoints using the test script
2. ✅ Integrate with your frontend application
3. ✅ Set up scheduled scraping
4. ✅ Monitor database growth and implement cleanup
5. ✅ Add authentication if needed

## 📞 Support

For detailed documentation, see [GOLD_SCRAPER_API.md](./GOLD_SCRAPER_API.md)

---

**Happy Scraping! 🎉**
