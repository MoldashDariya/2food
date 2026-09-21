-- rename old type
ALTER TYPE order_status RENAME TO order_status_old;

-- create new type
CREATE TYPE order_status AS ENUM (
    'PENDING',
    'PAID',
    'READY',
    'COLLECTED',
    'CANCELLED',
    'UNPAID',
    'EXPIRED',
    'REFUNDED',
    'PARTIALLY_REFUNDED',
    'DISPUTED'
);

-- map values
UPDATE orders 
SET status = 'PAID'
WHERE status::text = 'ACCEPTED';

UPDATE order_status_history
SET old_status = 'PAID'
WHERE old_status::text = 'ACCEPTED';

UPDATE order_status_history
SET new_status = 'PAID'
WHERE new_status::text = 'ACCEPTED';

-- cast values
ALTER TABLE orders 
    ALTER COLUMN status TYPE order_status 
    USING status::text::order_status;

ALTER TABLE order_status_history 
    ALTER COLUMN old_status TYPE order_status 
    USING old_status::text::order_status,
    ALTER COLUMN new_status TYPE order_status 
    USING new_status::text::order_status;

-- drop old type
DROP TYPE IF EXISTS order_status_old;