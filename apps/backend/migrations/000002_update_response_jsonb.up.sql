-- This migration ensures the response column is of type jsonb
-- Since the initial migration already creates it as jsonb, this is effectively a no-op
-- but is included for parity with the C# migration history
DO $$
BEGIN
    -- Check if the column exists and isn't jsonb
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'content' 
        AND column_name = 'response'
    ) THEN
        -- If response column exists and isn't jsonb, convert it
        IF (SELECT data_type FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = 'content' 
            AND column_name = 'response') != 'jsonb' THEN
            ALTER TABLE content 
            ALTER COLUMN response TYPE jsonb USING response::jsonb;
        END IF;
    END IF;
END $$;