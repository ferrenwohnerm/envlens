// Package template provides utilities for generating stub .env template files
// from an existing set of environment variables.
//
// A template replaces all values with a configurable placeholder (default:
// "CHANGEME"), making it safe to commit to version control as a reference for
// required configuration keys. Sensitive keys — those whose names contain
// patterns such as PASSWORD, SECRET, or TOKEN — are always redacted regardless
// of the IncludeValues option.
//
// Basic usage:
//
//	vars, _ := parser.ParseFile(".env.production")
//	out := template.Generate(vars, template.DefaultOptions())
//	fmt.Print(out)
//
// To write the template directly to a file:
//
//	err := template.WriteFile(".env.template", vars, template.DefaultOptions())
package template
