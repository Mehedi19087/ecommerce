CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    order_number VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL DEFAULT 'pending',
    payment_status VARCHAR(255) DEFAULT 'pending',
    total DECIMAL(10,2) NOT NULL,
    shipping_address TEXT NOT NULL DEFAULT '',
    customer_name VARCHAR(255) NOT NULL DEFAULT '',
    customer_phone VARCHAR(255) NOT NULL DEFAULT '',
    payment_method VARCHAR(255) NOT NULL DEFAULT '',
    notes TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE UNIQUE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);