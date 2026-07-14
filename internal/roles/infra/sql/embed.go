package sql

import "embed"

//go:embed 001_roles.sql
var Schema embed.FS
