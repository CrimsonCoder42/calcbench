package business

import (
	"fmt"
	"math"
)

// SalePriceAfterDiscount returns the price after a percentage discount: original × (1 − discountRate),
// rounded half-even to the cent (0.25 = 25 % off). The factor 1 − discountRate is formed in float64; the
// product is then rounded exactly. The discount itself is original − the result.
// Errors: ErrOutOfRange when originalPriceCents is outside ±2⁵³; ErrNegative when originalPriceCents < 0;
// ErrNotFinite when discountRate is NaN or ±Inf; ErrRate when discountRate is outside 0…1.
//
// ID:      FIN-BIZ-054
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  D. Kirk, Contemporary Mathematics, OpenStax (Rice University), 2023, §6.2 "Discounts, Markups, and Sales Tax" (sale price = original price × (1 − percent discount)), Example 6.10, https://openstax.org/books/contemporary-mathematics/pages/6-2-discounts-markups-and-sales-tax (read 2026-11-13)
// Tags:    discount, pricing, sale-price, percent-off
// Units:   originalPriceCents and the result in cents; discountRate decimal
// Example: $75.80 at 25 % off → $56.85 (Contemporary Mathematics §6.2, Example 6.10)
func SalePriceAfterDiscount(originalPriceCents int64, discountRate float64) (int64, error) {
	if err := checkNonNegativeCents("originalPriceCents", originalPriceCents); err != nil {
		return 0, err
	}
	if math.IsNaN(discountRate) || math.IsInf(discountRate, 0) {
		return 0, fmt.Errorf("%w: discountRate = %v", ErrNotFinite, discountRate)
	}
	if discountRate < 0 || discountRate > 1 {
		return 0, fmt.Errorf("%w: discountRate = %v, want 0…1", ErrRate, discountRate)
	}
	return scaleCents(originalPriceCents, 1-discountRate)
}
