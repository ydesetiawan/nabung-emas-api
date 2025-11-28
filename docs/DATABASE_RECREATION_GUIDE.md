# Database Recreation Guide

## Summary of Changes

### Removed Fields
- ✅ **scraped_at** - Removed from `gold_pricing_histories` table
- ✅ **unit** - Removed from scraper service (was never in database)

### Current Schema
The `gold_pricing_histories` table now has:
- `id` - SERIAL PRIMARY KEY
- `pricing_date` - DATE NOT NULL
- `gold_type` - VARCHAR(255) NOT NULL
- `buy_price` - BIGINT NOT NULL (stores Rupiah as integer)
- `sell_price` - BIGINT NOT NULL (stores Rupiah as integer)
- `source` - gold_source ENUM ('antam', 'usb')
- `created_at` - TIMESTAMP (auto-set on insert)
- `updated_at` - TIMESTAMP (auto-updated via trigger)

## Steps to Recreate Database

### Option 1: Full Reset (Recommended for Development)
```bash
# This will drop the database, recreate it, and run all migrations
make db-reset
```

### Option 2: Manual Steps
```bash
# 1. Drop existing database
dropdb nabung_emas

# 2. Create new database
createdb nabung_emas

# 3. Run migrations
make migrate
```

### Option 3: Individual Commands
```bash
# Drop database
psql -c "DROP DATABASE IF EXISTS nabung_emas;"

# Create database
psql -c "CREATE DATABASE nabung_emas;"

# Run migrations one by one
psql -d nabung_emas -f migrations/001_initial_schema.sql
psql -d nabung_emas -f migrations/002_add_token_blacklist.sql
psql -d nabung_emas -f migrations/003_create_gold_pricing_histories.sql
```

## Verification

After recreating the database, verify the schema:

```bash
# Check table structure
psql -d nabung_emas -c "\d+ gold_pricing_histories"
```

Expected output should show:
```
Column       | Type                     | Nullable | Default
-------------+--------------------------+----------+----------------------------
id           | integer                  | not null | nextval('gold_pricing...')
pricing_date | date                     | not null |
gold_type    | character varying(255)   | not null |
buy_price    | bigint                   | not null |
sell_price   | bigint                   | not null |
source       | gold_source              | not null |
created_at   | timestamp                |          | CURRENT_TIMESTAMP
updated_at   | timestamp                |          | CURRENT_TIMESTAMP
```

## Test the Application

After database recreation:

```bash
# 1. Start the server
make run

# 2. Test scraping (in another terminal)
curl -X POST http://localhost:8080/api/scrape

# 3. Check latest prices
curl http://localhost:8080/api/prices/latest
```

## Migration Files

Current migration files:
1. `001_initial_schema.sql` - User authentication tables
2. `002_add_token_blacklist.sql` - Token blacklist for logout
3. `003_create_gold_pricing_histories.sql` - Gold pricing table (updated)

**Note:** Migration `004_alter_price_columns_to_bigint.sql` has been deleted as it's no longer needed since we're recreating the database with the correct schema from the start.

## Changes Made to Code

### Files Modified:
1. ✅ `migrations/003_create_gold_pricing_histories.sql`
   - Removed `scraped_at` column
   - Removed `scraped_at` index
   - Removed `scraped_at` comment

2. ✅ `migrations/004_alter_price_columns_to_bigint.sql`
   - **DELETED** (no longer needed)

3. ✅ `internal/models/gold_pricing_history.go`
   - Removed `ScrapedAt` field from struct
   - Changed `BuyPrice` and `SellPrice` to `int64`

4. ✅ `internal/repositories/gold_pricing_history_repository.go`
   - Removed `scraped_at` from all SQL queries
   - Removed `&history.ScrapedAt` from all Scan operations
   - Removed `scraped_at = CURRENT_TIMESTAMP` from UPDATE queries

5. ✅ `internal/services/gold_scraper_service.go`
   - Removed `Unit` field from `ScrapedGoldData` struct
   - Removed all unit scraping and logging

## API Response Format

Prices are now returned as integers (Rupiah):
```json
{
  "id": 1,
  "pricing_date": "2025-11-28",
  "gold_type": "Logam Mulia 1 gram",
  "buy_price": 1164580,
  "sell_price": 1238915,
  "source": "antam",
  "created_at": "2025-11-28T10:00:00Z",
  "updated_at": "2025-11-28T10:00:00Z"
}
```

**Note:** `scraped_at` field is no longer present in the response.
