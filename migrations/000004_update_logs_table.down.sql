ALTER TABLE logs
DROP COLUMN IF EXISTS log_type;

ALTER TABLE logs
DROP COLUMN executed_by;

ALTER TABLE logs
RENAME COLUMN triggered_on TO executed_on;

ALTER TABLE logs
DROP COLUMN api_payload;

DROP TYPE IF EXISTS "log_types";