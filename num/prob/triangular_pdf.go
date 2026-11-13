package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkTriangular validates low < high, low ≤ mode ≤ high, all finite, and returns the width high − low
// (ErrOverflow when it exceeds the float64 range).
func checkTriangular(low, mode, high float64) (float64, error) {
	for _, v := range [...]struct {
		name string
		x    float64
	}{{"low", low}, {"mode", mode}, {"high", high}} {
		if err := checkFinite(v.name, v.x); err != nil {
			return 0, err
		}
	}
	if !(low < high) || mode < low || mode > high {
		return 0, fmt.Errorf("%w: low = %v, mode = %v, high = %v, want low ≤ mode ≤ high, low < high", ErrOutOfRange, low, mode, high)
	}
	w := high - low
	if math.IsInf(w, 0) {
		return 0, fmt.Errorf("%w: high − low = %v − %v", ErrOverflow, high, low)
	}
	return w, nil
}

// TriangularPDF returns the triangular density on [a, b] with mode c: 2(x − a)/((b − a)(c − a)) for
// a ≤ x < c, 2/(b − a) at x = c, 2(b − x)/((b − a)(b − c)) for c < x ≤ b, and 0 outside [a, b].
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange unless low ≤ mode ≤ high and
// low < high; ErrOverflow when high − low or the density exceeds the float64 range.
//
// ID:      NUM-PROB-123
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 26 "Uniform (Rectangular) Distributions", §9 (triangular distributions)
// Tags:    probability, distribution, triangular, pdf, continuous, pert, estimation, monte-carlo
// Units:   x, low, mode and high: one unit; result: probability per that unit
// Example: a = 0, c = 1, b = 2, x = 0.5 → 0.5
func TriangularPDF(x, low, mode, high float64) (float64, error) {
	w, err := checkTriangular(low, mode, high)
	if err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	var d float64
	switch {
	case x < low || x > high:
		return 0, nil
	case x < mode:
		d = 2 / w * ((x - low) / (mode - low))
	case x == mode:
		d = 2 / w
	default:
		d = 2 / w * ((high - x) / (high - mode))
	}
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density, width = %v", ErrOverflow, w)
	}
	return d, nil
}
