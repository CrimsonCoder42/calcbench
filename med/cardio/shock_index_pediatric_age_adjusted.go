package cardio

import "fmt"

// SIPACutoff returns the age-adjusted shock-index cut-off of SIPA for injured children: above 1.22 at ages
// 4–6, above 1.0 at 7–12, above 0.9 at 13–16. A shock index above the returned value is "elevated" in the
// paper's sense; the comparison is left to the caller.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 4 to under 17.
//
// ID:      MED-CARDIO-037
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Acker SN, Ross JT, Partrick DA, Tong S, Bensard DD. Pediatric specific shock index accurately identifies severely injured children. J Pediatr Surg. 2015;50(2):331-334. doi:10.1016/j.jpedsurg.2014.08.009. PMID 25638631. Abstract: "Cutoffs included SI>1.22 (age 4-6), >1.0 (7-12), and >0.9 (13-16)"
// Tags:    cardio, pediatrics, trauma, shock-index, sipa
// Units:   ageYears completed years (decimal allowed: 6.9 is in the 4–6 band); result unitless shock-index cut-off
// Range:   4 ≤ ageYears < 17 (543 injured children; the abstract defines no band below age 4)
// Example: age 5 → 1.22; age 10 → 1.0; age 15 → 0.9 (abstract)
func SIPACutoff(ageYears float64) (float64, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	switch {
	case ageYears < 4 || ageYears >= 17:
		return 0, fmt.Errorf("%w: ageYears = %v, want 4 ≤ ageYears < 17", ErrOutOfRange, ageYears)
	case ageYears < 7:
		return 1.22, nil
	case ageYears < 13:
		return 1.0, nil
	default:
		return 0.9, nil
	}
}
