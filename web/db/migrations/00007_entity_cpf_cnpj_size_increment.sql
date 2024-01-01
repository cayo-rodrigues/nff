-- +goose Up
-- +goose StatementBegin
ALTER TABLE entities ALTER COLUMN cpf_cnpj TYPE VARCHAR(18);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE entities ALTER COLUMN cpf_cnpj TYPE VARCHAR(14);
-- +goose StatementEnd

