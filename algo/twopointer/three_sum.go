package twopointer

import "slices"

// ThreeSum appends to dst every distinct triple of values (a, b, c), taken from three different
// positions, with a + b + c == target, and returns the extended slice. Each triple is in ascending
// order (a <= b <= c) and the triples are in ascending lexicographic order; triples equal as value
// sets appear once. Sums are compared exactly (in 128-bit arithmetic), so overflow neither hides nor
// invents a triple.
// values is read, never modified: the search runs on a sorted copy.
//
// ID:      ALG-TWOPTR-002
// Time:    O(n²)   n = len(values); O(n log n) to sort the copy, then a linear two-pointer scan per first element
// Space:   O(n) extra (the sorted copy, one allocation); O(t) for the result, t = number of triples ≤ n²
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.4, Exercise 1.4.15 "Faster 3-sum" (a quadratic 3-sum from the linear two-sum on a sorted array)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., booksite §1.4 "Analysis of Algorithms", Web Exercise "Sum of three" (answer: sort, then "for each a in A, scan B and C for a pair that sums to -a (when the sum is too small, advance in B, when the sum is too large, advance in C)"), https://algs4.cs.princeton.edu/14analysis/ (read 2026-11-12)
// Source:  Gajentaan, Overmars, "On a class of O(n²) problems in computational geometry", Computational Geometry: Theory and Applications, 1995, 5(3):165–185, doi:10.1016/0925-7721(95)00022-2 (the 3SUM problem and its quadratic algorithm)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", "2SUM problem", pp. 78–79, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    two-pointers, three-sum, 3sum, triplets, target-sum, sorted, integer
func ThreeSum[T integer](dst [][3]T, values []T, target T) [][3]T {
	n := len(values)
	if n < 3 {
		return dst
	}
	s := slices.Clone(values)
	slices.Sort(s)
	t := widen(target)
	for i := 0; i < n-2; i++ {
		if i > 0 && s[i] == s[i-1] {
			continue // the same first value would repeat the same triples
		}
		first := widen(s[i])
		lo, hi := i+1, n-1
		for lo < hi {
			switch c := first.add(widen(s[lo])).add(widen(s[hi])).compare(t); {
			case c < 0:
				lo++
			case c > 0:
				hi--
			default:
				dst = append(dst, [3]T{s[i], s[lo], s[hi]})
				lo++
				for lo < hi && s[lo] == s[lo-1] {
					lo++
				}
				hi--
				for lo < hi && s[hi] == s[hi+1] {
					hi--
				}
			}
		}
	}
	return dst
}
