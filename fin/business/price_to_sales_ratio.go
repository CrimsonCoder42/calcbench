package business

// PriceToSalesRatio returns market capitalization divided by revenue, as a decimal. It is not rounded. The
// source defines it per share (share price ÷ sales per share); multiplying both by the share count gives
// this whole-firm form, the same ratio.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when marketCapitalizationCents or
// revenueCents is negative; ErrZeroDenominator when revenueCents is 0.
//
// ID:      FIN-BIZ-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §11.1 "Multiple Approaches to Stock Valuation", https://openstax.org/books/principles-finance/pages/11-1-multiple-approaches-to-stock-valuation (read 2026-11-12)
// Tags:    ratio, market-value, valuation, price-to-sales
// Units:   marketCapitalizationCents and revenueCents in cents; the result is a dimensionless decimal
// Example: $50,000,000 ÷ $20,000,000 = 2.5 (hand-checked)
func PriceToSalesRatio(marketCapitalizationCents, revenueCents int64) (float64, error) {
	if err := checkNonNegativeCents("marketCapitalizationCents", marketCapitalizationCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("revenueCents", revenueCents); err != nil {
		return 0, err
	}
	return ratioCents(marketCapitalizationCents, "revenueCents", revenueCents)
}
