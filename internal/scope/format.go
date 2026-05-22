package scope

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable summary of the scope result to w.
func WriteText(w io.Writer, r Result) error {
	keys := sortedKeys(r.Vars)
	fmt.Fprintf(w, "Scoped vars (%d):\n", len(keys))
	for _, k := range keys {
		fmt.Fprintf(w, "  %s=%s\n", k, r.Vars[k])
	}
	if len(r.Dropped) > 0 {
		fmt.Fprintf(w, "Dropped (%d):\n", len(r.Dropped))
		for _, k := range r.Dropped {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}
	return nil
}

type jsonOutput struct {
	Vars    map[string]string `json:"vars"`
	Dropped []string          `json:"dropped"`
}

// WriteJSON writes the scope result as a JSON object to w.
func WriteJSON(w io.Writer, r Result) error {
	out := jsonOutput{
		Vars:    r.Vars,
		Dropped: r.Dropped,
	}
	if out.Vars == nil {
		out.Vars = map[string]string{}
	}
	if out.Dropped == nil {
		out.Dropped = []string{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
