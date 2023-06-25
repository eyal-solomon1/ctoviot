#!/bin/sh

set -e 

echo "Running DB migrations"

/app/migrate -path /app/db/migrations  -database "$DB_SOURCE" -verbose up


