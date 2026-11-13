package dp

import "fmt"

// MaxNonAdjacentSum returns the largest sum of a set of elements of values no two of which are
// adjacent (the maximum-weight independent set of a path). The empty set is allowed, so the result is
// at least 0 and negative elements are never taken. values is read, never modified.
// Errors: ErrOverflow when the sum does not fit in int64.
//
// ID:      ALG-DP-011
// Time:    O(n)   n = len(values)
// Space:   O(1), no allocation
// Source:  Roughgarden, Algorithms Illuminated, Part 3: Greedy Algorithms and Dynamic Programming, Soundlikeyourself Publishing, 2019, ch. 16, §16.1–16.2 (weighted independent set in path graphs: W_i = max(W_{i−1}, W_{i−2} + w_i))
// Source:  Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §6.1 (the same recurrence OPT(j) = max(v_j + OPT(p(j)), OPT(j−1)); on a path, p(j) = j−2)
// Tags:    dp, independent-set, path-graph, non-adjacent, house-robber
func MaxNonAdjacentSum(values []int64) (int64, error) {
	return maxNonAdjacent(values, 0)
}

// MaxNonAdjacentSumCircular returns the largest sum of a set of elements of values no two of which are
// adjacent when values is read as a ring: the last element is adjacent to the first (the maximum-weight
// independent set of a cycle). One element is a ring of one, adjacent to nothing. The empty set is
// allowed. values is read, never modified.
// The ring reduces to two paths: an independent set of the cycle cannot hold both values[0] and
// values[n−1], so it is an independent set of values[1:] or of values[:n−1], and every independent set
// of either path is one of the cycle.
// Errors: ErrOverflow when the sum does not fit in int64.
//
// ID:      ALG-DP-012
// Time:    O(n)   n = len(values); two passes
// Space:   O(1), no allocation
// Source:  Roughgarden, Algorithms Illuminated, Part 3, Soundlikeyourself Publishing, 2019, ch. 16, §16.2 (the path recurrence, applied to the two paths above; the reduction is proved in this comment)
// Tags:    dp, independent-set, cycle-graph, circular, non-adjacent, house-robber
func MaxNonAdjacentSumCircular(values []int64) (int64, error) {
	n := len(values)
	if n <= 1 {
		return maxNonAdjacent(values, 0)
	}
	withoutFirst, err := maxNonAdjacent(values[1:], 1)
	if err != nil {
		return 0, err
	}
	withoutLast, err := maxNonAdjacent(values[:n-1], 0)
	if err != nil {
		return 0, err
	}
	return max(withoutFirst, withoutLast), nil
}

// maxNonAdjacent runs the path recurrence; offset is added to indices in error messages.
func maxNonAdjacent(values []int64, offset int) (int64, error) {
	var skip, take int64 // best of values[:i] without and with values[i−1]
	for i, v := range values {
		withV, ok := addInt64(skip, v)
		if !ok {
			return 0, fmt.Errorf("%w: sum through values[%d]", ErrOverflow, i+offset)
		}
		skip, take = max(skip, take), withV
	}
	return max(skip, take), nil
}
