-- Drop trigger
DROP TRIGGER IF EXISTS set_updated_at ON public.perspectives;

-- Drop perspectives table
DROP TABLE IF EXISTS public.perspectives;

-- Drop users table
DROP TABLE IF EXISTS public.users;

-- Drop function
DROP FUNCTION IF EXISTS public.update_updated_at();

-- Drop custom domain
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'valid_integer_range') THEN
        DROP DOMAIN IF EXISTS valid_integer_range;
    END IF;
END $$;