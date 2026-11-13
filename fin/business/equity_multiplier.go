package business

// EquityMultiplier returns total assets divided by total stockholders' equity, as a decimal. It is not
// rounded. The source uses average total assets and average equity; pass averages to follow it. Negative
// equity is allowed and gives a negative multiplier.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when totalAssetsCents is negative;
// ErrZeroDenominator when totalEquityCents is 0.
//
// ID:      FIN-BIZ-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.6 "Profitability Ratios and the DuPont Method", https://openstax.org/books/principles-finance/pages/6-6-profitability-ratios-and-the-dupont-method (read 2026-11-12)
// Tags:    ratio, solvency, equity-multiplier, dupont, leverage
// Units:   totalAssetsCents and totalEquityCents in cents; the result is a dimensionless decimal
// Example: Clear Lake Sporting Goods $225,000 ÷ $95,000 = 2.3684…, printed 2.37 (Principles of Finance §6.6, DuPont)
func EquityMultiplier(totalAssetsCents, totalEquityCents int64) (float64, error) {
	if err := checkNonNegativeCents("totalAssetsCents", totalAssetsCents); err != nil {
		return 0, err
	}
	if err := checkCents("totalEquityCents", totalEquityCents); err != nil {
		return 0, err
	}
	return ratioCents(totalAssetsCents, "totalEquityCents", totalEquityCents)
}
