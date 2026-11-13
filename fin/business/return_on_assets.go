package business

// ReturnOnAssets returns net income divided by average total assets, as a decimal. It is not rounded. A net
// loss is allowed and gives a negative return.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when averageTotalAssetsCents is negative;
// ErrZeroDenominator when averageTotalAssetsCents is 0.
//
// ID:      FIN-BIZ-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.6 "Profitability Ratios and the DuPont Method", https://openstax.org/books/principles-finance/pages/6-6-profitability-ratios-and-the-dupont-method (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, profitability, return, roa
// Units:   netIncomeCents and averageTotalAssetsCents in cents; the result is a decimal per period
// Example: Clear Lake Sporting Goods $35,000 ÷ $225,000 = 0.1556, printed 0.16 (Principles of Finance §6.6)
func ReturnOnAssets(netIncomeCents, averageTotalAssetsCents int64) (float64, error) {
	if err := checkCents("netIncomeCents", netIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("averageTotalAssetsCents", averageTotalAssetsCents); err != nil {
		return 0, err
	}
	return ratioCents(netIncomeCents, "averageTotalAssetsCents", averageTotalAssetsCents)
}
