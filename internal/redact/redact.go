package redact

import (
	"regexp"
	"strings"
)

// Rule defines a single redaction rule applied to env var values.
type Rule struct {
	// Pattern is a regex matched against the key name.
	Pattern *regexp.Regexp
	// Replacement is the string used in place of the matched value.
	Replacement string
}

// Options controls the behaviour of Apply.
type Options struct {
	// Rules is an ordered list of redaction rules.
	Rules []Rule
	// DefaultReplacement is used when a key matches no rule but MaskAll is set.
	DefaultReplacement string
	// MaskAll redacts every value regardless of rules.
	MaskAll bool
}

// DefaultOptions returns sensible defaults that redact common secret patterns.
func DefaultOptions() Options {
	return Options{
		DefaultReplacement: "[REDACTED]",
		Rules: []Rule{
			{Pattern: regexp.MustCompile(`(?i)(password|passwd|secret|token|apikey|api_key|private_key|auth)`), Replacement: "[REDACTED]"},
			{Pattern: regexp.MustCompile(`(?i)(dsn|database_url|connection_string)`), Replacement: "[REDACTED-DSN]"},
		},
	}
}

// Apply returns a copy of vars with values redacted according to opts.
// Keys are matched case-insensitively against each rule's Pattern in order;
// the first match wins. If MaskAll is true every value is replaced with
// DefaultReplacement regardless of rules.
func Apply(vars map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = redactValue(k, v, opts)
	}
	return out
}

func redactValue(key, value string, opts Options) string {
	if opts.MaskAll {
		if opts.DefaultReplacement == "" {
			return "[REDACTED]"
		}
		return opts.DefaultReplacement
	}
	upper := strings.ToUpper(key)
	for _, rule := range opts.Rules {
		if rule.Pattern.MatchString(upper) {
			if rule.Replacement == "" {
				return "[REDACTED]"
			}
			return rule.Replacement
		}
	}
	return value
}
