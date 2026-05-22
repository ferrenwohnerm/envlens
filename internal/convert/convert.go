package convert

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format represents a supported output format for environment variable conversion.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
	FormatYAML   Format = "yaml"
)

// Options controls conversion behaviour.
type Options struct {
	SourceFormat Format
	TargetFormat Format
	SortKeys     bool
}

// DefaultOptions returns sensible defaults for Convert.
func DefaultOptions() Options {
	return Options{
		SourceFormat: FormatDotenv,
		TargetFormat: FormatJSON,
		SortKeys:     true,
	}
}

// Convert transforms env vars from one representation to another.
// vars is the already-parsed map; Convert serialises it into the target format.
func Convert(vars map[string]string, opts Options) (string, error) {
	switch opts.TargetFormat {
	case FormatDotenv:
		return toDotenv(vars, opts.SortKeys), nil
	case FormatJSON:
		return toJSON(vars)
	case FormatShell:
		return toShell(vars, opts.SortKeys), nil
	case FormatYAML:
		return toYAML(vars, opts.SortKeys), nil
	default:
		return "", fmt.Errorf("convert: unsupported target format %q", opts.TargetFormat)
	}
}

func toDotenv(vars map[string]string, sorted bool) string {
	var sb strings.Builder
	for _, k := range sortedOrNot(vars, sorted) {
		v := vars[k]
		if strings.ContainsAny(v, " \t\n") {
			fmt.Fprintf(&sb, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return sb.String()
}

func toJSON(vars map[string]string) (string, error) {
	b, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return "", fmt.Errorf("convert: json marshal: %w", err)
	}
	return string(b), nil
}

func toShell(vars map[string]string, sorted bool) string {
	var sb strings.Builder
	for _, k := range sortedOrNot(vars, sorted) {
		fmt.Fprintf(&sb, "export %s=%q\n", k, vars[k])
	}
	return sb.String()
}

func toYAML(vars map[string]string, sorted bool) string {
	var sb strings.Builder
	for _, k := range sortedOrNot(vars, sorted) {
		v := vars[k]
		if strings.ContainsAny(v, ":#{}[]|>&*!,") || strings.TrimSpace(v) != v {
			fmt.Fprintf(&sb, "%s: %q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s: %s\n", k, v)
		}
	}
	return sb.String()
}
