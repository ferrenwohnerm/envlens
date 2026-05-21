package validate

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

func severityPrefix(s string) string {
	switch s {
	case "error":
		return "[ERROR]  "
	case "warning":
		return "[WARN]   "
	default:
		return "[INFO]   "
	}
}

// WriteText writes findings in a human-readable text format.
func WriteText(w io.Writer, findings []Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "schema validation passed: no findings")
		return err
	}
	sorted := sortedFindings(findings)
	for _, f := range sorted {
		_, err := fmt.Fprintf(w, "%s%-30s %s\n", severityPrefix(f.Severity), f.Key, f.Message)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w, Summary(findings))
	return err
}

// WriteJSON writes findings as a JSON array.
func WriteJSON(w io.Writer, findings []Finding) error {
	type jsonFinding struct {
		Key      string `json:"key"`
		Message  string `json:"message"`
		Severity string `json:"severity"`
	}
	out := make([]jsonFinding, len(findings))
	for i, f := range findings {
		out[i] = jsonFinding{Key: f.Key, Message: f.Message, Severity: f.Severity}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func sortedFindings(findings []Finding) []Finding {
	copy := append([]Finding(nil), findings...)
	sort.Slice(copy, func(i, j int) bool {
		if copy[i].Severity != copy[j].Severity {
			return copy[i].Severity < copy[j].Severity
		}
		return copy[i].Key < copy[j].Key
	})
	return copy
}
