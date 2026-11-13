package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkLognormal validates mu finite and sigma finite and > 0.
func checkLognormal(mu, sigma float64) error {
	if err := checkFinite("mu", mu); err != nil {
		return err
	}
	return checkPositive("sigma", sigma)
}

// LognormalPDF returns the lognormal density f(x) = exp(−(ln x − μ)²/(2σ²)) / (xσ√(2π)) for x > 0, and 0
// for x ≤ 0. μ and σ are the mean and standard deviation of ln X; NIST's scale (the median) is m = e^μ,
// so a caller holding m passes math.Log(m).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when sigma ≤ 0; ErrOverflow when the
// density exceeds the float64 range.
//
// ID:      NUM-PROB-051
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12)
// Tags:    probability, distribution, lognormal, pdf, continuous
// Units:   x: the caller's unit, > 0; mu: natural log of that unit; sigma: dimensionless; result: probability per unit of x
// Example: μ = 0, σ = 1, x = 1 → 1/√(2π) ≈ 0.3989422804014327
func LognormalPDF(x, mu, sigma float64) (float64, error) {
	if err := checkLognormal(mu, sigma); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	z := standardize(math.Log(x), mu, sigma)
	var d float64
	if xs := x * sigma; xs >= 0x1p-1022 && reduce.AllFinite(xs) {
		d = standardNormalPDF(z) / xs
	} else {
		// x·σ is subnormal, zero or beyond the float64 range: combine in logs instead.
		d = math.Exp(-0.5*z*z-math.Log(x)-math.Log(sigma)) * invSqrt2Pi
	}
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v", ErrOverflow, x)
	}
	return d, nil
}
