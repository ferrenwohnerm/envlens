package inject

import (
	"fmt"
	"os"
	"sort"
)

// Options controls the behaviour of the Inject operation.
type Options struct {
	// Overwrite allows existing environment variables in the process to be
	// replaced by values from the provided map.
	Overwrite bool

	// DryRun skips calling os.Setenv and returns the planned changes only.
	DryRun bool

	// Keys restricts injection to the listed keys. When empty all keys are
	// injected.
	Keys []string
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		DryRun:    false,
	}
}

// Result describes the outcome of a single key injection.
type Result struct {
	Key     string
	Value   string
	Action  string // "set", "skipped", "overwritten"
}

// Inject applies the supplied env map to the current process environment and
// returns a slice of Results describing what happened for each key.
func Inject(env map[string]string, opts Options) ([]Result, error) {
	allowSet := buildAllowSet(opts.Keys)

	keys := sortedKeys(env)
	results := make([]Result, 0, len(keys))

	for _, k := range keys {
		if len(allowSet) > 0 {
			if _, ok := allowSet[k]; !ok {
				continue
			}
		}

		v := env[k]
		existing, exists := os.LookupEnv(k)
		_ = existing

		var action string
		switch {
		case !exists:
			action = "set"
		case opts.Overwrite:
			action = "overwritten"
		default:
			results = append(results, Result{Key: k, Value: v, Action: "skipped"})
			continue
		}

		if !opts.DryRun {
			if err := os.Setenv(k, v); err != nil {
				return results, fmt.Errorf("inject: failed to set %q: %w", k, err)
			}
		}

		results = append(results, Result{Key: k, Value: v, Action: action})
	}

	return results, nil
}

func buildAllowSet(keys []string) map[string]struct{} {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
