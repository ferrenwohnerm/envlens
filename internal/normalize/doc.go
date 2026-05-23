// Package normalize provides utilities for standardising the keys and values
// of environment variable maps before diffing, auditing, or exporting.
//
// Normalization operations include:
//
//   - Converting keys to UPPER_CASE
//   - Trimming leading/trailing whitespace from keys and values
//   - Replacing hyphens in key names with underscores
//
// None of the functions in this package mutate their inputs; all operations
// return a new Result containing the transformed map and a record of any keys
// that were renamed during the process.
//
// Example usage:
//
//	opts := normalize.DefaultOptions()
//	result := normalize.Apply(vars, opts)
//	fmt.Println(result.Vars)    // normalized map
//	fmt.Println(result.Renamed) // original -> normalized key mapping
package normalize
