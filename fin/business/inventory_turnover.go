package business

// InventoryTurnover returns cost of goods sold divided by average inventory, as a decimal. It is not
// rounded.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when costOfGoodsSoldCents or
// averageInventoryCents is negative; ErrZeroDenominator when averageInventoryCents is 0.
//
// ID:      FIN-BIZ-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.2 "Operating Efficiency Ratios", https://openstax.org/books/principles-finance/pages/6-2-operating-efficiency-ratios (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, efficiency, activity, inventory-turnover
// Units:   costOfGoodsSoldCents and averageInventoryCents in cents; the result is a dimensionless count of times per period
// Example: Clear Lake Sporting Goods $60,000 ÷ $37,500 = 1.6 (Principles of Finance §6.2)
func InventoryTurnover(costOfGoodsSoldCents, averageInventoryCents int64) (float64, error) {
	if err := checkNonNegativeCents("costOfGoodsSoldCents", costOfGoodsSoldCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("averageInventoryCents", averageInventoryCents); err != nil {
		return 0, err
	}
	return ratioCents(costOfGoodsSoldCents, "averageInventoryCents", averageInventoryCents)
}
