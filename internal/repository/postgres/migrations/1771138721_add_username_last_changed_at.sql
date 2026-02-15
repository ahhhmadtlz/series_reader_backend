-- +migrate Up
ALTER TABLE users ADD COLUMN username_last_changed_at TIMESTAMP;

-- +migrate Down
ALTER TABLE users DROP COLUMN username_last_changed_at;