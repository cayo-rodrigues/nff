CREATE TABLE IF NOT EXISTS entities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    user_type VARCHAR(64),
    cpf_cnpj VARCHAR(14),
    ie VARCHAR(13),
    email VARCHAR(128),
    password TEXT,
    postal_code VARCHAR(8),
    neighborhood VARCHAR(64),
    street_type VARCHAR(64),
    street_name VARCHAR(64),
    number VARCHAR(6)
);

