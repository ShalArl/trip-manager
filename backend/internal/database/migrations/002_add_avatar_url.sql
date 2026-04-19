-- Migration: Add avatar_url column
-- Description: Add avatar_url column to users table for storing avatar image URLs

-- Add the new avatar_url column
ALTER TABLE users
ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500);


-- Add index for avatar_url
CREATE INDEX IF NOT EXISTS idx_users_avatar_url ON users(avatar_url);

