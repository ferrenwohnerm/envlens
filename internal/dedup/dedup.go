package dedup

import "sort"

// DefaultOptions returns a safe default Options value.
func DefaultOptions() Options {
	return Options{
		PreferFirst: true,
	}
}

// Options controls how duplicate keys are resolved when merging
// multiple env maps into one.
type Options struct {
	// PreferFirst keeps the value from the first map that defines a key.
	// When false the last definition wins.
	PreferFirst bool

	// ErrorOnDuplicate returns an error instead of silently resolving
	// conflicts.
	ErrorOnDuplicate bool
}

// Result is returned by Run and carries both the deduplicated map and
// metadata about which keys were duplicated.
type Result struct {
	Vars       map[string]string
	Duplicates []string // sorted list of keys that appeared more than once
}

// Run merges the supplied maps in order, deduplicating keys according
// to opts. Maps are processed left-to-right so index 0 is "first".
func Run(maps []map[string]string, opts Options) (Result, error) {
	out := make(map[string]string)
	seen := make(map[string]bool)
	dupSet := make(map[string]struct{})

	for _, m := range maps {
		for k, v := range m {
			if seen[k] {
				dupSet[k] = struct{}{}
				if opts.ErrorOnDuplicate {
					return Result{}, &DuplicateKeyError{Key: k}
				}
				if !opts.PreferFirst {
					out[k] = v
				}
			} else {
				out[k] = v
				seen[k] = true
			}
		}
	}

	duplicates := make([]string, 0, len(dupSet))
	for k := range dupSet {
		duplicates = append(duplicates, k)
	}
	sort.Strings(duplicates)

	return Result{Vars: out, Duplicates: duplicates}, nil
}

// DuplicateKeyError is returned when ErrorOnDuplicate is set and a
// repeated key is encountered.
type DuplicateKeyError struct {
	Key string
}

func (e *DuplicateKeyError) Error() string {
	return "dedup: duplicate key: " + e.Key
}
