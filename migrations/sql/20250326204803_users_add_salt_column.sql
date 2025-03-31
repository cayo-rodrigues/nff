-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN salt TEXT NOT NULL DEFAULT '';
UPDATE users SET salt = hex(randomblob(16)) WHERE salt IS NULL OR salt = '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN salt;
-- +goose StatementEnd
