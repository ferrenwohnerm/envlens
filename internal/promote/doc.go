// Package promote implements the promotion workflow for envlens.
//
// Promotion copies environment variables from a source environment (e.g.
// staging) into a destination environment (e.g. production), with fine-grained
// control over which keys are allowed to be added or overwritten.
//
// # Basic usage
//
//	src, _ := parser.ParseFile("staging.env")
//	dst, _ := parser.ParseFile("production.env")
//
//	opts := promote.DefaultOptions()
//	opts.OverwriteExisting = false
//	opts.SkipKeys = map[string]bool{"DATABASE_URL": true}
//
//	result, err := promote.Promote(src, dst, opts)
//
// # Dry-run mode
//
// Set Options.DryRun = true to inspect what would change without modifying the
// destination map.  The returned Result.Changes slice still reflects every
// intended action.
//
// # Actions
//
// Each Change in the result carries one of three action labels:
//
//   - "add"       – key was absent in dst and has been (or would be) inserted.
//   - "overwrite" – key existed in dst with a different value and OverwriteExisting was true.
//   - "skip"      – key was excluded via SkipKeys or already present without OverwriteExisting.
package promote
