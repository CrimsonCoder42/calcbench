package series

import (
	"fmt"
	"math"
)

// GeometricSequenceTerm returns the n-th term of the geometric sequence with first term first and common
// ratio ratio: aₙ = a₁·rⁿ⁻¹, n counted from 1.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange when n < 1; ErrOverflow when the term is
// outside the float64 range.
//
// ID:      NUM-SERIES-022
// Time:    O(1)   (math.Pow)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §9.3 "Geometric Sequences", explicit formula aₙ = a₁rⁿ⁻¹
// Tags:    sequence, geometric-sequence, geometric-progression, nth-term, scalar
// Units:   first and the result in one unit; ratio dimensionless; n a count
// Example: a₁ = 2, r = 5 → a₁₂ = 2·5¹¹ = 97 656 250 (hand-checked)
func GeometricSequenceTerm(first, ratio float64, n int) (float64, error) {
	if !isFinite(first) || !isFinite(ratio) {
		return 0, fmt.Errorf("%w: first = %v, ratio = %v", ErrNotFinite, first, ratio)
	}
	if n < 1 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	if first == 0 {
		return 0, nil
	}
	r := first * math.Pow(ratio, float64(n-1))
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: term %d", ErrOverflow, n)
	}
	return r, nil
}

// GeometricSeriesSum returns the sum of the first n terms of the geometric sequence with first term first
// and common ratio ratio: Sₙ = a₁(1 − rⁿ)/(1 − r) for r ≠ 1 and n·a₁ for r = 1. For 1/2 < r < 3/2 the factor
// (rⁿ − 1)/(r − 1) is formed as expm1(n·ln r)/(r − 1), which keeps full accuracy when r is close to 1;
// elsewhere math.Pow is used, exact for small integer ratios.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange when n < 1; ErrOverflow when the sum is
// outside the float64 range.
//
// ID:      NUM-SERIES-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §9.4 "Series and Their Notations", sum of the first n terms of a geometric series Sₙ = a₁(1 − rⁿ)/(1 − r), r ≠ 1
// Tags:    series, geometric-series, geometric-progression, partial-sum, scalar
// Units:   first and the result in one unit; ratio dimensionless; n a count
// Example: a₁ = 2, r = 5, n = 12 → 2(5¹² − 1)/4 = 122 070 312 (hand-checked)
func GeometricSeriesSum(first, ratio float64, n int) (float64, error) {
	if !isFinite(first) || !isFinite(ratio) {
		return 0, fmt.Errorf("%w: first = %v, ratio = %v", ErrNotFinite, first, ratio)
	}
	if n < 1 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	if first == 0 {
		return 0, nil
	}
	fn := float64(n)
	var factor float64
	switch {
	case ratio == 1:
		factor = fn
	case ratio > 0.5 && ratio < 1.5:
		factor = math.Expm1(fn*math.Log(ratio)) / (ratio - 1)
	default:
		factor = (1 - math.Pow(ratio, fn)) / (1 - ratio)
	}
	r := first * factor
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: sum of %d terms", ErrOverflow, n)
	}
	return r, nil
}
