#!/usr/bin/env bash
# Regenerate the OpenAPI spec under ./docs/ from annotations in the source tree.
#
# Requires `swag` (https://github.com/swaggo/swag) on PATH:
#   go install github.com/swaggo/swag/cmd/swag@v1.16.3
set -euo pipefail

cd "$(dirname "$0")/.."

if ! command -v swag >/dev/null 2>&1; then
  echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@v1.16.3" >&2
  exit 1
fi

swag init \
  --generalInfo cmd/api/main.go \
  --output docs \
  --parseDependency \
  --parseInternal

echo "OK: docs/{docs.go,swagger.json,swagger.yaml} regenerated."
