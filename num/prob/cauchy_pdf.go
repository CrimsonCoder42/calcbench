package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkCauchy validates a Cauchy location x₀ (finite) and scale γ (finite, > 0).
func checkCauchy(location, scale float64) error {
	if err := checkFinite("location", location); err != nil {
		return err
	}
	return checkPositive("scale", scale)
}

// CauchyPDF returns the Cauchy density f(x) = 1/(πγ·(1 + z²)), z = (x − x₀)/γ, with x₀ the location and γ
// the scale; x₀ = 0, γ = 1 is the standard Cauchy. Its mean and variance do not exist, so there is no
// CauchyMeanVariance.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// density exceeds the float64 range (a scale below about 3e-309).
//
// ID:      NUM-PROB-091
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 16 "Cauchy Distribution", §1 (eq. 16.1)
// Tags:    probability, distribution, cauchy, lorentz, pdf, continuous, heavy-tail
// Units:   x, location and scale: one unit; result: probability per that unit
// Example: x₀ = 0, γ = 1, x = 0 → 1/π ≈ 0.3183098861837907; x = 1 → 1/(2π) ≈ 0.15915494309189535
func CauchyPDF(x, location, scale float64) (float64, error) {
	if err := checkCauchy(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	z := x/scale - location/scale // no overflow of x − x₀
	d := 1 / (math.Pi * scale * (1 + z*z))
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density for scale = %v", ErrOverflow, scale)
	}
	return d, nil
}
