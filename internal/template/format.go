package template

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Format represents the output format for template metadata.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// KeyEntry describes a single key in the template output.
type KeyEntry struct {
	Key       string `json:"key"`
	Sensitive bool   `json:"sensitive"`
	Value     string `json:"value"`
}

// WriteMetadata writes a human-readable or JSON summary of which keys are
// present in the template and whether they are considered sensitive.
func WriteMetadata(w io.Writer, vars map[string]string, opts Options, format Format) error {
	entries := buildEntries(vars, opts)
	switch format {
	case FormatJSON:
		return writeMetadataJSON(w, entries)
	default:
		return writeMetadataText(w, entries)
	}
}

func buildEntries(vars map[string]string, opts Options) []KeyEntry {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]KeyEntry, 0, len(keys))
	for _, k := range keys {
		sens := isSensitive(k)
		val := opts.Placeholder
		if !sens && opts.IncludeValues {
			val = vars[k]
		}
		entries = append(entries, KeyEntry{Key: k, Sensitive: sens, Value: val})
	}
	return entries
}

func writeMetadataText(w io.Writer, entries []KeyEntry) error {
	for _, e := range entries {
		sensTag := ""
		if e.Sensitive {
			sensTag = " [sensitive]"
		}
		if _, err := fmt.Fprintf(w, "%-30s = %-20s%s\n", e.Key, e.Value, sensTag); err != nil {
			return err
		}
	}
	return nil
}

func writeMetadataJSON(w io.Writer, entries []KeyEntry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
