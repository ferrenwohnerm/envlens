// Package filter provides post-processing utilities for envlens diff results.
//
// After comparing two environment files with the diff package, callers may
// wish to narrow the result set before rendering a report. This package
// exposes a single Apply function that accepts a map of diff.Result values
// and an Options struct, returning only the entries that match every
// active filter criterion.
//
// Supported filter criteria:
//
//   - Prefix     – retain only keys that start with a given prefix (case-insensitive)
//   - OnlyChanged – retain only keys whose values differ between the two files
//   - OnlyMissing – retain only keys absent from one of the two files
//   - ExcludeKeys – drop specific keys by exact name (case-insensitive)
//
// Multiple criteria are combined with AND semantics: a key must satisfy
// every active condition to appear in the output.
package filter
