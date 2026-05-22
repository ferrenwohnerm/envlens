package scope

import "sort"

// Options controls how scoping is applied.
type Options struct {
	// Prefixes restricts the result to keys matching any of the given prefixes.
	Prefixes []string
	// StripPrefix removes the matched prefix from the resulting key names.
	StripPrefix bool
	// CaseFold performs case-insensitive prefix matching.
	CaseFold bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		StripPrefix: false,
		CaseFold:    false,
	}
}

// Result holds the output of a scope operation.
type Result struct {
	// Vars contains the scoped key/value pairs.
	Vars map[string]string
	// Dropped lists keys that were excluded by the scope filter.
	Dropped []string
}

// Run filters env to only the keys matching opts.Prefixes.
// If opts.Prefixes is empty, all keys are retained.
func Run(env map[string]string, opts Options) Result {
	if len(opts.Prefixes) == 0 {
		vars := make(map[string]string, len(env))
		for k, v := range env {
			vars[k] = v
		}
		return Result{Vars: vars}
	}

	vars := make(map[string]string)
	var dropped []string

	for k, v := range env {
		matched, prefix := matchesAny(k, opts.Prefixes, opts.CaseFold)
		if !matched {
			dropped = append(dropped, k)
			continue
		}
		outKey := k
		if opts.StripPrefix && prefix != "" {
			outKey = k[len(prefix):]
		}
		vars[outKey] = v
	}

	sort.Strings(dropped)
	return Result{Vars: vars, Dropped: dropped}
}

func matchesAny(key string, prefixes []string, caseFold bool) (bool, string) {
	candidate := key
	if caseFold {
		candidate = toLower(key)
	}
	for _, p := range prefixes {
		pfx := p
		if caseFold {
			pfx = toLower(p)
		}
		if len(candidate) >= len(pfx) && candidate[:len(pfx)] == pfx {
			return true, p
		}
	}
	return false, ""
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
