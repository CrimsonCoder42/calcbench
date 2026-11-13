package invest

import "fmt"

// DividendYield returns the annual dividend per share divided by the price per share, as a decimal.
// For a holding-period return the price is the beginning price (OpenStax §15.1).
// Errors: ErrOutOfRange if an amount is above 2⁵³ cents; ErrNegative if an amount is negative;
// ErrZeroDenominator if pricePerShareCents is 0.
//
// ID:      FIN-INVEST-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §15.1 "Risk and Return to an Individual Asset" (dividend yield = dividends received / initial stock price), https://openstax.org/books/principles-finance/pages/15-1-risk-and-return-to-an-individual-asset
// Tags:    invest, stock, dividend, dividend-yield, yield
// Units:   amounts in cents per share; result decimal per year
// Example: OpenStax §15.1 3M: $5.91 dividends on $175.00 → 3.38 %; Target $2.68 on $128.74 → 0.0208
func DividendYield(annualDividendPerShareCents, pricePerShareCents int64) (float64, error) {
	if err := checkNonNegativeCents("annualDividendPerShareCents", annualDividendPerShareCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("pricePerShareCents", pricePerShareCents); err != nil {
		return 0, err
	}
	if pricePerShareCents == 0 {
		return 0, fmt.Errorf("%w: pricePerShareCents = 0", ErrZeroDenominator)
	}
	return float64(annualDividendPerShareCents) / float64(pricePerShareCents), nil
}
