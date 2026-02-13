-- Convert length column back to varchar
ALTER TABLE content 
ALTER COLUMN length TYPE varchar USING length::varchar;