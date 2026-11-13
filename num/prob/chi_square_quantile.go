package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"
)

// chiSquareQuantileMaxSteps caps the safeguarded Newton iteration (a defect guard: it converges in under
// 60 steps; every step either keeps a Newton step inside the bracket or halves the bracket).
const chiSquareQuantileMaxSteps = 400

// ChiSquareQuantile returns the u-quantile of the chi-square distribution with df degrees of freedom: the x
// with P(X ≤ x) = u, so u = 0.95 gives the upper critical value χ²(0.95, df) of a test at α = 0.05. It
// solves ChiSquareCDF(x, df) = u (for u > ½, ChiSquareSurvival(x, df) = 1 − u, which keeps the upper
// tail's accuracy) by Newton's method on a bracket: each step uses the density
// f(x) = x^(ν/2 − 1)·e^(−x/2)/(2^(ν/2)·Γ(ν/2)) and falls back to bisection (geometric once the lower end
// is positive) when the step leaves the bracket. It stops when a step changes x by at most 4 ulps.
// The first iterate is x₀ = 2·(u·Γ(ν/2 + 1))^(2/ν) (the leading term of the series, DLMF 8.7.1) when that
// is below df, else df; the result does not depend on it. Below x₀ = 1e-20, x₀ is the quantile to double
// precision (the next term is smaller by a factor of about x₀) and is returned; it is 0 when the quantile
// is below the smallest float64.
// u = 1 (an infinite quantile) is an error, not +Inf.
// Errors: ErrNotFinite when u or df is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 1, df ≤ 0 or df > 10⁶.
//
// ID:      NUM-PROB-065
// Time:    O(√df) per step, at most 400 steps (under 60 in practice)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.6 "Chi-Square Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3666.htm (read 2026-11-12) (pdf and cdf; the percent point function has no closed form)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.7.4 "Critical Values of the Chi-Square Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3674.htm (read 2026-11-12) (the table the tests reproduce)
// Source:  Burden, Faires, Numerical Analysis, 9th ed., Brooks/Cole, 2011, §2.3 (Newton's method) and §2.1 (bisection)
// Tags:    probability, distribution, chi-square, quantile, inverse-cdf, ppf, critical-value
// Units:   u a probability in [0, 1); df dimensionless; result dimensionless, ≥ 0
// Range:   0 < df ≤ 10⁶
// Example: e-Handbook §1.3.6.7.4: u = 0.95, df = 3 → 7.815; u = 0.95, df = 6 → 12.592; u = 0.99, df = 10 → 23.209 (table to 3 decimals)
func ChiSquareQuantile(u, df float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	if err := checkChiSquareDF(df); err != nil {
		return 0, err
	}
	if u == 0 {
		return 0, nil
	}
	a := df / 2
	upper := u > 0.5
	tail := u
	if upper {
		tail = 1 - u // exact for u in (½, 1) (Sterbenz)
	}
	// h(x) is increasing in x and 0 at the quantile.
	h := func(x float64) (float64, error) {
		if upper {
			q, err := special.RegularizedGammaQ(a, x/2)
			if err != nil {
				return 0, specialErr("chi-square quantile", err)
			}
			return tail - q, nil
		}
		p, err := special.RegularizedGammaP(a, x/2)
		if err != nil {
			return 0, specialErr("chi-square quantile", err)
		}
		return p - u, nil
	}

	// Bracket [lo, hi] with h(lo) < 0 ≤ h(hi).
	lo, hi := 0.0, math.Max(df, 1)
	for {
		v, err := h(hi)
		if err != nil {
			return 0, err
		}
		if v >= 0 {
			break
		}
		lo, hi = hi, 2*hi
		if math.IsInf(hi, 0) {
			return 0, fmt.Errorf("%w: quantile for u = %v, df = %v", ErrOverflow, u, df)
		}
	}

	x := df
	if !upper {
		// P(a, z) = zᵃ/Γ(a + 1)·(1 + O(z)) (DLMF 8.7.1, k = 0), so below z ≈ 1e-20 the leading term is the
		// quantile to double precision (and 0 when it underflows).
		lg, _ := math.Lgamma(a + 1)
		x0 := 2 * math.Exp((math.Log(u)+lg)/a)
		if x0 < 1e-20 {
			return x0, nil
		}
		if x0 < df {
			x = x0
		}
	}
	if !(x > lo && x < hi) {
		x = lo + (hi-lo)/2
	}
	for range chiSquareQuantileMaxSteps {
		v, err := h(x)
		if err != nil {
			return 0, err
		}
		if v == 0 {
			return x, nil
		}
		if v < 0 {
			lo = x
		} else {
			hi = x
		}
		// Density f(x) = (x/2)^(a−1)·e^(−x/2)/(2·Γ(a)), the chi-square pdf (NUM-PROB-096 ChiSquarePDF); it only
		// sizes the Newton step, so an error (none occurs for x > 0) falls back to bisection like f = 0.
		density, derr := ChiSquarePDF(x, df)
		if derr != nil {
			density = 0
		}
		next := x - v/density
		if density == 0 || !(next > lo && next < hi) {
			if lo > 0 {
				next = math.Sqrt(lo) * math.Sqrt(hi)
			} else {
				next = hi / 2
			}
		}
		if math.Abs(next-x) <= 4*0x1p-52*next || hi-lo <= 4*0x1p-52*hi {
			return next, nil
		}
		x = next
	}
	return x, nil
}
