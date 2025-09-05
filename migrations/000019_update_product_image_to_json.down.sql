-- Remove the image column
ALTER TABLE products 
DROP COLUMN IF EXISTS image;