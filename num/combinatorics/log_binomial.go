package combinatorics

import (
	"fmt"
	"math"
)

// logBinomialSumLimit is the largest m = min(k, n−k) summed term by term; above it the Stirling series
// is used. At m ≥ 65 the series is truncated after the x⁻⁷ term, whose successor is below 1e-16.
const logBinomialSumLimit = 64

// LogBinomial returns ln C(n,k), the natural log of the binomial coefficient, for 0 ≤ k ≤ n. It keeps
// full relative accuracy where the difference of three ln Γ values loses it (n = 10⁹, k = 3: that
// difference is off by about 1.5e-6, a relative error of 2.5e-8). This is the form binomial and hypergeometric
// probabilities need.
// Errors: ErrNegative when n < 0 or k < 0; ErrOutOfRange when k > n (C(n,k) = 0 has no log).
//
// ID:      NUM-COMBI-004
// Time:    O(min(k, n−k, 64))
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.3 eq. 26.3.1, §5.4 eq. 5.4.1, §5.11 eq. 5.11.1
// Tags:    combinatorics, binomial, log-space, probability
// Units:   n, k: counts of objects; result: natural log of a count of subsets
// Example: LogBinomial(100, 50) = 66.783841652017 (ln of the exact C(100,50))
func LogBinomial(n, k int) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return 0, fmt.Errorf("%w: k", ErrNegative)
	}
	if k > n {
		return 0, fmt.Errorf("%w: k > n", ErrOutOfRange)
	}
	m := min(k, n-k)
	rest := float64(n - m)
	if m <= logBinomialSumLimit {
		// C(n,m) = Π_{i=1..m} (n−m+i)/i = Π (1 + (n−m)/i).
		sum := 0.0
		for i := 1; i <= m; i++ {
			sum += math.Log1p(rest / float64(i))
		}
		return sum, nil
	}
	// ln x! = (x+½)ln x − x + ½ln 2π + s(x) (DLMF 5.11.1 with ln x! = ln x + ln Γ(x)). With a = m,
	// b = n−m ≥ a, N = a+b, the −x terms cancel and
	// ln C = a·ln(N/a) + b·log1p(a/b) + ½·ln(N/(a·b)) − ½ln 2π + s(N) − s(a) − s(b).
	a, b, total := float64(m), rest, float64(n)
	v := a*math.Log(total/a) + b*math.Log1p(a/b) + 0.5*math.Log(total/a/b) - 0.5*math.Log(2*math.Pi)
	return v + stirlingRemainder(total) - stirlingRemainder(a) - stirlingRemainder(b), nil
}

// stirlingRemainder returns the Stirling-series tail Σ B₂ₖ/(2k(2k−1)x^(2k−1)) of DLMF 5.11.1 for k = 1…4:
// 1/(12x) − 1/(360x³) + 1/(1260x⁵) − 1/(1680x⁷). Used only for x ≥ 65.
func stirlingRemainder(x float64) float64 {
	inv := 1 / x
	inv2 := inv * inv
	return inv * (1.0/12 - inv2*(1.0/360-inv2*(1.0/1260-inv2/1680)))
}
