-- +goose Up
-- +goose StatementBegin
-- Add list_id column to tasks table
ALTER TABLE tasks
    ADD COLUMN list_id BIGINT;

-- Set list_id to the user's default list for existing tasks
UPDATE tasks t
SET list_id = l.id
FROM lists l
WHERE t.user_id = l.user_id
    AND l.name = 'default';

-- Make list_id NOT NULL after populating it
ALTER TABLE tasks
    ALTER COLUMN list_id SET NOT NULL;

-- Add foreign key constraint
ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_list FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE;

-- Create index on list_id for efficient querying
CREATE INDEX idx_tasks_list_id ON tasks(list_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint and index
DROP INDEX IF EXISTS idx_tasks_list_id;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_list;

-- Drop the list_id column
ALTER TABLE tasks DROP COLUMN list_id;
-- +goose StatementEnd
