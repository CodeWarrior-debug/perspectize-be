DROP INDEX IF EXISTS idx_perspectives_primary_perspective_id;
DROP INDEX IF EXISTS idx_perspectives_related_ids;
ALTER TABLE perspectives DROP COLUMN IF EXISTS review;
ALTER TABLE perspectives DROP COLUMN IF EXISTS custom_fields;
ALTER TABLE perspectives DROP COLUMN IF EXISTS related_perspective_ids;
ALTER TABLE perspectives DROP COLUMN IF EXISTS primary_perspective_id;
