ALTER TABLE USERS
    ALTER COLUMN password_hash TYPE TEXT,
    ALTER COLUMN password_hash SET NOT NULL,
    DROP COLUMN IF EXISTS profile_url,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;

