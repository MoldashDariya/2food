ALTER TABLE users
    ADD COLUMN IF NOT EXISTS city_id INT REFERENCES cities(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS language VARCHAR(5) NOT NULL DEFAULT 'en',
    ADD COLUMN IF NOT EXISTS notification_settings JSONB NOT NULL DEFAULT '{"order_updates": true, "promotions": true}',
    ADD COLUMN IF NOT EXISTS dietary_preferences JSONB NOT NULL DEFAULT '[]';

CREATE INDEX IF NOT EXISTS idx_users_city_id ON users(city_id);