// Package scope provides functionality for narrowing an environment variable
// map to a specific subset defined by key prefixes.
//
// It is useful when working with multi-service .env files where each service
// owns a distinct prefix (e.g. APP_, DB_, CACHE_). Scoping lets callers
// extract only the relevant slice for further diffing, auditing, or exporting.
//
// Basic usage:
//
//	opts := scope.Options{
//		Prefixes:    []string{"APP_"},
//		StripPrefix: true,
//	}
//	result := scope.Run(env, opts)
//	// result.Vars contains keys with "APP_" stripped
//	// result.Dropped lists all excluded keys
package scope
