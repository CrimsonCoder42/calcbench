package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// stirlingRemainderTable holds δ(m) = ln m! − (m + ½)·ln m + m − ½·ln 2π for m = 0…15 (δ(0) is unused
// and set to 0). Computed for this package from the definition with 60-digit decimal arithmetic
// (doc/specs/s9_prob.md); the same values are tabulated in Loader's report.
var stirlingRemainderTable = [16]float64{
	0,
	8.1061466795327261070e-02,
	4.1340695955409297035e-02,
	2.7677925684998338357e-02,
	2.0790672103765093365e-02,
	1.6644691189821193139e-02,
	1.3876128823070748436e-02,
	1.1896709945891769528e-02,
	1.0411265261972096202e-02,
	9.2554621827127328548e-03,
	8.3305634333628707927e-03,
	7.5736754879518405903e-03,
	6.9428401072095299179e-03,
	6.4089941880042071432e-03,
	5.9513701127588474957e-03,
	5.5547335519628010525e-03,
}

// stirlingRemainder returns δ(m), the remainder of Stirling's formula, ln m! = (m + ½)ln m − m + ½ln 2π + δ(m),
// for an integer m ≥ 1 held in a float64. Above 15 it sums the Stirling series of DLMF 5.11.1,
// 1/(12m) − 1/(360m³) + 1/(1260m⁵) − 1/(1680m⁷) + 1/(1188m⁹), cut where the first omitted term is at
// most about 2⁻⁵³ (Loader's thresholds 35, 80 and 500). δ enters an exponent, so that absolute error is a
// relative error of the same size in the probability.
func stirlingRemainder(m float64) float64 {
	if m <= 15 {
		return stirlingRemainderTable[int(m)]
	}
	const s0, s1, s2, s3, s4 = 1.0 / 12, 1.0 / 360, 1.0 / 1260, 1.0 / 1680, 1.0 / 1188
	inv2 := 1 / (m * m)
	switch {
	case m > 500:
		return (s0 - s1*inv2) / m
	case m > 80:
		return (s0 - (s1-s2*inv2)*inv2) / m
	case m > 35:
		return (s0 - (s1-(s2-s3*inv2)*inv2)*inv2) / m
	}
	return (s0 - (s1-(s2-(s3-s4*inv2)*inv2)*inv2)*inv2) / m
}

// deviance returns D₀(x, M) = x·ln(x/M) + M − x for x ≥ 0, M > 0, Loader's bd0, without the cancellation
// of the direct form when x is close to M: there it sums the series
// (x − M)·v + 2x·Σⱼ v^(2j+1)/(2j+1), v = (x − M)/(x + M), which converges by a factor v² < 0.01 per term.
func deviance(x, m float64) float64 {
	if d := x - m; math.Abs(d) < 0.1*(x+m) {
		v := d / (x + m)
		s := d * v
		ej := 2 * x * v
		v2 := v * v
		for j := 3.0; j < 200; j += 2 { // |v| < 0.1: the terms fall below an ulp within 9 steps
			ej *= v2
			s1 := s + ej/j
			if s1 == s {
				return s1
			}
			s = s1
		}
		return s
	}
	r := x / m
	if r == 0 || !reduce.AllFinite(r) {
		// x/M under- or overflowed (subnormal M or x): the difference of logs is still finite.
		return x*(math.Log(x)-math.Log(m)) + m - x
	}
	return x*math.Log(r) + m - x
}

// binomialTermParts returns (lc, scale) with b(k; n, p) = exp(lc)·scale for 0 ≤ k ≤ n and 0 < p < 1, in
// Loader's saddle-point form: for 0 < k < n,
// lc = δ(n) − δ(k) − δ(n−k) − D₀(k, np) − D₀(n−k, nq), scale = √(n / (2π·k·(n−k))), q = 1 − p;
// k = 0 and k = n are q^n = exp(n·log1p(−p)) and p^n = exp(n·ln p). Keeping the two parts apart lets
// HypergeometricPMF combine three terms before any of them underflows.
func binomialTermParts(k, n int, p float64) (lc, scale float64) {
	nf := float64(n)
	switch k {
	case 0:
		return nf * math.Log1p(-p), 1
	case n:
		return nf * math.Log(p), 1
	}
	x, y := float64(k), float64(n-k)
	q := 1 - p
	lc = stirlingRemainder(nf) - stirlingRemainder(x) - stirlingRemainder(y) - deviance(x, nf*p) - deviance(y, nf*q)
	scale = math.Sqrt(nf / x / y / (2 * math.Pi))
	return lc, scale
}

// binomialTerm returns b(k; n, p) for 0 ≤ k ≤ n and 0 < p < 1, clamped to at most 1.
func binomialTerm(k, n int, p float64) float64 {
	lc, scale := binomialTermParts(k, n, p)
	return min(1, math.Exp(lc)*scale)
}

// checkBinomial validates n ≥ 0 and p a probability.
func checkBinomial(n int, p float64) error {
	if n < 0 {
		return fmt.Errorf("%w: n = %d, want >= 0", ErrOutOfRange, n)
	}
	return checkProbability("p", p)
}

// BinomialPMF returns P(X = k) = C(n, k)·pᵏ·(1 − p)ⁿ⁻ᵏ, the probability of exactly k successes in n
// independent trials with success probability p; 0 for k outside 0…n. It never forms C(n, k) or a power:
// it uses Loader's saddle-point form, exp(δ(n) − δ(k) − δ(n−k) − D₀(k, np) − D₀(n−k, nq))·√(n/(2πk(n−k))),
// with δ the Stirling remainder and D₀ the deviance: a few ulps near the mode, and in the tails a relative
// error of about 2⁻⁵³ times the size of the deviance terms, for every n (the difference of ln Γ values
// loses about log₁₀(n·ln n) digits everywhere). The result is clamped to at most 1 as the last step.
// Errors: ErrOutOfRange when n < 0; ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is
// outside [0, 1].
//
// ID:      NUM-PROB-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.18 "Binomial Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366i.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables)
// Tags:    probability, distribution, binomial, pmf, discrete, saddle-point
// Units:   k, n: counts of successes and trials; p and result: probabilities in [0, 1]
// Example: n = 10, p = 0.5, k = 5 → 252/1024 = 0.24609375
func BinomialPMF(k, n int, p float64) (float64, error) {
	if err := checkBinomial(n, p); err != nil {
		return 0, err
	}
	if k < 0 || k > n {
		return 0, nil
	}
	switch p {
	case 0:
		return boolProb(k == 0), nil
	case 1:
		return boolProb(k == n), nil
	}
	return binomialTerm(k, n, p), nil
}

// boolProb returns 1 for true and 0 for false: the PMF of a degenerate distribution.
func boolProb(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
