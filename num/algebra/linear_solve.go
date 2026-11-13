package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LinearSolve returns the x that solves slope·x + intercept = y, that is x = (y − intercept)/slope.
// A zero slope is never compared with a tolerance: only slope == 0 exactly is degenerate.
// Errors: ErrDegenerate when slope is 0 (no solution, or every x); ErrNotFinite when an input is NaN or
// ±Inf; ErrOverflow when y − intercept or x is not finite. x is 0 with any error.
//
// ID:      NUM-ALGEB-005
// Time:    O(1)
// Space:   O(1)
// Source:  Abramson, College Algebra, 2nd ed., OpenStax (Rice University), 2021, §4.1 "Linear Functions" (slope-intercept form), https://openstax.org/books/college-algebra-2e/pages/4-1-linear-functions (read 2026-11-12)
// Tags:    linear, solve, inverse, algebra
// Units:   y and intercept in the caller's y-unit; slope in y-unit per x-unit; the result in the x-unit
// Example: slope = 2, intercept = 3, y = 11 → 4 (hand-checked, the inverse of NUM-ALGEB-004's example)
func LinearSolve(slope, intercept, y float64) (float64, error) {
	if !reduce.AllFinite(slope) {
		return 0, fmt.Errorf("%w: slope = %v", ErrNotFinite, slope)
	}
	if !reduce.AllFinite(intercept) {
		return 0, fmt.Errorf("%w: intercept = %v", ErrNotFinite, intercept)
	}
	if !reduce.AllFinite(y) {
		return 0, fmt.Errorf("%w: y = %v", ErrNotFinite, y)
	}
	if slope == 0 {
		return 0, fmt.Errorf("%w: slope = 0", ErrDegenerate)
	}
	rise := y - intercept
	if !reduce.AllFinite(rise) {
		return 0, fmt.Errorf("%w: y − intercept", ErrOverflow)
	}
	x := rise / slope
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrOverflow)
	}
	return x, nil
}
