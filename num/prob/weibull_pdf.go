package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkWeibull validates a Weibull shape k and scale λ: both finite and > 0.
func checkWeibull(shape, scale float64) error {
	if err := checkPositive("shape", shape); err != nil {
		return err
	}
	return checkPositive("scale", scale)
}

// WeibullPDF returns the Weibull density f(x) = (k/λ)·(x/λ)^(k−1)·e^(−(x/λ)^k) for x ≥ 0 and 0 for x < 0,
// with k the shape and λ the scale. At x = 0 it is 0 for k > 1, 1/λ for k = 1, and unbounded
// (ErrOutOfRange) for k < 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0 or scale ≤ 0, or at
// x = 0 with shape < 1; ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-087
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 21 "Weibull Distributions", §1 (eq. 21.3)
// Tags:    probability, distribution, weibull, pdf, continuous, reliability, lifetime, wind-speed
// Units:   x and scale: one unit; shape dimensionless; result: probability per that unit
// Example: k = 2, λ = 1, x = 1 → 2e⁻¹ ≈ 0.7357588823428847
func WeibullPDF(x, shape, scale float64) (float64, error) {
	if err := checkWeibull(shape, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x < 0 {
		return 0, nil
	}
	if x == 0 {
		switch {
		case shape > 1:
			return 0, nil
		case shape == 1:
			return 1 / scale, nil
		default:
			return 0, fmt.Errorf("%w: density unbounded at x = 0 for shape = %v < 1", ErrOutOfRange, shape)
		}
	}
	z := x / scale
	lz := math.Log(z)
	d := shape / scale * math.Exp((shape-1)*lz-math.Exp(shape*lz))
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, shape = %v, scale = %v", ErrOverflow, x, shape, scale)
	}
	return d, nil
}
