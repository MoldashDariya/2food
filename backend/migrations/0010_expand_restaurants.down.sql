DROP TABLE IF EXISTS restaurant_cuisine_types;

ALTER TABLE restaurants
  DROP COLUMN IF EXISTS social_links,
  DROP COLUMN IF EXISTS working_hours,
  DROP COLUMN IF EXISTS longitude,
  DROP COLUMN IF EXISTS latitude,
  DROP COLUMN IF EXISTS website,
  DROP COLUMN IF EXISTS phone;
  DROP COLUMN IF EXISTS city_id;