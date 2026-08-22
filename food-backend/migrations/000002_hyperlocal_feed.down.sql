DROP INDEX IF EXISTS idx_foods_available_partner_created;
DROP INDEX IF EXISTS idx_food_partners_delivery_radius;
DROP INDEX IF EXISTS idx_food_partners_location;
DROP INDEX IF EXISTS idx_users_location;

ALTER TABLE foods
    DROP COLUMN IF EXISTS available,
    DROP COLUMN IF EXISTS category,
    DROP COLUMN IF EXISTS thumbnail_url,
    DROP COLUMN IF EXISTS image_url;

ALTER TABLE food_partners
    DROP COLUMN IF EXISTS menu_pdf,
    DROP COLUMN IF EXISTS menu_images,
    DROP COLUMN IF EXISTS food_images,
    DROP COLUMN IF EXISTS cover_image,
    DROP COLUMN IF EXISTS profile_image,
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS closing_time,
    DROP COLUMN IF EXISTS opening_time,
    DROP COLUMN IF EXISTS delivery_radius_km,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS latitude;

ALTER TABLE users
    DROP COLUMN IF EXISTS current_location_updated_at,
    DROP COLUMN IF EXISTS last_location_updated,
    DROP COLUMN IF EXISTS current_longitude,
    DROP COLUMN IF EXISTS current_latitude,
    DROP COLUMN IF EXISTS saved_longitude,
    DROP COLUMN IF EXISTS saved_latitude,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS latitude;
