package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkBetaShapes validates the two beta shapes α and β: both finite and > 0.
func checkBetaShapes(alpha, beta float64) error {
	if err := checkPositive("alpha", alpha); err != nil {
		return err
	}
	return checkPositive("beta", beta)
}

// BetaPDF returns the beta density f(x) = x^(α−1)·(1 − x)^(β−1)/B(α, β) on [0, 1], and 0 outside, evaluated
// as exp((α−1)·ln x + (β−1)·log1p(−x) − ln B(α, β)) with ln B from num/special (NUM-SPEC LogBeta). At an end
// point the density is 0 when the shape on that side exceeds 1, 1/B(α, β) when it equals 1, and unbounded
// (ErrOutOfRange) when it is below 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when α ≤ 0 or β ≤ 0, or at an
// unbounded end point; ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-073
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §1 (eq. 25.2 with a = 0, b = 1)
// Source:  NIST DLMF, Release 1.2.1, §5.12.1 (B(a, b))
// Tags:    probability, distribution, beta, pdf, continuous, proportion, bayesian-prior
// Units:   dimensionless; x a proportion in [0, 1]
// Example: α = 2, β = 3, x = 0.5 → 12·0.5·0.25 = 1.5
func BetaPDF(x, alpha, beta float64) (float64, error) {
	if err := checkBetaShapes(alpha, beta); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x < 0 || x > 1 {
		return 0, nil
	}
	lb, err := special.LogBeta(alpha, beta)
	if err != nil {
		return 0, specialErr("log beta", err)
	}
	endpoint := func(shape float64) (float64, error) {
		switch {
		case shape > 1:
			return 0, nil
		case shape == 1:
			return math.Exp(-lb), nil
		default:
			return 0, fmt.Errorf("%w: density unbounded at x = %v", ErrOutOfRange, x)
		}
	}
	if x == 0 {
		return endpoint(alpha)
	}
	if x == 1 {
		return endpoint(beta)
	}
	d := math.Exp((alpha-1)*math.Log(x) + (beta-1)*math.Log1p(-x) - lb)
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, alpha = %v, beta = %v", ErrOverflow, x, alpha, beta)
	}
	return d, nil
}
