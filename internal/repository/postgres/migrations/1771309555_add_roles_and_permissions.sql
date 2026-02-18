-- +migrate Up

-- Create ENUM types
CREATE TYPE user_role AS ENUM ('user', 'manager', 'admin');
CREATE TYPE subscription_tier AS ENUM ('free', 'premium');
CREATE TYPE collaboration_role AS ENUM ('viewer', 'editor', 'publisher', 'owner');

-- Add role and subscription_tier to users table
ALTER TABLE users
    ADD COLUMN role user_role NOT NULL DEFAULT 'user',
    ADD COLUMN subscription_tier subscription_tier NOT NULL DEFAULT 'free';

-- Create user_permissions table
CREATE TABLE IF NOT EXISTS user_permissions (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission  VARCHAR(100) NOT NULL,
    granted_by  INTEGER NOT NULL REFERENCES users(id),
    granted_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- One permission per user (no duplicates)
    CONSTRAINT uq_user_permission UNIQUE (user_id, permission),

    -- Only valid permissions allowed
    CONSTRAINT chk_valid_permission CHECK (
        permission IN (
            'moderate_comments',
            'unpublish_content',
            'manage_series_global',
            'manage_users',
            'view_analytics'
        )
    )
);

-- Add created_by and is_premium_only to series table
ALTER TABLE series
    ADD COLUMN created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN is_premium_only BOOLEAN NOT NULL DEFAULT false;

-- Create series_collaborators table
CREATE TABLE IF NOT EXISTS series_collaborators (
    series_id           INTEGER NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    user_id             INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    collaboration_role  collaboration_role NOT NULL DEFAULT 'viewer',
    added_by            INTEGER NOT NULL REFERENCES users(id),
    added_at            TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (series_id, user_id)
);

-- Enforce single owner per series
-- Only one collaborator can have 'owner' role per series
CREATE UNIQUE INDEX uq_series_single_owner
    ON series_collaborators(series_id)
    WHERE collaboration_role = 'owner';

-- Indexes
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_subscription_tier ON users(subscription_tier);
CREATE INDEX idx_user_permissions_user_id ON user_permissions(user_id);
CREATE INDEX idx_series_created_by ON series(created_by);
CREATE INDEX idx_series_collaborators_user_id ON series_collaborators(user_id);
CREATE INDEX idx_series_collaborators_series_id ON series_collaborators(series_id);

-- +migrate Down

DROP INDEX IF EXISTS uq_series_single_owner;
DROP INDEX IF EXISTS idx_series_collaborators_series_id;
DROP INDEX IF EXISTS idx_series_collaborators_user_id;
DROP INDEX IF EXISTS idx_series_created_by;
DROP INDEX IF EXISTS idx_user_permissions_user_id;
DROP INDEX IF EXISTS idx_users_subscription_tier;
DROP INDEX IF EXISTS idx_users_role;

DROP TABLE IF EXISTS series_collaborators;
DROP TABLE IF EXISTS user_permissions;

ALTER TABLE series
    DROP COLUMN IF EXISTS is_premium_only,
    DROP COLUMN IF EXISTS created_by;

ALTER TABLE users
    DROP COLUMN IF EXISTS subscription_tier,
    DROP COLUMN IF EXISTS role;

DROP TYPE IF EXISTS collaboration_role;
DROP TYPE IF EXISTS subscription_tier;
DROP TYPE IF EXISTS user_role;