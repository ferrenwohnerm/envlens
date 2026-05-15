package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// severityOrder defines the display / sort order for severity levels.
var severityOrder = map[Severity]int{
	SeverityCritical: 0,
	SeverityWarning:  1,
	SeverityInfo:     2,
}

// WriteText writes a human-readable audit report to w.
func WriteText(w io.Writer, r Report) error {
	findings := sortedFindings(r.Findings)

	fmt.Fprintf(w, "Audit Report — %s\n", r.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Total findings: %d\n\n", len(findings))

	if len(findings) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return nil
	}

	for _, f := range findings {
		prefix := severityPrefix(f.Severity)
		fmt.Fprintf(w, "%s %s\n", prefix, f.Message)
	}
	return nil
}

// WriteJSON writes the audit report as a JSON object to w.
func WriteJSON(w io.Writer, r Report) error {
	type jsonFinding struct {
		Key      string `json:"key"`
		Severity string `json:"severity"`
		Message  string `json:"message"`
	}
	type jsonReport struct {
		GeneratedAt string       `json:"generated_at"`
		Total       int          `json:"total_findings"`
		Findings    []jsonFinding `json:"findings"`
	}

	findings := sortedFindings(r.Findings)
	out := jsonReport{
		GeneratedAt: r.GeneratedAt.Format("2006-01-02T15:04:05Z"),
		Total:       len(findings),
	}
	for _, f := range findings {
		out.Findings = append(out.Findings, jsonFinding{
			Key:      f.Key,
			Severity: string(f.Severity),
			Message:  f.Message,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func severityPrefix(s Severity) string {
	switch s {
	case SeverityCritical:
		return "[CRITICAL]"
	case SeverityWarning:
		return "[WARNING] "
	default:
		return "[INFO]    "
	}
}

func sortedFindings(findings []Finding) []Finding {
	out := make([]Finding, len(findings))
	copy(out, findings)
	sort.Slice(out, func(i, j int) bool {
		oi := severityOrder[out[i].Severity]
		oj := severityOrder[out[j].Severity]
		if oi != oj {
			return oi < oj
		}
		return out[i].Key < out[j].Key
	})
	return out
}
