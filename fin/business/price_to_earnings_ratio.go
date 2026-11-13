package business

// PriceToEarningsRatio returns market price per share divided by earnings per share, as a decimal. It is not
// rounded. Negative earnings are allowed and give a negative ratio, which the caller decides how to report.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when pricePerShareCents is negative;
// ErrZeroDenominator when earningsPerShareCents is 0.
//
// ID:      FIN-BIZ-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.5 "Market Value Ratios", https://openstax.org/books/principles-finance/pages/6-5-market-value-ratios (read 2026-11-12)
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §11.1 "Multiple Approaches to Stock Valuation", https://openstax.org/books/principles-finance/pages/11-1-multiple-approaches-to-stock-valuation (read 2026-11-12)
// Tags:    ratio, market-value, valuation, price-earnings, pe
// Units:   pricePerShareCents and earningsPerShareCents in cents per share; the result is a dimensionless decimal
// Example: Clear Lake Sporting Goods $69.41 ÷ $3.75 = 18.509…, printed 18.51 (Principles of Finance §6.5)
func PriceToEarningsRatio(pricePerShareCents, earningsPerShareCents int64) (float64, error) {
	if err := checkNonNegativeCents("pricePerShareCents", pricePerShareCents); err != nil {
		return 0, err
	}
	if err := checkCents("earningsPerShareCents", earningsPerShareCents); err != nil {
		return 0, err
	}
	return ratioCents(pricePerShareCents, "earningsPerShareCents", earningsPerShareCents)
}
