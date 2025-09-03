CREATE TABLE payment_proofs (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    transaction_id VARCHAR(255),
    amount DECIMAL(10,2),
    payment_method VARCHAR(255),
    status VARCHAR(255) DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_payment_proofs_order_id ON payment_proofs(order_id);