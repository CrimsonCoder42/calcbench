package prob

import (
	"fmt"
	"math"
)

// checkWald returns ErrNotFinite or ErrOutOfRange unless mean > 0 and shape > 0, both finite.
func checkWald(mean, shape float64) error {
	if err := checkPositive("mean", mean); err != nil {
		return err
	}
	return checkPositive("shape", shape)
}

// WaldPDF returns the inverse Gaussian (Wald) density f(x) = √(λ/(2πx³)) · exp(−λ(x − μ)²/(2μ²x)) for x > 0,
// mean μ and shape λ; x ≤ 0 gives 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when mean ≤ 0 or shape ≤ 0;
// ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-138
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chhikara, Folks, The Inverse Gaussian Distribution: Theory, Methodology, and Applications, Marcel Dekker, 1989, ch. 2 (the density)
// Tags:    probability, distribution, wald, inverse-gaussian, pdf, continuous, first-passage-time
// Units:   x and mean: one unit; shape: the same unit; result: density per that unit
// Example: μ = 1, λ = 1, x = 1 → 1/√(2π) ≈ 0.3989422804014327
func WaldPDF(x, mean, shape float64) (float64, error) {
	if err := checkWald(mean, shape); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	r := (x - mean) / mean
	d := math.Sqrt(shape/(2*math.Pi*x)) / x * math.Exp(-shape*r*r/(2*x))
	if math.IsInf(d, 0) || math.IsNaN(d) {
		return 0, fmt.Errorf("%w: Wald density at x = %v", ErrOverflow, x)
	}
	return d, nil
}
