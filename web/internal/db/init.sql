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

CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(9),
    protocol VARCHAR(13),
    operation VARCHAR(7) NOT NULL,
    cfop INT NOT NULL,
    is_final_customer VARCHAR(4) NOT NULL,
    is_icms_contributor VARCHAR(6) NOT NULL,
    shipping DOUBLE PRECISION NOT NULL,
    add_shipping_to_total VARCHAR(4) NOT NULL,
    gta VARCHAR(16),
    sender_id BIGINT NOT NULL,
    recipient_id BIGINT NOT NULL,
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(256) DEFAULT 'Em andamento...',
    CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES entities(id) ON DELETE CASCADE,
    CONSTRAINT fk_recipient FOREIGN KEY (recipient_id) REFERENCES entities(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_items (
    id BIGSERIAL PRIMARY KEY,
    item_group VARCHAR(64) NOT NULL,
    description VARCHAR(128) NOT NULL,
    origin VARCHAR(64) NOT NULL,
    unity_of_measurement VARCHAR(8) NOT NULL,
    quantity DOUBLE PRECISION NOT NULL,
    value_per_unity DOUBLE PRECISION NOT NULL,
    invoice_id BIGINT NOT NULL,
    CONSTRAINT fk_invoice FOREIGN KEY(invoice_id) REFERENCES invoices(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_cancelings (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(9) NOT NULL,
    year INT NOT NULL,
    justification VARCHAR(128),
    entity_id BIGINT NOT NULL,
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(256) DEFAULT 'Em andamento...',
    CONSTRAINT fk_canceling_entity FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS metrics_history (
    id BIGSERIAL PRIMARY KEY,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    total_income VARCHAR(16) DEFAULT '0,00',
    total_expenses VARCHAR(16) DEFAULT '0,00',
    avg_income VARCHAR(16) DEFAULT '0,00',
    avg_expenses VARCHAR(16) DEFAULT '0,00',
    diff VARCHAR(16) DEFAULT '0,00',
    is_positive BOOLEAN DEFAULT false,
    total_records INT DEFAULT 0,
    positive_records INT DEFAULT 0,
    negative_records INT DEFAULT 0,
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(256) DEFAULT 'Em andamento...',
    entity_id BIGINT NOT NULL,
    CONSTRAINT fk_metrics_entity FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_printings (
    id BIGSERIAL PRIMARY KEY,
    invoice_id VARCHAR(13), -- number or protocol
    invoice_id_type VARCHAR(13),
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(256) DEFAULT 'Em andamento...',
    entity_id BIGINT NOT NULL,
    CONSTRAINT fk_printing_entity FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE
);
