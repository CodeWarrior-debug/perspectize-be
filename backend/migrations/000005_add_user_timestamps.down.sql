-- Remove trigger
DROP TRIGGER IF EXISTS set_updated_at_users ON public.users;

-- Remove timestamps from users table
ALTER TABLE public.users
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;
