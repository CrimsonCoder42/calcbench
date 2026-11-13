package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// DiscreteVariance returns Var X = Σᵢ pᵢ·(xᵢ − μ)², μ = E[X], for a discrete distribution. It uses the
// corrected two-pass algorithm: the computed mean μ̂, then Σ pᵢdᵢ² − (Σ pᵢdᵢ)² with dᵢ = xᵢ − μ̂ (Chan,
// Golub and LeVeque's eq. (1.7) with the weights 1/N replaced by pᵢ, Σ pᵢ = 1). The second term is zero in
// exact arithmetic and removes the error left by rounding μ̂. Both sums are compensated, and each pᵢdᵢ²
// keeps the rounding error of dᵢ² (FMA). The one-pass E[X²] − μ² is not used because it cancels
// catastrophically. The standard deviation is the caller's math.Sqrt of the result.
// Errors: as ExpectedValue; ErrOverflow when the variance exceeds the float64 range.
//
// ID:      NUM-PROB-019
// Time:    O(n)   n = len(outcomes)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.5 (variance)
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §1.9 "Computing the Sample Variance"
// Source:  Chan, Golub, LeVeque, "Algorithms for Computing the Sample Variance: Analysis and Recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115, eq. (1.7) (the corrected two-pass algorithm)
// Tags:    probability, variance, discrete, dispersion, compensated-summation
// Units:   outcomes: the caller's unit; probabilities: in [0, 1]; result: that unit squared
// Example: a fair die (1…6, each 1/6) → 35/12 ≈ 2.9166666666666665
func DiscreteVariance(outcomes, probabilities []float64) (float64, error) {
	if err := validateDiscrete(outcomes, probabilities); err != nil {
		return 0, err
	}
	mean, err := discreteMean(outcomes, probabilities)
	if err != nil {
		return 0, err
	}
	var sq, dev compensatedSum // Σ pᵢdᵢ² and Σ pᵢdᵢ
	for i, x := range outcomes {
		p := probabilities[i]
		d := x - mean
		d2 := d * d
		sq.addProduct(d2, p)
		sq.comp += math.FMA(d, d, -d2) * p // the rounding error of dᵢ², weighted
		dev.addProduct(d, p)
	}
	c := dev.total()
	v := sq.total() - c*c
	if !reduce.AllFinite(v) {
		return 0, fmt.Errorf("%w: variance", ErrOverflow)
	}
	return max(0, v), nil
}
