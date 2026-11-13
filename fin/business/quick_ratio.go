package business

// QuickRatio returns (cash + marketable securities + receivables) divided by current liabilities (the
// acid-test ratio), as a decimal. Inventory is not in the numerator. The three-term sum is exact in int64;
// when it exceeds 2⁵³ its conversion to float64 adds one rounding before the division.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative;
// ErrZeroDenominator when currentLiabilitiesCents is 0.
//
// ID:      FIN-BIZ-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.3 "Liquidity Ratios" (quick ratio = cash + short-term investments + accounts receivable over current liabilities), https://openstax.org/books/principles-finance/pages/6-3-liquidity-ratios (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, liquidity, quick-ratio, acid-test, balance-sheet
// Units:   every input in cents (marketableSecuritiesCents = short-term investments); the result is a dimensionless decimal
// Example: Clear Lake Sporting Goods ($110,000 + $20,000 + $30,000) ÷ $100,000 = 1.6 (Principles of Finance §6.3)
func QuickRatio(cashCents, marketableSecuritiesCents, receivablesCents, currentLiabilitiesCents int64) (float64, error) {
	if err := checkNonNegativeCents("cashCents", cashCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("marketableSecuritiesCents", marketableSecuritiesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("receivablesCents", receivablesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("currentLiabilitiesCents", currentLiabilitiesCents); err != nil {
		return 0, err
	}
	quick := cashCents + marketableSecuritiesCents + receivablesCents // ≤ 3·2⁵³, no int64 overflow
	return ratioCents(quick, "currentLiabilitiesCents", currentLiabilitiesCents)
}
