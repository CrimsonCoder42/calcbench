package business

// DebtRatio returns total liabilities divided by total assets (debt-to-assets), as a decimal. It is not
// rounded.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when totalLiabilitiesCents or
// totalAssetsCents is negative; ErrZeroDenominator when totalAssetsCents is 0.
//
// ID:      FIN-BIZ-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.4 "Solvency Ratios", https://openstax.org/books/principles-finance/pages/6-4-solvency-ratios (read 2026-11-12)
// Tags:    ratio, solvency, debt-ratio, debt-to-assets, leverage
// Units:   totalLiabilitiesCents and totalAssetsCents in cents; the result is a dimensionless decimal
// Example: Clear Lake Sporting Goods ($100,000 + $50,000) ÷ $250,000 = 0.6 (Principles of Finance §6.4)
func DebtRatio(totalLiabilitiesCents, totalAssetsCents int64) (float64, error) {
	if err := checkNonNegativeCents("totalLiabilitiesCents", totalLiabilitiesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("totalAssetsCents", totalAssetsCents); err != nil {
		return 0, err
	}
	return ratioCents(totalLiabilitiesCents, "totalAssetsCents", totalAssetsCents)
}
