package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RayleighPDF returns the Rayleigh density f(x) = (x/σ²)·e^(−x²/(2σ²)) for x ≥ 0 and 0 below, σ the scale.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// density exceeds the float64 range.
//
// ID:      NUM-PROB-119
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 18 "Chi-Square Distributions Including Chi and Rayleigh", §10 (Rayleigh distribution: chi with 2 degrees of freedom, scaled)
// Tags:    probability, distribution, rayleigh, pdf, continuous, wind-speed, signal-envelope
// Units:   x and scale: one unit; result: probability per that unit
// Example: σ = 1, x = 1 → e^(−½) ≈ 0.6065306597126334
func RayleighPDF(x, scale float64) (float64, error) {
	if err := checkPositive("scale", scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	z := x / scale
	if !reduce.AllFinite(z) {
		return 0, nil
	}
	d := z / scale * math.Exp(-0.5*z*z)
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, scale = %v", ErrOverflow, x, scale)
	}
	return d, nil
}
