package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkGamma validates a gamma shape k and scale θ: both finite and > 0.
func checkGamma(shape, scale float64) error {
	if err := checkPositive("shape", shape); err != nil {
		return err
	}
	return checkPositive("scale", scale)
}

// GammaPDF returns the gamma density f(x) = x^(k−1)·e^(−x/θ)/(Γ(k)·θ^k) for x > 0 and 0 for x < 0, with k
// the shape and θ the scale (a caller holding the rate β passes θ = 1/β). It is evaluated as
// exp((k−1)·ln(x/θ) − x/θ − ln Γ(k))/θ, so it neither overflows nor underflows in its intermediate terms.
// At x = 0 the density is 0 for k > 1 and 1/θ for k = 1; for k < 1 it is unbounded (ErrOutOfRange).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0, scale ≤ 0, or x = 0
// with shape < 1; ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-069
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 17 "Gamma Distributions", §1 (eq. 17.1)
// Source:  NIST DLMF, Release 1.2.1, §5.2 (Γ) and §8.2.4 (P(a, x) = γ(a, x)/Γ(a), whose derivative this is)
// Tags:    probability, distribution, gamma, pdf, continuous, erlang, waiting-time
// Units:   x and scale: one unit (time, length, …); shape dimensionless; result: probability per that unit
// Example: k = 2, θ = 1, x = 1 → e⁻¹ ≈ 0.36787944117144233
func GammaPDF(x, shape, scale float64) (float64, error) {
	if err := checkGamma(shape, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	switch {
	case x < 0:
		return 0, nil
	case x == 0:
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
	lg, _ := math.Lgamma(shape)
	d := math.Exp((shape-1)*math.Log(z)-z-lg) / scale
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, shape = %v, scale = %v", ErrOverflow, x, shape, scale)
	}
	return d, nil
}
