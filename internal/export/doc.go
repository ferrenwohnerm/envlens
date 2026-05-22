// Package export renders a map of environment variables to various
// output formats suitable for consumption by shells, CI systems, or
// configuration management tools.
//
// Supported formats:
//
//   - dotenv  – standard KEY=VALUE lines, values quoted when necessary
//   - json    – pretty-printed JSON object
//   - shell   – export KEY=VALUE lines ready for sourcing in bash/zsh
//
// Example usage:
//
//	opts := export.DefaultOptions()
//	opts.Format = export.FormatShell
//	if err := export.Write(vars, opts, os.Stdout); err != nil {
//		log.Fatal(err)
//	}
package export
