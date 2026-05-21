// Package validate provides schema-based validation for environment variable maps.
//
// A Schema maps key names to Rule definitions. Each Rule can mark a key as
// required and optionally enforce a regex pattern on its value. Keys present
// in the environment but absent from the schema are reported as warnings,
// allowing operators to detect undocumented configuration drift.
//
// Basic usage:
//
//	schema := validate.Schema{
//		"APP_ENV": {Required: true, Pattern: `^(production|staging)$`},
//		"PORT":    {Required: true, Pattern: `^\d+$`},
//	}
//	findings := validate.Run(env, schema)
//	validate.WriteText(os.Stdout, findings)
//
// Output formats:
//   - WriteText: human-readable prefixed lines
//   - WriteJSON: structured JSON array
package validate
