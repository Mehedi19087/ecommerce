-- Drop the column and recreate as JSON
ALTER TABLE products DROP COLUMN image;
ALTER TABLE products ADD COLUMN image json DEFAULT '[]'::json;