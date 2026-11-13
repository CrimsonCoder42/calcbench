package invest

import "fmt"

// CapitalGainsYield returns the change in price over the holding period divided by the beginning
// price, (P₁ − P₀)/P₀, as a decimal; a fall in price gives a negative yield.
// Variants: return on investment, (final − cost)/cost (FIN-BIZ-033 was its duplicate).
// Errors: ErrOutOfRange if a price is above 2⁵³ cents; ErrNegative if a price is negative;
// ErrZeroDenominator if beginningPriceCents is 0.
//
// ID:      FIN-INVEST-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §15.1 "Risk and Return to an Individual Asset" (capital gain yield), https://openstax.org/books/principles-finance/pages/15-1-risk-and-return-to-an-individual-asset
// Tags:    invest, return, capital-gains-yield, price-appreciation, stock, roi, return-on-investment
// Units:   prices in cents; result decimal for the holding period
// Example: OpenStax §15.1 3M: $175.00 → $169.72 → −3.02 %; Target $128.74 → $176.53 → 0.3712
func CapitalGainsYield(beginningPriceCents, endingPriceCents int64) (float64, error) {
	if err := checkNonNegativeCents("beginningPriceCents", beginningPriceCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("endingPriceCents", endingPriceCents); err != nil {
		return 0, err
	}
	if beginningPriceCents == 0 {
		return 0, fmt.Errorf("%w: beginningPriceCents = 0", ErrZeroDenominator)
	}
	return float64(endingPriceCents-beginningPriceCents) / float64(beginningPriceCents), nil
}
