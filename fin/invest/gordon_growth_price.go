package invest

import "fmt"

// GordonGrowthPrice returns the value of a share whose dividend grows at a constant rate forever,
// P₀ = D₁/(r − g), rounded half-even to the cent. D₁ is the NEXT dividend (one period from now);
// from the last dividend paid, pass D₀·(1 + g). requiredReturn and growthRate are per the same period.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value except
// when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents. (n = 1)
// Errors: ErrOutOfRange if |nextDividendCents| > 2⁵³; ErrNegative if nextDividendCents < 0; ErrRate
// if a rate is NaN, ±Inf or ≤ −1; ErrDomain if requiredReturn ≤ growthRate (no finite value);
// ErrOverflow if the price exceeds 2⁵³ cents.
//
// ID:      FIN-INVEST-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §11.2 "Dividend Discount Models (DDMs)" (constant growth / Gordon growth model, Stock Value = D1/(r − g)), https://openstax.org/books/principles-finance/pages/11-2-dividend-discount-models-ddms
// Source:  Gordon and Shapiro, "Capital Equipment Analysis: The Required Rate of Profit", Management Science 3(1):102–110, 1956, doi:10.1287/mnsc.3.1.102 (the original; DOI confirmed on Crossref)
// Tags:    invest, stock, valuation, dividend-discount-model, gordon-growth, constant-growth, terminal-value
// Units:   nextDividendCents and the result in cents per share; requiredReturn and growthRate decimal per period
// Example: OpenStax §11.2: D₀ $5.00, g 4 %, r 8 % → D₁ $5.20 → $130.00
func GordonGrowthPrice(nextDividendCents int64, requiredReturn, growthRate float64) (int64, error) {
	if err := checkNonNegativeCents("nextDividendCents", nextDividendCents); err != nil {
		return 0, err
	}
	if err := checkRate("requiredReturn", requiredReturn); err != nil {
		return 0, err
	}
	if err := checkRate("growthRate", growthRate); err != nil {
		return 0, err
	}
	if requiredReturn <= growthRate {
		return 0, fmt.Errorf("%w: requiredReturn = %v must exceed growthRate = %v", ErrDomain, requiredReturn, growthRate)
	}
	return roundCents(float64(nextDividendCents) / (requiredReturn - growthRate))
}
