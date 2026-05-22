package resolve

import (
	"fmt"
	"os"
	"strings"
)

// Options controls how variable resolution behaves.
type Options struct {
	// FallbackToEnv allows falling back to the process environment
	// when a variable reference is not found in the provided map.
	FallbackToEnv bool

	// ErrorOnMissing causes Resolve to return an error if a referenced
	// variable cannot be resolved from any source.
	ErrorOnMissing bool
}

// DefaultOptions returns sensible defaults for resolution.
func DefaultOptions() Options {
	return Options{
		FallbackToEnv:  false,
		ErrorOnMissing: false,
	}
}

// Resolve expands ${VAR} and $VAR references within values of the given map
// using other entries in the same map, and optionally the process environment.
// It returns a new map with all resolvable references expanded.
func Resolve(vars map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(vars))

	for k, v := range vars {
		expanded, err := expandValue(v, vars, opts)
		if err != nil {
			return nil, fmt.Errorf("resolving key %q: %w", k, err)
		}
		result[k] = expanded
	}

	return result, nil
}

// expandValue replaces all variable references in s with their resolved values.
func expandValue(s string, vars map[string]string, opts Options) (string, error) {
	var missingKeys []string

	result := os.Expand(s, func(key string) string {
		if val, ok := vars[key]; ok {
			return val
		}
		if opts.FallbackToEnv {
			if val, ok := os.LookupEnv(key); ok {
				return val
			}
		}
		missingKeys = append(missingKeys, key)
		return ""
	})

	if opts.ErrorOnMissing && len(missingKeys) > 0 {
		return "", fmt.Errorf("unresolved references: %s", strings.Join(missingKeys, ", "))
	}

	return result, nil
}
