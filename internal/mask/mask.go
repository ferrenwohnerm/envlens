package mask

import "strings"

// sensitivePatterns holds substrings that indicate a key holds sensitive data.
var sensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

// Options controls masking behaviour.
type Options struct {
	// MaskAll replaces every value regardless of key name.
	MaskAll bool
	// Placeholder is the string used in place of a masked value.
	// Defaults to "***" when empty.
	Placeholder string
}

// Apply returns a copy of env with sensitive values replaced by the
// placeholder defined in opts.
func Apply(env map[string]string, opts Options) map[string]string {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.MaskAll || isSensitive(k) {
			out[k] = placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// isSensitive reports whether key contains any known sensitive substring.
func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
