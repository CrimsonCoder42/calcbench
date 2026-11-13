package special

import (
	"fmt"
	"math"
)

// zetaBernoulli holds B₂, B₄, …, B₂₀ divided by (2k)! (DLMF Table 24.2.1), the Euler–Maclaurin weights.
var zetaBernoulli = [10]float64{
	1.0 / 6 / 2,
	-1.0 / 30 / 24,
	1.0 / 42 / 720,
	-1.0 / 30 / 40320,
	5.0 / 66 / 3628800,
	-691.0 / 2730 / 479001600,
	7.0 / 6 / 87178291200,
	-3617.0 / 510 / 20922789888000,
	43867.0 / 798 / 6402373705728000,
	-174611.0 / 330 / 2432902008176640000,
}

// RiemannZeta returns the Riemann zeta function ζ(s) = Σ n^(−s) for real s > 1. For s < 20 it uses
// the Euler–Maclaurin form DLMF 25.2.9 with N = 12 and ten Bernoulli terms (remainder below 1e-20
// relative on 1 < s < 20); for s ≥ 20 the series 25.2.1 itself, whose terms fall below ε after a few.
//
// ID:      NUM-SPEC-008
// Time:    O(1)   at most 12 powers and 10 correction terms, or about 7 series terms for s ≥ 20
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §25.2 (eq. 25.2.1 series, eq. 25.2.9 Euler–Maclaurin), Table 24.2.1 (Bernoulli numbers)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, Table 23.3 (ζ(n) values for the tests)
// Tags:    special-function, zeta, riemann-zeta, series, zipf
// Units:   dimensionless
// Range:   s > 1 and finite. Errors: ErrNotFinite, ErrOutOfRange (s ≤ 1; ζ has a pole at 1).
// Example: ζ(2) = π²/6 = 1.6449340668482264 (DLMF 25.6.1)
func RiemannZeta(s float64) (float64, error) {
	if math.IsNaN(s) || math.IsInf(s, 0) {
		return 0, fmt.Errorf("%w: s = %v", ErrNotFinite, s)
	}
	if s <= 1 {
		return 0, fmt.Errorf("%w: s = %v, want > 1", ErrOutOfRange, s)
	}
	if s >= 20 {
		sum := 1.0
		for n := 2.0; n < 100; n++ {
			t := math.Pow(n, -s)
			sum += t
			if t < epsilon*sum {
				break
			}
		}
		return sum, nil
	}
	const bigN = 12.0
	sum := 0.0
	for k := bigN - 1; k >= 1; k-- { // smallest terms first
		sum += math.Pow(k, -s)
	}
	nPow := math.Pow(bigN, -s) // N^(−s)
	sum += nPow                // the k = N term of the finite sum
	sum += nPow*bigN/(s-1) - nPow/2
	// Σ B₂ₖ/(2k)! · s(s+1)…(s+2k−2) · N^(1−s−2k)
	poch := s            // (s)_(2k−1) for k = 1
	power := nPow / bigN // N^(−s−1) = N^(1−s−2k) for k = 1
	for k := 0; k < len(zetaBernoulli); k++ {
		sum += zetaBernoulli[k] * poch * power
		j := float64(2*k + 1) // advance (s)_(2k−1) to (s)_(2k+1)
		poch *= (s + j) * (s + j + 1)
		power /= bigN * bigN
	}
	return sum, nil
}
