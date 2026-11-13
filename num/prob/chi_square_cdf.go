package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/special"
)

// chiSquareMaxDF is the largest degrees of freedom accepted: the series and the continued fraction of
// num/special's incomplete gamma then need at most a few thousand terms, inside its iteration cap
// (1000 + 20·√a terms) and inside the range it is checked over (a ≤ 5·10⁵).
const chiSquareMaxDF = 1e6

// checkChiSquareDF validates the degrees of freedom: finite, > 0 and ≤ chiSquareMaxDF.
func checkChiSquareDF(df float64) error {
	if err := checkPositive("df", df); err != nil {
		return err
	}
	if df > chiSquareMaxDF {
		return fmt.Errorf("%w: df = %v, want ≤ %v", ErrOutOfRange, df, chiSquareMaxDF)
	}
	return nil
}

// ChiSquareCDF returns P(X ≤ x) for a chi-square variable with df degrees of freedom:
// F(x) = γ(df/2, x/2)/Γ(df/2), the regularized lower incomplete gamma function P(df/2, x/2) (num/special
// NUM-SPEC-001 RegularizedGammaP); 0 for x ≤ 0. The smaller tail is computed directly (series below
// df/2 + 1, continued fraction above), so a CDF near 0 keeps its relative accuracy; for a CDF near 1 use
// ChiSquareSurvival, which keeps that of 1 − F.
// Accuracy: in the lower tail a relative error of a few times 10⁻¹⁵ for small df, growing to about 10⁻¹³
// at df = 10⁶ (Loader's saddle-point form for the factor zᵃe^(−z)/Γ(a + 1) keeps it from growing with
// a·ln z); above x ≈ df + 2, F = 1 − Q to an absolute error of a few 2⁻⁵³.
// Errors: ErrNotFinite when x or df is NaN or ±Inf; ErrOutOfRange when df ≤ 0 or df > 10⁶ (and, wrapping
// special.ErrNoConvergence, if the incomplete gamma iteration ever failed to converge).
//
// ID:      NUM-PROB-063
// Time:    O(√df) worst case (terms of the series or the continued fraction near x ≈ df); O(1) when x is far from df
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.6 "Chi-Square Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3666.htm (read 2026-11-12) (cdf F(x) = γ(ν/2, x/2)/Γ(ν/2), x ≥ 0)
// Source:  NIST Digital Library of Mathematical Functions, §8.2 eqs. 8.2.1–8.2.5 (γ, Γ, P, Q), §8.7 eq. 8.7.1 (series), §8.9 eq. 8.9.2 (continued fraction), https://dlmf.nist.gov/8 (read 2026-11-12)
// Source:  Thompson, Barnett, "Coulomb and Bessel functions of complex arguments and order", Journal of Computational Physics, 1986, 64(2):490–509, doi:10.1016/0021-9991(86)90046-X (the modified Lentz method)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf, §2 (the saddle-point form of zᵃe^(−z)/Γ(a), inside NUM-SPEC-001)
// Tags:    probability, distribution, chi-square, cdf, incomplete-gamma, goodness-of-fit, continuous
// Units:   x and df dimensionless; result a probability in [0, 1]
// Range:   0 < df ≤ 10⁶
// Example: e-Handbook §1.3.6.7.4: df = 3, x = 7.815 → 0.950 (table to 3 decimals); df = 2, x = 2 → 1 − e⁻¹ ≈ 0.6321205588285577
func ChiSquareCDF(x, df float64) (float64, error) {
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if err := checkChiSquareDF(df); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	p, err := special.RegularizedGammaP(df/2, x/2)
	if err != nil {
		return 0, specialErr("chi-square cdf", err)
	}
	return p, nil
}
