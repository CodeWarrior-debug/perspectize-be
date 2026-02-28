-- Remove the [anonymous] sentinel user.
-- Content previously attributed to [anonymous] will need manual reassignment.
DELETE FROM public.users WHERE username = '[anonymous]';
