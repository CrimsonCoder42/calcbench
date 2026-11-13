package prob

import (
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"
)

// NoncentralFCDF returns P(X ≤ x) for the noncentral F distribution with df1 and df2 degrees of freedom
// and noncentrality λ: with y = df1·x/(df1·x + df2), F(x) = Σⱼ e^(−λ/2)(λ/2)ʲ/j! · I_y(df1/2 + j, df2/2),
// I the regularized incomplete beta function (num/special RegularizedBeta). As in Lenth's AS 226 it needs
// only two incomplete beta values: I and its decrement u = I_y(a, b) − I_y(a + 1, b) =
// Γ(a + b)/(Γ(a + 1)Γ(b))·yᵃ(1 − y)ᵇ are formed at the Poisson mode j₀ = ⌊λ/2⌋ (u as a difference of two I,
// or of the two complements I_(1−y)(b, ·) when y is above the mean a/(a + b)), and the sum walks both
// ways with I(a + 1) = I(a) − u, u(a + 1) = u(a)·y(a + b)/(a + 1) (Lenth's recurrences, summed outward from
// the mode rather than from j = 0), stopping on a geometric bound of the rest below 2⁻⁵³ of the sum. The
// error is absolute, about one rounding per term. λ = 0 is FDistributionCDF; 0 for x ≤ 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df1 ≤ 0, df2 ≤ 0, λ < 0 or
// λ > 10⁵ (and, wrapping special.ErrNoConvergence, if the incomplete beta failed to converge).
//
// ID:      NUM-PROB-154
// Time:    O(√λ + 1) mixture terms plus two incomplete beta evaluations
// Space:   O(1), no allocation
// Source:  Lenth, "Algorithm AS 226: Computing Noncentral Beta Probabilities", Applied Statistics, 1987, 36(2):241–244, doi:10.2307/2347558 (the mixture of incomplete betas and its recurrences; the noncentral F via y = df1·x/(df1·x + df2))
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 30 "Noncentral F-Distributions" (the Poisson mixture)
// Tags:    probability, distribution, f-distribution, noncentral, cdf, continuous, power-analysis, anova
// Units:   x, df1, df2, lambda dimensionless; result a probability in [0, 1]
// Range:   df1, df2 > 0, 0 ≤ λ ≤ 10⁵
// Example: λ = 0 reduces to the central F: df1 = 2, df2 = 2, x = 1 → 0.5; the mixture with df1 + 2j, x·df1/(df1 + 2j) is the test oracle
func NoncentralFCDF(x, df1, df2, lambda float64) (float64, error) {
	if err := checkFDistribution(df1, df2); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if err := checkNoncentrality(lambda); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	if lambda == 0 {
		return FDistributionCDF(x, df1, df2)
	}
	mu := lambda / 2
	den := df1*x + df2
	y, cy := df1*x/den, df2/den // cy = 1 − y without cancellation
	if y >= 1 || math.IsInf(den, 0) {
		return 1, nil // x so large that y rounds to 1: every I_y is 1
	}
	b := df2 / 2
	j0 := int(math.Floor(mu))
	a := df1/2 + float64(j0)
	var base, dec float64
	if y <= a/(a+b) {
		i0, err := special.RegularizedBeta(y, a, b)
		if err != nil {
			return 0, specialErr("noncentral F cdf", err)
		}
		i1, err := special.RegularizedBeta(y, a+1, b)
		if err != nil {
			return 0, specialErr("noncentral F cdf", err)
		}
		base, dec = i0, i0-i1
	} else {
		c0, err := special.RegularizedBeta(cy, b, a)
		if err != nil {
			return 0, specialErr("noncentral F cdf", err)
		}
		c1, err := special.RegularizedBeta(cy, b, a+1)
		if err != nil {
			return 0, specialErr("noncentral F cdf", err)
		}
		base, dec = 1-c0, c1-c0
	}
	return poissonMixtureCDF(mu, df1/2, j0, base, max(0, dec), func(a float64) float64 { return y * (a + b) / (a + 1) })
}
