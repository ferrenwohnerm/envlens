// Package convert provides utilities for converting parsed environment variable
// maps between different serialisation formats.
//
// Supported target formats:
//
//   - dotenv  — KEY=VALUE pairs, values with whitespace are double-quoted
//   - json    — a JSON object mapping string keys to string values
//   - shell   — export KEY="VALUE" statements suitable for sourcing in bash/zsh
//   - yaml    — simple YAML mapping (key: value), special characters are quoted
//
// Basic usage:
//
//	vars, _ := parser.ParseFile(".env.production")
//	out, err := convert.Convert(vars, convert.Options{
//		TargetFormat: convert.FormatJSON,
//		SortKeys:     true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Print(out)
package convert
