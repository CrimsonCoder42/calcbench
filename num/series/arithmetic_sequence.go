package series

import (
	"fmt"
	"math"
)

// ArithmeticSequenceTerm returns the n-th term of the arithmetic sequence with first term first and
// common difference difference: aₙ = a₁ + (n − 1)·d, n counted from 1.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange when n < 1; ErrOverflow when the term is
// outside the float64 range.
//
// ID:      NUM-SERIES-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §9.2 "Arithmetic Sequences", explicit formula aₙ = a₁ + d(n − 1)
// Tags:    sequence, arithmetic-sequence, arithmetic-progression, nth-term, scalar
// Units:   first, difference and the result in one unit; n a count
// Example: a₁ = 2, d = 5 → a₂₀ = 2 + 19·5 = 97 (hand-checked)
func ArithmeticSequenceTerm(first, difference float64, n int) (float64, error) {
	if !isFinite(first) || !isFinite(difference) {
		return 0, fmt.Errorf("%w: first = %v, difference = %v", ErrNotFinite, first, difference)
	}
	if n < 1 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	r := math.FMA(float64(n-1), difference, first)
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: term %d", ErrOverflow, n)
	}
	return r, nil
}

// ArithmeticSeriesSum returns the sum of the first n terms of the arithmetic sequence with first term
// first and common difference difference: Sₙ = n(a₁ + aₙ)/2 = n·a₁ + d·n(n − 1)/2.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange when n < 1; ErrOverflow when the sum is
// outside the float64 range.
//
// ID:      NUM-SERIES-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §9.4 "Series and Their Notations", sum of the first n terms of an arithmetic series Sₙ = n(a₁ + aₙ)/2
// Tags:    series, arithmetic-series, arithmetic-progression, partial-sum, scalar
// Units:   first, difference and the result in one unit; n a count
// Example: a₁ = 2, d = 5, n = 20 → 20·(2 + 97)/2 = 990 (hand-checked)
func ArithmeticSeriesSum(first, difference float64, n int) (float64, error) {
	if !isFinite(first) || !isFinite(difference) {
		return 0, fmt.Errorf("%w: first = %v, difference = %v", ErrNotFinite, first, difference)
	}
	if n < 1 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	fn := float64(n)
	// n·a₁ + d·(n(n − 1)/2); n(n − 1)/2 is formed in float64 (exact for n ≤ 2²⁶).
	r := math.FMA(fn*(fn-1)/2, difference, fn*first)
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: sum of %d terms", ErrOverflow, n)
	}
	return r, nil
}

func isFinite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }
