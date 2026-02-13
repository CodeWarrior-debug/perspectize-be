-- Revert to text type (though this is unlikely to be needed)
ALTER TABLE content 
ALTER COLUMN response TYPE text USING response::text;