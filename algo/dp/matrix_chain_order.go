package dp

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MatrixChainOrder returns the fewest scalar multiplications needed to compute the product
// A1·A2·…·An, where matrix Ai is dims[i−1] × dims[i], and a full parenthesisation that achieves it,
// written as in CLRS: "((A1(A2A3))((A4A5)A6))"; a single matrix is "A1". Ties go to the smallest
// split point k. dims is read, never modified.
// Hu and Shing's O(n log n) algorithm (SIAM J. Comput. 11(2):362–373, 1982, doi:10.1137/0211028, and
// Part II, 1984) is faster for long chains and is not built here (catalog/dp.md gives the reason).
// Errors: ErrEmpty when len(dims) < 2 (no matrix); ErrNotPositive when a dimension is ≤ 0; ErrTooLarge
// when the two n×n tables exceed 128 MiB (n above about 3300); ErrOverflow when the fewest
// multiplications exceed 2⁶³−1 (a split whose cost overflows is never the best, so only a true
// overflow of the answer is reported).
//
// ID:      ALG-DP-019
// Time:    O(n³)   n = len(dims) − 1 matrices
// Space:   O(n²) extra (the cost table m and the split table s, two allocations), O(n) recursion depth while writing; O(n log n) bytes for the result string
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.2, eq. (15.7), MATRIX-CHAIN-ORDER and PRINT-OPTIMAL-PARENS, Figure 15.5 (the six-matrix example)
// Tags:    dp, matrix-chain, parenthesization, interval-dp, optimisation
func MatrixChainOrder(dims []int) (multiplications int64, parens string, err error) {
	if len(dims) < 2 {
		return 0, "", fmt.Errorf("%w: len(dims) = %d, need at least 2", ErrEmpty, len(dims))
	}
	for i, d := range dims {
		if d <= 0 {
			return 0, "", fmt.Errorf("%w: dims[%d] = %d", ErrNotPositive, i, d)
		}
	}
	n := len(dims) - 1
	if !tableFits(n, n, 8+4) {
		return 0, "", fmt.Errorf("%w: %d matrices", ErrTooLarge, n)
	}
	// m[i*n+j] = fewest multiplications for A(i+1)..A(j+1), 0-based i ≤ j; math.MaxInt64 = overflowed.
	// Every cost is non-negative, so an overflowing split exceeds every split that fits.
	const overflowed = math.MaxInt64
	m := make([]int64, n*n)
	split := make([]int32, n*n)
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			best, bestK := int64(overflowed), i
			for k := i; k < j; k++ {
				left, right := m[i*n+k], m[(k+1)*n+j]
				if left == overflowed || right == overflowed {
					continue
				}
				p, ok := mulPositiveInt64(int64(dims[i]), int64(dims[k+1]))
				if ok {
					p, ok = mulPositiveInt64(p, int64(dims[j+1]))
				}
				if !ok {
					continue
				}
				q, ok := addInt64(left, right)
				if ok {
					q, ok = addInt64(q, p)
				}
				if ok && q < best {
					best, bestK = q, k
				}
			}
			m[i*n+j], split[i*n+j] = best, int32(bestK)
		}
	}
	if m[n-1] == overflowed {
		return 0, "", fmt.Errorf("%w: multiplications for %d matrices", ErrOverflow, n)
	}
	var b strings.Builder
	b.Grow(n * (3 + len(strconv.Itoa(n))))
	writeParens(&b, split, n, 0, n-1)
	return m[n-1], b.String(), nil
}

// writeParens writes PRINT-OPTIMAL-PARENS(s, i+1, j+1). Recursion depth ≤ j−i+1 ≤ n.
func writeParens(b *strings.Builder, split []int32, n, i, j int) {
	if i == j {
		b.WriteByte('A')
		b.WriteString(strconv.Itoa(i + 1))
		return
	}
	k := int(split[i*n+j])
	b.WriteByte('(')
	writeParens(b, split, n, i, k)
	writeParens(b, split, n, k+1, j)
	b.WriteByte(')')
}
