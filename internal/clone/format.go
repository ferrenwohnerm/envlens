package clone

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Summary holds a human-readable account of what a clone operation did.
type Summary struct {
	Added    []string `json:"added"`
	Skipped  []string `json:"skipped"`
	Replaced []string `json:"replaced"`
}

// BuildSummary compares before and after maps to produce a Summary.
func BuildSummary(before, after map[string]string, src map[string]string, opts Options) Summary {
	s := Summary{}
	for srcKey := range src {
		destKey := applyPrefixReplace(srcKey, opts.PrefixReplace)
		_, existed := before[destKey]
		if !existed {
			s.Added = append(s.Added, destKey)
		} else if opts.Overwrite {
			s.Replaced = append(s.Replaced, destKey)
		} else {
			s.Skipped = append(s.Skipped, destKey)
		}
	}
	sort.Strings(s.Added)
	sort.Strings(s.Skipped)
	sort.Strings(s.Replaced)
	return s
}

// WriteText writes a human-readable clone summary to w.
func WriteText(w io.Writer, s Summary) {
	for _, k := range s.Added {
		fmt.Fprintf(w, "+ %s\n", k)
	}
	for _, k := range s.Replaced {
		fmt.Fprintf(w, "~ %s\n", k)
	}
	for _, k := range s.Skipped {
		fmt.Fprintf(w, "= %s (skipped)\n", k)
	}
}

// WriteJSON writes the clone summary as JSON to w.
func WriteJSON(w io.Writer, s Summary) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
