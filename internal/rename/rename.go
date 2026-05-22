// Package rename provides utilities for renaming or remapping environment
// variable keys across one or more env maps. It supports bulk prefix
// substitution, explicit key mappings, and dry-run inspection.
package rename

import "fmt"

// Options controls the behaviour of the Rename operation.
type Options struct {
	// Mapping is an explicit old-key → new-key map.
	Mapping map[string]string

	// OldPrefix and NewPrefix, when both non-empty, replace a leading prefix
	// on every matching key.
	OldPrefix string
	NewPrefix string

	// ErrorOnMissing causes Rename to return an error when an explicit mapping
	// key is not present in the source env.
	ErrorOnMissing bool
}

// Result describes a single rename operation that was applied (or would be
// applied in a dry-run).
type Result struct {
	OldKey string
	NewKey string
	Value  string
}

// Rename applies the rename rules described by opts to env and returns the
// transformed map together with a log of every rename that occurred.
// The original map is never modified.
func Rename(env map[string]string, opts Options) (map[string]string, []Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var results []Result

	// Explicit key mapping takes priority.
	for oldKey, newKey := range opts.Mapping {
		val, ok := out[oldKey]
		if !ok {
			if opts.ErrorOnMissing {
				return nil, nil, fmt.Errorf("rename: key %q not found in env", oldKey)
			}
			continue
		}
		if oldKey == newKey {
			continue
		}
		delete(out, oldKey)
		out[newKey] = val
		results = append(results, Result{OldKey: oldKey, NewKey: newKey, Value: val})
	}

	// Prefix substitution (applied after explicit mapping).
	if opts.OldPrefix != "" && opts.NewPrefix != "" {
		for k, v := range out {
			if len(k) >= len(opts.OldPrefix) && k[:len(opts.OldPrefix)] == opts.OldPrefix {
				newKey := opts.NewPrefix + k[len(opts.OldPrefix):]
				if newKey == k {
					continue
				}
				delete(out, k)
				out[newKey] = v
				results = append(results, Result{OldKey: k, NewKey: newKey, Value: v})
			}
		}
	}

	return out, results, nil
}
