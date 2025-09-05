ALTER TABLE products ADD COLUMN sub_category_id INTEGER REFERENCES sub_categories(id);
ALTER TABLE products ADD COLUMN sub_sub_category_id INTEGER REFERENCES sub_sub_categories(id);

CREATE INDEX idx_products_sub_category_id ON products(sub_category_id);
CREATE INDEX idx_products_sub_sub_category_id ON products(sub_sub_category_id);