-- Migration: Create Galeri24 Gold Pricing Schema
-- Description: Creates table and enum for storing gold pricing data from Galeri24

-- Create enum type for gold sources (vendors)
DO $$ BEGIN
    CREATE TYPE gold_source AS ENUM (
        'GALERI_24',
        'DINAR_G24',
        'BABY_GALERI_24',
        'ANTAM',
        'UBS',
        'ANTAM_MULIA_RETRO',
        'ANTAM_NON_PEGADAIAN',
        'LOTUS_ARCHI',
        'UBS_DISNEY',
        'UBS_ELSA',
        'UBS_ANNA',
        'UBS_MICKEY_FULLBODY',
        'LOTUS_ARCHI_GIFT',
        'UBS_HELLO_KITTY',
        'BABY_SERIES_TUMBUHAN',
        'BABY_SERIES_INVESTASI',
        'BATIK_SERIES'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create gold_pricing_histories table
CREATE TABLE IF NOT EXISTS gold_pricing_histories (
    id SERIAL PRIMARY KEY,
    pricing_date DATE NOT NULL,
    gold_type VARCHAR(255) NOT NULL,  -- e.g., "0.5", "1", "2", "5", "10", "25", "50", "100", "250", "500", "1000"
    buy_price VARCHAR(50) NOT NULL,    -- Harga Buyback (e.g., "Rp1.132.000")
    sell_price VARCHAR(50) NOT NULL,   -- Harga Jual (e.g., "Rp1.271.000")
    source gold_source NOT NULL,
    scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pricing_date, gold_type, source)  -- Prevent duplicate entries for same date, type, and source
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_pricing_date ON gold_pricing_histories(pricing_date DESC);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_gold_type ON gold_pricing_histories(gold_type);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_source ON gold_pricing_histories(source);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_scraped_at ON gold_pricing_histories(scraped_at DESC);

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_date_source ON gold_pricing_histories(pricing_date DESC, source);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_type_source ON gold_pricing_histories(gold_type, source);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_latest ON gold_pricing_histories(gold_type, source, pricing_date DESC);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_gold_pricing_histories_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_gold_pricing_histories_updated_at
    BEFORE UPDATE ON gold_pricing_histories
    FOR EACH ROW
    EXECUTE FUNCTION update_gold_pricing_histories_updated_at();

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'gold_pricing_histories_unique_date_type_source'
    ) THEN
        ALTER TABLE gold_pricing_histories 
        ADD CONSTRAINT gold_pricing_histories_unique_date_type_source 
        UNIQUE (pricing_date, gold_type, source);
    END IF;
END $$;

-- Create index on pricing_date if it doesn't exist
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_pricing_date 
ON gold_pricing_histories(pricing_date DESC);

-- Create composite index for date and source if it doesn't exist
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_date_source 
ON gold_pricing_histories(pricing_date DESC, source);


-- Add comments for documentation
COMMENT ON TABLE gold_pricing_histories IS 'Stores historical gold pricing data from various vendors (primarily Galeri24)';
COMMENT ON COLUMN gold_pricing_histories.pricing_date IS 'The date when the prices were published by the vendor';
COMMENT ON COLUMN gold_pricing_histories.gold_type IS 'Weight of gold in grams (0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000)';
COMMENT ON COLUMN gold_pricing_histories.buy_price IS 'Harga Buyback - Price at which vendor buys back gold';
COMMENT ON COLUMN gold_pricing_histories.sell_price IS 'Harga Jual - Price at which vendor sells gold';
COMMENT ON COLUMN gold_pricing_histories.source IS 'Vendor/brand name';
COMMENT ON COLUMN gold_pricing_histories.scraped_at IS 'Timestamp when data was scraped from website';
