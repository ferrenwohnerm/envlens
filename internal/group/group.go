// Package group provides utilities for grouping environment variable
// keys by a shared prefix, delimiter, or custom classifier function.
package group

import (
	"sort"
	"strings"
)

// Options controls how grouping is performed.
type Options struct {
	// Delimiter separates the group prefix from the rest of the key.
	// Defaults to "_" if empty.
	Delimiter string

	// MaxDepth is the number of delimiter-separated segments used to
	// form the group name. 1 means only the first segment.
	// Defaults to 1.
	MaxDepth int

	// Ungrouped is the label used for keys that have no delimiter.
	// Defaults to "(ungrouped)".
	Ungrouped string
}

// DefaultOptions returns an Options struct with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Delimiter: "_",
		MaxDepth:  1,
		Ungrouped: "(ungrouped)",
	}
}

// Result holds the grouped output.
type Result struct {
	// Groups maps a group name to the keys (and their values) that
	// belong to it.
	Groups map[string]map[string]string

	// Order is a sorted slice of group names for deterministic output.
	Order []string
}

// Run partitions env into groups according to opts.
func Run(env map[string]string, opts Options) Result {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = 1
	}
	if opts.Ungrouped == "" {
		opts.Ungrouped = "(ungrouped)"
	}

	groups := make(map[string]map[string]string)

	for k, v := range env {
		label := groupLabel(k, opts.Delimiter, opts.MaxDepth, opts.Ungrouped)
		if groups[label] == nil {
			groups[label] = make(map[string]string)
		}
		groups[label][k] = v
	}

	order := make([]string, 0, len(groups))
	for g := range groups {
		order = append(order, g)
	}
	sort.Strings(order)

	return Result{Groups: groups, Order: order}
}

// groupLabel derives the group name for a key.
func groupLabel(key, delimiter string, maxDepth int, ungrouped string) string {
	parts := strings.Split(key, delimiter)
	if len(parts) <= 1 {
		return ungrouped
	}
	end := maxDepth
	if end > len(parts)-1 {
		end = len(parts) - 1
	}
	return strings.Join(parts[:end], delimiter)
}
