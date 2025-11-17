-- +goose Up
-- Add uuid extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Add uuid columns to users table
ALTER TABLE users ADD COLUMN uuid VARCHAR(255) UNIQUE;
UPDATE users SET uuid = uuid_generate_v4()::text WHERE uuid IS NULL;
ALTER TABLE users ALTER COLUMN uuid SET NOT NULL;
CREATE INDEX idx_users_uuid ON users(uuid);

-- Add uuid columns to lists table
ALTER TABLE lists ADD COLUMN uuid VARCHAR(255) UNIQUE;
UPDATE lists SET uuid = uuid_generate_v4()::text WHERE uuid IS NULL;
ALTER TABLE lists ALTER COLUMN uuid SET NOT NULL;
CREATE INDEX idx_lists_uuid ON lists(uuid);

-- Add uuid columns to tasks table
ALTER TABLE tasks ADD COLUMN uuid VARCHAR(255) UNIQUE;
UPDATE tasks SET uuid = uuid_generate_v4()::text WHERE uuid IS NULL;
ALTER TABLE tasks ALTER COLUMN uuid SET NOT NULL;
CREATE INDEX idx_tasks_uuid ON tasks(uuid);

-- +goose Down
-- Remove uuid columns and indexes
DROP INDEX IF EXISTS idx_tasks_uuid;
ALTER TABLE tasks DROP COLUMN IF EXISTS uuid;

DROP INDEX IF EXISTS idx_lists_uuid;
ALTER TABLE lists DROP COLUMN IF EXISTS uuid;

DROP INDEX IF EXISTS idx_users_uuid;
ALTER TABLE users DROP COLUMN IF EXISTS uuid;

-- Note: We don't drop the uuid-ossp extension as it might be used by other tables
