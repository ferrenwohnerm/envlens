// Package snapshot provides functionality for saving, loading, and listing
// named snapshots of environment variable maps. Snapshots are persisted as
// JSON files in a configurable directory, identified by a label string.
//
// Typical usage:
//
//	// Save a snapshot
//	err := snapshot.Save("/tmp/envlens", "production-2024-01-15", vars)
//
//	// Load a snapshot later
//	vars, err := snapshot.Load("/tmp/envlens", "production-2024-01-15")
//
//	// List all available snapshots
//	labels, err := snapshot.List("/tmp/envlens")
//
// Snapshot filenames are derived from the label with a .json extension.
// Labels must not contain path separators or be empty.
package snapshot
