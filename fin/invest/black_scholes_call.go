package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/prob"
)

// normalCDF returns N(x), the standard normal CDF, from NUM-PROB-048 prob.NormalCDF (½·erfc(−x/√2)).
// prob.NormalCDF rejects ±Inf, which d₁, d₂ reach when σ√T is infinite; there N is its limit, 1 at +Inf
// and 0 at −Inf (the values ½·erfc gives). x is never NaN here (blackScholesTerms returns ErrOverflow).
func normalCDF(x float64) float64 {
	if math.IsInf(x, 0) {
		if x > 0 {
			return 1
		}
		return 0
	}
	p, _ := prob.NormalCDF(x, 0, 1) // mean 0, sd 1 and finite x: no error is possible
	return p
}

// blackScholesTerms checks the Black-Scholes inputs and returns S, K·e^(−rT) and d₁, d₂. When
// σ√T = 0 the option is worth its discounted intrinsic value and degenerate is true (d₁, d₂ unused).
func blackScholesTerms(spotCents, strikeCents int64, riskFreeRatePerYear, volatilityPerYear, years float64) (s, kDisc, d1, d2 float64, degenerate bool, err error) {
	for _, a := range [...]struct {
		name  string
		cents int64
	}{{"spotCents", spotCents}, {"strikeCents", strikeCents}} {
		if err := checkNonNegativeCents(a.name, a.cents); err != nil {
			return 0, 0, 0, 0, false, err
		}
		if a.cents == 0 {
			return 0, 0, 0, 0, false, fmt.Errorf("%w: %s = 0 (ln(S/K) undefined)", ErrDomain, a.name)
		}
	}
	if err := checkRate("riskFreeRatePerYear", riskFreeRatePerYear); err != nil {
		return 0, 0, 0, 0, false, err
	}
	if err := checkFinite("volatilityPerYear", volatilityPerYear); err != nil {
		return 0, 0, 0, 0, false, err
	}
	if volatilityPerYear < 0 {
		return 0, 0, 0, 0, false, fmt.Errorf("%w: volatilityPerYear = %v", ErrNegative, volatilityPerYear)
	}
	if err := checkFinite("years", years); err != nil {
		return 0, 0, 0, 0, false, err
	}
	if years < 0 {
		return 0, 0, 0, 0, false, fmt.Errorf("%w: years = %v", ErrNegative, years)
	}
	s = float64(spotCents)
	k := float64(strikeCents)
	kDisc = k * math.Exp(-riskFreeRatePerYear*years)
	if math.IsInf(kDisc, 0) || math.IsNaN(kDisc) {
		return 0, 0, 0, 0, false, fmt.Errorf("%w: discounted strike %v is beyond float64", ErrOverflow, kDisc)
	}
	sd := volatilityPerYear * math.Sqrt(years)
	if sd == 0 {
		return s, kDisc, 0, 0, true, nil
	}
	// d₁,₂ = ln(S/K)/sd + r·T/sd ± sd/2, term by term so σ² is never formed (σ up to MaxFloat64);
	// r·T/sd is r·(√T/σ). An infinite sd gives d₁ = +Inf, d₂ = −Inf: call → S, put → K·e^(−rT).
	base := math.Log(s/k)/sd + riskFreeRatePerYear*(math.Sqrt(years)/volatilityPerYear)
	d1 = base + sd/2
	d2 = base - sd/2
	if math.IsNaN(d1) || math.IsNaN(d2) {
		return 0, 0, 0, 0, false, fmt.Errorf("%w: d1/d2 terms are beyond float64 (σ = %v, T = %v)", ErrOverflow, volatilityPerYear, years)
	}
	return s, kDisc, d1, d2, false, nil
}

// BlackScholesCall returns the Black-Scholes value of a European call on a share that pays no
// dividends, S·N(d₁) − K·e^(−rT)·N(d₂), rounded half-even to the cent, with
// d₁ = [ln(S/K) + (r + σ²/2)T]/(σ√T) and d₂ = d₁ − σ√T. r is the continuously compounded riskless
// rate per year, σ the volatility per year, T the time to expiry in years. When σ√T = 0 (T = 0 or
// σ = 0) the value is the limit max(S − K·e^(−rT), 0) (at T = 0 the paper's boundary condition (8));
// as σ√T → ∞ it tends to S. d₁ and d₂ are formed without σ², so any finite σ is accepted.
// N is computed as ½·erfc(−x/√2). Assumptions (paper p. 640): constant r and σ, log-normal prices,
// no dividends, exercise only at expiry, no costs.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value except
// when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents. (n ≈ 20 operations)
// Errors: ErrOutOfRange if spotCents or strikeCents > 2⁵³; ErrNegative if spotCents, strikeCents,
// volatilityPerYear or years < 0; ErrDomain if spotCents or strikeCents is 0; ErrRate if
// riskFreeRatePerYear is NaN, ±Inf or ≤ −1; ErrNotFinite if volatilityPerYear or years is NaN or ±Inf;
// ErrOverflow if the value exceeds 2⁵³ cents or d₁, d₂ cannot be formed in float64.
//
// ID:      FIN-INVEST-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Black and Scholes, "The Pricing of Options and Corporate Liabilities", Journal of Political Economy 81(3):637–654, 1973, eq. (13), p. 644, doi:10.1086/260062
// Tags:    invest, option, derivative, black-scholes, call, european-option, valuation
// Units:   spotCents, strikeCents and the result in cents; riskFreeRatePerYear decimal per year (continuous); volatilityPerYear decimal per √year; years years
// Example: hand-checked; same inputs as Hull's textbook example (section not confirmed): S $42, K $40, r 10 %, σ 20 %, T 0.5 → 475.942 cents → $4.76
func BlackScholesCall(spotCents, strikeCents int64, riskFreeRatePerYear, volatilityPerYear, years float64) (int64, error) {
	s, kDisc, d1, d2, degenerate, err := blackScholesTerms(spotCents, strikeCents, riskFreeRatePerYear, volatilityPerYear, years)
	if err != nil {
		return 0, err
	}
	if degenerate {
		return roundCents(math.Max(s-kDisc, 0))
	}
	return roundCents(s*normalCDF(d1) - kDisc*normalCDF(d2))
}
