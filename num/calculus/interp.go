package calculus

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/algo/search"
)

// Interp returns the value at x of the piecewise-linear function through the table points (xp[i], fp[i]):
// it finds the interval xp[j] ≤ x ≤ xp[j+1] by binary search and returns fp[j] + t·(fp[j+1] − fp[j]),
// t = (x − xp[j])/(xp[j+1] − xp[j]). At a table abscissa it returns that table value exactly. x outside
// [xp[0], xp[n−1]] is an error, never a silent extrapolation or clamp. xp and fp are read, never modified.
// Errors: ErrLengthMismatch when len(fp) ≠ len(xp); ErrTooFew when the table has fewer than 2 points;
// ErrNotFinite when x or a table value is NaN or ±Inf; ErrNotAscending when xp is not strictly ascending;
// ErrOutOfRange when x is outside the table.
//
// ID:      NUM-CALC-004
// Time:    O(n)   n = len(xp); the table check is O(n), the search O(log n)
// Space:   O(1) extra, no allocation
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge, 2007, §3.1 (bracketing search, Locate) and §3.3 eq. 3.3.1 (linear interpolation in a table)
// Tags:    calculus, interpolation, linear-interpolation, table, lookup, sampled-data
// Units:   x and xp in one unit, fp in any unit; result in the unit of fp
// Example: table (1, 3), (2, 2), (3, 0) at x = 2.5 → 1 (hand-checked)
func Interp(x float64, xp, fp []float64) (float64, error) {
	n := len(xp)
	if len(fp) != n {
		return 0, fmt.Errorf("%w: len(fp) = %d, len(xp) = %d", ErrLengthMismatch, len(fp), n)
	}
	if n < 2 {
		return 0, fmt.Errorf("%w: %d table points, want ≥ 2", ErrTooFew, n)
	}
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	if err := checkFinite("xp", xp); err != nil {
		return 0, err
	}
	if err := checkFinite("fp", fp); err != nil {
		return 0, err
	}
	for i := 1; i < n; i++ {
		if !(xp[i] > xp[i-1]) {
			return 0, fmt.Errorf("%w: xp[%d] = %v ≤ xp[%d] = %v", ErrNotAscending, i, xp[i], i-1, xp[i-1])
		}
	}
	if x < xp[0] || x > xp[n-1] {
		return 0, fmt.Errorf("%w: x = %v, want %v..%v", ErrOutOfRange, x, xp[0], xp[n-1])
	}
	j := min(search.UpperBound(xp, x)-1, n-2) // xp[j] ≤ x < xp[j+1], or j = n−2 at x = xp[n−1]
	if x == xp[j] {
		return fp[j], nil
	}
	if x == xp[j+1] {
		return fp[j+1], nil
	}
	t := x/2 - xp[j]/2 // halves: no overflow for any finite x, xp
	t /= xp[j+1]/2 - xp[j]/2
	return lerp(fp[j], fp[j+1], t), nil
}

// lerp returns a + t·(b − a) for t in [0, 1], without overflow when b − a would overflow.
func lerp(a, b, t float64) float64 {
	if d := b - a; reduce.AllFinite(d) {
		return a + t*d
	}
	return (1-t)*a + t*b
}
