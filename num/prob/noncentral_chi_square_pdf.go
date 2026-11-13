package prob

import (
	"fmt"
	"math"
)

// noncentralMaxLambda bounds the noncentrality λ: the Poisson mixture then spans a few thousand terms and
// the incomplete gamma and beta arguments stay inside the range num/special is checked over.
const noncentralMaxLambda = 1e5

// noncentralMaxTerms caps each direction of a Poisson-mixture sum (a defect guard: the sums stop after
// O(√λ + 1) terms).
const noncentralMaxTerms = 1 << 20

// checkNoncentrality validates λ finite and in [0, 10⁵].
func checkNoncentrality(lambda float64) error {
	if err := checkFinite("lambda", lambda); err != nil {
		return err
	}
	if lambda < 0 || lambda > noncentralMaxLambda {
		return fmt.Errorf("%w: lambda = %v, want 0..1e5", ErrOutOfRange, lambda)
	}
	return nil
}

// logGammaTerm returns ln(zᵃe^(−z)/Γ(a + 1)) for a ≥ 0, z > 0: above a = 15 in Loader's saddle-point form
// −δ(a) − D₀(a, z) − ½·ln(2πa) (δ the Stirling remainder, D₀ the deviance), which has no a·ln a
// cancellation; at or below 15 the direct form. With a an integer and z = μ it is the log Poisson weight.
func logGammaTerm(a, z float64) float64 {
	if a > 15 {
		return -stirlingRemainder(a) - deviance(a, z) - 0.5*math.Log(2*math.Pi*a)
	}
	lg, _ := math.Lgamma(a + 1)
	return a*math.Log(z) - z - lg
}

// NoncentralChiSquarePDF returns the density at x of the noncentral chi-square distribution with df
// degrees of freedom and noncentrality λ, the Poisson mixture f(x) = Σⱼ e^(−λ/2)(λ/2)ʲ/j! · g_(df+2j)(x),
// g_ν the central chi-square density. The terms are log-concave in j, so the sum starts at the largest
// one, j* = ⌈(√((ν/2 − 1)² + λx) − ν/2 − 1)/2⌉, ν = df (where the term ratio (λ/2)/(j + 1)·(x/2)/(ν/2 + j) crosses
// 1), formed in log space by Loader's saddle-point form, and walks both ways by that ratio, stopping when
// the geometric bound on the rest is below 2⁻⁵³ of the sum. λ = 0 is ChiSquarePDF. At x = 0 the density is
// 0 for df > 2, e^(−λ/2)/2 for df = 2, and unbounded for df < 2 (ErrOutOfRange); 0 for x < 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df ≤ 0, df > 10⁶, λ < 0, λ > 10⁵,
// or x = 0 with df < 2.
//
// ID:      NUM-PROB-151
// Time:    O(√(λ·x) + 1) terms in the worst case, O(√λ + 1) near the bulk
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 29 "Noncentral χ² Distributions" (the Poisson mixture of central χ² densities)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf, §2 (the saddle-point form of the Poisson and gamma terms)
// Tags:    probability, distribution, chi-square, noncentral, pdf, continuous, power-analysis
// Units:   x, df, lambda dimensionless; result a density
// Range:   0 < df ≤ 10⁶, 0 ≤ λ ≤ 10⁵
// Example: df = 1: f(x) = (φ(√x − √λ) + φ(√x + √λ))/(2√x) (X = (Z + √λ)²); λ = 1, x = 1 → (φ(0) + φ(2))/2 ≈ 0.2264666234573104
func NoncentralChiSquarePDF(x, df, lambda float64) (float64, error) {
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if err := checkChiSquareDF(df); err != nil {
		return 0, err
	}
	if err := checkNoncentrality(lambda); err != nil {
		return 0, err
	}
	switch {
	case x < 0:
		return 0, nil
	case lambda == 0:
		return ChiSquarePDF(x, df)
	case x == 0:
		switch {
		case df > 2:
			return 0, nil
		case df == 2:
			return math.Exp(-lambda/2) / 2, nil
		}
		return 0, fmt.Errorf("%w: density unbounded at x = 0 for df = %v < 2", ErrOutOfRange, df)
	}
	mu, z, a0 := lambda/2, x/2, df/2
	jr := (math.Sqrt((a0-1)*(a0-1)+4*mu*z) - (a0 + 1)) / 2
	j := 0
	if jr > 0 {
		j = int(math.Min(math.Ceil(jr), 1<<50))
	}
	jf := float64(j)
	a := a0 + jf
	t0 := math.Exp(logGammaTerm(jf, mu) + logGammaTerm(a, z) + math.Log(a/z) - math.Ln2)
	var s compensatedSum
	s.add(t0)
	// Upward: term(j+1)/term(j) = (μ/(j + 1))·(z/(a₀ + j)), nonincreasing in j and ≤ 1 from j* on.
	t := t0
	for i := 0; ; i++ {
		if i >= noncentralMaxTerms {
			return 0, fmt.Errorf("%w: noncentral chi-square pdf did not converge for x = %v, df = %v, lambda = %v", ErrOutOfRange, x, df, lambda)
		}
		k := jf + float64(i)
		r := mu / (k + 1) * (z / (a0 + k))
		t *= r
		s.add(t)
		if t == 0 || (r < 1 && t*r/(1-r) <= 0x1p-53*s.total()) {
			break
		}
	}
	// Downward: term(k−1)/term(k) = (k/μ)·((a₀ + k − 1)/z), ≤ 1 from j* down.
	t = t0
	for k := jf; k > 0; k-- {
		r := k / mu * ((a0 + k - 1) / z)
		t *= r
		s.add(t)
		if t == 0 || (r < 1 && t*r/(1-r) <= 0x1p-53*s.total()) {
			break
		}
	}
	return s.total(), nil
}
