package twopointer

import "math/bits"

// integer is the set of Go's built-in integer types and the types defined on them.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// TwoSumSorted returns indexes i < j with sorted[i] + sorted[j] == target and true, or (-1, -1, false)
// when no such pair exists. Of all such pairs it returns the one with the smallest i and, for that i,
// the largest j. sorted must be in ascending order; on unsorted input the result is unspecified but
// nothing panics. Sums are compared exactly (in 128-bit arithmetic), so no pair is missed or invented
// by overflow.
// sorted is read, never modified.
//
// ID:      ALG-TWOPTR-001
// Time:    O(n)   n = len(sorted); at most n−1 pointer moves
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", "2SUM problem", pp. 78–79, https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.4, Exercise 1.4.15 "Faster 3-sum" (a linear two-sum on a sorted array)
// Tags:    two-pointers, two-sum, pair-sum, sorted, integer
func TwoSumSorted[T integer](sorted []T, target T) (i, j int, ok bool) {
	lo, hi := 0, len(sorted)-1
	t := widen(target)
	for lo < hi {
		switch c := widen(sorted[lo]).add(widen(sorted[hi])).compare(t); {
		case c < 0:
			lo++
		case c > 0:
			hi--
		default:
			return lo, hi, true
		}
	}
	return -1, -1, false
}

// int128 is a two's-complement 128-bit integer (hi signed, lo unsigned); it holds the exact sum of up to
// 2⁶³ values of any 64-bit integer type.
type int128 struct {
	hi int64
	lo uint64
}

// widen returns v as an int128, exactly.
func widen[T integer](v T) int128 {
	var zero T
	if ^zero < zero { // signed T
		x := int64(v)
		return int128{hi: x >> 63, lo: uint64(x)}
	}
	return int128{lo: uint64(v)}
}

func (a int128) add(b int128) int128 {
	lo, carry := bits.Add64(a.lo, b.lo, 0)
	return int128{hi: a.hi + b.hi + int64(carry), lo: lo}
}

// compare returns -1, 0 or +1 as a is less than, equal to or greater than b.
func (a int128) compare(b int128) int {
	switch {
	case a.hi < b.hi:
		return -1
	case a.hi > b.hi:
		return 1
	case a.lo < b.lo:
		return -1
	case a.lo > b.lo:
		return 1
	}
	return 0
}
