package special

import (
	"fmt"
	"math"
)

// BesselI0 returns the modified Bessel function of the first kind of order zero, I₀(x). I₀ is even.
// For |x| ≤ 20 it sums the ascending series DLMF 10.25.2 (all terms positive, so no cancellation);
// for |x| > 20 it uses the asymptotic expansion DLMF 10.40.1, whose truncation error there is below
// e^(−2|x|) < 1e-17 relative. e^|x| is formed as two halves so the result overflows only when I₀ does.
//
// ID:      NUM-SPEC-006
// Time:    O(|x|) for |x| ≤ 20 (at most about 60 series terms); O(1) beyond (at most 30 asymptotic terms)
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §10.25 (eq. 10.25.2 with ν = 0), §10.40 (eq. 10.40.1, a_k(0) from eq. 10.17.1)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eqs. 9.6.12 and 9.7.1; Table 9.8
// Tags:    special-function, bessel, modified-bessel, i0, kaiser-window
// Units:   dimensionless
// Range:   x finite; |x| up to about 713.98 (beyond, ErrOverflow). Errors: ErrNotFinite, ErrOverflow.
// Example: I₀(1) = 1.2660658777520082 (A&S Table 9.8: e⁻¹·I₀(1) = 0.4657596077)
func BesselI0(x float64) (float64, error) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	ax := math.Abs(x)
	if ax <= 20 {
		q := ax * ax / 4
		term := 1.0
		sum := 1.0
		for k := 1.0; k < 200; k++ {
			term *= q / (k * k)
			sum += term
			if term < sum*epsilon {
				break
			}
		}
		return sum, nil
	}
	// DLMF 10.40.1: I₀(z) ~ e^z/√(2πz) Σ a_k(0)/z^k, a_k(0) = ((1)(9)…((2k−1)²))/(k! 8^k).
	sum := 1.0
	term := 1.0
	for k := 1.0; k <= 30; k++ {
		next := term * (2*k - 1) * (2*k - 1) / (k * 8 * ax)
		if next < sum*epsilon || next > term {
			break
		}
		term = next
		sum += term
	}
	half := math.Exp(ax / 2)
	v := half * (half * sum / math.Sqrt(2*math.Pi*ax))
	if math.IsInf(v, 0) {
		return 0, fmt.Errorf("%w: I0(%v)", ErrOverflow, x)
	}
	return v, nil
}
