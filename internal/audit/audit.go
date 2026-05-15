package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/envlens/internal/diff"
)

// Severity represents the risk level of an audit finding.
type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// Finding represents a single audit observation.
type Finding struct {
	Key      string
	Severity Severity
	Message  string
}

// Report holds the full audit output for a diff run.
type Report struct {
	GeneratedAt time.Time
	Findings    []Finding
}

// sensitivePatterns are substrings that indicate a key may hold a secret.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "KEY",
}

// Run analyses a slice of diff results and returns an audit Report.
func Run(results []diff.Result) Report {
	report := Report{GeneratedAt: time.Now()}

	for _, r := range results {
		upper := strings.ToUpper(r.Key)

		switch r.Status {
		case diff.StatusMissingInB:
			sev := SeverityWarning
			if isSensitive(upper) {
				sev = SeverityCritical
			}
			report.Findings = append(report.Findings, Finding{
				Key:      r.Key,
				Severity: sev,
				Message:  fmt.Sprintf("key %q is present in source but missing in target", r.Key),
			})

		case diff.StatusMissingInA:
			report.Findings = append(report.Findings, Finding{
				Key:      r.Key,
				Severity: SeverityInfo,
				Message:  fmt.Sprintf("key %q exists only in target", r.Key),
			})

		case diff.StatusChanged:
			sev := SeverityInfo
			if isSensitive(upper) {
				sev = SeverityWarning
			}
			report.Findings = append(report.Findings, Finding{
				Key:      r.Key,
				Severity: sev,
				Message:  fmt.Sprintf("key %q value differs between source and target", r.Key),
			})
		}
	}

	return report
}

// isSensitive returns true if the uppercased key contains a sensitive pattern.
func isSensitive(upperKey string) bool {
	for _, p := range sensitivePatterns {
		if strings.Contains(upperKey, p) {
			return true
		}
	}
	return false
}
