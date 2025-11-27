-- Create enum type for source
DO $$ BEGIN
    CREATE TYPE gold_source AS ENUM ('antam', 'usb');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Gold Pricing Histories table
CREATE TABLE IF NOT EXISTS gold_pricing_histories (
    id SERIAL PRIMARY KEY,
    pricing_date DATE NOT NULL,
    gold_type VARCHAR(255) NOT NULL,
    buy_price VARCHAR(50) NOT NULL,
    sell_price VARCHAR(50) NOT NULL,
    source gold_source NOT NULL,
    scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pricing_date, gold_type, source)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_pricing_date ON gold_pricing_histories(pricing_date DESC);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_gold_type ON gold_pricing_histories(gold_type);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_source ON gold_pricing_histories(source);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_scraped_at ON gold_pricing_histories(scraped_at);

-- Composite indexes
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_date_type ON gold_pricing_histories(pricing_date DESC, gold_type);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_date_source ON gold_pricing_histories(pricing_date DESC, source);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_type_source ON gold_pricing_histories(gold_type, source);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_latest ON gold_pricing_histories(gold_type, source, pricing_date DESC);

-- Create trigger function for updated_at
CREATE OR REPLACE FUNCTION update_gold_pricing_histories_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
CREATE TRIGGER trigger_update_gold_pricing_histories_updated_at
    BEFORE UPDATE ON gold_pricing_histories
    FOR EACH ROW
    EXECUTE FUNCTION update_gold_pricing_histories_updated_at();

-- Add comments
COMMENT ON TABLE gold_pricing_histories IS 'Stores historical gold pricing data from various sources';
COMMENT ON COLUMN gold_pricing_histories.pricing_date IS 'The date when the prices were published';
COMMENT ON COLUMN gold_pricing_histories.gold_type IS 'Type/weight of gold (e.g., 1 gram, 5 gram, etc.)';
COMMENT ON COLUMN gold_pricing_histories.buy_price IS 'Buyback price (calculated as 94% of sell price)';
COMMENT ON COLUMN gold_pricing_histories.sell_price IS 'Selling price from the source';
COMMENT ON COLUMN gold_pricing_histories.source IS 'Source of the pricing data (antam, usb)';
COMMENT ON COLUMN gold_pricing_histories.scraped_at IS 'Timestamp when data was scraped';
