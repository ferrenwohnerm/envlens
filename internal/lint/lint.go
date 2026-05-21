package lint

import (
	"fmt"
	"strings"
)

// Rule represents a linting rule applied to environment variable keys or values.
type Rule string

const (
	RuleUppercaseKeys    Rule = "uppercase-keys"
	RuleNoEmptyValues    Rule = "no-empty-values"
	RuleNoWhitespaceKeys Rule = "no-whitespace-keys"
	RuleNoDuplicateKeys  Rule = "no-duplicate-keys"
)

// Finding represents a single lint violation.
type Finding struct {
	Rule    Rule
	Key     string
	Message string
}

// Options controls which rules are applied during linting.
type Options struct {
	UppercaseKeys    bool
	NoEmptyValues    bool
	NoWhitespaceKeys bool
	NoDuplicateKeys  bool
}

// DefaultOptions returns an Options struct with all rules enabled.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:    true,
		NoEmptyValues:    true,
		NoWhitespaceKeys: true,
		NoDuplicateKeys:  true,
	}
}

// Run applies lint rules to the provided environment map and returns any findings.
func Run(env map[string]string, opts Options) []Finding {
	var findings []Finding
	seen := make(map[string]int)

	for key, value := range env {
		seen[key]++

		if opts.UppercaseKeys && key != strings.ToUpper(key) {
			findings = append(findings, Finding{
				Rule:    RuleUppercaseKeys,
				Key:     key,
				Message: fmt.Sprintf("key %q should be uppercase (got %q)", strings.ToUpper(key), key),
			})
		}

		if opts.NoEmptyValues && strings.TrimSpace(value) == "" {
			findings = append(findings, Finding{
				Rule:    RuleNoEmptyValues,
				Key:     key,
				Message: fmt.Sprintf("key %q has an empty or blank value", key),
			})
		}

		if opts.NoWhitespaceKeys && strings.ContainsAny(key, " \t") {
			findings = append(findings, Finding{
				Rule:    RuleNoWhitespaceKeys,
				Key:     key,
				Message: fmt.Sprintf("key %q contains whitespace characters", key),
			})
		}
	}

	if opts.NoDuplicateKeys {
		for key, count := range seen {
			if count > 1 {
				findings = append(findings, Finding{
					Rule:    RuleNoDuplicateKeys,
					Key:     key,
					Message: fmt.Sprintf("key %q appears %d times", key, count),
				})
			}
		}
	}

	return findings
}
