DROP INDEX IF EXISTS idx_unique_active_display_number;

ALTER TABLE orders 
    DROP COLUMN IF EXISTS display_number;