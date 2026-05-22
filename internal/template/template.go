package template

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Options controls how template generation behaves.
type Options struct {
	// IncludeValues includes actual values as comments next to placeholders.
	IncludeValues bool
	// Placeholder is the string used as the value stub.
	Placeholder string
}

// DefaultOptions returns sensible defaults for template generation.
func DefaultOptions() Options {
	return Options{
		IncludeValues: false,
		Placeholder:   "CHANGEME",
	}
}

// Generate produces a template .env file from the given key-value map.
// Sensitive keys are always redacted; non-sensitive keys may optionally
// retain their original values as inline comments.
func Generate(vars map[string]string, opts Options) string {
	if opts.Placeholder == "" {
		opts.Placeholder = DefaultOptions().Placeholder
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := vars[k]
		if isSensitive(k) {
			fmt.Fprintf(&sb, "%s=%s\n", k, opts.Placeholder)
		} else if opts.IncludeValues {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, opts.Placeholder)
		}
	}
	return sb.String()
}

// WriteFile writes a generated template to the given file path.
func WriteFile(path string, vars map[string]string, opts Options) error {
	content := Generate(vars, opts)
	return os.WriteFile(path, []byte(content), 0o644)
}

// isSensitive returns true if the key name suggests it holds a secret.
func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	sensitivePatterns := []string{
		"PASSWORD", "PASSWD", "SECRET", "TOKEN",
		"API_KEY", "APIKEY", "PRIVATE_KEY", "CREDENTIALS",
		"AUTH", "DSN", "DATABASE_URL",
	}
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
