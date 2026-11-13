package hashing

import (
	"errors"
	"fmt"
)

// ErrOutOfRange: a value or a size parameter is outside the range the function accepts.
var ErrOutOfRange = errors.New("hashing: value out of range")

// BinCount returns how many times each integer 0…bins-1 occurs in values, as a dense slice: counts[v] is
// the number of i with values[i] == v. It is a direct-address table, so no hashing is needed; the caller
// chooses bins (numpy's default is max(values)+1), which bounds the allocation. bins = 0 with empty values
// gives an empty, non-nil slice. values is read, never modified.
// Errors: ErrOutOfRange when bins < 0, or when a value is negative or ≥ bins (wrapped with its index).
//
// ID:      ALG-HASH-014
// Time:    O(n + bins)   n = len(values)
// Space:   O(bins) for the result (one allocation); O(1) extra
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §11.1 (direct-address tables) and §8.2 (counting sort, step 1: the count array)
// Tags:    count, histogram, bincount, direct-address, tally, frequency, dense
func BinCount(values []int, bins int) ([]int, error) {
	if bins < 0 {
		return nil, fmt.Errorf("%w: bins = %d, want ≥ 0", ErrOutOfRange, bins)
	}
	counts := make([]int, bins)
	for i, v := range values {
		if v < 0 || v >= bins {
			return nil, fmt.Errorf("%w: values[%d] = %d, want 0..%d", ErrOutOfRange, i, v, bins-1)
		}
		counts[v]++
	}
	return counts, nil
}
