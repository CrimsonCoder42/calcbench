package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LinearRegression returns the intercept and slope of the least-squares straight line y = intercept +
// slope·x through the points (x[i], y[i]): slope = Σ(x − x̄)(y − ȳ)/Σ(x − x̄)², intercept = ȳ − slope·x̄.
// The centred sums use the corrected two-pass form, compensated, never the one-pass Σxy − n·x̄·ȳ.
// It is a descriptive fit: no standard errors or tests. x and y are read, never modified.
// Errors: ErrEmpty when x is empty; ErrLengthMismatch when len(y) ≠ len(x); ErrTooFew when there is one
// point; ErrNotFinite when a value is NaN or ±Inf; ErrZeroVariance when every x is equal; ErrOverflow
// when the slope or the intercept exceeds the float64 range.
//
// ID:      NUM-STATS-009
// Time:    O(n)   n = len(x); three passes (means, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §4.4.3.1 "Least Squares", https://www.itl.nist.gov/div898/handbook/pmd/section4/pmd431.htm (read 2026-11-12)
// Tags:    linear-regression, least-squares, trend-line, slope, intercept, fit, descriptive
// Units:   x in unit X; y in unit Y; intercept in unit Y; slope in unit Y per unit X
// Example: NIST StRD Linear Least Squares, Norris (36 points) → intercept −0.262323073774029, slope 1.00211681802045 (certified)
func LinearRegression(x, y []float64) (intercept, slope float64, err error) {
	if err := checkPaired(x, y); err != nil {
		return 0, 0, err
	}
	xMean, yMean := meanOf(x), meanOf(y)
	sxx, _, sxy, ex, ey := centredProducts(x, y, xMean, yMean)
	if sxx == 0 {
		return 0, 0, fmt.Errorf("%w: every x equals %v", ErrZeroVariance, x[0])
	}
	slope = math.Ldexp(sxy/sxx, ey-ex) // undo the scales: 2^(ex+ey)/4^ex
	if !reduce.AllFinite(slope) {
		return 0, 0, fmt.Errorf("%w: slope", ErrOverflow)
	}
	intercept = math.FMA(-slope, xMean, yMean)
	if !reduce.AllFinite(intercept) {
		return 0, 0, fmt.Errorf("%w: intercept", ErrOverflow)
	}
	return intercept, slope, nil
}

// checkPaired validates two paired samples: ErrEmpty, ErrLengthMismatch, ErrTooFew (fewer than 2
// pairs) and ErrNotFinite.
func checkPaired(x, y []float64) error {
	if len(x) == 0 {
		return fmt.Errorf("%w: x", ErrEmpty)
	}
	if len(y) != len(x) {
		return fmt.Errorf("%w: len(y) = %d, len(x) = %d", ErrLengthMismatch, len(y), len(x))
	}
	if len(x) < 2 {
		return fmt.Errorf("%w: len(x) = %d, want ≥ 2", ErrTooFew, len(x))
	}
	if err := checkFinite("x", x); err != nil {
		return err
	}
	return checkFinite("y", y)
}
