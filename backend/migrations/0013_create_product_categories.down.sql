DROP INDEX IF EXISTS idx_products_category_id;

ALTER TABLE products DROP COLUMN IF EXISTS category_id;

DROP TABLE IF EXISTS product_categories;

CREATE TABLE IF NOT EXISTS product_categories (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    icon_url      TEXT,
    sort_order    INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE products ADD COLUMN category_id UUID REFERENCES product_categories(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_product_categories_restaurant_id ON product_categories(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);