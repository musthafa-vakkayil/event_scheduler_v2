CREATE TYPE log_types AS ENUM ('TEST_EVENT', 'API_EVENT', 'SCHEDULED_EVENT');

ALTER table logs
ADD COLUMN log_type log_types NOT NULL DEFAULT 'API_EVENT';

ALTER TABLE logs
ADD COLUMN executed_by varchar NOT NULL;

ALTER TABLE logs
RENAME COLUMN executed_on TO triggered_on;

ALTER TABLE logs
ADD COLUMN api_payload JSON;