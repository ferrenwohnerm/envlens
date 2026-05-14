// Package filter provides utilities for filtering diff results
// based on key patterns, status types, and other criteria.
package filter

import (
	"strings"

	"github.com/user/envlens/internal/diff"
)

// Options holds the configuration for filtering diff results.
type Options struct {
	// Prefix restricts results to keys with this prefix (case-insensitive).
	Prefix string
	// OnlyChanged, when true, includes only keys whose values changed.
	OnlyChanged bool
	// OnlyMissing, when true, includes only keys missing from one side.
	OnlyMissing bool
	// ExcludeKeys is a set of exact key names to omit from results.
	ExcludeKeys []string
}

// Apply returns a new map containing only the entries from results that
// satisfy all conditions expressed in opts.
func Apply(results map[string]diff.Result, opts Options) map[string]diff.Result {
	excluded := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]diff.Result)
	for key, result := range results {
		if _, skip := excluded[strings.ToUpper(key)]; skip {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(opts.Prefix)) {
			continue
		}
		if opts.OnlyChanged && result.Status != diff.StatusChanged {
			continue
		}
		if opts.OnlyMissing && result.Status != diff.StatusOnlyInA && result.Status != diff.StatusOnlyInB {
			continue
		}
		out[key] = result
	}
	return out
}
