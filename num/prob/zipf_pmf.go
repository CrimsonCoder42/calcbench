package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"
)

// checkZipf returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange unless the exponent is > 1.
func checkZipf(a float64) error {
	if err := checkFinite("a", a); err != nil {
		return err
	}
	if a <= 1 {
		return fmt.Errorf("%w: a = %v, want > 1", ErrOutOfRange, a)
	}
	return nil
}

// zipfZeta returns ζ(a) through special.RiemannZeta, re-wrapping its error.
func zipfZeta(a float64) (float64, error) {
	z, err := special.RiemannZeta(a)
	if err != nil {
		return 0, specialErr("RiemannZeta", err)
	}
	return z, nil
}

// ZipfPMF returns P(X = k) = k^(−a)/ζ(a), k = 1, 2, …, the Zipf (zeta) distribution with exponent a > 1;
// k < 1 gives 0. ζ is special.RiemannZeta (NUM-SPEC-008).
// Errors: ErrNotFinite when a is NaN or ±Inf; ErrOutOfRange when a ≤ 1.
//
// ID:      NUM-PROB-144
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 11, §11.2.20 (zeta and Zipf distributions)
// Tags:    probability, distribution, zipf, zeta, power-law, pmf, discrete, rank-frequency
// Units:   k: a rank (≥ 1); a: dimensionless; result: a probability in [0, 1]
// Example: a = 2, k = 1 → 1/ζ(2) = 6/π² ≈ 0.6079271018540267
func ZipfPMF(k int, a float64) (float64, error) {
	if err := checkZipf(a); err != nil {
		return 0, err
	}
	if k < 1 {
		return 0, nil
	}
	z, err := zipfZeta(a)
	if err != nil {
		return 0, err
	}
	return math.Pow(float64(k), -a) / z, nil
}
