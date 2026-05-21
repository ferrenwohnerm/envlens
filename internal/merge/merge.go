// Package merge provides functionality for merging two environment variable
// maps, with configurable conflict resolution strategies.
package merge

import "fmt"

// Strategy defines how conflicts are resolved when a key exists in both maps.
type Strategy int

const (
	// PreferA keeps the value from map A on conflict.
	PreferA Strategy = iota
	// PreferB keeps the value from map B on conflict.
	PreferB
	// ErrorOnConflict returns an error if any key exists in both maps with
	// different values.
	ErrorOnConflict
)

// Options controls the behaviour of the Merge operation.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns a sensible default: prefer values from map B so that
// the "overlay" file wins, which mirrors typical staging→production promotion.
func DefaultOptions() Options {
	return Options{Strategy: PreferB}
}

// Result holds the merged environment variables together with metadata about
// which keys were overridden during the merge.
type Result struct {
	Vars      map[string]string
	Overrides []string // keys whose value was overridden by the winning map
}

// Merge combines envA and envB according to opts and returns a Result.
// Neither input map is mutated.
func Merge(envA, envB map[string]string, opts Options) (Result, error) {
	out := make(map[string]string, len(envA)+len(envB))
	var overrides []string

	for k, v := range envA {
		out[k] = v
	}

	for k, vB := range envB {
		vA, exists := out[k]
		if !exists {
			out[k] = vB
			continue
		}
		if vA == vB {
			continue
		}
		switch opts.Strategy {
		case PreferA:
			// keep vA already in out; record the attempted override
			overrides = append(overrides, k)
		case PreferB:
			out[k] = vB
			overrides = append(overrides, k)
		case ErrorOnConflict:
			return Result{}, fmt.Errorf("merge conflict on key %q: %q vs %q", k, vA, vB)
		}
	}

	return Result{Vars: out, Overrides: overrides}, nil
}
