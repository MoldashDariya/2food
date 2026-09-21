DROP INDEX IF EXISTS idx_users_city_id;

ALTER TABLE users
    DROP COLUMN IF EXISTS city_id,
    DROP COLUMN IF EXISTS language,
    DROP COLUMN IF EXISTS notification_settings,
    DROP COLUMN IF EXISTS dietary_preferences;