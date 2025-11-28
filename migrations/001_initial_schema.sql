-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Password reset tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Type Pockets (Categories)
CREATE TABLE type_pockets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    color VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pockets
CREATE TABLE pockets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type_pocket_id UUID NOT NULL REFERENCES type_pockets(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    aggregate_total_price DECIMAL(15, 2) DEFAULT 0,
    aggregate_total_weight DECIMAL(10, 3) DEFAULT 0,
    target_weight DECIMAL(10, 3),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_pocket_name_per_user UNIQUE(user_id, name)
);

-- Transactions
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pocket_id UUID NOT NULL REFERENCES pockets(id) ON DELETE CASCADE,
    transaction_date DATE NOT NULL,
    brand VARCHAR(50) NOT NULL,
    weight DECIMAL(10, 3) NOT NULL CHECK (weight >= 0.1 AND weight <= 1000),
    price_per_gram DECIMAL(15, 2) NOT NULL CHECK (price_per_gram >= 1000 AND price_per_gram <= 10000000),
    total_price DECIMAL(15, 2) NOT NULL,
    description TEXT,
    receipt_image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- User Settings
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language VARCHAR(10) DEFAULT 'en',
    theme VARCHAR(20) DEFAULT 'light',
    currency VARCHAR(10) DEFAULT 'IDR',
    email_notifications BOOLEAN DEFAULT TRUE,
    push_notifications BOOLEAN DEFAULT FALSE,
    price_alerts BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_settings_per_user UNIQUE(user_id)
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_pockets_user_id ON pockets(user_id);
CREATE INDEX idx_pockets_type_pocket_id ON pockets(type_pocket_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_pocket_id ON transactions(pocket_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_brand ON transactions(brand);

CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_type_pockets_updated_at BEFORE UPDATE ON type_pockets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pockets_updated_at BEFORE UPDATE ON pockets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update pocket aggregates
CREATE OR REPLACE FUNCTION update_pocket_aggregates()
RETURNS TRIGGER AS $$
DECLARE
    v_pocket_id UUID;
BEGIN
    -- Determine which pocket to update
    IF TG_OP = 'DELETE' THEN
        v_pocket_id := OLD.pocket_id;
    ELSE
        v_pocket_id := NEW.pocket_id;
    END IF;

    -- Update pocket aggregates
    UPDATE pockets
    SET 
        aggregate_total_weight = COALESCE((
            SELECT SUM(weight) FROM transactions WHERE pocket_id = v_pocket_id
        ), 0),
        aggregate_total_price = COALESCE((
            SELECT SUM(total_price) FROM transactions WHERE pocket_id = v_pocket_id
        ), 0),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_pocket_id;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update pocket aggregates
CREATE TRIGGER trigger_update_pocket_aggregates
AFTER INSERT OR UPDATE OR DELETE ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_pocket_aggregates();

-- Seed data for type_pockets
INSERT INTO type_pockets (name, description, icon, color) VALUES
    ('Emergency Fund', 'Savings for emergency situations', 'heroicons:shield-check', 'blue'),
    ('Wedding', 'Savings for wedding expenses', 'heroicons:heart', 'pink'),
    ('Investment', 'General investment savings', 'heroicons:chart-bar', 'gold'),
    ('Education', 'Savings for education expenses', 'heroicons:academic-cap', 'green'),
    ('Retirement', 'Savings for retirement', 'heroicons:home', 'purple'),
    ('Vacation', 'Savings for vacation and travel', 'heroicons:globe-alt', 'cyan'),
    ('Business', 'Savings for business capital', 'heroicons:briefcase', 'orange');
