package snapshot

import (
	"fmt"

	"github.com/user/envlens/internal/diff"
)

// CompareResult holds the result of comparing two named snapshots.
type CompareResult struct {
	// LabelA is the name of the first snapshot.
	LabelA string
	// LabelB is the name of the second snapshot.
	LabelB string
	// Results contains the diff entries between the two snapshots.
	Results []diff.Result
}

// Compare loads two snapshots by label from the given directory and returns
// a CompareResult containing the diff between them. Returns an error if
// either snapshot cannot be loaded.
func Compare(dir, labelA, labelB string) (*CompareResult, error) {
	if labelA == "" || labelB == "" {
		return nil, fmt.Errorf("snapshot labels must not be empty")
	}

	varsA, err := Load(dir, labelA)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot %q: %w", labelA, err)
	}

	varsB, err := Load(dir, labelB)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot %q: %w", labelB, err)
	}

	results := diff.Compare(varsA, varsB)

	return &CompareResult{
		LabelA:  labelA,
		LabelB:  labelB,
		Results: results,
	}, nil
}
