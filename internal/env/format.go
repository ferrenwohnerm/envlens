package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// WriteText writes vars to w in a human-readable key=value format, one entry
// per line, keys sorted alphabetically.
func WriteText(w io.Writer, vars map[string]string) error {
	for _, k := range SortedKeys(vars) {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, vars[k]); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes vars to w as a JSON object with sorted keys.
func WriteJSON(w io.Writer, vars map[string]string) error {
	// Use a sorted slice of entries so output is deterministic.
	type kv struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	entries := make([]kv, 0, len(vars))
	for _, k := range SortedKeys(vars) {
		entries = append(entries, kv{Key: k, Value: vars[k]})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

// WriteSummary writes a one-line summary of the variable map to w.
func WriteSummary(w io.Writer, vars map[string]string) error {
	keys := SortedKeys(vars)
	_, err := fmt.Fprintf(w, "%d variable(s): %s\n", len(keys), strings.Join(keys, ", "))
	return err
}
