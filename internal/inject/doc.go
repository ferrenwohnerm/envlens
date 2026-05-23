// Package inject provides functionality for applying a map of environment
// variables to the current process environment.
//
// It supports fine-grained control over which keys are injected, whether
// existing values may be overwritten, and a dry-run mode that plans changes
// without modifying the process environment.
//
// Typical usage:
//
//	env, _ := parser.ParseFile(".env.staging")
//	opts := inject.DefaultOptions()
//	opts.Overwrite = true
//	results, err := inject.Inject(env, opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, r := range results {
//		fmt.Printf("%s: %s\n", r.Key, r.Action)
//	}
package inject
