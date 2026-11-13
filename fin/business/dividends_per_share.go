package business

// DividendsPerShare returns total cash dividends ÷ shares outstanding (issued less treasury shares), in
// cents per share, rounded half-even with exact integer arithmetic.
// Errors: ErrOutOfRange when totalDividendsCents is outside ±2⁵³ or sharesOutstanding > 2⁵³; ErrNegative
// when an input is negative; ErrZeroDenominator when sharesOutstanding is 0.
//
// ID:      FIN-BIZ-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §14.3 "Record Transactions and the Effects on Financial Statements for Cash Dividends, Property Dividends, Stock Dividends, and Stock Splits" (total cash dividend = per-share amount × shares outstanding), https://openstax.org/books/principles-financial-accounting/pages/14-3-record-transactions-and-the-effects-on-financial-statements-for-cash-dividends-property-dividends-stock-dividends-and-stock-splits (read 2026-11-12)
// Tags:    per-share, dividends, dividends-per-share
// Units:   totalDividendsCents in cents; sharesOutstanding a share count; the result in cents per share
// Example: La Cantina $4,600 ÷ 9,200 shares = $0.50 (Principles of Accounting Vol. 1 §14.3: 9,200 × $0.50 = $4,600)
func DividendsPerShare(totalDividendsCents, sharesOutstanding int64) (int64, error) {
	if err := checkNonNegativeCents("totalDividendsCents", totalDividendsCents); err != nil {
		return 0, err
	}
	if err := checkShares("sharesOutstanding", sharesOutstanding); err != nil {
		return 0, err
	}
	return divRoundHalfEven(totalDividendsCents, sharesOutstanding), nil
}
