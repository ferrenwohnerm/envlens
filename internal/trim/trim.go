// Package trim provides utilities for removing keys from environment
// variable maps based on prefix, suffix, or explicit key lists.
package trim

import "strings"

// Options controls which keys are removed during a Trim operation.
type Options struct {
	// Keys is an explicit list of key names to remove.
	Keys []string

	// Prefix removes any key whose name starts with the given string.
	Prefix string

	// Suffix removes any key whose name ends with the given string.
	Suffix string

	// DryRun returns a report of what would be removed without mutating
	// the input map.
	DryRun bool
}

// Result holds the outcome of a Trim operation.
type Result struct {
	// Removed contains the keys that were (or would be) removed.
	Removed []string

	// Vars is the resulting map after removal (nil when DryRun is true).
	Vars map[string]string
}

// Apply removes keys from vars according to opts and returns a Result.
// The original map is never mutated; a new map is always returned unless
// DryRun is enabled.
func Apply(vars map[string]string, opts Options) Result {
	removed := []string{}

	should := func(key string) bool {
		if opts.Prefix != "" && strings.HasPrefix(key, opts.Prefix) {
			return true
		}
		if opts.Suffix != "" && strings.HasSuffix(key, opts.Suffix) {
			return true
		}
		for _, k := range opts.Keys {
			if k == key {
				return true
			}
		}
		return false
	}

	for k := range vars {
		if should(k) {
			removed = append(removed, k)
		}
	}

	if opts.DryRun {
		return Result{Removed: removed, Vars: nil}
	}

	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if !should(k) {
			out[k] = v
		}
	}
	return Result{Removed: removed, Vars: out}
}
