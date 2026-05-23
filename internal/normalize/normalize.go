package normalize

import (
	"strings"
)

// Options controls how normalization is applied to environment variable maps.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// TrimValues removes leading and trailing whitespace from values.
	TrimValues bool
	// TrimKeys removes leading and trailing whitespace from keys.
	TrimKeys bool
	// ReplaceHyphens replaces hyphens in keys with underscores.
	ReplaceHyphens bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:  true,
		TrimValues:     true,
		TrimKeys:       true,
		ReplaceHyphens: true,
	}
}

// Result holds the output of a normalization run.
type Result struct {
	// Vars is the normalized map.
	Vars map[string]string
	// Renamed records keys that changed during normalization (original -> normalized).
	Renamed map[string]string
}

// Apply normalizes the provided environment variable map according to opts.
// It never mutates the input map.
func Apply(vars map[string]string, opts Options) Result {
	out := make(map[string]string, len(vars))
	renamed := make(map[string]string)

	for k, v := range vars {
		newKey := k

		if opts.TrimKeys {
			newKey = strings.TrimSpace(newKey)
		}
		if opts.ReplaceHyphens {
			newKey = strings.ReplaceAll(newKey, "-", "_")
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}

		if newKey != k {
			renamed[k] = newKey
		}
		out[newKey] = v
	}

	return Result{Vars: out, Renamed: renamed}
}
