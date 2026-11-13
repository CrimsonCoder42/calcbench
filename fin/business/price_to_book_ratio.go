package business

// PriceToBookRatio returns market price per share divided by book value per share (market-to-book), as a
// decimal. It is not rounded. A negative book value is allowed and gives a negative ratio.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when pricePerShareCents is negative;
// ErrZeroDenominator when bookValuePerShareCents is 0.
//
// ID:      FIN-BIZ-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §11.1 "Multiple Approaches to Stock Valuation", https://openstax.org/books/principles-finance/pages/11-1-multiple-approaches-to-stock-valuation (read 2026-11-12)
// Tags:    ratio, market-value, valuation, price-to-book, market-to-book
// Units:   pricePerShareCents and bookValuePerShareCents in cents per share; the result is a dimensionless decimal
// Example: $69.41 ÷ $10.00 = 6.941 (hand-checked from the Clear Lake numbers of Principles of Finance §6.5; §11.1 has no numbers)
func PriceToBookRatio(pricePerShareCents, bookValuePerShareCents int64) (float64, error) {
	if err := checkNonNegativeCents("pricePerShareCents", pricePerShareCents); err != nil {
		return 0, err
	}
	if err := checkCents("bookValuePerShareCents", bookValuePerShareCents); err != nil {
		return 0, err
	}
	return ratioCents(pricePerShareCents, "bookValuePerShareCents", bookValuePerShareCents)
}
