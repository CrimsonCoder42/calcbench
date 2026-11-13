package cardio

import "fmt"

// LACEIndex returns the LACE index for death or unplanned readmission within 30 days of discharge, 0–19:
// L, length of stay in days (< 1: 0; 1: 1; 2: 2; 3: 3; 4–6: 4; 7–13: 5; ≥ 14: 7); A, acute (emergent)
// admission: 3; C, Charlson comorbidity index score (0: 0; 1: 1; 2: 2; 3: 3; ≥ 4: 5); E, emergency department
// visits in the previous 6 months (0–3: that number; ≥ 4: 4).
// Errors: ErrOutOfRange when a count is negative.
//
// ID:      MED-CARDIO-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  van Walraven C, Dhalla IA, Bell C, et al. Derivation and validation of an index to predict early death or unplanned readmission after discharge from hospital to the community. CMAJ. 2010;182(6):551-557. doi:10.1503/cmaj.091117. PMID 20194559. PMC2845681, Table 3 (points for L, A, C, E)
// Tags:    readmission, hospital-discharge, lace, charlson, risk-score
// Units:   lengthOfStayDays whole days; charlsonScore points (computed by the caller with the paper's Charlson weights); edVisits6Months count; result points 0–19
// Range:   all counts ≥ 0 (4,812 adults discharged from 11 Ontario hospitals)
// Example: stay 5 days, emergent admission, Charlson 2, one ED visit → 4 + 3 + 2 + 1 = 10 (HAND-DERIVED from Table 3; Table 4 gives 12.2% expected risk at 10)
func LACEIndex(lengthOfStayDays int, acuteAdmission bool, charlsonScore, edVisits6Months int) (int, error) {
	if lengthOfStayDays < 0 || charlsonScore < 0 || edVisits6Months < 0 {
		return 0, fmt.Errorf("%w: lengthOfStayDays = %d, charlsonScore = %d, edVisits6Months = %d, want all ≥ 0",
			ErrOutOfRange, lengthOfStayDays, charlsonScore, edVisits6Months)
	}
	var l int
	switch {
	case lengthOfStayDays <= 3:
		l = lengthOfStayDays
	case lengthOfStayDays <= 6:
		l = 4
	case lengthOfStayDays <= 13:
		l = 5
	default:
		l = 7
	}
	a := 0
	if acuteAdmission {
		a = 3
	}
	c := min(charlsonScore, 3)
	if charlsonScore >= 4 {
		c = 5
	}
	e := min(edVisits6Months, 4)
	return l + a + c + e, nil
}
