package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// invSqrt2Pi is 1/√(2π), correctly rounded.
const invSqrt2Pi = 0.398942280401432677939946059934381868475858631164934657665926

// checkNormal validates mean finite and sd finite and > 0.
func checkNormal(mean, sd float64) error {
	if err := checkFinite("mean", mean); err != nil {
		return err
	}
	return checkPositive("sd", sd)
}

// standardNormalPDF returns φ(z) = exp(−z²/2)/√(2π). For |z| ≥ 5 it splits z = hi + lo with hi on a
// 2⁻¹⁶ grid, so hi² is exact and exp(−z²/2) = exp(−hi²/2)·exp(−lo·(hi + lo/2)) keeps full relative
// accuracy in the tail instead of the z²·u error of rounding z² first. ±Inf gives 0.
func standardNormalPDF(z float64) float64 {
	a := math.Abs(z)
	if a < 5 {
		return math.Exp(-0.5*z*z) * invSqrt2Pi
	}
	if a > 40 {
		return 0 // exp(−800) underflows
	}
	hi := math.Round(a*65536) / 65536
	lo := a - hi
	return math.Exp(-0.5*hi*hi) * math.Exp(-lo*(hi+0.5*lo)) * invSqrt2Pi
}

// standardize returns z = (x − mean)/sd. A difference beyond the float64 range gives ±Inf, which every
// caller maps to the correct limit (density 0, probability 0 or 1).
func standardize(x, mean, sd float64) float64 {
	return (x - mean) / sd
}

// NormalPDF returns the normal density f(x) = exp(−z²/2)/(σ√(2π)), z = (x − μ)/σ. The standard normal is
// mean 0, sd 1. z is formed before squaring, so (x − μ)² never overflows.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when sd ≤ 0; ErrOverflow when the
// density exceeds the float64 range (sd below about 4e-309).
//
// ID:      NUM-PROB-047
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.4 (normal random variables)
// Tags:    probability, distribution, normal, gaussian, pdf, continuous, bell-curve
// Units:   x, mean, sd: one unit; result: probability per that unit
// Example: x = 0, mean 0, sd 1 → 1/√(2π) ≈ 0.3989422804014327
func NormalPDF(x, mean, sd float64) (float64, error) {
	if err := checkNormal(mean, sd); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	d := standardNormalPDF(standardize(x, mean, sd)) / sd
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density for sd = %v", ErrOverflow, sd)
	}
	return d, nil
}
