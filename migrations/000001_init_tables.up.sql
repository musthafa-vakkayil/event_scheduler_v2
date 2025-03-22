CREATE TYPE event_types AS ENUM ('API', 'SCHEDULED');
CREATE TYPE status_types AS ENUM ('SUCCESS', 'FAILED');

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "type" event_types NOT NULL,
  "api_end_point" varchar NOT NULL,
  "api_method" varchar,
  "api_request_body" JSON,
  "created_by" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "executed_at" timestamptz DEFAULT null
);

CREATE TABLE "logs" (
  "id" bigserial PRIMARY KEY,
  "event_id" bigint NOT NULL,
  "executed_on" timestamptz NOT NULL,
  "status" status_types NOT NULL,
  "is_archived" bool DEFAULT false
);

ALTER TABLE "events" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("username");

ALTER TABLE "logs" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");