package hashing

import (
	"errors"
	"fmt"
)

// ErrNegativeGap: maxGap is less than 0.
var ErrNegativeGap = errors.New("hashing: maxGap must not be negative")

// ContainsNearbyDuplicate reports whether two equal values sit at most maxGap positions apart: some
// i < j with values[i] == values[j] and j − i <= maxGap. maxGap = 0 always gives false. A NaN equals
// nothing, so it never forms a duplicate. values is read, never modified.
// Errors: ErrNegativeGap when maxGap < 0.
//
// ID:      ALG-HASH-013
// Time:    O(n) expected   n = len(values)
// Space:   O(min(n, maxGap)) extra for the set of the last maxGap values
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary INSERT, SEARCH, DELETE)
// Tags:    hash, set, sliding-window, duplicate, nearby, window
func ContainsNearbyDuplicate[T comparable](values []T, maxGap int) (bool, error) {
	if maxGap < 0 {
		return false, fmt.Errorf("%w: maxGap = %d", ErrNegativeGap, maxGap)
	}
	// The window holds values[j-maxGap : j], all distinct (a repeat inside it would already have returned).
	window := make(map[T]struct{}, min(len(values), maxGap))
	for j, v := range values {
		if j > maxGap {
			delete(window, values[j-maxGap-1])
		}
		if _, ok := window[v]; ok {
			return true, nil
		}
		if maxGap > 0 {
			window[v] = struct{}{}
		}
	}
	return false, nil
}
