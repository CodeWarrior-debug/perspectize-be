-- Convert length column from varchar to integer
ALTER TABLE content 
ALTER COLUMN length TYPE integer USING length::integer;