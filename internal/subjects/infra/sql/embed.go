package sql

import "embed"

//go:embed 001_subjects.sql
var Schema embed.FS
