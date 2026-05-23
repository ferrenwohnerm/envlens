// Package dedup merges multiple environment variable maps into a single
// map, detecting and resolving duplicate keys.
//
// When the same key appears in more than one source map, the caller can
// choose between three strategies:
//
//   - PreferFirst (default) — the value from the earliest map wins.
//   - PreferLast            — the value from the latest map wins.
//   - ErrorOnDuplicate      — Run returns a *DuplicateKeyError immediately.
//
// The Result type records which keys were duplicated so callers can
// surface warnings to the user without aborting the operation.
//
// Example:
//
//	result, err := dedup.Run(
//		[]map[string]string{staging, production},
//		dedup.DefaultOptions(),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Duplicates)
package dedup
