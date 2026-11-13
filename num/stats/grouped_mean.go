package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// GroupedMean returns the estimated mean of grouped data: Σ f·X / N, with X the class mark (midpoint)
// of each class, f its count and N = Σ f. Class j (0-based) runs from edges[j] to edges[j+1]. It is an
// estimate that treats every observation as sitting at its class mark, not the data's own mean.
// The sum is compensated. edges and counts are read, never modified.
// Errors: ErrEmpty when counts is empty; ErrLengthMismatch when len(edges) ≠ len(counts) + 1;
// ErrNotFinite when an edge is NaN or ±Inf; ErrNotAscending when edges are not strictly ascending;
// ErrNegative when a count is negative; ErrZeroTotal when every count is 0; ErrOverflow when the count
// total exceeds the int range.
//
// ID:      NUM-STATS-005
// Time:    O(b)   b = len(counts); two passes
// Space:   O(1) extra, no allocation
// Source:  Illowsky, Dean, Introductory Statistics, 2nd ed., OpenStax, 2023, §2.5 "Measures of the Center of the Data", "Calculating the Mean of Grouped Frequency Tables", https://openstax.org/books/introductory-statistics-2e/pages/2-5-measures-of-the-center-of-the-data (read 2026-11-12)
// Tags:    grouped-data, mean, class-mark, frequency, location, central-tendency, descriptive
// Units:   edges in the unit of the data; counts dimensionless; result in the unit of the data
// Example: classes [0,10) [10,20) [20,30) with counts 2, 5, 3 → (2·5 + 5·15 + 3·25)/10 = 16 (hand-checked)
func GroupedMean(edges []float64, counts []int) (float64, error) {
	total, err := checkGrouped(edges, counts)
	if err != nil {
		return 0, err
	}
	n := float64(total)
	var s compensatedSum
	for j, f := range counts {
		s.add(float64(f) * midpoint(edges[j], edges[j+1]))
	}
	if t := s.total(); reduce.AllFinite(t) {
		return t / n, nil
	}
	// A product or the sum overflowed: weight each class mark by its fraction f/N instead. The result
	// lies between the first and last class marks, so it is finite.
	s = compensatedSum{}
	for j, f := range counts {
		s.add(float64(f) / n * midpoint(edges[j], edges[j+1]))
	}
	t := s.total()
	if math.IsInf(t, 0) {
		return math.Copysign(math.MaxFloat64, t), nil
	}
	return t, nil
}

// checkGrouped validates grouped data (edges of len(counts)+1 classes, strictly ascending and finite;
// counts non-negative with a positive total) and returns the total count.
func checkGrouped(edges []float64, counts []int) (int, error) {
	if len(counts) == 0 {
		return 0, fmt.Errorf("%w: counts", ErrEmpty)
	}
	if len(edges) != len(counts)+1 {
		return 0, fmt.Errorf("%w: len(edges) = %d, want len(counts) + 1 = %d", ErrLengthMismatch, len(edges), len(counts)+1)
	}
	if err := checkEdges(edges); err != nil {
		return 0, err
	}
	return countTotal(counts)
}
