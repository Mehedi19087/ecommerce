CREATE TABLE sub_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE sub_sub_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    sub_category_id INTEGER NOT NULL REFERENCES sub_categories(id) ON DELETE CASCADE,
    product_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_sub_categories_category_id ON sub_categories(category_id);
CREATE INDEX idx_sub_categories_deleted_at ON sub_categories(deleted_at);
CREATE INDEX idx_sub_sub_categories_sub_category_id ON sub_sub_categories(sub_category_id);
CREATE INDEX idx_sub_sub_categories_deleted_at ON sub_sub_categories(deleted_at);