-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN password_digest VARCHAR(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN password_digest;
-- +goose StatementEnd
