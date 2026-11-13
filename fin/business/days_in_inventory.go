package business

import "fmt"

// DaysInInventory returns the number of days of cost of goods sold held in inventory:
// daysInPeriod × inventory ÷ cost of goods sold, which equals daysInPeriod ÷ inventory turnover
// (FIN-BIZ-016). Principles of Finance §19.1 uses average inventory (via the turnover); §6.2 and Principles
// of Accounting Appendix A use ending inventory: pass the figure your convention uses. At most two
// roundings; the result is not rounded to whole days.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when a money input is negative;
// ErrPeriods when daysInPeriod < 1; ErrZeroDenominator when costOfGoodsSoldCents is 0.
//
// ID:      FIN-BIZ-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §19.1 "What Is Working Capital?" (inventory conversion period = 365 ÷ inventory turnover), https://openstax.org/books/principles-finance/pages/19-1-what-is-working-capital (read 2026-11-12)
// Source:  Dahlquist and Knight, Principles of Finance, 2022, §6.2 "Operating Efficiency Ratios" (days' sales in inventory = ending inventory ÷ cost of goods sold × 365), https://openstax.org/books/principles-finance/pages/6-2-operating-efficiency-ratios (read 2026-11-12)
// Tags:    ratio, efficiency, activity, inventory, days-in-inventory, inventory-conversion-period
// Units:   averageInventoryCents and costOfGoodsSoldCents in cents; daysInPeriod in days; the result in days
// Example: Clear Lake Sporting Goods $40,000 ÷ $60,000 × 365 = 243.33, printed 243 days (Principles of Finance §6.2)
func DaysInInventory(averageInventoryCents, costOfGoodsSoldCents int64, daysInPeriod int) (float64, error) {
	if err := checkNonNegativeCents("averageInventoryCents", averageInventoryCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("costOfGoodsSoldCents", costOfGoodsSoldCents); err != nil {
		return 0, err
	}
	if daysInPeriod < 1 {
		return 0, fmt.Errorf("%w: daysInPeriod = %d, want ≥ 1", ErrPeriods, daysInPeriod)
	}
	if costOfGoodsSoldCents == 0 {
		return 0, fmt.Errorf("%w: costOfGoodsSoldCents = 0", ErrZeroDenominator)
	}
	return float64(daysInPeriod) * float64(averageInventoryCents) / float64(costOfGoodsSoldCents), nil
}
