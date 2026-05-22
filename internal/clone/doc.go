// Package clone provides functionality for copying environment variables
// from one map (source) into another (destination), with optional prefix
// renaming, overwrite control, and dry-run support.
//
// A common use-case is promoting a staging environment's variables into a
// production template while rewriting key prefixes:
//
//	opts := clone.DefaultOptions()
//	opts.PrefixReplace = map[string]string{"STAGING_": "PROD_"}
//	opts.DryRun = true
//
//	preview, err := clone.Clone(stagingVars, prodVars, opts)
//
// Use BuildSummary to understand what was added, replaced, or skipped, and
// WriteText / WriteJSON to render that summary for the user.
package clone
