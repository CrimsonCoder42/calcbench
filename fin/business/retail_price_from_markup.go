package business

import (
	"fmt"
	"math"
)

// scaleCents returns cents × factor rounded half-even to the cent. The product's rounding error is
// recovered with FMA, so the rounding is exact for the float64 factor as given (roundHalfEvenExact).
// Errors: ErrOverflow when the result exceeds 2⁵³ cents.
func scaleCents(cents int64, factor float64) (int64, error) {
	c := float64(cents)
	p := c * factor
	if math.Abs(p) > maxExactCents {
		return 0, fmt.Errorf("%w: %v cents is above 2^53", ErrOverflow, p)
	}
	return roundHalfEvenExact(p, math.FMA(c, factor, -p), 0), nil
}

// RetailPriceFromMarkup returns the selling price of an item that costs costCents with a markup on cost of
// markupRate (0.15 = 15 %): cost × (1 + markupRate), rounded half-even to the cent. The factor 1 + markupRate
// is formed in float64; the product is then rounded exactly.
// Errors: ErrOutOfRange when costCents is outside ±2⁵³; ErrNegative when costCents < 0; ErrNotFinite when
// markupRate is NaN or ±Inf; ErrRate when markupRate < 0; ErrOverflow when the price exceeds 2⁵³ cents.
//
// ID:      FIN-BIZ-053
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  D. Kirk, Contemporary Mathematics, OpenStax (Rice University), 2023, §6.2 "Discounts, Markups, and Sales Tax" (retail price = cost × (1 + percent markup)), Example 6.16, https://openstax.org/books/contemporary-mathematics/pages/6-2-discounts-markups-and-sales-tax (read 2026-11-13)
// Tags:    markup, pricing, retail-price, cost
// Units:   costCents and the result in cents; markupRate decimal
// Example: cost $62.00, markup 15 % → $71.30 (Contemporary Mathematics §6.2, Example 6.16)
func RetailPriceFromMarkup(costCents int64, markupRate float64) (int64, error) {
	if err := checkNonNegativeCents("costCents", costCents); err != nil {
		return 0, err
	}
	if math.IsNaN(markupRate) || math.IsInf(markupRate, 0) {
		return 0, fmt.Errorf("%w: markupRate = %v", ErrNotFinite, markupRate)
	}
	if markupRate < 0 {
		return 0, fmt.Errorf("%w: markupRate = %v, want ≥ 0", ErrRate, markupRate)
	}
	return scaleCents(costCents, 1+markupRate)
}
