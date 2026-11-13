package business

import "fmt"

// BreakEvenUnits returns the number of units at which total contribution margin equals fixed costs:
// fixed costs ÷ (price per unit − variable cost per unit). The result is not rounded up to whole units;
// the caller decides. The subtraction is exact; the result has one rounding.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative; ErrDomain
// when pricePerUnitCents ≤ variableCostPerUnitCents (no positive contribution margin, so no break-even).
//
// ID:      FIN-BIZ-042
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 2: Managerial Accounting, OpenStax (Rice University), 2019, §3.2 "Calculate a Break-Even Point in Units and Dollars" (break-even point in units = total fixed costs ÷ contribution margin per unit), https://openstax.org/books/principles-managerial-accounting/pages/3-2-calculate-a-break-even-point-in-units-and-dollars (read 2026-11-12)
// Tags:    cost-volume-profit, break-even, contribution-margin, managerial
// Units:   fixedCostsCents in cents; pricePerUnitCents and variableCostPerUnitCents in cents per unit; the result in units
// Example: Hicks Manufacturing $18,000 ÷ ($100 − $20) = 225 units (Principles of Accounting Vol. 2 §3.2)
func BreakEvenUnits(fixedCostsCents, pricePerUnitCents, variableCostPerUnitCents int64) (float64, error) {
	if err := checkNonNegativeCents("fixedCostsCents", fixedCostsCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("pricePerUnitCents", pricePerUnitCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("variableCostPerUnitCents", variableCostPerUnitCents); err != nil {
		return 0, err
	}
	margin := pricePerUnitCents - variableCostPerUnitCents
	if margin <= 0 {
		return 0, fmt.Errorf("%w: pricePerUnitCents − variableCostPerUnitCents = %d, want > 0", ErrDomain, margin)
	}
	return float64(fixedCostsCents) / float64(margin), nil
}
