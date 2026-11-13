package enumerate

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrNonPositive: a value that must be at least 1 is zero or negative.
	ErrNonPositive = errors.New("enumerate: value must be positive")
	// ErrNegative: a count or target is negative.
	ErrNegative = errors.New("enumerate: value is negative")
	// ErrTooLarge: a count is so large that the generator's buffer would exceed maxAllocBytes.
	ErrTooLarge = errors.New("enumerate: count too large")
	// ErrBitsRange: a bit width is outside [0, 64].
	ErrBitsRange = errors.New("enumerate: bits outside [0, 64]")
)

// CombinationSum returns an iterator over every multiset of candidates (each candidate usable any number
// of times) whose sum is target, each yielded in non-decreasing order, the multisets in lexicographic
// order: candidates [2 3 5], target 8 yields [2 2 2 2] [2 3 3] [3 5]. Duplicate candidates are merged.
// target 0 yields one empty combination. The yielded slice is reused between steps: copy it to keep it,
// never modify it. candidates is copied before returning and never modified.
// Errors: ErrNonPositive when a candidate is ≤ 0 (a zero or negative part gives infinitely many
// combinations); ErrNegative when target < 0.
//
// ID:      ALG-ENUM-007
// Time:    O(c log c) to sort, then O(1) per search-tree node; the tree has at most (output + dead ends) × d nodes, d = target / min(candidates)   c = len(candidates)
// Space:   O(c + d) extra: the sorted candidates and the explicit stack (no recursion)
// Source:  Knuth, The Art of Computer Programming, Vol. 4B: Combinatorial Algorithms, Part 2, Addison-Wesley, 2022, §7.2.2, Algorithm B (basic backtrack)
// Source:  Skiena, The Algorithm Design Manual, 2nd ed., Springer, 2008, §7.1 (backtracking)
// Tags:    backtracking, combination-sum, partition, enumeration, iterator
func CombinationSum(candidates []int, target int) (iter.Seq[[]int], error) {
	if target < 0 {
		return nil, fmt.Errorf("%w: target = %d", ErrNegative, target)
	}
	for i, c := range candidates {
		if c <= 0 {
			return nil, fmt.Errorf("%w: candidates[%d] = %d", ErrNonPositive, i, c)
		}
	}
	cands := slices.Clone(candidates)
	slices.Sort(cands)
	cands = slices.Compact(cands)
	return func(yield func([]int) bool) {
		var pick []int // indices into cands, non-decreasing: the path from the root
		var buf []int  // cands[pick[i]], the yielded combination
		rem := target  // target − sum(buf), never negative
		next := 0      // the smallest candidate index the next part may use
		for {
			if rem == 0 {
				if !yield(buf) {
					return
				}
			} else if next < len(cands) && cands[next] <= rem {
				pick = append(pick, next)
				buf = append(buf, cands[next])
				rem -= cands[next]
				continue // the same candidate may be used again
			}
			// Backtrack: drop the last part and try the next larger candidate in its place.
			if len(pick) == 0 {
				return
			}
			last := pick[len(pick)-1]
			pick = pick[:len(pick)-1]
			buf = buf[:len(buf)-1]
			rem += cands[last]
			next = last + 1
		}
	}, nil
}
