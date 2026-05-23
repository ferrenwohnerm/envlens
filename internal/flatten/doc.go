// Package flatten normalises nested or dot-separated environment variable keys
// into a flat, uppercase format suitable for use in standard .env files or
// process environments.
//
// Keys that use dots ("db.host"), forward-slashes ("db/host"), or mixed
// separators are converted to a single configurable delimiter (default "_").
// An optional prefix can be prepended to every output key.
//
// Example:
//
//	input:  {"db.host": "localhost", "db.port": "5432"}
//	output: {"DB_HOST": "localhost", "DB_PORT": "5432"}
//
// Duplicate output keys (after normalisation) are rejected with an error to
// prevent silent data loss.
package flatten
