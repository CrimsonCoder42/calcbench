package matrix

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/prefix"
)

// signed is the signed integer types MaxRectangleSumNoLargerThanK accepts; sums are carried in int64.
type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// MaxRectangleSumNoLargerThanK returns the largest sum of a non-empty axis-aligned sub-rectangle of grid
// that is ≤ k. Pairs of lines along the shorter side are fixed in turn; the cells between them are summed
// into one line, and the best "P[j] − P[i] ≤ k, i < j" over that line's prefix sums P is found by a
// bottom-up merge sort that examines the pairs split across each merge with two pointers.
// grid is read, never modified.
// Errors: ErrRagged; ErrEmpty (no cells); ErrNotFound (every sub-rectangle sums to more than k);
// ErrOverflow when a line or prefix sum leaves int64, or the answer does not fit T.
//
// ID:      ALG-MATRIX-006
// Time:    O(s² · L log L)   s = min(rows, columns), L = max(rows, columns)
// Space:   O(L) extra (three int64 lines, allocated once)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8, Problems (the two-dimensional maximum-sum subarray, reduced to one dimension over pairs of rows)
// Source:  Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §5.3 "Counting Inversions" (Merge-and-Count: the pairs i < j split across sorted halves, examined with two pointers while merging)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.4 (merge sort); prefix sums via PrefixSums (ALG-PREFIX-001)
// Tags:    matrix, max-sum-rectangle, prefix-sum, merge-sort, divide-and-conquer
func MaxRectangleSumNoLargerThanK[T signed](grid [][]T, k T) (T, error) {
	m, n, err := gridShape(grid)
	if err != nil {
		return 0, err
	}
	if m == 0 || n == 0 {
		return 0, fmt.Errorf("%w: %d×%d", ErrEmpty, m, n)
	}
	// Iterate pairs along the short side; cell(i, j) reads the grid with i along the short side.
	short, long := m, n
	cell := func(i, j int) int64 { return int64(grid[i][j]) }
	if m > n {
		short, long = n, m
		cell = func(i, j int) int64 { return int64(grid[j][i]) }
	}
	line := make([]int64, long)
	sums := make([]int64, 0, long+1)
	tmp := make([]int64, long+1)
	bound := int64(k)
	var best int64
	found, overflowCandidate := false, false
	for top := range short {
		clear(line)
		for bottom := top; bottom < short; bottom++ {
			for j := range long {
				s, ok := add64(line[j], cell(bottom, j))
				if !ok {
					return 0, fmt.Errorf("%w: column sum at line %d, position %d", ErrOverflow, bottom, j)
				}
				line[j] = s
			}
			if sums, err = prefix.PrefixSums(sums[:0], line); err != nil {
				return 0, fmt.Errorf("%w: prefix sums: %w", ErrOverflow, err)
			}
			b, ok, ovf := bestDiffAtMost(sums, tmp, bound)
			overflowCandidate = overflowCandidate || ovf
			if ok && (!found || b > best) {
				best, found = b, true
				if best == bound {
					return T(best), nil
				}
			}
		}
	}
	if !found {
		if overflowCandidate {
			return 0, fmt.Errorf("%w: every sum ≤ k is below the int64 range", ErrOverflow)
		}
		return 0, fmt.Errorf("%w: k = %d", ErrNotFound, bound)
	}
	if int64(T(best)) != best {
		return 0, fmt.Errorf("%w: answer %d does not fit the element type", ErrOverflow, best)
	}
	return T(best), nil
}

// add64 returns x+y and whether it did not overflow.
func add64(x, y int64) (int64, bool) {
	s := x + y
	return s, (y >= 0) == (s >= x)
}

// sub64 returns x−y and whether it did not overflow.
func sub64(x, y int64) (int64, bool) {
	d := x - y
	return d, (y >= 0) == (d <= x)
}

// atMost reports whether r − l ≤ k exactly, even when r − l overflows int64.
func atMost(r, l, k int64) bool {
	if d, ok := sub64(r, l); ok {
		return d <= k
	}
	return r < l // the true difference is below MinInt64
}

// bestDiffAtMost returns the largest p[j] − p[i] ≤ k over i < j (ok = false when no pair qualifies) and
// whether some qualifying pair's difference was below the int64 range. It sorts p ascending, using tmp
// (len ≥ len(p)) as the merge buffer: a bottom-up merge sort, with every run [lo, mid) holding indices
// before every index of [mid, hi), so the pairs split across one merge are exactly i in the left run, j
// in the right one; for each right value r (ascending) the best left value is the smallest l ≥ r − k,
// and that pointer only moves right as r grows.
func bestDiffAtMost(p, tmp []int64, k int64) (best int64, ok, overflowed bool) {
	n := len(p)
	best = math.MinInt64
	src, dst := p, tmp[:n]
	for width := 1; width < n; width *= 2 {
		for lo := 0; lo < n; lo += 2 * width {
			mid := min(lo+width, n)
			hi := min(lo+2*width, n)
			// Cross pairs.
			ptr := lo
			for _, r := range src[mid:hi] {
				for ptr < mid && !atMost(r, src[ptr], k) {
					ptr++
				}
				if ptr == mid {
					break
				}
				if d, exact := sub64(r, src[ptr]); exact {
					if !ok || d > best {
						best, ok = d, true
					}
				} else {
					overflowed = true
				}
			}
			// Merge src[lo:mid] and src[mid:hi] into dst[lo:hi].
			i, j, o := lo, mid, lo
			for i < mid && j < hi {
				if src[j] < src[i] {
					dst[o] = src[j]
					j++
				} else {
					dst[o] = src[i]
					i++
				}
				o++
			}
			o += copy(dst[o:], src[i:mid])
			copy(dst[o:], src[j:hi])
		}
		src, dst = dst, src
	}
	return best, ok, overflowed
}
