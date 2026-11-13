package business

import (
	"fmt"
	"math"
)

// EconomicOrderQuantity returns the lot size that minimises the sum of set-up (ordering) cost and carrying
// cost for a steady demand, √(2·D·S ÷ H), in units, not rounded (Harris's own example truncates to whole
// units; the caller decides). D is annualDemandUnits per year, S is orderCostCents per order, H is
// holdingCostPerUnitPerYearCents. Harris writes it as √(240·M·S ÷ C) with M units per month, C the unit
// cost and a 10 % yearly charge for interest and depreciation, which is this formula with D = 12·M and
// H = C ÷ 10. It is computed as √(2·S ÷ H)·√D: S ÷ H is one rounding and the two square roots cannot
// overflow, so the result is finite for every finite D and within a few ulp of the exact value.
// Errors: ErrNotFinite when annualDemandUnits is NaN or ±Inf; ErrNegative when annualDemandUnits,
// orderCostCents or holdingCostPerUnitPerYearCents is negative; ErrOutOfRange when a money input is outside
// ±2⁵³; ErrZeroDenominator when holdingCostPerUnitPerYearCents is 0.
//
// ID:      FIN-BIZ-043
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  F. W. Harris, "What Quantity to Make at Once", ch. IV in Operation and Costs (The Library of Factory Management, vol. V; Shaw Factory Management Series), A. W. Shaw Company, Chicago, 1915, pp. 47–52 (pp. 49–51: symbols M, C, S, X; "the square root of (240 MS divided by C)"), full text HathiTrust uc2.ark:/13960/t02z13n82 read at https://distantreader.org/stacks/trust/uc2/uc2_ark+=13960=t02z13n82.txt (read 2026-11-12)
// Source:  D. Erlenkotter, "Ford Whitman Harris's Economical Lot Size Model", UCLA Anderson School, 2013, eScholarship qt4s8482p2, https://escholarship.org/content/qt4s8482p2/qt4s8482p2_noSplash_5dfa07046f9562ca39dee86651320d0b.pdf (read 2026-11-12; identifies the 1915 chapter as Harris's EOQ model, first printed in Factory, The Magazine of Management 10(2):135–136, 152, 1913, reprinted in Operations Research 38(6):947–950, 1990, doi:10.1287/opre.38.6.947, DOI confirmed on Crossref)
// Source:  D. Erlenkotter, "Ford Whitman Harris's economical lot size model", International Journal of Production Economics 155:12–15, 2014, doi:10.1016/j.ijpe.2013.12.008 (DOI confirmed on Crossref 2026-11-12)
// Tags:    inventory, eoq, economic-order-quantity, lot-size, harris, working-capital
// Units:   annualDemandUnits in units per year; orderCostCents in cents per order; holdingCostPerUnitPerYearCents in cents per unit per year; the result in units per order
// Example: Harris p. 51: M = 1,000 units a month, set-up $2, unit cost $0.10 → "2,190 units"; here D = 12,000, S = 200¢, H = 1¢ (10 % of 10¢) → √4,800,000 = 2,190.89…
func EconomicOrderQuantity(annualDemandUnits float64, orderCostCents, holdingCostPerUnitPerYearCents int64) (float64, error) {
	if math.IsNaN(annualDemandUnits) || math.IsInf(annualDemandUnits, 0) {
		return 0, fmt.Errorf("%w: annualDemandUnits = %v", ErrNotFinite, annualDemandUnits)
	}
	if annualDemandUnits < 0 {
		return 0, fmt.Errorf("%w: annualDemandUnits = %v", ErrNegative, annualDemandUnits)
	}
	if err := checkNonNegativeCents("orderCostCents", orderCostCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("holdingCostPerUnitPerYearCents", holdingCostPerUnitPerYearCents); err != nil {
		return 0, err
	}
	costRatio, err := ratioCents(orderCostCents, "holdingCostPerUnitPerYearCents", holdingCostPerUnitPerYearCents)
	if err != nil {
		return 0, err
	}
	// 2·costRatio ≤ 2⁵⁴ and √D ≤ 2⁵¹²: the product is finite.
	return math.Sqrt(2*costRatio) * math.Sqrt(annualDemandUnits), nil
}
