package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported env vars.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
)

// Options controls how the export is rendered.
type Options struct {
	Format Format
	Sorted bool
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatDotenv,
		Sorted: true,
	}
}

// Write serialises env vars to w in the requested format.
func Write(vars map[string]string, opts Options, w io.Writer) error {
	keys := keys(vars, opts.Sorted)

	switch opts.Format {
	case FormatDotenv:
		return writeDotenv(vars, keys, w)
	case FormatJSON:
		return writeJSON(vars, keys, w)
	case FormatShell:
		return writeShell(vars, keys, w)
	default:
		return fmt.Errorf("export: unknown format %q", opts.Format)
	}
}

func writeDotenv(vars map[string]string, keys []string, w io.Writer) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, quote(vars[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(vars map[string]string, keys []string, w io.Writer) error {
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = vars[k]
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}

func writeShell(vars map[string]string, keys []string, w io.Writer) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, quote(vars[k])); err != nil {
			return err
		}
	}
	return nil
}

func quote(v string) string {
	if strings.ContainsAny(v, " \t\n#$\"\'\\") {
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return v
}

func keys(vars map[string]string, sorted bool) []string {
	out := make([]string, 0, len(vars))
	for k := range vars {
		out = append(out, k)
	}
	if sorted {
		sort.Strings(out)
	}
	return out
}
