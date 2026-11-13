package invest

import "fmt"

// HoldingPeriodYield returns the total return over a holding period, (P₁ − P₀ + D)/P₀ as a decimal:
// the capital gain plus the income received, divided by the beginning value. It equals the capital
// gains yield plus the dividend yield, and is the total stock return. incomeCents may be negative
// (a net charge). The numerator is computed exactly in int64 before the one division.
// Errors: ErrOutOfRange if an amount is outside ±2⁵³ cents; ErrNegative if beginningValueCents or
// endingValueCents < 0; ErrZeroDenominator if beginningValueCents is 0.
//
// ID:      FIN-INVEST-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §15.1 "Risk and Return to an Individual Asset" (Total Percent Return = Dividend Yield + Capital Gain Yield), https://openstax.org/books/principles-finance/pages/15-1-risk-and-return-to-an-individual-asset
// Tags:    invest, return, holding-period-return, holding-period-yield, total-return, stock
// Units:   all amounts in cents; result decimal for the holding period (not annualized)
// Example: OpenStax §15.1 Target: bought $128.74, sold $176.53, dividends $2.68 → $50.47, 39.20 %
func HoldingPeriodYield(beginningValueCents, endingValueCents, incomeCents int64) (float64, error) {
	if err := checkNonNegativeCents("beginningValueCents", beginningValueCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("endingValueCents", endingValueCents); err != nil {
		return 0, err
	}
	if err := checkCents("incomeCents", incomeCents); err != nil {
		return 0, err
	}
	if beginningValueCents == 0 {
		return 0, fmt.Errorf("%w: beginningValueCents = 0", ErrZeroDenominator)
	}
	gain := endingValueCents - beginningValueCents + incomeCents // |gain| ≤ 3·2⁵³, no int64 overflow
	return float64(gain) / float64(beginningValueCents), nil
}
