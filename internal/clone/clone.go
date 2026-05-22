package clone

import (
	"fmt"
	"strings"
)

// Options controls how a clone operation behaves.
type Options struct {
	// PrefixReplace maps an old prefix to a new prefix for all cloned keys.
	// e.g. {"STAGING_": "PROD_"}
	PrefixReplace map[string]string

	// Overwrite allows existing keys in Dst to be replaced.
	Overwrite bool

	// DryRun returns the would-be result without modifying Dst.
	DryRun bool
}

// DefaultOptions returns a safe, non-destructive Options.
func DefaultOptions() Options {
	return Options{
		PrefixReplace: map[string]string{},
		Overwrite:     false,
		DryRun:        false,
	}
}

// Clone copies keys from src into dst, optionally renaming prefixes.
// When DryRun is true the returned map is a preview and dst is not mutated.
func Clone(src, dst map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(dst))
	for k, v := range dst {
		result[k] = v
	}

	for srcKey, srcVal := range src {
		destKey := applyPrefixReplace(srcKey, opts.PrefixReplace)

		if _, exists := result[destKey]; exists && !opts.Overwrite {
			continue
		}

		result[destKey] = srcVal
	}

	if opts.DryRun {
		return result, nil
	}

	for k, v := range result {
		dst[k] = v
	}

	return result, nil
}

func applyPrefixReplace(key string, replacements map[string]string) string {
	for old, newPrefix := range replacements {
		if strings.HasPrefix(key, old) {
			return fmt.Sprintf("%s%s", newPrefix, strings.TrimPrefix(key, old))
		}
	}
	return key
}
