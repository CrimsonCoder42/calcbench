package search

import (
	"errors"
	"fmt"
)

// Errors returned by SearchPredicate. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrBounds: lo is greater than hi.
	ErrBounds = errors.New("search: lo > hi")
	// ErrNilPredicate: the predicate is nil.
	ErrNilPredicate = errors.New("search: predicate is nil")
)

// SearchPredicate returns the smallest x in [lo, hi) for which pred(x) is true, or hi when pred is false
// on all of it. pred must be monotone on [lo, hi): false for every x below some point and true from it
// on. This is binary search on the answer: x can be a capacity, a speed or a day number rather than a
// slice index. pred is called only on values in [lo, hi), at most ⌈log2(hi-lo+1)⌉ times. Any int range
// works, including lo = math.MinInt and hi = math.MaxInt (the midpoint cannot overflow).
// A non-monotone pred is a precondition violation: the result is still in [lo, hi] and nothing panics.
// Errors: ErrBounds when lo > hi; ErrNilPredicate when pred is nil.
//
// ID:      ALG-SEARCH-008
// Time:    O(log(hi-lo)) calls to pred
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (the invariant that brackets the first position where the test holds)
// Tags:    search, binary-search, predicate, monotone, answer-search, bisection
func SearchPredicate(lo, hi int, pred func(int) bool) (int, error) {
	if pred == nil {
		return lo, fmt.Errorf("%w: pred", ErrNilPredicate)
	}
	if lo > hi {
		return lo, fmt.Errorf("%w: lo = %d, hi = %d", ErrBounds, lo, hi)
	}
	// Invariant: pred is false below lo and true from hi on (hi itself is never passed to pred).
	for lo < hi {
		// hi-lo can exceed MaxInt; as a uint it is exact, and half of it fits an int.
		mid := lo + int(uint(hi-lo)/2)
		if pred(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo, nil
}
