package respiratory

import (
	"fmt"
	"math"
)

// cigarettesPerPack is the pack size in the USPSTF definition: "20 cigarettes—1 pack" (USPSTF 2021, Clinician Summary).
const cigarettesPerPack = 20

// PackYears returns cumulative smoking exposure in pack-years: packs smoked per day on average × years smoked.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative input; ErrDomain when the product overflows.
//
// ID:      MED-RESP-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  US Preventive Services Task Force. Lung Cancer: Screening. Final Recommendation Statement, March 9, 2021, Clinician Summary, "How to implement this recommendation?": "One pack-year is the equivalent of smoking an average of 20 cigarettes—1 pack—per day for a year." https://www.uspreventiveservicestaskforce.org/uspstf/recommendation/lung-cancer-screening (read 2026-11-12)
// Source:  US Preventive Services Task Force. Screening for Lung Cancer: US Preventive Services Task Force Recommendation Statement. JAMA. 2021;325(10):962-970. doi:10.1001/jama.2021.1117. PMID 33687470. Recommendation statement, eligibility: "a 20 pack-year smoking history"
// Tags:    respiratory, smoking, pack-years, tobacco, exposure
// Units:   packsPerDay packs per day (average); yearsSmoked years; result pack-years
// Range:   packsPerDay ≥ 0; yearsSmoked ≥ 0 (the definition states no limits)
// Example: 1 pack/day for 20 years → 20 pack-years (the definition: 1 pack per day for a year is 1 pack-year)
func PackYears(packsPerDay, yearsSmoked float64) (float64, error) {
	if err := checkClosedRange("packsPerDay", packsPerDay, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	if err := checkClosedRange("yearsSmoked", yearsSmoked, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	py := packsPerDay * yearsSmoked
	if math.IsInf(py, 0) {
		return 0, fmt.Errorf("%w: packsPerDay × yearsSmoked overflows", ErrDomain)
	}
	return py, nil
}

// PackYearsFromCigarettes returns pack-years from cigarettes per day: (cigarettesPerDay / 20) × yearsSmoked,
// 20 cigarettes being one pack in the USPSTF definition.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative input; ErrDomain when the product overflows.
//
// ID:      MED-RESP-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  US Preventive Services Task Force. Lung Cancer: Screening. Final Recommendation Statement, March 9, 2021, Clinician Summary: "an average of 20 cigarettes—1 pack—per day for a year" https://www.uspreventiveservicestaskforce.org/uspstf/recommendation/lung-cancer-screening (read 2026-11-12)
// Tags:    respiratory, smoking, pack-years, tobacco, exposure, cigarettes
// Units:   cigarettesPerDay cigarettes per day (average); yearsSmoked years; result pack-years
// Range:   cigarettesPerDay ≥ 0; yearsSmoked ≥ 0
// Example: 10 cigarettes/day for 30 years → 0.5 × 30 = 15 pack-years (hand-derived from the definition)
func PackYearsFromCigarettes(cigarettesPerDay, yearsSmoked float64) (float64, error) {
	if err := checkClosedRange("cigarettesPerDay", cigarettesPerDay, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	return PackYears(cigarettesPerDay/cigarettesPerPack, yearsSmoked)
}
