ALTER TABLE restaurants
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS website TEXT,
    ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS working_hours JSONB,
    ADD COLUMN IF NOT EXISTS social_links JSONB,
    ADD COLUMN IF NOT EXISTS city_id INT REFERENCES cities(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS restaurant_cuisine_types (
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    cuisine_type_id INT NOT NULL REFERENCES cuisine_types(id),
    PRIMARY KEY (restaurant_id, cuisine_type_id)
);

CREATE INDEX IF NOT EXISTS idx_restaurants_city_id ON restaurants(city_id);
CREATE INDEX IF NOT EXISTS idx_rct_cuisine_type_id ON restaurant_cuisine_types(cuisine_type_id);