// Package redact provides value-level redaction for environment variable maps.
//
// It applies an ordered list of regex-based Rules to key names; the first
// matching rule's Replacement string is substituted for the original value.
// A MaskAll flag is available to unconditionally redact every entry, which is
// useful when producing logs or audit trails that must never contain secrets.
//
// Typical usage:
//
//	opts := redact.DefaultOptions()
//	safe := redact.Apply(vars, opts)
//
// DefaultOptions ships with rules that match common secret key patterns such
// as PASSWORD, TOKEN, API_KEY, SECRET, and DATABASE_URL.
package redact
