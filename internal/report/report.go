// Package report provides formatting and output utilities for
// displaying environment variable diff results to the user.
package report

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envlens/internal/diff"
)

// Format controls how the report is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Write renders a diff result to the given writer in the specified format.
func Write(w io.Writer, result diff.Result, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, result)
	default:
		return writeText(w, result)
	}
}

func writeText(w io.Writer, result diff.Result) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	keys := collectKeys(result)
	sort.Strings(keys)

	for _, k := range keys {
		if entry, ok := result[k]; ok {
			switch entry.Status {
			case diff.OnlyInA:
				fmt.Fprintf(tw, "- %s\t%s\t(only in A)\n", k, entry.ValueA)
			case diff.OnlyInB:
				fmt.Fprintf(tw, "+ %s\t%s\t(only in B)\n", k, entry.ValueB)
			case diff.Changed:
				fmt.Fprintf(tw, "~ %s\t%s -> %s\n", k, entry.ValueA, entry.ValueB)
			case diff.Identical:
				fmt.Fprintf(tw, "  %s\t%s\n", k, entry.ValueA)
			}
		}
	}

	return tw.Flush()
}

func writeJSON(w io.Writer, result diff.Result) error {
	keys := collectKeys(result)
	sort.Strings(keys)

	fmt.Fprintln(w, "{")
	for i, k := range keys {
		entry := result[k]
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(w, `  %q: {"status": %q, "value_a": %q, "value_b": %q}%s`+"\n",
			k, entry.Status, entry.ValueA, entry.ValueB, comma)
	}
	fmt.Fprintln(w, "}")
	return nil
}

func collectKeys(result diff.Result) []string {
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	return keys
}
