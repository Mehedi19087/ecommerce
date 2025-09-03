CREATE TABLE visitor_logs (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(45),
    country VARCHAR(100),
    region VARCHAR(100),
    city VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);