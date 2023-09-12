CREATE TABLE IF NOT EXISTS addresses (
    id BIGSERIAL PRIMARY KEY,
    postal_code VARCHAR(8),
    street_type VARCHAR(64),
    street_name VARCHAR(64),
    number VARCHAR(6),
    neighborhood VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS entities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    user_type VARCHAR(64),
    cpf_cnpj VARCHAR(14),
    ie VARCHAR(13),
    email VARCHAR(128),
    password TEXT,
    address_id BIGINT,
    FOREIGN KEY(address_id) REFERENCES addresses(id)
);

