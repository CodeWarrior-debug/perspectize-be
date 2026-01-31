-- Add timestamps to users table
ALTER TABLE public.users
    ADD COLUMN created_at timestamptz NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_at timestamptz NOT NULL DEFAULT NOW();

-- Create trigger for users table
CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
