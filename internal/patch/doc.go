// Package patch provides a lightweight instruction-based mutation layer for
// environment variable maps.
//
// An instruction set is an ordered list of [Instruction] values, each
// describing one of three operations:
//
//   - OpSet    – add or overwrite a key with a given value.
//   - OpUnset  – remove a key from the map.
//   - OpRename – rename a key while preserving its value.
//
// Instructions are applied in order against a shallow copy of the input map,
// so the original is never mutated. This makes it safe to use patch alongside
// other envlens pipeline stages without defensive copying at the call site.
//
// Example:
//
//	instructions := []patch.Instruction{
//		{Op: patch.OpSet,    Key: "APP_ENV",   Value: "production"},
//		{Op: patch.OpRename, Key: "DB_HOST",   To: "DATABASE_HOST"},
//		{Op: patch.OpUnset,  Key: "DEBUG"},
//	}
//	out, err := patch.Apply(env, instructions, patch.DefaultOptions())
package patch
