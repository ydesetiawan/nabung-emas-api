-- Create enum type for gold category
DO $$ BEGIN
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
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add category column to gold_pricing_histories table
ALTER TABLE gold_pricing_histories 
ADD COLUMN IF NOT EXISTS category gold_category DEFAULT 'emas_batangan';

-- Create index for category
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_category ON gold_pricing_histories(category);

-- Create composite index for category and pricing_date
CREATE INDEX IF NOT EXISTS idx_gold_pricing_histories_category_date ON gold_pricing_histories(category, pricing_date DESC);

-- Add comment
COMMENT ON COLUMN gold_pricing_histories.category IS 'Category of gold/silver product';
