package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"
)

// poissonMixtureCDF returns Σⱼ w_j·B_j, w_j = e^(−μ)μʲ/j! the Poisson weights, for a base sequence that
// falls by a decrement: B_(j+1) = B_j − u_j, u_(j+1) = u_j·ratioUp(a₀ + j), with B and u given at the
// Poisson mode j₀ = ⌊μ⌋ (Ding 1992's recurrences, summed outward from the mode). 0 ≤ B ≤ 1 decreases in
// j, so the rest of the upward sum is below w_j·B_j·r/(1 − r), r = μ/(j + 2), and of the downward sum below
// w_j·r/(1 − r), r = j/μ; each stops when its bound is under 2⁻⁵³ of the sum.
func poissonMixtureCDF(mu, a0 float64, j0 int, base0, dec0 float64, ratioUp func(a float64) float64) (float64, error) {
	w0 := poissonTerm(float64(j0), mu)
	var s compensatedSum
	s.add(w0 * base0)
	w, b, u := w0, base0, dec0
	for j := j0 + 1; ; j++ {
		if j-j0 > noncentralMaxTerms {
			return 0, fmt.Errorf("%w: Poisson mixture did not converge for mean %v", ErrOutOfRange, mu)
		}
		b = max(0, b-u)
		u *= ratioUp(a0 + float64(j-1))
		w *= mu / float64(j)
		term := w * b
		s.add(term)
		r := mu / float64(j+1)
		if term == 0 || term*r/(1-r) <= 0x1p-53*s.total() {
			break
		}
	}
	w, b, u = w0, base0, dec0
	for j := j0 - 1; j >= 0; j-- {
		if ru := ratioUp(a0 + float64(j)); ru > 0 {
			u /= ru
		}
		b = min(1, b+u)
		w *= float64(j+1) / mu
		s.add(w * b)
		r := float64(j) / mu
		if w == 0 || w*r/(1-r) <= 0x1p-53*s.total() {
			break
		}
	}
	return min(1, max(0, s.total())), nil
}

// NoncentralChiSquareCDF returns P(X ≤ x) for the noncentral chi-square distribution with df degrees of
// freedom and noncentrality λ: F(x) = Σⱼ e^(−λ/2)(λ/2)ʲ/j! · P(df/2 + j, x/2), P the regularized lower
// incomplete gamma function (num/special NUM-SPEC-001). As in Ding's AS 275 it needs only two incomplete
// gamma values: P and its decrement t = P(a, z) − P(a + 1, z) = zᵃe^(−z)/Γ(a + 1) are formed at the Poisson
// mode j₀ = ⌊λ/2⌋ (t as a difference of two P, or of two Q when z ≥ a + 1), and the sum walks both ways
// with P(a + 1) = P(a) − t, t(a + 1) = t(a)·z/(a + 1), stopping on a geometric bound of the rest below 2⁻⁵³
// of the sum. The error is absolute, about one rounding per term (a few thousand terms at λ = 10⁵), so a
// tiny lower-tail probability has no relative accuracy. λ = 0 is ChiSquareCDF; 0 for x ≤ 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df ≤ 0, df > 10⁶, λ < 0 or λ > 10⁵.
//
// ID:      NUM-PROB-152
// Time:    O(√λ + √df) (the two incomplete gamma values and O(√λ + 1) mixture terms)
// Space:   O(1), no allocation
// Source:  Ding, "Algorithm AS 275: Computing the Non-Central χ² Distribution Function", Applied Statistics, 1992, 41(2):478–482, doi:10.2307/2347584 (the mixture with the recurrences for P and its decrement)
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 29 "Noncentral χ² Distributions" (the Poisson mixture)
// Tags:    probability, distribution, chi-square, noncentral, cdf, continuous, power-analysis
// Units:   x, df, lambda dimensionless; result a probability in [0, 1]
// Range:   0 < df ≤ 10⁶, 0 ≤ λ ≤ 10⁵
// Example: df = 1: F(x) = Φ(√x − √λ) − Φ(−√x − √λ); λ = 1, x = 1 → Φ(0) − Φ(−2) ≈ 0.4772498680518208
func NoncentralChiSquareCDF(x, df, lambda float64) (float64, error) {
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if err := checkChiSquareDF(df); err != nil {
		return 0, err
	}
	if err := checkNoncentrality(lambda); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	if lambda == 0 {
		return ChiSquareCDF(x, df)
	}
	mu, z := lambda/2, x/2
	j0 := int(math.Floor(mu))
	a := df/2 + float64(j0)
	var base, dec float64
	if z < a+1 {
		p0, err := special.RegularizedGammaP(a, z)
		if err != nil {
			return 0, specialErr("noncentral chi-square cdf", err)
		}
		p1, err := special.RegularizedGammaP(a+1, z)
		if err != nil {
			return 0, specialErr("noncentral chi-square cdf", err)
		}
		base, dec = p0, p0-p1
	} else {
		q0, err := special.RegularizedGammaQ(a, z)
		if err != nil {
			return 0, specialErr("noncentral chi-square cdf", err)
		}
		q1, err := special.RegularizedGammaQ(a+1, z)
		if err != nil {
			return 0, specialErr("noncentral chi-square cdf", err)
		}
		base, dec = 1-q0, q1-q0
	}
	return poissonMixtureCDF(mu, df/2, j0, base, max(0, dec), func(a float64) float64 { return z / (a + 1) })
}
