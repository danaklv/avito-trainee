#!/bin/bash
set -e

echo "Waiting for postgres..."

until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do
  sleep 1
done

echo "Postgres is up - running migrations..."

for file in migrations/*.sql; do
  echo "Applying migration: $file"
  psql postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME -f "$file"
done

echo "Starting application..."
./service
