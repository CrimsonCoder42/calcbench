package prob

import "github.com/CrimsonCoder42/calcbench/num/special"

// ChiSquareSurvival returns P(X > x) = 1 − F(x) for a chi-square variable with df degrees of freedom, the
// regularized upper incomplete gamma function Q(df/2, x/2) (num/special NUM-SPEC-002 RegularizedGammaQ);
// 1 for x ≤ 0. It is the p-value of a
// chi-square statistic x (NIST e-Handbook §1.3.5.15: reject when x exceeds the 1 − α quantile), computed
// directly in the upper tail (x ≥ df + 2), so a p-value of 10⁻³⁰⁰ keeps its relative accuracy; below
// df + 2 it is 1 − F, at least e⁻² ≈ 0.135 when df ≥ 2, with an absolute error of a few 2⁻⁵³ when df < 2.
// Errors: ErrNotFinite when x or df is NaN or ±Inf; ErrOutOfRange when df ≤ 0 or df > 10⁶.
//
// ID:      NUM-PROB-064
// Time:    O(√df) worst case (near x ≈ df); O(1) when x is far from df
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.6 "Chi-Square Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3666.htm (read 2026-11-12)
// Source:  NIST Digital Library of Mathematical Functions, §8.2 eq. 8.2.4 (Q(a, z) = Γ(a, z)/Γ(a)), §8.9 eq. 8.9.2 (continued fraction), https://dlmf.nist.gov/8 (read 2026-11-12)
// Tags:    probability, distribution, chi-square, survival, p-value, upper-tail, incomplete-gamma
// Units:   x and df dimensionless; result a probability in [0, 1]
// Range:   0 < df ≤ 10⁶
// Example: e-Handbook §1.3.6.7.4: df = 6, x = 12.592 → 0.050 (table to 3 decimals); df = 2, x = 2 → e⁻¹ ≈ 0.36787944117144233
func ChiSquareSurvival(x, df float64) (float64, error) {
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if err := checkChiSquareDF(df); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 1, nil
	}
	q, err := special.RegularizedGammaQ(df/2, x/2)
	if err != nil {
		return 0, specialErr("chi-square survival", err)
	}
	return q, nil
}
