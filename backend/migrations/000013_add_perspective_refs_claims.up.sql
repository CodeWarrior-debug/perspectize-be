-- Add perspective reference columns
ALTER TABLE perspectives ADD COLUMN primary_perspective_id INTEGER REFERENCES perspectives(id) ON DELETE SET NULL;
ALTER TABLE perspectives ADD COLUMN related_perspective_ids INTEGER[] DEFAULT '{}';
ALTER TABLE perspectives ADD COLUMN custom_fields JSONB DEFAULT '{}';
ALTER TABLE perspectives ADD COLUMN review TEXT;

-- GIN index for reverse lookups on related_perspective_ids
CREATE INDEX idx_perspectives_related_ids ON perspectives USING GIN (related_perspective_ids);

-- Index on primary_perspective_id for FK lookups
CREATE INDEX idx_perspectives_primary_perspective_id ON perspectives(primary_perspective_id);
