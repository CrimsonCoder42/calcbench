package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LinearEvaluate returns y = slope·x + intercept (slope-intercept form of a line), computed with one
// fused multiply-add so the result is rounded once.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOverflow when y is not finite. y is 0 with any
// error.
//
// ID:      NUM-ALGEB-004
// Time:    O(1)
// Space:   O(1)
// Source:  Abramson, College Algebra, 2nd ed., OpenStax (Rice University), 2021, §4.1 "Linear Functions" (slope-intercept form), https://openstax.org/books/college-algebra-2e/pages/4-1-linear-functions (read 2026-11-12)
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, clause 5.4.1 (fusedMultiplyAdd, one rounding)
// Tags:    linear, line, slope-intercept, algebra
// Units:   x in the caller's x-unit; slope in y-unit per x-unit; intercept and the result in the y-unit
// Example: slope = 2, intercept = 3, x = 4 → 11 (hand-checked)
func LinearEvaluate(slope, intercept, x float64) (float64, error) {
	if !reduce.AllFinite(slope) {
		return 0, fmt.Errorf("%w: slope = %v", ErrNotFinite, slope)
	}
	if !reduce.AllFinite(intercept) {
		return 0, fmt.Errorf("%w: intercept = %v", ErrNotFinite, intercept)
	}
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	y := math.FMA(slope, x, intercept)
	if !reduce.AllFinite(y) {
		return 0, fmt.Errorf("%w: y", ErrOverflow)
	}
	return y, nil
}
