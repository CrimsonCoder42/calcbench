package dp

import (
	"fmt"
	"math"
)

// Knapsack01 returns the largest total value of a set of items, each used at most once, whose total
// weight is at most capacity, and the indices of one such set in ascending order (empty and non-nil
// when nothing is taken). Item i has weights[i] and values[i]. Of several optimal sets it returns the
// one the backward trace finds: item i is taken exactly when the table recorded it as an improvement
// at the remaining capacity. Inputs are read, never modified.
// The table is W = min(capacity, Σ weights) + 1 columns wide; one bit per item and column records the
// choices.
// Errors: ErrLengthMismatch when len(weights) ≠ len(values); ErrNegative when capacity, a weight or a
// value is negative; ErrTooLarge when n·W bits plus W int64s exceed 128 MiB; ErrOverflow when a total
// value does not fit in int64.
//
// ID:      ALG-DP-006
// Time:    O(n·W)   n = len(weights), W = min(capacity, Σ weights) (pseudo-polynomial; the problem is NP-hard)
// Space:   O(W) int64s plus n·W bits extra (two allocations); O(n) for the result
// Source:  Kellerer, Pferschy, Pisinger, Knapsack Problems, Springer, 2004, doi:10.1007/978-3-540-24777-7, §2.3 (dynamic programming by the Bellman recursion over capacities)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 16.2-2 (an O(nW) dynamic program for the 0-1 knapsack problem)
// Tags:    dp, knapsack, 0-1-knapsack, subset, optimisation, pseudo-polynomial
func Knapsack01(weights []int, values []int64, capacity int) (best int64, items []int, err error) {
	if err := checkItems(weights, values, capacity); err != nil {
		return 0, nil, err
	}
	n := len(weights)
	// A capacity above the total weight changes nothing; capping it bounds the table by the input.
	total := 0
	for _, w := range weights {
		if total > math.MaxInt-w {
			total = math.MaxInt
			break
		}
		total += w
	}
	capacity = min(capacity, total)
	cols := capacity + 1
	words := (cols + 63) / 64 // words of the choice bitset per item
	if cols <= 0 || !tableFits(n, words, 8) || !tableFits(1, cols, 8) ||
		!tableFits(1, n*words+cols, 8) {
		return 0, nil, fmt.Errorf("%w: %d items × capacity %d", ErrTooLarge, n, capacity)
	}
	row := make([]int64, cols)      // row[c] = best value within capacity c using items 0..i
	took := make([]uint64, n*words) // bit c of item i: item i improves capacity c
	for i, w := range weights {
		v := values[i]
		bitsOfI := took[i*words : (i+1)*words]
		for c := capacity; c >= w; c-- {
			cand, ok := addInt64(row[c-w], v)
			if !ok {
				return 0, nil, fmt.Errorf("%w: total value with item %d", ErrOverflow, i)
			}
			if cand > row[c] {
				row[c] = cand
				bitsOfI[c/64] |= 1 << (c % 64)
			}
		}
	}
	items = []int{}
	for i, c := n-1, capacity; i >= 0; i-- {
		if took[i*words+c/64]>>(c%64)&1 == 1 {
			items = append(items, i)
			c -= weights[i]
		}
	}
	for l, r := 0, len(items)-1; l < r; l, r = l+1, r-1 {
		items[l], items[r] = items[r], items[l]
	}
	return row[capacity], items, nil
}

// UnboundedKnapsack returns the largest total value of a multiset of items, each usable any number of
// times, whose total weight is at most capacity, and counts[i], how many copies of item i one such
// multiset takes. Item i has weights[i] and values[i]. An item of weight 0 and value 0 is never
// taken. Inputs are read, never modified.
// Errors: ErrLengthMismatch when len(weights) ≠ len(values); ErrNegative when capacity, a weight or a
// value is negative; ErrUnbounded when an item has weight 0 and a positive value; ErrTooLarge when the
// two tables of capacity+1 entries exceed 128 MiB; ErrOverflow when a total value does not fit in int64.
//
// ID:      ALG-DP-007
// Time:    O(n·W)   n = len(weights), W = capacity (pseudo-polynomial)
// Space:   O(W) extra (two tables of W+1 entries, two allocations); O(n) for the result
// Source:  Kellerer, Pferschy, Pisinger, Knapsack Problems, Springer, 2004, doi:10.1007/978-3-540-24777-7, ch. 8 (the unbounded knapsack problem; dynamic programming over capacities)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.1 (the same recurrence with one item per length: rod cutting)
// Tags:    dp, knapsack, unbounded-knapsack, integer-knapsack, optimisation, pseudo-polynomial
func UnboundedKnapsack(weights []int, values []int64, capacity int) (best int64, counts []int, err error) {
	if err := checkItems(weights, values, capacity); err != nil {
		return 0, nil, err
	}
	for i, w := range weights {
		if w == 0 && values[i] > 0 {
			return 0, nil, fmt.Errorf("%w: item %d has weight 0 and value %d", ErrUnbounded, i, values[i])
		}
	}
	if capacity == math.MaxInt || !tableFits(capacity+1, 2, 8) {
		return 0, nil, fmt.Errorf("%w: capacity %d", ErrTooLarge, capacity)
	}
	row := make([]int64, capacity+1)  // row[c] = best value within capacity c
	choice := make([]int, capacity+1) // item taken last at capacity c, or -1: capacity c−1 is as good
	choice[0] = -1
	for c := 1; c <= capacity; c++ {
		row[c], choice[c] = row[c-1], -1
		for i, w := range weights {
			if w == 0 || w > c {
				continue
			}
			cand, ok := addInt64(row[c-w], values[i])
			if !ok {
				return 0, nil, fmt.Errorf("%w: total value at capacity %d", ErrOverflow, c)
			}
			if cand > row[c] {
				row[c], choice[c] = cand, i
			}
		}
	}
	counts = make([]int, len(weights))
	for c := capacity; c > 0; {
		if i := choice[c]; i >= 0 {
			counts[i]++
			c -= weights[i]
		} else {
			c--
		}
	}
	return row[capacity], counts, nil
}

// checkItems validates the shared knapsack inputs.
func checkItems(weights []int, values []int64, capacity int) error {
	if len(weights) != len(values) {
		return fmt.Errorf("%w: len(weights) = %d, len(values) = %d", ErrLengthMismatch, len(weights), len(values))
	}
	if capacity < 0 {
		return fmt.Errorf("%w: capacity = %d", ErrNegative, capacity)
	}
	for i, w := range weights {
		if w < 0 {
			return fmt.Errorf("%w: weights[%d] = %d", ErrNegative, i, w)
		}
		if values[i] < 0 {
			return fmt.Errorf("%w: values[%d] = %d", ErrNegative, i, values[i])
		}
	}
	return nil
}
