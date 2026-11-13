package nutrition

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ISSN 2017 position statement point 2: "1.4–2.0 g protein/kg body weight/day (g/kg/d)".
const (
	issnProteinLowGPerKg  = 1.4
	issnProteinHighGPerKg = 2.0
)

// ProteinRangeISSN returns the low and high ends, in g/day, of the overall daily protein intake range of the
// International Society of Sports Nutrition 2017 position stand for exercising individuals:
// 1.4 × body weight and 2.0 × body weight. Not rounded. The stand's narrower ranges (2.3–3.1 g/kg/d during
// hypocaloric periods) are not computed.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      LIF-NUTR-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Jäger R, Kerksick CM, Campbell BI, et al. International Society of Sports Nutrition Position Stand: protein and exercise. J Int Soc Sports Nutr. 2017;14:20. doi:10.1186/s12970-017-0177-8. PMID 28642676. PMC5477153, position statement point 2: "an overall daily protein intake in the range of 1.4–2.0 g protein/kg body weight/day (g/kg/d)"
// Tags:    nutrition, protein, sports-nutrition, exercise, issn
// Units:   weightKg kg body weight; results g/day
// Range:   weightKg > 0 (the stand states no weight or age limits)
// Example: 70 kg → 98–140 g/day (HAND-DERIVED: 1.4 × 70, 2.0 × 70)
func ProteinRangeISSN(weightKg float64) (lowGPerDay, highGPerDay float64, err error) {
	if err := checkFinitePositive("weightKg", weightKg); err != nil {
		return 0, 0, err
	}
	hi := issnProteinHighGPerKg * weightKg
	if !reduce.AllFinite(hi) {
		return 0, 0, fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return issnProteinLowGPerKg * weightKg, hi, nil
}
