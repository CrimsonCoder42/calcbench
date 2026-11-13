package special

import (
	"errors"
	"fmt"
	"math"
)

// Sentinel errors for package special. Every error a function returns wraps one of these with the
// offending parameter's name and value.
var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("special: value is NaN or infinite")
	// ErrOutOfRange: an input lies outside the function's domain.
	ErrOutOfRange = errors.New("special: argument out of range")
	// ErrOverflow: the true result exceeds the largest float64.
	ErrOverflow = errors.New("special: result overflows float64")
	// ErrNoConvergence: an iteration did not reach full precision within its bound.
	ErrNoConvergence = errors.New("special: iteration did not converge")
)

// tiny guards the modified Lentz algorithm against a zero denominator (Numerical Recipes 3rd ed. §5.2).
const tiny = 1e-300

// iterationCap bounds a series or continued fraction whose term count grows like √a (Numerical
// Recipes 3rd ed. §6.2): 1000 terms plus 20 per unit of √a.
func iterationCap(a float64) int {
	return 1000 + int(20*math.Sqrt(a))
}

// RegularizedGammaP returns the regularized lower incomplete gamma function P(a, x) = γ(a, x)/Γ(a).
// x < a+1 uses the power series of DLMF 8.7.1; otherwise P = 1 − Q with Q from the continued fraction.
// P(a, 0) = 0. Accuracy: about 1e-14 relative for every a (checked to a = 5e5); above a = 15 the prefactor
// e^(−x)·x^a/Γ(a) uses Loader's saddle-point form, so it does not lose log10(a) digits to ln Γ(a).
//
// ID:      NUM-SPEC-001
// Time:    O(√a + log(1/ε))   a = shape, ε = 2⁻⁵² (terms of the series or the continued fraction)
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §8.2 (eq. 8.2.4 P = γ/Γ), §8.7 (eq. 8.7.1 series), §8.9 (eq. 8.9.2 continued fraction)
// Source:  Press, Teukolsky, Vetterling & Flannery, Numerical Recipes, 3rd ed., Cambridge, 2007, §6.2 (gser, gcf)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf, §2 and Appendix A (stirlerr, bd0: the saddle-point prefactor); NIST DLMF §5.11 eq. 5.11.1 (Stirling series)
// Tags:    special-function, incomplete-gamma, regularized, gamma, chi-square, poisson
// Units:   dimensionless
// Range:   a > 0 and finite; x ≥ 0 and finite. Errors: ErrNotFinite, ErrOutOfRange, ErrNoConvergence.
// Example: P(1, 1) = 1 − e⁻¹ = 0.6321205588285577 (DLMF 8.4.4 with a = 1)
func RegularizedGammaP(a, x float64) (float64, error) {
	if err := checkGammaArgs(a, x); err != nil {
		return 0, err
	}
	if x == 0 {
		return 0, nil
	}
	if x < a+1 {
		return gammaSeries(a, x)
	}
	q, err := gammaContinuedFraction(a, x)
	if err != nil {
		return 0, err
	}
	return 1 - q, nil
}

// RegularizedGammaQ returns the regularized upper incomplete gamma function Q(a, x) = Γ(a, x)/Γ(a)
// = 1 − P(a, x). x ≥ a+1 uses the continued fraction of DLMF 8.9.2 directly, so a tiny Q keeps its
// relative accuracy; otherwise Q = 1 − P from the series. Q(a, 0) = 1. Accuracy: about 1e-14 relative
// for every a (checked to a = 5e5), the prefactor in Loader's saddle-point form above a = 15.
//
// ID:      NUM-SPEC-002
// Time:    O(√a + log(1/ε))   a = shape, ε = 2⁻⁵²
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §8.2 (eq. 8.2.4), §8.7 (eq. 8.7.1), §8.9 (eq. 8.9.2); modified Lentz evaluation per Numerical Recipes 3rd ed. §5.2
// Source:  Press et al., Numerical Recipes, 3rd ed., Cambridge, 2007, §6.2
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf, §2 and Appendix A (stirlerr, bd0); NIST DLMF §5.11 eq. 5.11.1
// Tags:    special-function, incomplete-gamma, regularized, gamma, survival, chi-square
// Units:   dimensionless
// Range:   a > 0 and finite; x ≥ 0 and finite. Errors: ErrNotFinite, ErrOutOfRange, ErrNoConvergence.
// Example: Q(1, 1) = e⁻¹ = 0.36787944117144233
func RegularizedGammaQ(a, x float64) (float64, error) {
	if err := checkGammaArgs(a, x); err != nil {
		return 0, err
	}
	if x == 0 {
		return 1, nil
	}
	if x < a+1 {
		p, err := gammaSeries(a, x)
		if err != nil {
			return 0, err
		}
		return 1 - p, nil
	}
	return gammaContinuedFraction(a, x)
}

func checkGammaArgs(a, x float64) error {
	if math.IsNaN(a) || math.IsInf(a, 0) {
		return fmt.Errorf("%w: a = %v", ErrNotFinite, a)
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	if a <= 0 {
		return fmt.Errorf("%w: a = %v, want > 0", ErrOutOfRange, a)
	}
	if x < 0 {
		return fmt.Errorf("%w: x = %v, want ≥ 0", ErrOutOfRange, x)
	}
	return nil
}

// gammaPrefactor returns e^(−x)·x^a/Γ(a) for a > 0, x > 0 (underflow to 0 is the true value). Above
// a = 15 it uses Loader's saddle-point form: with Γ(a + 1) = √(2πa)·aᵃe^(−a)·e^δ(a) (δ the Stirling
// remainder, DLMF 5.11.1), x^a e^(−x)/Γ(a) = √(a/(2π))·exp(−δ(a) − D₀(a, x)), D₀(a, x) = a·ln(a/x) + x − a
// the deviance (Loader 2000, stirlerr and bd0). That avoids the cancellation of a·ln x − x − ln Γ(a), whose
// terms are of size a·ln a and lose about log₁₀(a) digits. At or below 15 the direct form is exact to a
// few ulps.
func gammaPrefactor(a, x float64) float64 {
	if a > 15 {
		return math.Sqrt(a/(2*math.Pi)) * math.Exp(-stirlingRemainder(a)-gammaDeviance(a, x))
	}
	lg, _ := math.Lgamma(a)
	return math.Exp(a*math.Log(x) - x - lg)
}

// stirlingRemainder returns δ(m) = ln Γ(m + 1) − (m + ½)·ln m + m − ½·ln 2π for real m > 15, from the
// Stirling series of DLMF 5.11.1: δ(m) = 1/(12m) − 1/(360m³) + 1/(1260m⁵) − 1/(1680m⁷) + 1/(1188m⁹) − …,
// cut where the first omitted term is at most about 2⁻⁵³·δ (Loader's thresholds 35, 80 and 500). δ enters
// an exponent, so its absolute error is a relative error of the same size in the prefactor.
func stirlingRemainder(m float64) float64 {
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

// gammaDeviance returns D₀(a, x) = a·ln(a/x) + x − a ≥ 0 for a > 0, x > 0 (Loader's bd0). When a and x
// are within 10% of their sum it sums the cancellation-free series
// D₀ = (a − x)·v + 2a·Σⱼ v^(2j+1)/(2j+1), v = (a − x)/(a + x), |v| < 0.1, so each term is under
// a hundredth of the one before and the loop ends within about 9 steps.
func gammaDeviance(a, x float64) float64 {
	if d := a - x; math.Abs(d) < 0.1*(a+x) {
		v := d / (a + x)
		s := d * v
		ej := 2 * a * v
		v2 := v * v
		for j := 3.0; j < 200; j += 2 {
			ej *= v2
			s1 := s + ej/j
			if s1 == s {
				return s1
			}
			s = s1
		}
		return s
	}
	r := a / x
	if r == 0 || math.IsInf(r, 0) {
		// a/x under- or overflowed: the difference of logarithms is still finite.
		return a*(math.Log(a)-math.Log(x)) + x - a
	}
	return a*math.Log(r) + x - a
}

// gammaSeries sums DLMF 8.7.1: P(a,x) = e^(−x) x^a/Γ(a) · Σ x^k/(a(a+1)…(a+k)).
func gammaSeries(a, x float64) (float64, error) {
	ap := a
	term := 1 / a
	sum := term
	limit := iterationCap(a)
	for n := 0; n < limit; n++ {
		ap++
		term *= x / ap
		sum += term
		if math.Abs(term) < math.Abs(sum)*epsilon {
			return math.Min(1, sum*gammaPrefactor(a, x)), nil
		}
	}
	return 0, fmt.Errorf("%w: series for P(a, x), a = %v, x = %v", ErrNoConvergence, a, x)
}

// gammaContinuedFraction evaluates Q(a,x) from the Legendre continued fraction DLMF 8.9.2 in its even
// form (Numerical Recipes 3rd ed. eq. 6.2.7) by the modified Lentz method.
func gammaContinuedFraction(a, x float64) (float64, error) {
	b := x + 1 - a
	c := 1 / tiny
	d := 1 / b
	h := d
	limit := iterationCap(a)
	for i := 1; i <= limit; i++ {
		an := -float64(i) * (float64(i) - a)
		b += 2
		d = an*d + b
		if math.Abs(d) < tiny {
			d = tiny
		}
		c = b + an/c
		if math.Abs(c) < tiny {
			c = tiny
		}
		d = 1 / d
		del := d * c
		h *= del
		if math.Abs(del-1) <= epsilon {
			return math.Min(1, gammaPrefactor(a, x)*h), nil
		}
	}
	return 0, fmt.Errorf("%w: continued fraction for Q(a, x), a = %v, x = %v", ErrNoConvergence, a, x)
}

// epsilon is the float64 unit roundoff 2⁻⁵² used as the stopping test of every iteration here.
const epsilon = 0x1p-52
