-- +goose Up
-- +goose StatementBegin

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS entities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    user_type TEXT,
    cpf_cnpj TEXT,
    ie TEXT,
    email TEXT,
    password TEXT,
    postal_code TEXT,
    neighborhood TEXT,
    street_type TEXT,
    street_name TEXT,
    number TEXT,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    number TEXT,
    protocol TEXT,
    operation TEXT NOT NULL,
    cfop INTEGER NOT NULL,
    is_final_customer TEXT NOT NULL,
    is_icms_contributor TEXT NOT NULL,
    shipping REAL NOT NULL,
    add_shipping_to_total TEXT NOT NULL,
    gta TEXT,
    invoice_pdf TEXT DEFAULT '',
    req_status TEXT DEFAULT 'pending',
    req_msg TEXT DEFAULT 'Em andamento...',
    sender_id INTEGER NOT NULL,
    recipient_id INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES entities(id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES entities(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_group TEXT NOT NULL,
    description TEXT NOT NULL,
    origin TEXT NOT NULL,
    unity_of_measurement TEXT NOT NULL,
    quantity REAL NOT NULL,
    value_per_unity REAL NOT NULL,
    invoice_id INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_cancelings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    invoice_number TEXT NOT NULL,
    year INTEGER NOT NULL,
    justification TEXT,
    req_status TEXT DEFAULT 'pending',
    req_msg TEXT DEFAULT 'Em andamento...',
    entity_id INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_printings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    invoice_id TEXT,
    invoice_id_type TEXT,
    invoice_pdf TEXT DEFAULT '',
    req_status TEXT DEFAULT 'pending',
    req_msg TEXT DEFAULT 'Em andamento...',
    entity_id INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS metrics_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    total_income TEXT DEFAULT '0,00',
    total_expenses TEXT DEFAULT '0,00',
    avg_income TEXT DEFAULT '0,00',
    avg_expenses TEXT DEFAULT '0,00',
    diff TEXT DEFAULT '0,00',
    is_positive INTEGER DEFAULT 0,
    total_records INTEGER DEFAULT 0,
    positive_records INTEGER DEFAULT 0,
    negative_records INTEGER DEFAULT 0,
    req_status TEXT DEFAULT 'pending',
    req_msg TEXT DEFAULT 'Em andamento...',
    entity_id INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics_history;
DROP TABLE IF EXISTS invoices_printings;
DROP TABLE IF EXISTS invoices_cancelings;
DROP TABLE IF EXISTS invoices_items;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS entities;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
