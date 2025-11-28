#!/bin/bash

echo "🔍 Verifying SQL queries and Scan operations..."
echo ""

echo "📊 Checking SELECT queries (should have 9 columns):"
grep -n "SELECT.*pricing_date.*gold_type.*buy_price.*sell_price.*source" internal/repositories/gold_pricing_history_repository.go | while read line; do
    echo "  $line"
    # Count commas in the SELECT statement
    column_count=$(echo "$line" | grep -o "," | wc -l)
    column_count=$((column_count + 1))
    echo "    → Columns: $column_count"
done

echo ""
echo "📝 Checking Scan operations (should have 9 arguments):"
grep -A 10 "rows.Scan\|QueryRow.*Scan" internal/repositories/gold_pricing_history_repository.go | grep -B 1 "&history" | grep -c "&history"

echo ""
echo "✅ Verification complete!"
