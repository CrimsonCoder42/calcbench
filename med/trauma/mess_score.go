package trauma

import "fmt"

// MESSScore returns the Mangled Extremity Severity Score: skeletal/soft-tissue injury (low energy 1, medium 2,
// high 3, very high energy 4) + limb ischemia (pulse reduced or absent but perfusion normal 1; pulseless,
// paresthesias, diminished capillary refill 2; cool, paralyzed, insensate, numb 3; doubled for ischemia > 6 hours)
// + shock (systolic always > 90 mmHg 0; transiently hypotensive 1; persistent hypotension 2) + age (< 30 years 0;
// 30–50 1; > 50 2). The ischemia item has no 0 row in the sources; 1–3 are accepted.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an item outside its printed points or a negative age.
//
// ID:      MED-TRAUMA-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johansen K, Daines M, Howey T, Helfet D, Hansen ST Jr. Objective criteria accurately predict amputation following lower extremity trauma. J Trauma. 1990;30(5):568-572. doi:10.1097/00005373-199005000-00007. PMID 2342140 (primary; not opened). Abstract: "a MESS value greater than or equal to 7 predicted amputation with 100% accuracy"
// Source:  Topal AE, Eren MN, Celik Y. Lower extremity arterial injuries over a six-year period: outcomes, risk factors, and management. Vasc Health Risk Manag. 2010;6:1103-1110. doi:10.2147/VHRM.S15316. PMID 21191430. PMC3004513, MESS table (cites Johansen 1990)
// Source:  Okereke I, Abdelfatah E. Limb Salvage Versus Amputation for the Mangled Extremity: Factors Affecting Decision-Making and Outcomes. Cureus. 2022;14(8):e28153. doi:10.7759/cureus.28153. PMID 36148195. PMC9482680, MESS table ("Score doubled for ischemia > 6 hours"; cites Johansen 1990)
// Verified: secondary (Topal 2010; Okereke 2022)
// Tags:    trauma, extremity, amputation, mess, score
// Units:   skeletal, ischemia, shock printed points; ischemiaHours h; ageYears years; result points
// Range:   skeletal 1–4; ischemia 1–3; shock 0–2; ischemiaHours ≥ 0; ageYears ≥ 0
// Example: high energy 3, pulseless 2 for 8 h (→ 4), transient hypotension 1, age 40 (1) → 9 (HAND-DERIVED)
func MESSScore(skeletal, ischemia int, ischemiaHours float64, shock int, ageYears float64) (int, error) {
	for _, e := range [...]error{
		finite("ischemiaHours", ischemiaHours), finite("ageYears", ageYears),
		atLeastZero("ischemiaHours", ischemiaHours), atLeastZero("ageYears", ageYears),
	} {
		if e != nil {
			return 0, e
		}
	}
	switch {
	case skeletal < 1 || skeletal > 4:
		return 0, fmt.Errorf("%w: skeletal = %d, want 1–4", ErrOutOfRange, skeletal)
	case ischemia < 1 || ischemia > 3:
		return 0, fmt.Errorf("%w: ischemia = %d, want 1–3", ErrOutOfRange, ischemia)
	case shock < 0 || shock > 2:
		return 0, fmt.Errorf("%w: shock = %d, want 0–2", ErrOutOfRange, shock)
	}
	isch := ischemia
	if ischemiaHours > 6 {
		isch *= 2
	}
	age := 0
	switch {
	case ageYears > 50:
		age = 2
	case ageYears >= 30:
		age = 1
	}
	return skeletal + isch + shock + age, nil
}
