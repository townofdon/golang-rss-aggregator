#!/bin/bash

# run this from the project root

source "${PWD}/bin/_utils.sh"

DB_URL=$(read_env_var DB_URL)

cd sql/schema

echo "Running migrations..."

goose postgres $DB_URL up

cd - > /dev/null

echo "Generating DB transport..."

sqlc generate

echo "✅ All done!"
