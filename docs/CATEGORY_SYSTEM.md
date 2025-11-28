# Category System Documentation

## Overview
The gold pricing system now supports categorization of products to better organize different types of gold and silver items.

## Categories

| Category Enum | Display Name | Description |
|---------------|--------------|-------------|
| `emas_batangan` | Emas Batangan | Standard gold bars |
| `emas_batangan_gift_series` | Emas Batangan Gift Series | Gift series gold bars |
| `emas_batangan_selamat_idul_fitri` | Emas Batangan Selamat Idul Fitri | Idul Fitri special edition gold bars |
| `emas_batangan_imlek` | Emas Batangan Imlek | Chinese New Year special edition gold bars |
| `emas_batangan_batik_seri_iii` | Emas Batangan Batik Seri III | Batik Series III gold bars |
| `perak_murni` | Perak Murni | Pure silver |
| `perak_heritage` | Perak Heritage | Heritage silver collection |
| `liontin_batik_seri_iii` | Liontin Batik Seri III | Batik Series III pendants |

## Database Schema

### Enum Type
```sql
CREATE TYPE gold_category AS ENUM (
    'emas_batangan',
    'emas_batangan_gift_series',
    'emas_batangan_selamat_idul_fitri',
    'emas_batangan_imlek',
    'emas_batangan_batik_seri_iii',
    'perak_murni',
    'perak_heritage',
    'liontin_batik_seri_iii'
);
```

### Table Column
```sql
ALTER TABLE gold_pricing_histories 
ADD COLUMN category gold_category DEFAULT 'emas_batangan';
```

## Go Model

### Type Definition
```go
type GoldCategory string

const (
    GoldCategoryEmasBatangan                 GoldCategory = "emas_batangan"
    GoldCategoryEmasBatanganGiftSeries       GoldCategory = "emas_batangan_gift_series"
    GoldCategoryEmasBatanganSelamatIdulFitri GoldCategory = "emas_batangan_selamat_idul_fitri"
    GoldCategoryEmasBatanganImlek            GoldCategory = "emas_batangan_imlek"
    GoldCategoryEmasBatanganBatikSeriIII     GoldCategory = "emas_batangan_batik_seri_iii"
    GoldCategoryPerakMurni                   GoldCategory = "perak_murni"
    GoldCategoryPerakHeritage                GoldCategory = "perak_heritage"
    GoldCategoryLiontinBatikSeriIII          GoldCategory = "liontin_batik_seri_iii"
)
```

### Usage in Struct
```go
type GoldPricingHistory struct {
    ID          int          `json:"id" db:"id"`
    PricingDate time.Time    `json:"pricing_date" db:"pricing_date"`
    GoldType    string       `json:"gold_type" db:"gold_type"`
    BuyPrice    int64        `json:"buy_price" db:"buy_price"`
    SellPrice   int64        `json:"sell_price" db:"sell_price"`
    Source      GoldSource   `json:"source" db:"source"`
    Category    GoldCategory `json:"category" db:"category"`
    CreatedAt   time.Time    `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}
```

## API Usage

### Filtering by Category
```bash
# Get all gold bars
curl "http://localhost:8080/api/prices?category=emas_batangan"

# Get gift series
curl "http://localhost:8080/api/prices?category=emas_batangan_gift_series"

# Get pure silver
curl "http://localhost:8080/api/prices?category=perak_murni"
```

### API Response Example
```json
{
  "success": true,
  "message": "Successfully retrieved gold prices",
  "count": 5,
  "data": [
    {
      "id": 1,
      "pricing_date": "2025-11-28",
      "gold_type": "Logam Mulia 1 gram",
      "buy_price": 1164580,
      "sell_price": 1238915,
      "source": "antam",
      "category": "emas_batangan",
      "created_at": "2025-11-28T10:00:00Z",
      "updated_at": "2025-11-28T10:00:00Z"
    }
  ]
}
```

## Category Mapping

When scraping data, you can determine the category based on the product name:

```go
func DetermineCategoryFromGoldType(goldType string) models.GoldCategory {
    goldTypeLower := strings.ToLower(goldType)
    
    // Check for special editions
    if strings.Contains(goldTypeLower, "gift") {
        return models.GoldCategoryEmasBatanganGiftSeries
    }
    if strings.Contains(goldTypeLower, "idul fitri") || strings.Contains(goldTypeLower, "lebaran") {
        return models.GoldCategoryEmasBatanganSelamatIdulFitri
    }
    if strings.Contains(goldTypeLower, "imlek") || strings.Contains(goldTypeLower, "chinese new year") {
        return models.GoldCategoryEmasBatanganImlek
    }
    if strings.Contains(goldTypeLower, "batik") {
        if strings.Contains(goldTypeLower, "liontin") || strings.Contains(goldTypeLower, "pendant") {
            return models.GoldCategoryLiontinBatikSeriIII
        }
        return models.GoldCategoryEmasBatanganBatikSeriIII
    }
    
    // Check for silver
    if strings.Contains(goldTypeLower, "perak") || strings.Contains(goldTypeLower, "silver") {
        if strings.Contains(goldTypeLower, "heritage") {
            return models.GoldCategoryPerakHeritage
        }
        return models.GoldCategoryPerakMurni
    }
    
    // Default to standard gold bars
    return models.GoldCategoryEmasBatangan
}
```

## Migration

### For New Installations
The category column is included in migration `003_create_gold_pricing_histories.sql`.

```bash
make db-reset
```

### For Existing Databases
Run migration `004_add_category_to_gold_pricing.sql`:

```bash
psql -d nabung_emas -f migrations/004_add_category_to_gold_pricing.sql
```

## Default Behavior

- **Default Category**: If no category is specified, `emas_batangan` is used as the default
- **Repository**: Automatically sets default category if not provided
- **API**: Category is optional in create requests

## Display Names

Use the `GetDisplayName()` method to get human-readable names:

```go
category := models.GoldCategoryEmasBatanganGiftSeries
displayName := category.GetDisplayName() // Returns: "Emas Batangan Gift Series"
```

## Validation

The `IsValid()` method checks if a category is valid:

```go
category := models.GoldCategory("emas_batangan")
if category.IsValid() {
    // Category is valid
}
```

## Benefits

1. **Better Organization**: Group products by type
2. **Easier Filtering**: Filter API results by category
3. **Analytics**: Track sales/prices by category
4. **User Experience**: Display products in organized sections
5. **Extensibility**: Easy to add new categories in the future
