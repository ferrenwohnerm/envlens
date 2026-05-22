// Package resolve expands variable references within environment variable maps.
//
// It supports ${VAR} and $VAR syntax, resolving references against other
// entries in the same map. Optionally, unresolved references can fall back
// to the host process environment, enabling composition of partial configs
// with ambient runtime values.
//
// Example usage:
//
//	vars := map[string]string{
//		"HOST":     "db.internal",
//		"PORT":     "5432",
//		"DATABASE_URL": "postgres://${HOST}:${PORT}/app",
//	}
//
//	resolved, err := resolve.Resolve(vars, resolve.DefaultOptions())
//	// resolved["DATABASE_URL"] == "postgres://db.internal:5432/app"
package resolve
