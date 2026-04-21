ALTER TABLE users RENAME COLUMN avatar_url TO avatar_key;
DROP INDEX IF EXISTS idx_users_avatar_url;