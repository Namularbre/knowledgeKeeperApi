// Package sql exposes the auth bootstrap schema as an embedded string so the
// binary stays self-contained (no external SQL files at runtime).
package sql

import "embed"

//go:embed 001_users.sql
//go:embed 002_users.sql
var Schema embed.FS
