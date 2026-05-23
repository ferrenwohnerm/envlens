// Package diff compares two maps of environment variables and
// produces a structured result describing keys that are only in A,
// only in B, present in both with the same value, or present in both
// with different values.
//
// Usage:
//
//	result := diff.Compare(mapA, mapB)
//	for _, r := range result {
//		fmt.Println(r.Key, r.Status)
//	}
package diff
