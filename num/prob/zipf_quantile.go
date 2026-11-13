package prob

import (
	"fmt"
	"math"
)

// zipfQuantileMaxK is the largest rank ZipfQuantile searches; beyond it the quantile is ErrOutOfRange.
const zipfQuantileMaxK = 1 << 62

// ZipfQuantile returns the smallest rank k ≥ 1 with P(X ≤ k) ≥ u for the Zipf distribution with exponent
// a. Up to k = 1000 it walks the compensated partial sums H(k, a)/ζ(a) exactly as ZipfCDF (NUM-PROB-145)
// forms them; above, it brackets by doubling and bisects on 1 − T(k + 1)/ζ(a), T the Euler–Maclaurin tail,
// again ZipfCDF's own expression, so ZipfCDF(k) ≥ u > ZipfCDF(k − 1) for the result. u = 0 returns 1.
// Errors: ErrNotFinite when u or a is NaN or ±Inf; ErrProbability when u is outside [0, 1]; ErrOutOfRange
// when a ≤ 1, when u = 1 (no finite quantile), or when the quantile exceeds 2⁶² (a very close to 1).
//
// ID:      NUM-PROB-147
// Time:    O(min(k, 1000) + log k)   k = the result
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 11, §11.2.20 (zeta and Zipf distributions; inversion of the CDF by search)
// Source:  NIST DLMF, Release 1.2.1, eq. 2.10.1 (Euler–Maclaurin summation) for the tail above k = 1000
// Tags:    probability, distribution, zipf, zeta, power-law, quantile, inverse-cdf, discrete
// Units:   u: a probability; a: dimensionless; result: a rank (≥ 1)
// Example: a = 2: u = 0.7 → 2 (F(1) = 6/π² ≈ 0.6079 < 0.7 ≤ F(2) ≈ 0.7599)
func ZipfQuantile(u, a float64) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkZipf(a); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	z, err := zipfZeta(a)
	if err != nil {
		return 0, err
	}
	var s compensatedSum
	for k := 1; k <= zipfDirectTerms; k++ {
		s.add(math.Pow(float64(k), -a))
		if min(1, s.total()/z) >= u {
			return k, nil
		}
	}
	if a > 100 {
		return zipfDirectTerms + 1, nil // ZipfCDF is 1 above k = 1000 for a > 100
	}
	cdf := func(k int) float64 { return max(0, 1-zipfTail(float64(k)+1, a)/z) }
	lo, hi := zipfDirectTerms, 2*zipfDirectTerms // cdf(lo) < u
	for cdf(hi) < u {
		if hi >= zipfQuantileMaxK/2 {
			return 0, fmt.Errorf("%w: quantile for u = %v, a = %v exceeds 2^62", ErrOutOfRange, u, a)
		}
		lo, hi = hi, 2*hi
	}
	for hi-lo > 1 { // cdf(lo) < u ≤ cdf(hi)
		mid := lo + (hi-lo)/2
		if cdf(mid) >= u {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi, nil
}
