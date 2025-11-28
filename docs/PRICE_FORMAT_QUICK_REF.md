# Quick Reference: Price Format Changes

## Summary
Changed `buy_price` and `sell_price` from `VARCHAR(50)` to `BIGINT` for better data integrity and performance.

## What Changed

| Component | Before | After |
|-----------|--------|-------|
| **Database Type** | `VARCHAR(50)` | `BIGINT` |
| **Go Type** | `string` | `int64` |
| **Storage Format** | `"Rp1.234.567"` | `1234567` |
| **JSON Response** | `"1234567"` | `1234567` |

## Migration Commands

### New Installation
```bash
make migrate
```

### Existing Database
```bash
# Backup first!
pg_dump nabung_emas > backup.sql

# Run migration
psql -d nabung_emas -f migrations/004_alter_price_columns_to_bigint.sql
```

### Development Reset
```bash
make db-reset
```

## API Response Example

**Before:**
```json
{
  "buy_price": "1164580",
  "sell_price": "1238915"
}
```

**After:**
```json
{
  "buy_price": 1164580,
  "sell_price": 1238915
}
```

## Frontend Integration

### JavaScript
```javascript
const formatRupiah = (amount) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount);
};

// Usage
formatRupiah(1238915); // "Rp1.238.915"
```

### Go (if needed for display)
```go
import "golang.org/x/text/language"
import "golang.org/x/text/message"

func FormatRupiah(amount int64) string {
    p := message.NewPrinter(language.Indonesian)
    return p.Sprintf("Rp%d", amount)
}
```

## Benefits

✅ **Type Safety** - Database enforces numeric values  
✅ **Calculations** - Can use SUM, AVG, MIN, MAX in SQL  
✅ **Sorting** - Proper numeric sorting  
✅ **Performance** - Better indexing and storage  
✅ **Precision** - No floating-point issues  

## Files Modified

1. `migrations/003_create_gold_pricing_histories.sql` - Schema definition
2. `migrations/004_alter_price_columns_to_bigint.sql` - Migration script (new)
3. `internal/models/gold_pricing_history.go` - Model types
4. `internal/repositories/gold_pricing_history_repository.go` - Calculation logic
5. `internal/services/gold_scraper_service.go` - Price parsing

## Verification

```bash
# Check column types
psql -d nabung_emas -c "\d+ gold_pricing_histories"

# View sample data
psql -d nabung_emas -c "SELECT id, gold_type, buy_price, sell_price FROM gold_pricing_histories LIMIT 5;"

# Test API
curl http://localhost:8080/api/prices/latest
```
