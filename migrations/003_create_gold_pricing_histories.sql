-- Create enum type for source
DO $$ BEGIN
    CREATE TYPE gold_source AS ENUM ('antam', 'usb');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Gold Pricing Histories table
CREATE TABLE IF NOT EXISTS gold_pricing_histories (
    id SERIAL PRIMARY KEY,
    gold_type VARCHAR(255) NOT NULL,
    buy_price VARCHAR(50) NOT NULL,
    sell_price VARCHAR(50) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    source gold_source NOT NULL,
    scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_scraped_at ON gold_pricing_histories(scraped_at);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_gold_type ON gold_pricing_histories(gold_type);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_source ON gold_pricing_histories(source);
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_gold_type_source ON gold_pricing_histories(gold_type, source);

-- Create a composite index for getting latest prices
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_latest ON gold_pricing_histories(gold_type, source, scraped_at DESC);
