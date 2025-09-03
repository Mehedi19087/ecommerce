CREATE TABLE payment_proofs (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    transaction_id VARCHAR(255) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    amount NUMERIC NOT NULL,
    screenshot TEXT NOT NULL,
    sender_number VARCHAR(50) NOT NULL,
    sender_name VARCHAR(100) NOT NULL,
    payment_date VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    admin_notes TEXT,
    reviewed_by INTEGER,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_payment_proofs_order_id ON payment_proofs(order_id);