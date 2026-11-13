package business

// GrossProfitMargin returns (net sales − cost of goods sold) divided by net sales, as a decimal
// (0.39 = 39 %). The subtraction is exact; the result has one rounding and is not rounded further. Cost
// above sales gives a negative margin.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative;
// ErrZeroDenominator when netSalesCents is 0.
//
// ID:      FIN-BIZ-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §6.6 "Describe and Prepare Multi-Step and Simple Income Statements for Merchandising Companies" (gross profit margin ratio), https://openstax.org/books/principles-financial-accounting/pages/6-6-describe-and-prepare-multi-step-and-simple-income-statements-for-merchandising-companies (read 2026-11-12)
// Tags:    ratio, profitability, margin, gross-margin, gross-profit
// Units:   netSalesCents and costOfGoodsSoldCents in cents; the result is a dimensionless decimal
// Example: California Business Solutions ($293,500 − $180,000) ÷ $293,500 = 0.3867, printed 0.39 (Principles of Accounting Vol. 1 §6.6)
func GrossProfitMargin(netSalesCents, costOfGoodsSoldCents int64) (float64, error) {
	if err := checkNonNegativeCents("netSalesCents", netSalesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("costOfGoodsSoldCents", costOfGoodsSoldCents); err != nil {
		return 0, err
	}
	return ratioCents(netSalesCents-costOfGoodsSoldCents, "netSalesCents", netSalesCents)
}
