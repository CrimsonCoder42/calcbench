package twopointer

// RemoveDuplicatesInPlace moves the first element of every run of equal adjacent elements of values to
// the front, in order, and returns how many there are (the new length). On sorted input that removes
// every duplicate. values IS MODIFIED: values[:newLen] holds the kept elements; values[newLen:] keeps
// leftover elements in no promised order (it is not cleared). Equality is ==, so a floating-point NaN
// is never a duplicate, and -0 and +0 are duplicates (the first one is kept).
// The standard library's slices.Compact does the same and also clears the tail; this form returns the
// length and skips that pass.
//
// ID:      ALG-TWOPTR-003
// Time:    O(n)   n = len(values); n−1 comparisons, at most n−1 writes
// Space:   O(1), no allocation
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.5 "Applications", "Duplicates" (sort, then one pass noting equal values that are now adjacent), https://algs4.cs.princeton.edu/25applications/ (read 2026-11-12)
// Source:  The Go standard library, package slices, func Compact (replaces consecutive runs of equal elements with a single copy), https://pkg.go.dev/slices#Compact (read 2026-11-12)
// Tags:    two-pointers, dedup, unique, compact, in-place, sorted
func RemoveDuplicatesInPlace[T comparable](values []T) int {
	if len(values) == 0 {
		return 0
	}
	w := 1 // values[:w] are the kept elements; values[w-1] is the last one kept
	for r := 1; r < len(values); r++ {
		if values[r] != values[w-1] {
			values[w] = values[r]
			w++
		}
	}
	return w
}
