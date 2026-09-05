-- Thin onboarding state for checklist coach (version, displayNextSession, completedAt).
-- New-user default: show coach (version 0, displayNextSession true).
ALTER TABLE users
    ADD COLUMN onboarding JSONB NOT NULL DEFAULT '{"version": 0, "displayNextSession": true, "completedAt": null}'::jsonb;

-- Existing users (including sentinels): already seen at CURRENT_INTRO_VERSION = 1.
UPDATE users
SET onboarding = jsonb_build_object(
    'version', 1,
    'displayNextSession', false,
    'completedAt', to_char(NOW() AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
);
