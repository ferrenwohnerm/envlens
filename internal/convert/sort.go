package convert

import "sort"

// sortedOrNot returns map keys in sorted order when sorted is true,
// otherwise returns them in arbitrary (map iteration) order.
func sortedOrNot(vars map[string]string, sorted bool) []string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	if sorted {
		sort.Strings(keys)
	}
	return keys
}
