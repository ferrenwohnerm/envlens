package flatten

import (
	"fmt"
	"strings"
)

// DefaultOptions returns a safe default configuration for flattening nested
// key structures into a flat env-style map.
func DefaultOptions() Options {
	return Options{
		Delimiter:    "_",
		Prefix:       "",
		UppercaseKeys: true,
	}
}

// Options controls the behaviour of the Flatten operation.
type Options struct {
	// Delimiter is placed between segments when joining nested keys.
	Delimiter string

	// Prefix is prepended to every output key.
	Prefix string

	// UppercaseKeys converts all output keys to upper-case.
	UppercaseKeys bool
}

// Run takes a map that may contain dot-separated or delimiter-separated keys
// (e.g. "db.host", "db.port") and returns a new flat map where each key is
// normalised according to opts.
//
// Duplicate output keys (after normalisation) are rejected with an error so
// that callers are never silently surprised by a lost value.
func Run(env map[string]string, opts Options) (map[string]string, error) {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}

	out := make(map[string]string, len(env))

	for k, v := range env {
		flat := flattenKey(k, opts)
		if _, exists := out[flat]; exists {
			return nil, fmt.Errorf("flatten: duplicate output key %q (from input key %q)", flat, k)
		}
		out[flat] = v
	}

	return out, nil
}

func flattenKey(key string, opts Options) string {
	// Normalise any dots or forward-slashes to the chosen delimiter.
	replacer := strings.NewReplacer(
		".", opts.Delimiter,
		"/", opts.Delimiter,
	)
	result := replacer.Replace(key)

	if opts.UppercaseKeys {
		result = strings.ToUpper(result)
	}

	if opts.Prefix != "" {
		result = opts.Prefix + result
	}

	return result
}
