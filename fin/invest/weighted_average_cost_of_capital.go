package invest

import "fmt"

// WeightedAverageCostOfCapital returns the firm's WACC with debt and common equity (no preferred),
// E/(D+E)·re + D/(D+E)·rd·(1 − T), as a decimal per period. Values are market values in cents; the
// weights are formed exactly from them.
// Errors: ErrOutOfRange if a value is above 2⁵³ cents; ErrNegative if a value is negative;
// ErrZeroDenominator if both values are 0; ErrRate if costOfEquity or costOfDebt is NaN, ±Inf or ≤ −1;
// ErrNotFinite if taxRate is NaN or ±Inf; ErrDomain if taxRate is outside [0, 1]; ErrOverflow if
// the result exceeds float64.
//
// ID:      FIN-INVEST-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §17.3 "Calculating the Weighted Average Cost of Capital" (WACC = D%×rd(1 − T) + E%×re), https://openstax.org/books/principles-finance/pages/17-3-calculating-the-weighted-average-cost-of-capital
// Tags:    invest, corporate-finance, wacc, cost-of-capital, discount-rate, capital-structure
// Units:   equityValueCents and debtValueCents in cents; costOfEquity, costOfDebt decimal per period; taxRate decimal; result decimal per period
// Example: OpenStax §17.1–17.3 Bluebonnet: E $15,000,000, D $4,850,000, re 13.4 %, rd 6.312 %, T 21 % → 11.33 % (exact weights 0.1134430)
func WeightedAverageCostOfCapital(equityValueCents, debtValueCents int64, costOfEquity, costOfDebt, taxRate float64) (float64, error) {
	if err := checkNonNegativeCents("equityValueCents", equityValueCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("debtValueCents", debtValueCents); err != nil {
		return 0, err
	}
	if equityValueCents == 0 && debtValueCents == 0 {
		return 0, fmt.Errorf("%w: equityValueCents + debtValueCents = 0", ErrZeroDenominator)
	}
	if err := checkRate("costOfEquity", costOfEquity); err != nil {
		return 0, err
	}
	if err := checkRate("costOfDebt", costOfDebt); err != nil {
		return 0, err
	}
	if err := checkFinite("taxRate", taxRate); err != nil {
		return 0, err
	}
	if taxRate < 0 || taxRate > 1 {
		return 0, fmt.Errorf("%w: taxRate = %v, want 0 … 1", ErrDomain, taxRate)
	}
	e, d := float64(equityValueCents), float64(debtValueCents)
	return finiteResult((e*costOfEquity + d*costOfDebt*(1-taxRate)) / (e + d))
}
