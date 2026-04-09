// Package migrations embeds versioned SQL migration files for use with goose.
// The embedded FS is consumed by data.Connect to apply schema and seed
// migrations at startup without requiring migration files on the filesystem
// at runtime.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
