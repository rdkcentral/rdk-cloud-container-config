/**
 * Copyright 2025 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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
