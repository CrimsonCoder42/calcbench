package prob

import (
	"fmt"
	"math"
)

// WaldMeanVariance returns the mean μ and variance μ³/λ of the inverse Gaussian (Wald) distribution.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when mean ≤ 0 or shape ≤ 0;
// ErrOverflow when the variance exceeds the float64 range.
//
// ID:      NUM-PROB-141
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chhikara, Folks, The Inverse Gaussian Distribution, Marcel Dekker, 1989, ch. 2 (moments: E X = μ, Var X = μ³/λ)
// Tags:    probability, distribution, wald, inverse-gaussian, mean, variance, moments, continuous
// Units:   mean, shape: one unit; variance: that unit squared
// Example: μ = 2, λ = 4 → mean 2, variance 2
func WaldMeanVariance(mean, shape float64) (float64, float64, error) {
	if err := checkWald(mean, shape); err != nil {
		return 0, 0, err
	}
	v := mean * mean / shape * mean
	if math.IsInf(v, 0) {
		return 0, 0, fmt.Errorf("%w: Wald variance for mean = %v, shape = %v", ErrOverflow, mean, shape)
	}
	return mean, v, nil
}
