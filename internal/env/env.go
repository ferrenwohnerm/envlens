package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		Overwrite:    false,
		FailOnMissing: false,
	}
}

// Options controls the behaviour of Inject and Extract.
type Options struct {
	// Overwrite allows existing process environment variables to be replaced.
	Overwrite bool
	// FailOnMissing returns an error when a requested key is absent from the map.
	FailOnMissing bool
}

// Inject writes the key/value pairs in vars into the process environment.
// Keys already present in the environment are skipped unless Overwrite is set.
func Inject(vars map[string]string, opts Options) error {
	for k, v := range vars {
		if !opts.Overwrite {
			if _, exists := os.LookupEnv(k); exists {
				continue
			}
		}
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("env: setenv %q: %w", k, err)
		}
	}
	return nil
}

// Extract reads the keys listed in keys from the process environment and
// returns them as a map. If FailOnMissing is set, any absent key causes an
// error.
func Extract(keys []string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		v, ok := os.LookupEnv(k)
		if !ok && opts.FailOnMissing {
			return nil, fmt.Errorf("env: key %q not found in process environment", k)
		}
		if ok {
			out[k] = v
		}
	}
	return out, nil
}

// Snapshot captures all current process environment variables whose keys
// match prefix (empty prefix matches all).
func Snapshot(prefix string) map[string]string {
	out := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if prefix == "" || strings.HasPrefix(parts[0], prefix) {
			out[parts[0]] = parts[1]
		}
	}
	return out
}

// SortedKeys returns the keys of vars in ascending order.
func SortedKeys(vars map[string]string) []string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
