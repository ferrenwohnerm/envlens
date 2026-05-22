package promote

import (
	"fmt"
	"sort"
)

// Options controls the behaviour of the Promote operation.
type Options struct {
	// DryRun reports what would change without modifying dst.
	DryRun bool
	// SkipKeys is a set of keys that must never be copied from src to dst.
	SkipKeys map[string]bool
	// OverwriteExisting allows values already present in dst to be replaced.
	OverwriteExisting bool
}

// DefaultOptions returns a safe, non-destructive Options value.
func DefaultOptions() Options {
	return Options{
		DryRun:            false,
		SkipKeys:          map[string]bool{},
		OverwriteExisting: false,
	}
}

// Change describes a single key-level change produced by Promote.
type Change struct {
	Key      string
	OldValue string // empty when the key did not exist in dst
	NewValue string
	Action   string // "add" | "overwrite" | "skip"
}

// Result is returned by Promote and summarises what happened (or would happen).
type Result struct {
	Changes []Change
}

// Promote copies keys from src into dst according to opts.
// When DryRun is true dst is never modified; Changes still reflects what
// would have occurred.
func Promote(src, dst map[string]string, opts Options) (Result, error) {
	if src == nil {
		return Result{}, fmt.Errorf("promote: src map must not be nil")
	}
	if dst == nil {
		return Result{}, fmt.Errorf("promote: dst map must not be nil")
	}

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var changes []Change
	for _, k := range keys {
		if opts.SkipKeys[k] {
			changes = append(changes, Change{Key: k, NewValue: src[k], Action: "skip"})
			continue
		}

		old, exists := dst[k]
		switch {
		case !exists:
			if !opts.DryRun {
				dst[k] = src[k]
			}
			changes = append(changes, Change{Key: k, NewValue: src[k], Action: "add"})
		case opts.OverwriteExisting && old != src[k]:
			if !opts.DryRun {
				dst[k] = src[k]
			}
			changes = append(changes, Change{Key: k, OldValue: old, NewValue: src[k], Action: "overwrite"})
		default:
			// key exists and overwrite is disabled — treat as skip
			changes = append(changes, Change{Key: k, OldValue: old, NewValue: src[k], Action: "skip"})
		}
	}

	return Result{Changes: changes}, nil
}
