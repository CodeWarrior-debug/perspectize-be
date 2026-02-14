-- Create custom domain for rating fields
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'valid_integer_range') THEN
        CREATE DOMAIN valid_integer_range AS integer
        CHECK (VALUE BETWEEN 0 AND 10000);
    END IF;
END $$;

-- Create users table
CREATE TABLE public.users (
    id serial NOT NULL,
    username varchar(24) NOT NULL,
    email text NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY(id),
    CONSTRAINT users_unique_email UNIQUE(email),
    CONSTRAINT users_unique_username UNIQUE(username)
);

-- Create perspectives table
CREATE TABLE public.perspectives (
    id serial NOT NULL,
    claim varchar(255) NOT NULL,
    user_id integer NOT NULL,
    content_id integer NULL,
    "like" text NULL,
    quality valid_integer_range NULL,
    agreement valid_integer_range NULL,
    importance valid_integer_range NULL,
    confidence valid_integer_range NULL,
    privacy text NULL DEFAULT 'public',
    parts integer[] NULL,
    category text NULL,
    labels text[] NULL,
    description text NULL,
    review_status text NULL,
    categorized_ratings jsonb[] NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT perspectives_pk PRIMARY KEY(id),
    CONSTRAINT perspectives_unique_user_claims UNIQUE(claim, user_id),
    CONSTRAINT perspectives_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT perspectives_content_fk FOREIGN KEY (content_id) REFERENCES public.content(id)
);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION public.update_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$function$;

-- Create trigger for perspectives table
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON public.perspectives
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();