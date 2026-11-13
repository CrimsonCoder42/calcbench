package stats

import (
	"cmp"
	"math"
	"slices"
)

// countedValue is one distinct value and its frequency.
type countedValue struct {
	value float64
	count int64
}

// FrequencyWeightedMedian returns the median of data given as values with their frequencies (a frequency
// table): the median of the N = Σ fᵢ data with each value repeated fᵢ times, that is the middle value
// when N is odd and the mean of the two middle values when N is even. Values need not be sorted or
// distinct; zero frequencies are ignored. It equals NUM-STATS-012 on the expanded data, without expanding
// it: a weighted quickselect finds the value at rank ⌊(N − 1)/2⌋ (0-based) and one pass finds the next
// order statistic. values and counts are read, never modified.
// Errors: ErrEmpty when values is empty; ErrLengthMismatch when len(counts) ≠ len(values); ErrNotFinite
// when a value is NaN or ±Inf; ErrNegative when a count is negative; ErrZeroTotal when every count is 0;
// ErrOverflow when the total count exceeds the int64 range.
//
// ID:      NUM-STATS-037
// Time:    O(m) on typical input, O(m log m) worst case   m = len(values); quickselect with median-of-3 pivots within a 6m work budget, then a sort of the remaining range
// Space:   O(m) extra (one allocation, a copy of the pairs with non-zero counts)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12) (median Y((N+1)/2) for odd N, (Y(N/2) + Y(N/2+1))/2 for even N)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §9.2 (RANDOMIZED-SELECT, partition-based selection) and Problem 9-2 (weighted median by selection)
// Source:  Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience, 1997, 27(8):983–993 (the work budget and fallback)
// Tags:    median, weighted-median, frequency-table, counts, selection, location, descriptive
// Units:   values in any one unit; counts dimensionless; result in the unit of values
// Example: OpenStax Introductory Statistics 2e Example 2.32 ages 9, 9.5, 10, 10.5, 11, 11.5 with frequencies 1, 2, 4, 4, 6, 3 → N = 20, cumulative 1, 3, 7, 11, 17, 20, so the 10th and 11th values are both 10.5 → 10.5 (hand-checked)
func FrequencyWeightedMedian(values []float64, counts []int) (float64, error) {
	total, err := checkFrequencies(values, counts)
	if err != nil {
		return 0, err
	}
	pairs := make([]countedValue, 0, len(values))
	for i, v := range values {
		if counts[i] > 0 {
			pairs = append(pairs, countedValue{v, int64(counts[i])})
		}
	}
	lowRank := (total - 1) / 2 // 0-based rank of the lower middle value
	lower := weightedSelect(pairs, lowRank)
	if total%2 == 1 {
		return lower, nil
	}
	// Upper middle: lower itself if its copies reach rank lowRank + 1, else the least value above it.
	atMost := int64(0)
	next := math.Inf(1)
	for _, p := range pairs {
		if p.value <= lower {
			atMost += p.count
		} else if p.value < next {
			next = p.value
		}
	}
	if atMost > lowRank+1 {
		return lower, nil
	}
	return midpoint(lower, next), nil
}

// weightedSelectInsertionMax is the range length at or below which the selection sorts and scans.
const weightedSelectInsertionMax = 16

// weightedSelect returns the value at 0-based rank r of the multiset in which each pairs[i].value appears
// pairs[i].count times (counts > 0, 0 ≤ r < Σ counts). It reorders pairs. Quickselect with a median-of-3
// pivot and a three-way partition by value; the partitioned lengths are budgeted at selectWorkBudget·m,
// beyond which the remaining range is sorted and scanned.
func weightedSelect(pairs []countedValue, r int64) float64 {
	budget := selectWorkBudget * len(pairs)
	work := 0
	a := pairs
	for len(a) > weightedSelectInsertionMax && work <= budget {
		work += len(a)
		pivot := floatMedianOfThree(a[0].value, a[len(a)/2].value, a[len(a)-1].value)
		// Three-way partition: a[:lt] < pivot, a[lt:gt] == pivot, a[gt:] > pivot.
		lt, i, gt := 0, 0, len(a)
		var wLess, wEq int64
		for i < gt {
			switch v := a[i].value; {
			case v < pivot:
				wLess += a[i].count
				a[lt], a[i] = a[i], a[lt]
				lt++
				i++
			case v > pivot:
				gt--
				a[gt], a[i] = a[i], a[gt]
			default:
				wEq += a[i].count
				i++
			}
		}
		switch {
		case r < wLess:
			a = a[:lt]
		case r < wLess+wEq:
			return pivot
		default:
			r -= wLess + wEq
			a = a[gt:]
		}
	}
	slices.SortFunc(a, func(x, y countedValue) int { return cmp.Compare(x.value, y.value) })
	for _, p := range a {
		if r < p.count {
			return p.value
		}
		r -= p.count
	}
	return a[len(a)-1].value // unreachable for 0 ≤ r < Σ counts
}
