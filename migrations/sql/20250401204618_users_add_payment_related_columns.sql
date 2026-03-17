-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN is_blocked BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE users ADD COLUMN has_chosen_payment_plan BOOLEAN NOT NULL DEFAULT false;

UPDATE users SET is_blocked = true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_blocked;
ALTER TABLE users DROP COLUMN has_chosen_payment_plan;
-- +goose StatementEnd
