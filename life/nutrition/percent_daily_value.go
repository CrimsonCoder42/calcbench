package nutrition

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PercentDailyValue returns the percent Daily Value of a nutrient with a Daily Reference Value (DRV), expressed
// to the nearest whole percent: round(amount / DRV × 100). The amount and the DRV are in the same unit (g or
// mg); nothing is converted. A value exactly halfway between two whole percents rounds up (12.5 → 13), as the
// FDA Food Labeling Guide's rounding rules state (D-073; the regulation itself states no tie rule).
// Vitamins and minerals use the coarser increments of 21 CFR 101.9(c)(8)(iii) and are not covered here.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      LIF-NUTR-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  21 CFR 101.9(d)(7)(ii), eCFR current text read 2026-11-12: "The percent shall be calculated by dividing either the amount declared on the label for each nutrient or the actual amount of each nutrient (i.e., before rounding) by the DRV for the nutrient"; "expressed to the nearest whole percent"
// Source:  U.S. Food and Drug Administration. Guidance for Industry: A Food Labeling Guide. January 2013, Appendix H "Rounding the Values According to FDA Rounding Rules", p. 129: "When rounding % DV for nutrients other than vitamins and minerals, when the % DV values fall exactly halfway between two whole numbers or higher (e.g., 2.5 to 2.99), the values round up (e.g., 3 %)."
// Tags:    nutrition, labeling, daily-value, percent-dv, drv, cfr
// Units:   amountSameUnit and dailyValueSameUnit in one unit, g or mg (the DRVs of 101.9(c)(9) are in g or mg); result whole percent
// Range:   amountSameUnit ≥ 0; dailyValueSameUnit > 0
// Example: sodium 1150 mg of the 2,300 mg DRV → 50; total carbohydrate 30 g of 275 g → 10.909… → 11 (HAND-DERIVED; DRVs from 101.9(c)(9), adults and children ≥ 4 years)
func PercentDailyValue(amountSameUnit, dailyValueSameUnit float64) (float64, error) {
	if err := checkFinite("amountSameUnit", amountSameUnit); err != nil {
		return 0, err
	}
	if err := checkFinitePositive("dailyValueSameUnit", dailyValueSameUnit); err != nil {
		return 0, err
	}
	if !(amountSameUnit >= 0) {
		return 0, fmt.Errorf("%w: amountSameUnit = %v, want ≥ 0", ErrOutOfRange, amountSameUnit)
	}
	r := math.Round(amountSameUnit / dailyValueSameUnit * 100) // amount ≥ 0: half away from zero is half up
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return r, nil
}
