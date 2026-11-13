package economics

import "fmt"

// GDPExpenditure returns gross domestic product by the expenditure approach,
// consumption + investment + government + (exports − imports), in cents. The sum is exact integer
// arithmetic. Consumption, government, exports and imports must be non-negative; investment may be
// negative (inventories can fall). The result may be negative.
// Errors: ErrOutOfRange when an input is outside ±2⁵³ cents; ErrNegative when consumption, government,
// exports or imports is negative; ErrOverflow when the result is outside ±2⁵³ cents.
//
// ID:      FIN-ECON-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Greenlaw, Shapiro and MacDonald, Principles of Macroeconomics, 3rd ed., OpenStax (Rice University), 2022, §6.1 "Measuring the Size of the Economy: Gross Domestic Product" ("GDP = C + I + G + (X – M)"), Table 6.1, https://openstax.org/books/principles-macroeconomics-3e/pages/6-1-measuring-the-size-of-the-economy-gross-domestic-product (read 2026-11-13)
// Tags:    economics, gdp, expenditure-approach, national-accounts, macroeconomics
// Units:   every input and the result in cents (of the same currency and period)
// Example: Table 6.1 (U.S. 2020, trillions): 14.0 + 3.6 + 3.9 + (2.1 − 2.7) = 20.9
func GDPExpenditure(consumptionCents, investmentCents, governmentCents, exportsCents, importsCents int64) (int64, error) {
	for _, c := range []struct {
		name     string
		v        int64
		signedOK bool
	}{
		{"consumptionCents", consumptionCents, false},
		{"investmentCents", investmentCents, true},
		{"governmentCents", governmentCents, false},
		{"exportsCents", exportsCents, false},
		{"importsCents", importsCents, false},
	} {
		if c.v > maxExactCents || c.v < -maxExactCents {
			return 0, fmt.Errorf("%w: %s = %d is outside ±2^53 cents", ErrOutOfRange, c.name, c.v)
		}
		if c.v < 0 && !c.signedOK {
			return 0, fmt.Errorf("%w: %s = %d", ErrNegative, c.name, c.v)
		}
	}
	// Each |term| ≤ 2⁵³, so the int64 sum of five terms cannot overflow.
	gdp := consumptionCents + investmentCents + governmentCents + exportsCents - importsCents
	if gdp > maxExactCents || gdp < -maxExactCents {
		return 0, fmt.Errorf("%w: GDP %d cents is outside ±2^53", ErrOverflow, gdp)
	}
	return gdp, nil
}
