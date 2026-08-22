ALTER TABLE users
    ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS saved_latitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS saved_longitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS current_latitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS current_longitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS last_location_updated TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS current_location_updated_at TIMESTAMPTZ;

ALTER TABLE food_partners
    ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS delivery_radius_km DOUBLE PRECISION NOT NULL DEFAULT 6,
    ADD COLUMN IF NOT EXISTS opening_time VARCHAR(5),
    ADD COLUMN IF NOT EXISTS closing_time VARCHAR(5),
    ADD COLUMN IF NOT EXISTS rating DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS profile_image TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS cover_image TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS food_images TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS menu_images TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS menu_pdf TEXT NOT NULL DEFAULT '';

ALTER TABLE foods
    ADD COLUMN IF NOT EXISTS image_url TEXT,
    ADD COLUMN IF NOT EXISTS thumbnail_url TEXT,
    ADD COLUMN IF NOT EXISTS category VARCHAR(100),
    ADD COLUMN IF NOT EXISTS available BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_users_location ON users (latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_food_partners_location ON food_partners (latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_food_partners_delivery_radius ON food_partners (delivery_radius_km);
CREATE INDEX IF NOT EXISTS idx_foods_available_partner_created ON foods (available, food_partner_id, created_at DESC);
