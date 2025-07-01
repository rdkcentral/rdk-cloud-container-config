#!/bin/bash

set -e

# Wait for Cassandra to become available
until cqlsh cassandra -e "describe keyspaces"; do
  echo "Waiting for cqlsh to connect..."; sleep 5;
done

# Run schema file
echo "Running schema..."
cd /docker-entrypoint-initdb.d
cqlsh cassandra -f db_init.cql

# Mark schema as applied
echo "Schema done."
touch /db/schema_ready

if touch /db/schema_ready; then
  echo "File created successfully."
else
  echo "Failed to create file." >&2
  exit 1
fi
