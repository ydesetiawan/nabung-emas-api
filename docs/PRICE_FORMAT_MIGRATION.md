# Price Format Migration - Best Practices Implementation

## Overview
This document describes the migration from string-based price storage to integer-based (BIGINT) storage for `buy_price` and `sell_price` fields in the gold pricing system.

## Why This Change?

### Problems with String Storage
1. **No Type Safety**: Strings can contain invalid data
2. **Calculation Issues**: Cannot perform database-level calculations or aggregations
3. **Sorting Problems**: String sorting doesn't work correctly for numbers
4. **Storage Inefficiency**: Strings take more space than integers
5. **Precision Issues**: Inconsistent formatting can lead to data integrity problems

### Benefits of BIGINT Storage
1. **Type Safety**: Database enforces numeric values only
2. **Calculations**: Enable SUM, AVG, MIN, MAX operations at database level
3. **Proper Sorting**: Numeric sorting works correctly
4. **Storage Efficiency**: BIGINT uses 8 bytes vs VARCHAR(50) which can use up to 50+ bytes
5. **Precision**: No floating-point precision issues (storing as smallest unit - Rupiah)
6. **Indexing**: Better index performance for numeric types

## Changes Made

### 1. Database Schema (`migrations/003_create_gold_pricing_histories.sql`)
**Before:**
```sql
buy_price VARCHAR(50) NOT NULL,
sell_price VARCHAR(50) NOT NULL,
```

**After:**
```sql
buy_price BIGINT NOT NULL,
sell_price BIGINT NOT NULL,
```

### 2. Go Model (`internal/models/gold_pricing_history.go`)
**Before:**
```go
BuyPrice    string     `json:"buy_price" db:"buy_price"`
SellPrice   string     `json:"sell_price" db:"sell_price"`
```

**After:**
```go
BuyPrice    int64      `json:"buy_price" db:"buy_price"`
SellPrice   int64      `json:"sell_price" db:"sell_price"`
```

### 3. Repository (`internal/repositories/gold_pricing_history_repository.go`)
**Before:**
```go
func calculateBuyPrice(sellPrice string) string {
    // Complex string parsing and formatting
    // ~30 lines of code
}
```

**After:**
```go
func calculateBuyPrice(sellPrice int64) int64 {
    // Simple calculation
    return int64(float64(sellPrice) * 0.94)
}
```

### 4. Scraper Service (`internal/services/gold_scraper_service.go`)
Added `parsePrice` function to convert scraped string prices to int64:
```go
func parsePrice(priceStr string) (int64, error) {
    // Handles formats like "Rp1.234.567" or "1234567"
    cleaned := cleanPrice(priceStr)
    var result int64
    _, err := fmt.Sscanf(cleaned, "%d", &result)
    return result, err
}
```

## Data Format

### Storage Format
- **Type**: BIGINT (64-bit integer)
- **Unit**: Rupiah (smallest currency unit)
- **Example**: 1,234,567 Rp is stored as `1234567`

### JSON API Response
The API will return prices as integers:
```json
{
  "id": 1,
  "gold_type": "Logam Mulia 1 gram",
  "buy_price": 1164580,
  "sell_price": 1238915,
  "source": "antam"
}
```

### Frontend Display
Frontend applications should format the integer values for display:
```javascript
// JavaScript example
const formatRupiah = (amount) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount);
};

// Usage
formatRupiah(1238915); // Output: "Rp1.238.915"
```

## Migration Steps

### For New Installations
1. Run migrations in order:
   ```bash
   make migrate
   ```

### For Existing Databases
1. **Backup your database first!**
   ```bash
   pg_dump nabung_emas > backup_$(date +%Y%m%d_%H%M%S).sql
   ```

2. Run the migration script:
   ```bash
   psql -d nabung_emas -f migrations/004_alter_price_columns_to_bigint.sql
   ```

3. Verify the migration:
   ```bash
   psql -d nabung_emas -c "SELECT id, gold_type, buy_price, sell_price FROM gold_pricing_histories LIMIT 5;"
   ```

### Full Reset (Development Only)
If you want to start fresh:
```bash
make db-reset
```

## Testing

### 1. Test Scraping
```bash
curl -X POST http://localhost:8080/api/scrape
```

### 2. Test Price Retrieval
```bash
curl http://localhost:8080/api/prices/latest
```

### 3. Verify Data Types
```bash
psql -d nabung_emas -c "\d+ gold_pricing_histories"
```

Expected output should show:
```
Column      | Type   | ...
------------+--------+-----
buy_price   | bigint | ...
sell_price  | bigint | ...
```

## Rollback Plan

If you need to rollback (not recommended after data is in production):

```sql
-- Create rollback migration
ALTER TABLE gold_pricing_histories 
ADD COLUMN buy_price_old VARCHAR(50),
ADD COLUMN sell_price_old VARCHAR(50);

UPDATE gold_pricing_histories
SET 
    sell_price_old = 'Rp' || sell_price::TEXT,
    buy_price_old = 'Rp' || buy_price::TEXT;

ALTER TABLE gold_pricing_histories 
DROP COLUMN buy_price,
DROP COLUMN sell_price;

ALTER TABLE gold_pricing_histories 
RENAME COLUMN buy_price_old TO buy_price,
RENAME COLUMN sell_price_old TO sell_price;
```

## Best Practices Going Forward

1. **Always store monetary values as integers** (smallest currency unit)
2. **Format for display only in the presentation layer** (frontend)
3. **Use database-level calculations** when possible (SUM, AVG, etc.)
4. **Validate input** before converting to integer
5. **Document the unit** (Rupiah, cents, etc.) in comments and documentation

## References

- [PostgreSQL Numeric Types](https://www.postgresql.org/docs/current/datatype-numeric.html)
- [Stripe: Storing Money in Databases](https://stripe.com/docs/currencies#zero-decimal)
- [Martin Fowler: Money Pattern](https://martinfowler.com/eaaCatalog/money.html)
