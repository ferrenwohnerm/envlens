// Package lint provides static analysis rules for environment variable files.
//
// It inspects the keys and values of a parsed env map and reports violations
// according to configurable rules. Rules include:
//
//   - uppercase-keys: all keys must be fully uppercase
//   - no-empty-values: values must not be empty or whitespace-only
//   - no-whitespace-keys: keys must not contain spaces or tabs
//   - no-duplicate-keys: keys must appear only once in the source file
//
// Example usage:
//
//	env, _ := parser.ParseFile("staging.env")
//	findings := lint.Run(env, lint.DefaultOptions())
//	for _, f := range findings {
//		fmt.Printf("[%s] %s: %s\n", f.Rule, f.Key, f.Message)
//	}
package lint
