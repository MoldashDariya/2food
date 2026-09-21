-- rename old type
ALTER TYPE order_status RENAME TO order_status_old;

-- create new type
CREATE TYPE order_status AS ENUM (
    'PENDING',
    'PAID',
    'READY',
    'COLLECTED',
    'CANCELLED'
);

-- map values
UPDATE orders 
SET status = 'CANCELLED' 
WHERE status::text IN ('EXPIRED', 'REFUNDED');

UPDATE order_status_history 
SET old_status = 'CANCELLED' 
WHERE old_status::text IN ('EXPIRED', 'REFUNDED');

UPDATE order_status_history 
SET new_status = 'CANCELLED' 
WHERE new_status::text IN ('EXPIRED', 'REFUNDED');

UPDATE orders
SET status = 'COLLECTED'
WHERE status::text IN ('UNPAID', 'PARTIALLY_REFUNDED', 'DISPUTED');

UPDATE order_status_history
SET old_status = 'COLLECTED'
WHERE old_status::text IN ('UNPAID', 'PARTIALLY_REFUNDED', 'DISPUTED');

UPDATE order_status_history
SET new_status = 'COLLECTED'
WHERE new_status::text IN ('UNPAID', 'PARTIALLY_REFUNDED', 'DISPUTED');

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
DROP TYPE order_status_old;