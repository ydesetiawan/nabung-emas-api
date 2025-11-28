CREATE TABLE gold_pricing_histories (
    id SERIAL PRIMARY KEY,
    pricing_date TIMESTAMP NOT NULL,
    gold_type VARCHAR(50) NOT NULL,
    base_price BIGINT NOT NULL,
    buy_price BIGINT NOT NULL,
    sell_price BIGINT NOT NULL,
    include_tax BOOLEAN DEFAULT TRUE,
    source VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_pricing_history UNIQUE(pricing_date, source, gold_type, category)
);

-- Add indexes for common queries
CREATE INDEX idx_gold_pricing_histories_date ON gold_pricing_histories(pricing_date);
CREATE INDEX idx_gold_pricing_histories_source ON gold_pricing_histories(source);
CREATE INDEX idx_gold_pricing_histories_category ON gold_pricing_histories(category);
CREATE INDEX idx_gold_pricing_histories_gold_type ON gold_pricing_histories(gold_type);

-- Trigger for updated_at
CREATE TRIGGER update_gold_pricing_histories_updated_at BEFORE UPDATE ON gold_pricing_histories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
