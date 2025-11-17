-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'active';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN status;
-- +goose StatementEnd
