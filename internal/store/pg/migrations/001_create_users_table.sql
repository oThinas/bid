-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(50) UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash BYTEA NOT NULL,
  bio TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
