package business

// DebtToEquityRatio returns total liabilities divided by total stockholders' equity, as a decimal. It is not
// rounded. Negative equity is allowed and gives a negative ratio, which the caller must read as "liabilities
// exceed assets".
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when totalLiabilitiesCents is negative;
// ErrZeroDenominator when totalEquityCents is 0.
//
// ID:      FIN-BIZ-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.4 "Solvency Ratios", https://openstax.org/books/principles-finance/pages/6-4-solvency-ratios (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Source:  U.S. Securities and Exchange Commission, "Beginners' Guide to Financial Statements", https://www.sec.gov/about/reports-publications/beginners-guide-financial-statements (read 2026-11-12)
// Tags:    ratio, solvency, debt-to-equity, leverage
// Units:   totalLiabilitiesCents and totalEquityCents in cents; the result is a dimensionless decimal
// Example: Clear Lake Sporting Goods $150,000 ÷ $100,000 = 1.5 (Principles of Finance §6.4)
func DebtToEquityRatio(totalLiabilitiesCents, totalEquityCents int64) (float64, error) {
	if err := checkNonNegativeCents("totalLiabilitiesCents", totalLiabilitiesCents); err != nil {
		return 0, err
	}
	if err := checkCents("totalEquityCents", totalEquityCents); err != nil {
		return 0, err
	}
	return ratioCents(totalLiabilitiesCents, "totalEquityCents", totalEquityCents)
}
