// Package intersect provides utilities for finding common keys and
// shared values between two environment variable maps.
package intersect

import (
	"sort"
	"strings"
)

// Result holds the outcome of an intersection operation between two env maps.
type Result struct {
	// SharedKeys contains keys present in both maps (values may differ).
	SharedKeys []string
	// MatchingPairs contains keys where both key and value are identical.
	MatchingPairs map[string]string
	// DivergentKeys contains keys present in both maps but with different values.
	DivergentKeys []string
}

// Options controls the behaviour of Run.
type Options struct {
	// CaseFold treats keys as case-insensitive when matching.
	CaseFold bool
}

// DefaultOptions returns the recommended default Options.
func DefaultOptions() Options {
	return Options{CaseFold: false}
}

// Run compares two env maps and returns keys/values they share,
// as well as keys present in both but holding different values.
func Run(a, b map[string]string, opts Options) Result {
	lookup := buildLookup(b, opts.CaseFold)

	res := Result{
		MatchingPairs: make(map[string]string),
	}

	for k, va := range a {
		nk := normalise(k, opts.CaseFold)
		vb, exists := lookup[nk]
		if !exists {
			continue
		}
		res.SharedKeys = append(res.SharedKeys, k)
		if va == vb {
			res.MatchingPairs[k] = va
		} else {
			res.DivergentKeys = append(res.DivergentKeys, k)
		}
	}

	sort.Strings(res.SharedKeys)
	sort.Strings(res.DivergentKeys)
	return res
}

func buildLookup(m map[string]string, fold bool) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[normalise(k, fold)] = v
	}
	return out
}

func normalise(s string, fold bool) string {
	if fold {
		return strings.ToLower(s)
	}
	return s
}
