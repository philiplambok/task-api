-- +goose Up
-- Add default values for uuid columns to auto-generate UUIDs on insert
ALTER TABLE users ALTER COLUMN uuid SET DEFAULT uuid_generate_v4()::text;
ALTER TABLE lists ALTER COLUMN uuid SET DEFAULT uuid_generate_v4()::text;
ALTER TABLE tasks ALTER COLUMN uuid SET DEFAULT uuid_generate_v4()::text;

-- +goose Down
-- Remove default values
ALTER TABLE tasks ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE lists ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE users ALTER COLUMN uuid DROP DEFAULT;
