ALTER TABLE orders 
    ADD COLUMN IF NOT EXISTS display_number VARCHAR(10);

-- safe backfill
UPDATE orders 
SET display_number = left(id::text, 6) 
WHERE display_number IS NULL;

ALTER TABLE orders 
    ALTER COLUMN display_number SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_active_display_number 
ON orders(restaurant_id, display_number) 
WHERE status IN ('PENDING', 'PAID', 'ACCEPTED', 'READY');