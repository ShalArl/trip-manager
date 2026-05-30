ALTER TABLE trips
    DROP COLUMN user_name,
    DROP COLUMN user_email,
    DROP COLUMN user_avatar_key;

ALTER TABLE accommodations
    DROP COLUMN user_name,
    DROP COLUMN user_email,
    DROP COLUMN user_avatar_key;

ALTER TABLE transports
    DROP COLUMN user_name,
    DROP COLUMN user_email,
    DROP COLUMN user_avatar_key;