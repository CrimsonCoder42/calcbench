package business

// BookValuePerShare returns (total equity − preferred equity) ÷ total shares outstanding, in cents per
// share, rounded half-even with exact integer arithmetic. Negative equity gives a negative book value.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³ or sharesOutstanding > 2⁵³; ErrNegative when
// preferredEquityCents or sharesOutstanding is negative; ErrZeroDenominator when sharesOutstanding is 0;
// ErrOverflow when the result is outside ±2⁵³ cents.
//
// ID:      FIN-BIZ-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.5 "Market Value Ratios" (book value per share), https://openstax.org/books/principles-finance/pages/6-5-market-value-ratios (read 2026-11-12)
// Tags:    per-share, market-value, book-value, book-value-per-share
// Units:   totalEquityCents and preferredEquityCents in cents; sharesOutstanding a share count; the result in cents per share
// Example: Clear Lake Sporting Goods ($100,000 − $20,000) ÷ 8,000 = $10 (Principles of Finance §6.5)
func BookValuePerShare(totalEquityCents, preferredEquityCents, sharesOutstanding int64) (int64, error) {
	if err := checkCents("totalEquityCents", totalEquityCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("preferredEquityCents", preferredEquityCents); err != nil {
		return 0, err
	}
	if err := checkShares("sharesOutstanding", sharesOutstanding); err != nil {
		return 0, err
	}
	bvps := divRoundHalfEven(totalEquityCents-preferredEquityCents, sharesOutstanding)
	if err := checkResultCents("book value per share", bvps); err != nil {
		return 0, err
	}
	return bvps, nil
}
