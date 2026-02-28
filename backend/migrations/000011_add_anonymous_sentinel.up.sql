-- Insert "[anonymous]" sentinel user for content added without attribution.
-- Distinct from "[deleted]" (orphaned on user deletion) and "[system]" (pre-existing content).
INSERT INTO public.users (username, email, role)
VALUES ('[anonymous]', 'anonymous@system.internal', 'sentinel')
ON CONFLICT (username) DO NOTHING;
