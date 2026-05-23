package patch

import "fmt"

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"
	OpUnset  Op = "unset"
	OpRename Op = "rename"
)

// Instruction describes a single patch operation to apply to an env map.
type Instruction struct {
	Op    Op
	Key   string
	Value string // used by OpSet
	To    string // used by OpRename
}

// DefaultOptions returns a safe default Options value.
func DefaultOptions() Options {
	return Options{
		ErrorOnMissing: false,
	}
}

// Options controls patch behaviour.
type Options struct {
	// ErrorOnMissing causes Apply to return an error when an OpUnset or
	// OpRename instruction references a key that does not exist in env.
	ErrorOnMissing bool
}

// Apply executes each Instruction against a copy of env and returns the
// resulting map. The original map is never mutated.
func Apply(env map[string]string, instructions []Instruction, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for _, ins := range instructions {
		switch ins.Op {
		case OpSet:
			out[ins.Key] = ins.Value

		case OpUnset:
			if _, ok := out[ins.Key]; !ok && opts.ErrorOnMissing {
				return nil, fmt.Errorf("patch: unset: key %q not found", ins.Key)
			}
			delete(out, ins.Key)

		case OpRename:
			val, ok := out[ins.Key]
			if !ok {
				if opts.ErrorOnMissing {
					return nil, fmt.Errorf("patch: rename: key %q not found", ins.Key)
				}
				continue
			}
			delete(out, ins.Key)
			out[ins.To] = val

		default:
			return nil, fmt.Errorf("patch: unknown op %q", ins.Op)
		}
	}

	return out, nil
}
