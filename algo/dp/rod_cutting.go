package dp

import (
	"fmt"
	"math"
)

// RodCutting returns the largest revenue from cutting a rod of rodLength units into pieces of whole
// lengths and selling them, where a piece of length L sells for pricesByLength[L−1] (pieces longer
// than len(pricesByLength) cannot be sold; uncut leftover earns nothing), and the piece lengths of one
// best cut in the order the trace finds them (first piece = the best first cut of the whole rod).
// Ties go to the shortest first piece. rodLength 0 earns 0 with no pieces (empty, non-nil).
// pricesByLength is read, never modified.
// Errors: ErrNegative when rodLength or a price is negative; ErrTooLarge when the two tables of
// rodLength+1 entries exceed 128 MiB; ErrOverflow when a revenue does not fit in int64.
//
// ID:      ALG-DP-018
// Time:    O(n·min(n, m))   n = rodLength, m = len(pricesByLength)
// Space:   O(n) extra (two tables of n+1 entries, two allocations); O(n) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.1, eq. (15.2) r_n = max_{1≤i≤n}(p_i + r_{n−i}), BOTTOM-UP-CUT-ROD and EXTENDED-BOTTOM-UP-CUT-ROD, Figure 15.1 (the price table)
// Tags:    dp, rod-cutting, unbounded-knapsack, optimisation, revenue
func RodCutting(pricesByLength []int64, rodLength int) (revenue int64, pieces []int, err error) {
	if rodLength < 0 {
		return 0, nil, fmt.Errorf("%w: rodLength = %d", ErrNegative, rodLength)
	}
	for i, p := range pricesByLength {
		if p < 0 {
			return 0, nil, fmt.Errorf("%w: pricesByLength[%d] = %d", ErrNegative, i, p)
		}
	}
	if rodLength == math.MaxInt || !tableFits(rodLength+1, 2, 8) {
		return 0, nil, fmt.Errorf("%w: rodLength = %d", ErrTooLarge, rodLength)
	}
	best := make([]int64, rodLength+1) // best[j] = r_j
	first := make([]int, rodLength+1)  // first[j] = s_j, the first piece of a best cut of j; 0 = sell nothing more
	for j := 1; j <= rodLength; j++ {
		best[j] = best[j-1] // a unit of waste: a length with no price earns nothing
		for i := 1; i <= min(j, len(pricesByLength)); i++ {
			cand, ok := addInt64(pricesByLength[i-1], best[j-i])
			if !ok {
				return 0, nil, fmt.Errorf("%w: revenue for length %d", ErrOverflow, j)
			}
			if cand > best[j] || (first[j] == 0 && cand == best[j] && cand > 0) {
				best[j], first[j] = cand, i
			}
		}
	}
	pieces = []int{}
	for j := rodLength; j > 0; {
		if i := first[j]; i > 0 {
			pieces = append(pieces, i)
			j -= i
		} else {
			j-- // wasted unit
		}
	}
	return best[rodLength], pieces, nil
}
