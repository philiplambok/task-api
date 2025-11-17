-- +goose Up
-- +goose StatementBegin
-- Rename due_at to deadline and change type from TIMESTAMP to DATE
ALTER TABLE tasks
    ADD COLUMN deadline DATE;

-- Copy data from due_at to deadline (casting timestamp to date)
UPDATE tasks
    SET deadline = due_at::DATE
    WHERE due_at IS NOT NULL;

-- Drop the old column
ALTER TABLE tasks
    DROP COLUMN due_at;

-- Drop the old index
DROP INDEX IF EXISTS idx_tasks_due_at;

-- Create new index on deadline
CREATE INDEX idx_tasks_deadline ON tasks(deadline);

-- Add is_done column with default false
ALTER TABLE tasks
    ADD COLUMN is_done BOOLEAN NOT NULL DEFAULT FALSE;

-- Create index on is_done for efficient filtering
CREATE INDEX idx_tasks_is_done ON tasks(is_done);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove is_done column and its index
DROP INDEX IF EXISTS idx_tasks_is_done;
ALTER TABLE tasks DROP COLUMN is_done;

-- Rename deadline back to due_at and change type back to TIMESTAMP
ALTER TABLE tasks
    ADD COLUMN due_at TIMESTAMP;

-- Copy data back (casting date to timestamp)
UPDATE tasks
    SET due_at = deadline::TIMESTAMP
    WHERE deadline IS NOT NULL;

-- Drop the new column
ALTER TABLE tasks
    DROP COLUMN deadline;

-- Drop the new index and recreate the old one
DROP INDEX IF EXISTS idx_tasks_deadline;
CREATE INDEX idx_tasks_due_at ON tasks(due_at);
-- +goose StatementEnd
