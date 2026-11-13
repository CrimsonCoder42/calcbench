package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LaplacePDF returns the Laplace (double exponential) density f(x) = e^(−|x − μ|/b)/(2b), μ the location
// and b the scale.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// density exceeds the float64 range.
//
// ID:      NUM-PROB-107
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 24 "Laplace (Double Exponential) Distributions", §1 (definition)
// Tags:    probability, distribution, laplace, double-exponential, pdf, continuous, lasso
// Units:   x, location and scale: one unit; result: probability per that unit
// Example: μ = 0, b = 1, x = 1 → e⁻¹/2 ≈ 0.18393972058572117
func LaplacePDF(x, location, scale float64) (float64, error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	d := math.Exp(-math.Abs((x-location)/scale)) / (2 * scale)
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, scale = %v", ErrOverflow, x, scale)
	}
	return d, nil
}
