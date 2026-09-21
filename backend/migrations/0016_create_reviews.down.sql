DROP TABLE IF EXISTS reviews;

ALTER TABLE restaurants
    DROP COLUMN IF EXISTS rating_avg,
    DROP COLUMN IF EXISTS rating_count;

ALTER TABLE orders
    DROP COLUMN IF EXISTS reviewed_at;