#!/usr/bin/env bash

set -e

if [ -z "$POSTGRES_URL" ]; then
  export POSTGRES_URL="postgres://postgres:postgres@localhost:5432/musicapp?sslmode=disable"
fi

goose -dir migrations postgres "$POSTGRES_URL" "$1"
