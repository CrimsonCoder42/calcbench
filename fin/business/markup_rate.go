package business

import "fmt"

// MarkupRate returns the markup on cost, (selling price − cost) ÷ cost, as a decimal (0.15 = 15 %). It is
// not rounded. A selling price below cost gives a negative rate.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative;
// ErrZeroDenominator when costCents is 0.
//
// ID:      FIN-BIZ-052
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  D. Kirk, Contemporary Mathematics, OpenStax (Rice University), 2023, §6.2 "Discounts, Markups, and Sales Tax" (markup = percent markup × cost; retail price = cost + markup), Example 6.16, https://openstax.org/books/contemporary-mathematics/pages/6-2-discounts-markups-and-sales-tax (read 2026-11-13)
// Tags:    markup, pricing, retail-price, cost, margin
// Units:   costCents and sellingPriceCents in cents; result decimal
// Example: cost $62.00 sold at $71.30 → 0.15 (Contemporary Mathematics §6.2, Example 6.16)
func MarkupRate(costCents, sellingPriceCents int64) (float64, error) {
	if err := checkNonNegativeCents("costCents", costCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("sellingPriceCents", sellingPriceCents); err != nil {
		return 0, err
	}
	if costCents == 0 {
		return 0, fmt.Errorf("%w: costCents = 0", ErrZeroDenominator)
	}
	return float64(sellingPriceCents-costCents) / float64(costCents), nil
}
