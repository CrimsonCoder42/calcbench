package invest

import "math"

// BlackScholesPut returns the Black-Scholes value of a European put on a share that pays no
// dividends, −S·N(−d₁) + K·e^(−rT)·N(−d₂), rounded half-even to the cent (d₁, d₂ as in
// BlackScholesCall). It is evaluated directly, not as call − S + K·e^(−rT), so a deep
// out-of-the-money put keeps its digits; the two agree by put-call parity (paper eq. (25)). When
// σ√T = 0 the value is the limit max(K·e^(−rT) − S, 0); as σ√T → ∞ it tends to K·e^(−rT).
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value except
// when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents. (n ≈ 20 operations)
// Errors: ErrOutOfRange if spotCents or strikeCents > 2⁵³; ErrNegative if spotCents, strikeCents,
// volatilityPerYear or years < 0; ErrDomain if spotCents or strikeCents is 0; ErrRate if
// riskFreeRatePerYear is NaN, ±Inf or ≤ −1; ErrNotFinite if volatilityPerYear or years is NaN or ±Inf;
// ErrOverflow if the value exceeds 2⁵³ cents or d₁, d₂ cannot be formed in float64.
//
// ID:      FIN-INVEST-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Black and Scholes, "The Pricing of Options and Corporate Liabilities", Journal of Political Economy 81(3):637–654, 1973, eqs. (25)–(27), p. 647, doi:10.1086/260062
// Tags:    invest, option, derivative, black-scholes, put, european-option, valuation, put-call-parity
// Units:   spotCents, strikeCents and the result in cents; riskFreeRatePerYear decimal per year (continuous); volatilityPerYear decimal per √year; years years
// Example: hand-checked; same inputs as Hull's textbook example (section not confirmed): S $42, K $40, r 10 %, σ 20 %, T 0.5 → 80.860 cents → $0.81
func BlackScholesPut(spotCents, strikeCents int64, riskFreeRatePerYear, volatilityPerYear, years float64) (int64, error) {
	s, kDisc, d1, d2, degenerate, err := blackScholesTerms(spotCents, strikeCents, riskFreeRatePerYear, volatilityPerYear, years)
	if err != nil {
		return 0, err
	}
	if degenerate {
		return roundCents(math.Max(kDisc-s, 0))
	}
	return roundCents(-s*normalCDF(-d1) + kDisc*normalCDF(-d2))
}
