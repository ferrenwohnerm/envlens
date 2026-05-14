// Package parser provides utilities for reading and parsing .env files
// used by the envlens CLI tool.
//
// Supported formats:
//
//	# Comment lines are ignored
//	KEY=VALUE
//	KEY="quoted value"
//	KEY='single quoted value'
//
// Blank lines and lines beginning with '#' are silently skipped.
// Each file is parsed into an EnvMap (map[string]string) which can
// then be passed to the diff engine for comparison across environments.
//
// Example usage:
//
//	staging, err := parser.ParseFile(".env.staging")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	production, err := parser.ParseFile(".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
package parser
