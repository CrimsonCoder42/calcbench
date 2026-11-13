package anthropometry

import (
	"fmt"
	"math"
)

// IdealBodyMassTraubKichen returns the ideal body mass of a child, in kg, by the Traub–Kichen formula:
// 2.396 × e^(0.01863 × heightCm), fitted to the 50th-percentile weight for height of NCHS data, ages 1–17.
// ageYears is used only to enforce the source's age range; the formula itself takes height alone.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond
// float64 (heightCm above about 38,100).
//
// ID:      MED-ANTHRO-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Traub SL, Kichen L. Estimating ideal body mass in children. Am J Hosp Pharm. 1983;40(1):107-110. doi:10.1093/ajhp/40.1.107. PMID 6823980. Abstract: "IBM(kg) = 2.396e0.01863(ht)", ages 1–17, the 50th-percentile weight for height from NCHS data
// Tags:    anthropometry, ideal-body-weight, pediatric, traub-kichen
// Units:   heightCm cm; ageYears years (decimal); result kg
// Range:   1 ≤ ageYears ≤ 17 (the source's ages); heightCm > 0 (no height limits stated)
// Example: 100 cm → 2.396 × e^1.863 (HAND-DERIVED; the test writes the formula out)
func IdealBodyMassTraubKichen(heightCm, ageYears float64) (float64, error) {
	if err := checkFinite("heightCm", heightCm); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkPositive("heightCm", heightCm); err != nil {
		return 0, err
	}
	if ageYears < 1 || ageYears > 17 {
		return 0, fmt.Errorf("%w: ageYears = %v, want 1 ≤ ageYears ≤ 17", ErrOutOfRange, ageYears)
	}
	r := 2.396 * math.Exp(0.01863*heightCm)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}
