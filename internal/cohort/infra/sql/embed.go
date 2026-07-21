package sql

import "embed"

//go:embed 001_cohorts.sql
var Schema embed.FS
