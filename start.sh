#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migrations -database "postgresql://postgres:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" -verbose up

echo "start the app"
exec "$@"