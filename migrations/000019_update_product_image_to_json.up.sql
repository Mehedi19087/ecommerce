-- Add image column as JSON array (since it doesn't exist)
ALTER TABLE products 
ADD COLUMN image JSON DEFAULT '[]'::json;