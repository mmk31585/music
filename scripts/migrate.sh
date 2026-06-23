#!/usr/bin/env bash

set -e

if [ -z "$POSTGRES_URL" ]; then
  export POSTGRES_URL="postgres://postgres:postgres@localhost:5432/musicapp?sslmode=disable"
fi

# Extract DB name from URL and build a URL to the 'postgres' maintenance database
# URL format: postgres://user:pass@host:port/dbname?params
DB_NAME="${POSTGRES_URL##*/}"
DB_NAME="${DB_NAME%%\?*}"
POSTGRES_ADMIN_URL="${POSTGRES_URL/$DB_NAME/postgres}"

# Create the database if it doesn't exist
psql "$POSTGRES_ADMIN_URL" -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" | grep -q 1 || \
  psql "$POSTGRES_ADMIN_URL" -c "CREATE DATABASE $DB_NAME"

goose -dir migrations postgres "$POSTGRES_URL" "$1"
