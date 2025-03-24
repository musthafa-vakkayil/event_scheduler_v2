CREATE TYPE event_types AS ENUM ('API', 'SCHEDULED');

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT NULL
);

CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "type" event_types NOT NULL,
  "api_end_point" varchar,
  "api_method" varchar,
  "api_request_body" JSON,
  "created_by" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "executed_at" timestamp DEFAULT null
);

CREATE TABLE "logs" (
  "id" bigserial PRIMARY KEY,
  "event_id" bigint NOT NULL,
  "executed_on" timestamp NOT NULL,
  "status" varchar NOT NULL,
  "is_archived" bool DEFAULT false 
);

ALTER TABLE "events" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("username");

ALTER TABLE "logs" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");