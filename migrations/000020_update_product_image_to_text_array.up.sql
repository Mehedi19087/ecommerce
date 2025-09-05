-- Drop the column and recreate as text array
ALTER TABLE products DROP COLUMN image;
ALTER TABLE products ADD COLUMN image text[] DEFAULT '{}';