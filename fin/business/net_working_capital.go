package business

// NetWorkingCapital returns current assets minus current liabilities, in cents, exactly. A negative
// result means current liabilities exceed current assets.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative. (Both
// inputs are in [0, 2⁵³], so the difference cannot overflow.)
//
// ID:      FIN-BIZ-037
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis" (working capital), https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §19.1 "What Is Working Capital?", https://openstax.org/books/principles-finance/pages/19-1-what-is-working-capital (read 2026-11-12)
// Tags:    liquidity, working-capital, net-working-capital, balance-sheet
// Units:   currentAssetsCents, currentLiabilitiesCents and the result in cents
// Example: Banyan Goods $200,000 − $100,000 = $100,000 (Principles of Accounting Vol. 1, Appendix A)
func NetWorkingCapital(currentAssetsCents, currentLiabilitiesCents int64) (int64, error) {
	if err := checkNonNegativeCents("currentAssetsCents", currentAssetsCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("currentLiabilitiesCents", currentLiabilitiesCents); err != nil {
		return 0, err
	}
	return currentAssetsCents - currentLiabilitiesCents, nil
}
