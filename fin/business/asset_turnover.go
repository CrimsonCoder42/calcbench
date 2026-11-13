package business

// AssetTurnover returns net sales divided by average total assets (total asset turnover), as a decimal. It
// is not rounded.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when netSalesCents or
// averageTotalAssetsCents is negative; ErrZeroDenominator when averageTotalAssetsCents is 0.
//
// ID:      FIN-BIZ-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.2 "Operating Efficiency Ratios", https://openstax.org/books/principles-finance/pages/6-2-operating-efficiency-ratios (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, efficiency, activity, asset-turnover, dupont
// Units:   netSalesCents and averageTotalAssetsCents in cents; the result is a dimensionless count of times per period
// Example: Clear Lake Sporting Goods $120,000 ÷ $225,000 = 0.533…, printed 0.53 (Principles of Finance §6.2)
func AssetTurnover(netSalesCents, averageTotalAssetsCents int64) (float64, error) {
	if err := checkNonNegativeCents("netSalesCents", netSalesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("averageTotalAssetsCents", averageTotalAssetsCents); err != nil {
		return 0, err
	}
	return ratioCents(netSalesCents, "averageTotalAssetsCents", averageTotalAssetsCents)
}
