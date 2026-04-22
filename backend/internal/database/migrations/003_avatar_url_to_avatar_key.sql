-- Rename avatar_url to avatar_key, ensuring avatar_url doesn't exist after migration
DO $$
BEGIN
  -- Case 1: avatar_url exists, avatar_key doesn't -> rename
  IF EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_name='users' AND column_name='avatar_url'
  ) AND NOT EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_name='users' AND column_name='avatar_key'
  ) THEN
    ALTER TABLE users RENAME COLUMN avatar_url TO avatar_key;
  -- Case 2: Both columns exist -> drop avatar_url (avatar_key already has the data)
  ELSIF EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_name='users' AND column_name='avatar_url'
  ) AND EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_name='users' AND column_name='avatar_key'
  ) THEN
    ALTER TABLE users DROP COLUMN avatar_url;
  END IF;
END $$;

DROP INDEX IF EXISTS idx_users_avatar_url;