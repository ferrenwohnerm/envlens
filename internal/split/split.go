// Package split provides functionality for splitting a flat env map
// into multiple named groups based on prefix rules.
package split

import "strings"

// Options controls the behaviour of the Split operation.
type Options struct {
	// StripPrefix removes the matched prefix from keys in each group.
	StripPrefix bool
	// IncludeUnmatched collects keys that matched no prefix into a group
	// keyed by the empty string.
	IncludeUnmatched bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		StripPrefix:      true,
		IncludeUnmatched: false,
	}
}

// Result holds the outcome of a Split call.
type Result struct {
	// Groups maps each prefix (or "" for unmatched) to its env vars.
	Groups map[string]map[string]string
	// Unmatched contains keys that did not match any prefix.
	Unmatched []string
}

// Run partitions vars into groups according to prefixes.
// The first matching prefix wins; prefixes are tested in the order supplied.
func Run(vars map[string]string, prefixes []string, opts Options) Result {
	groups := make(map[string]map[string]string, len(prefixes))
	for _, p := range prefixes {
		groups[p] = make(map[string]string)
	}

	var unmatched []string

	for k, v := range vars {
		matched := false
		for _, p := range prefixes {
			if strings.HasPrefix(k, p) {
				key := k
				if opts.StripPrefix {
					key = strings.TrimPrefix(k, p)
				}
				groups[p][key] = v
				matched = true
				break
			}
		}
		if !matched {
			unmatched = append(unmatched, k)
			if opts.IncludeUnmatched {
				if groups[""] == nil {
					groups[""] = make(map[string]string)
				}
				groups[""][k] = v
			}
		}
	}

	return Result{
		Groups:    groups,
		Unmatched: unmatched,
	}
}
