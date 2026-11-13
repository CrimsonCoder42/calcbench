package business

// OperatingProfitMargin returns income from operations divided by net revenues, as a decimal (0.25 = 25 %).
// It is not rounded. A negative operating income is allowed and gives a negative margin.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when netSalesCents is negative;
// ErrZeroDenominator when netSalesCents is 0.
//
// ID:      FIN-BIZ-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  U.S. Securities and Exchange Commission, "Beginners' Guide to Financial Statements", https://www.sec.gov/about/reports-publications/beginners-guide-financial-statements (read 2026-11-12)
// Tags:    ratio, profitability, margin, operating-margin
// Units:   operatingIncomeCents and netSalesCents in cents; the result is a dimensionless decimal
// Example: $25,000 ÷ $100,000 = 0.25 (hand-checked; the source gives no numbers)
func OperatingProfitMargin(operatingIncomeCents, netSalesCents int64) (float64, error) {
	if err := checkCents("operatingIncomeCents", operatingIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("netSalesCents", netSalesCents); err != nil {
		return 0, err
	}
	return ratioCents(operatingIncomeCents, "netSalesCents", netSalesCents)
}
