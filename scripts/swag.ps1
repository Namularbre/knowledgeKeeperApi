# Regenerate the OpenAPI spec under ./docs/ from annotations in the source tree.
#
# Requires `swag` (https://github.com/swaggo/swag) on PATH:
#   go install github.com/swaggo/swag/cmd/swag@latest

$ErrorActionPreference = 'Stop'

Set-Location (Split-Path $PSScriptRoot -Parent)

if (-not (Get-Command swag -ErrorAction SilentlyContinue)) {
  Write-Error "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest" -ErrorAction Stop
  exit 1
}

& swag init `
  --generalInfo cmd/api/main.go `
  --output docs `
  --parseDependency `
  --parseInternal

Write-Host "OK: docs/{docs.go,swagger.json,swagger.yaml} regenerated."