-- +migrate Up
-- Create users table
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(20) NOT NULL UNIQUE,
  phone_number VARCHAR(20) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  avatar_url TEXT,
  bio TEXT,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

-- Create indexes
CREATE INDEX  idx_users_phone_number ON users(phone_number);
CREATE INDEX  idx_users_username ON users(username);
CREATE INDEX  idx_users_created_at ON users(created_at DESC);


-- Add constraints
ALTER TABLE users ADD CONSTRAINT chk_users_phone_number
    CHECK(phone_number ~ '^09[0-9]{9}$');


--+migrate StatementBegin
CREATE OR REPLACE FUNCTION update_series_updated_at()
RETURN TRIGGER AS $$
BEGIN
    NEW.updated_at=CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();
 -- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_users_updated_at();
-- +migrate StatementEnd

DROP TABLE IF EXISTS users;