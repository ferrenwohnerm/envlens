package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Schema defines expected keys and optional validation rules.
type Schema map[string]Rule

// Rule describes constraints for a single environment variable.
type Rule struct {
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Finding represents a single schema validation result.
type Finding struct {
	Key      string
	Message  string
	Severity string // "error" | "warning"
}

// Run validates env vars against the provided schema.
// It reports missing required keys and pattern mismatches.
func Run(env map[string]string, schema Schema) []Finding {
	var findings []Finding

	for key, rule := range schema {
		val, exists := env[key]

		if rule.Required && !exists {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "required key is missing",
				Severity: "error",
			})
			continue
		}

		if exists && rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, val)
			if err != nil {
				findings = append(findings, Finding{
					Key:      key,
					Message:  fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
					Severity: "warning",
				})
				continue
			}
			if !matched {
				findings = append(findings, Finding{
					Key:      key,
					Message:  fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
					Severity: "error",
				})
			}
		}
	}

	// Warn about keys present in env but not defined in schema.
	for key := range env {
		if _, defined := schema[key]; !defined {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key not defined in schema",
				Severity: "warning",
			})
		}
	}

	return findings
}

// Summary returns a human-readable summary line.
func Summary(findings []Finding) string {
	if len(findings) == 0 {
		return "schema validation passed: no findings"
	}
	errors, warnings := 0, 0
	for _, f := range findings {
		if strings.ToLower(f.Severity) == "error" {
			errors++
		} else {
			warnings++
		}
	}
	return fmt.Sprintf("schema validation: %d error(s), %d warning(s)", errors, warnings)
}
