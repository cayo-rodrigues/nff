-- +goose Up
-- +goose StatementBegin
ALTER TABLE entities ADD COLUMN other_ies TEXT[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE entities DROP COLUMN other_ies;
-- +goose StatementEnd
