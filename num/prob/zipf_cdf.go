package prob

import "math"

// zipfDirectTerms is the largest k whose Zipf CDF is summed term by term.
const zipfDirectTerms = 1000

// zipfTail returns Σᵢ₌ₙ^∞ i^(−a) for n > 1000 and 1 < a ≤ 100 by Euler–Maclaurin:
// n^(1−a)/(a − 1) + n^(−a)/2 + a·n^(−a−1)/12 − a(a+1)(a+2)·n^(−a−3)/720 + a(a+1)⋯(a+4)·n^(−a−5)/30240.
// The first omitted term is below 1e-13 of the result for these n and a.
func zipfTail(n, a float64) float64 {
	f := math.Pow(n, -a)
	inv := 1 / n
	t := n*f/(a-1) + f/2 + a*f*inv/12
	g := a * (a + 1) * (a + 2) * f * inv * inv * inv
	t -= g / 720
	t += g * (a + 3) * (a + 4) * inv * inv / 30240
	return t
}

// ZipfCDF returns P(X ≤ k) = H(k, a)/ζ(a), H(k, a) = Σᵢ₌₁ᵏ i^(−a), for the Zipf distribution. Up to k = 1000
// the sum is compensated term by term; above, P(X ≤ k) = 1 − Σᵢ₌ₖ₊₁^∞ i^(−a)/ζ(a), the tail by
// Euler–Maclaurin, and for a > 100 the tail beyond 1000 is below 10⁻³⁰⁰, so the result is 1. k < 1 gives 0.
// Errors: ErrNotFinite when a is NaN or ±Inf; ErrOutOfRange when a ≤ 1.
//
// ID:      NUM-PROB-145
// Time:    O(min(k, 1000))
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 11, §11.2.20 (zeta and Zipf distributions), summed
// Source:  NIST DLMF, Release 1.2.1, eq. 2.10.1 (Euler–Maclaurin summation) for the tail
// Tags:    probability, distribution, zipf, zeta, power-law, cdf, discrete
// Units:   k: a rank; a: dimensionless; result: a probability in [0, 1]
// Example: a = 2, k = 2 → (1 + 1/4)·6/π² ≈ 0.7599088773175333
func ZipfCDF(k int, a float64) (float64, error) {
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
	if k <= zipfDirectTerms {
		var s compensatedSum
		for i := 1; i <= k; i++ {
			s.add(math.Pow(float64(i), -a))
		}
		return min(1, s.total()/z), nil
	}
	if a > 100 {
		return 1, nil
	}
	return max(0, 1-zipfTail(float64(k)+1, a)/z), nil
}
