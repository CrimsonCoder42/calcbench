package window

import (
	"errors"
	"fmt"
)

// ErrNotPositive: a value that must be positive is zero or negative.
var ErrNotPositive = errors.New("window: value must be positive")

// integer is the set of Go's built-in integer types and the types defined on them.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// MinLengthSubarraySumAtLeast returns the shortest non-empty run values[start : start+length] whose sum
// is at least target, as (start, length). Every value must be positive. Among shortest runs it returns
// the one that starts first. When no run reaches target it returns (0, 0); when target <= 0 any single
// value does, so it returns (0, 1) for non-empty values.
// Two pointers move right only. The code never forms a sum: it tracks how much the window still needs
// (target − sum, while short) or has to spare (sum − target, once reached), both between 0 and target or
// the entering value, so no input overflows.
// values is read, never modified.
// Errors: ErrNotPositive when a value is <= 0 (its index is in the message).
//
// ID:      ALG-WINDOW-011
// Time:    O(n)   n = len(values); each pointer moves at most n times
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", "Subarray sum", p. 77 (positive values; the right pointer advances while the sum is short, the left one after; O(n) because both move one way), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    window, sliding-window, two-pointers, variable-window, shortest, subarray-sum, integer
func MinLengthSubarraySumAtLeast[T integer](values []T, target T) (start, length int, err error) {
	for i, v := range values {
		if v <= 0 {
			return 0, 0, fmt.Errorf("%w: values[%d] = %v", ErrNotPositive, i, v)
		}
	}
	if len(values) == 0 {
		return 0, 0, nil
	}
	if target <= 0 {
		return 0, 1, nil
	}
	need := target // the window values[lo:i] sums to target − need, with need > 0
	lo := 0
	for i, v := range values {
		if v < need {
			need -= v
			continue
		}
		// values[lo:i+1] sums to target + spare. Drop from the left while it still reaches target.
		spare := v - need
		for spare >= values[lo] {
			spare -= values[lo]
			lo++
		}
		if length == 0 || i+1-lo < length {
			start, length = lo, i+1-lo
		}
		// Drop one more so the window falls short again: it then needs values[lo] − spare > 0.
		need = values[lo] - spare
		lo++
	}
	return start, length, nil
}
