package calculus

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Gradient returns dy/dx at every sample of the tabulated function y(x), by second-order finite
// differences on the grid x as given (it need not be evenly spaced). At an interior point, with
// h₋ = x[i] − x[i−1] and h₊ = x[i+1] − x[i], the three-point weights are
// −h₊/(h₋(h₋+h₊)), (h₊−h₋)/(h₋h₊), h₋/(h₊(h₋+h₊)) for y[i−1], y[i], y[i+1]; at each end the one-sided
// three-point weights of the same order are used. The result is exact for every quadratic. With exactly
// two samples both entries are the one slope (y[1] − y[0])/(x[1] − x[0]), first order.
// The result is written into dst when cap(dst) ≥ len(y), else a new slice is allocated. dst must not
// overlap y or x. y and x are read, never modified.
// Errors: ErrLengthMismatch when len(y) ≠ len(x); ErrTooFew when there are fewer than 2 samples;
// ErrNotFinite when a value is NaN or ±Inf; ErrNotAscending when x is not strictly ascending;
// ErrOverflow when a step or a derivative is beyond float64.
//
// ID:      NUM-CALC-003
// Time:    O(n)   n = len(y); one pass
// Space:   O(1) extra; O(n) for the result unless dst has the capacity (then no allocation)
// Source:  Fornberg, "Generation of finite difference formulas on arbitrarily spaced grids", Mathematics of Computation, 1988, 51(184):699–706, doi:10.1090/S0025-5718-1988-0935077-0 (the algorithm, first derivative, three nodes; Table 1 for the uniform case)
// Tags:    calculus, derivative, gradient, finite-difference, nonuniform-grid, sampled-data
// Units:   x in any unit, y in any unit; result in unit(y)/unit(x)
// Example: y = x² at x = 0, 1, 3 → 0, 2, 6 (exact for a quadratic; hand-checked)
func Gradient(dst, y, x []float64) ([]float64, error) {
	n := len(y)
	if len(x) != n {
		return nil, fmt.Errorf("%w: len(y) = %d, len(x) = %d", ErrLengthMismatch, n, len(x))
	}
	if n < 2 {
		return nil, fmt.Errorf("%w: %d samples, want ≥ 2", ErrTooFew, n)
	}
	if err := checkFinite("x", x); err != nil {
		return nil, err
	}
	if err := checkFinite("y", y); err != nil {
		return nil, err
	}
	for i := 1; i < n; i++ {
		if !(x[i] > x[i-1]) {
			return nil, fmt.Errorf("%w: x[%d] = %v ≤ x[%d] = %v", ErrNotAscending, i, x[i], i-1, x[i-1])
		}
		if !reduce.AllFinite(x[i] - x[i-1]) {
			return nil, fmt.Errorf("%w: x[%d] − x[%d]", ErrOverflow, i, i-1)
		}
	}
	if cap(dst) < n {
		dst = make([]float64, n)
	}
	dst = dst[:n]
	if n == 2 {
		d := slope(y[0], y[1], x[1]-x[0])
		dst[0], dst[1] = d, d
		return checkResult(dst)
	}
	// Each formula is written as a combination of divided differences (weights summing to zero exactly),
	// which is the Fornberg three-point stencil rearranged: fewer roundings and no y·weight overflow.
	for i := 1; i < n-1; i++ {
		hm, hp := x[i]-x[i-1], x[i+1]-x[i]
		dm, dp := slope(y[i-1], y[i], hm), slope(y[i], y[i+1], hp)
		// (h₊·d₋ + h₋·d₊)/(h₋ + h₊), a weighted mean of the two slopes.
		w := hm / (hm + hp)
		dst[i] = dm*(1-w) + dp*w
	}
	{
		h1, h2 := x[1]-x[0], x[2]-x[1]
		d1, d2 := slope(y[0], y[1], h1), slope(y[1], y[2], h2)
		// d₁ − h₁·(d₂ − d₁)/(h₁ + h₂)
		dst[0] = d1 - (d2-d1)*(h1/(h1+h2))
	}
	{
		h1, h2 := x[n-2]-x[n-3], x[n-1]-x[n-2]
		d1, d2 := slope(y[n-3], y[n-2], h1), slope(y[n-2], y[n-1], h2)
		// d₂ + h₂·(d₂ − d₁)/(h₁ + h₂)
		dst[n-1] = d2 + (d2-d1)*(h2/(h1+h2))
	}
	return checkResult(dst)
}

// slope returns (b − a)/h, computed as b/h − a/h when b − a overflows.
func slope(a, b, h float64) float64 {
	if d := b - a; reduce.AllFinite(d) {
		return d / h
	}
	return b/h - a/h
}

// checkResult returns dst, or ErrOverflow naming the first entry beyond float64.
func checkResult(dst []float64) ([]float64, error) {
	for i, v := range dst {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: derivative at sample %d", ErrOverflow, i)
		}
	}
	return dst, nil
}
