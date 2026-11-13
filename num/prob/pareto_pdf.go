package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkPareto validates a Pareto scale x_m and shape α: both finite and > 0.
func checkPareto(scale, shape float64) error {
	if err := checkPositive("scale", scale); err != nil {
		return err
	}
	return checkPositive("shape", shape)
}

// ParetoPDF returns the Pareto (first kind) density f(x) = α·x_mᵅ/x^(α+1) = (α/x)·(x_m/x)ᵅ for x ≥ x_m and 0
// below, x_m the scale (minimum) and α the shape (tail index).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0 or shape ≤ 0;
// ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-115
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 20 "Pareto Distributions", §3 (Pareto distribution of the first kind)
// Tags:    probability, distribution, pareto, power-law, pdf, continuous, heavy-tail, income
// Units:   x and scale: one unit; shape dimensionless; result: probability per that unit
// Example: x_m = 1, α = 2, x = 2 → 2·1/8 = 0.25
func ParetoPDF(x, scale, shape float64) (float64, error) {
	if err := checkPareto(scale, shape); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x < scale {
		return 0, nil
	}
	d := shape / x * math.Pow(scale/x, shape)
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, scale = %v, shape = %v", ErrOverflow, x, scale, shape)
	}
	return d, nil
}
