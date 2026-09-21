DROP TABLE IF EXISTS product_categories CASCADE;

CREATE TABLE IF NOT EXISTS product_categories (
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(100) NOT NULL UNIQUE,
    label VARCHAR(100) NOT NULL,
    icon_url TEXT
);

ALTER TABLE products DROP COLUMN IF EXISTS category_id;
ALTER TABLE products ADD COLUMN category_id INT REFERENCES product_categories(id) ON DELETE SET NULL;

DROP INDEX IF EXISTS idx_product_categories_restaurant_id;
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);