// Package mask provides utilities for redacting sensitive values from
// environment variable maps before they are printed or stored.
//
// Sensitivity is determined by matching well-known substrings (e.g. PASSWORD,
// TOKEN, API_KEY) against the upper-cased key name. Callers may also opt-in to
// masking every value unconditionally via Options.MaskAll.
//
// Example usage:
//
//	masked := mask.Apply(env, mask.Options{
//		Placeholder: "<redacted>",
//	})
//
// The original map is never modified; Apply always returns a new map.
package mask
