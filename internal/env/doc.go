// Package env provides utilities for reading from and writing to the process
// environment at runtime.
//
// It complements the parser package (which handles .env files on disk) by
// operating directly on os.Environ and os.Setenv/os.LookupEnv.
//
// Key functions:
//
//   - Inject   – push a map of key/value pairs into the process environment.
//   - Extract  – pull a named set of keys out of the process environment.
//   - Snapshot – capture a filtered view of the current process environment.
//
// Format helpers (WriteText, WriteJSON, WriteSummary) are provided for
// rendering captured variable maps to any io.Writer.
package env
