package critical

// RochesterInputs are the Rochester criteria items (Jaskiewicz 1994 abstract).
type RochesterInputs struct {
	AppearsWell       bool    // "appear well"
	PreviouslyHealthy bool    // "were previously healthy"
	NoFocalInfection  bool    // "have no focal infection"
	WBCE9PerL         float64 // white blood cell count, 10⁹ cells/L: low-risk "5.0-15.0 x 10(9) cells/L"
	BandsE9PerL       float64 // band form count, 10⁹ cells/L: low-risk "≤ 1.5 x 10(9) cells/L"
	UrineWBCPerHPF    float64 // WBC per high power field, spun urine sediment: low-risk "≤ 10"
	Diarrhea          bool    // the stool criterion applies only "if diarrhea"
	StoolWBCPerHPF    float64 // WBC per high power field, stool smear: low-risk "≤ 5" (read only when Diarrhea)
}

// RochesterLowRisk reports whether a febrile infant 60 days of age or younger meets every Rochester low-risk
// criterion: appears well; previously healthy; no focal infection; WBC 5.0–15.0 × 10⁹/L (both ends included);
// band form count ≤ 1.5 × 10⁹/L; ≤ 10 WBC per high power field of spun urine sediment; and, if diarrhea, ≤ 5 WBC
// per high power field of a stool smear.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for ageDays outside 0–60 or a negative count.
//
// ID:      MED-CRIT-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Jaskiewicz JA, McCarthy CA, Richardson AC, et al. Febrile infants at low risk for serious bacterial infection—an appraisal of the Rochester criteria and implications for management. Pediatrics. 1994;94(3):390-396. PMID 8065869. Abstract: "appear well; were previously healthy; have no focal infection; have WBC count 5.0-15.0 x 10(9) cells/L ... band form count < or = 1.5 x 10(9) cells/L ... < or = 10 WBC per high power field on microscopic examination of spun urine sediment, and < or = 5 WBC per high power field on microscopic examination of a stool smear (if diarrhea)"
// Source:  Dagan R, Powell KR, Hall CB, Menegus MA. Identification of infants unlikely to have serious bacterial infection although hospitalized for suspected sepsis. J Pediatr. 1985;107(6):855-860. doi:10.1016/s0022-3476(85)80175-x. PMID 4067741 (original criteria: "between 5000 and 15,000 white blood cells/mm3, had less than 1500 bands/mm3")
// Tags:    critical, pediatrics, febrile infant, rochester criteria, serious bacterial infection, decision rule
// Units:   ageDays days; WBCE9PerL, BandsE9PerL 10⁹ cells/L; UrineWBCPerHPF, StoolWBCPerHPF cells per high power field; result boolean
// Range:   0 ≤ ageDays ≤ 60 ("< or = 60 days of age"); every count ≥ 0
// Example: well, healthy, no focus, WBC 15.0, bands 1.5, urine 10, no diarrhea → true; WBC 4.9 → false (from the abstract's criteria)
func RochesterLowRisk(ageDays float64, in RochesterInputs) (bool, error) {
	for _, e := range [...]error{
		finite("ageDays", ageDays), finite("WBCE9PerL", in.WBCE9PerL), finite("BandsE9PerL", in.BandsE9PerL),
		finite("UrineWBCPerHPF", in.UrineWBCPerHPF), finite("StoolWBCPerHPF", in.StoolWBCPerHPF),
		within("ageDays", ageDays, 0, 60), atLeastZero("WBCE9PerL", in.WBCE9PerL),
		atLeastZero("BandsE9PerL", in.BandsE9PerL), atLeastZero("UrineWBCPerHPF", in.UrineWBCPerHPF),
		atLeastZero("StoolWBCPerHPF", in.StoolWBCPerHPF),
	} {
		if e != nil {
			return false, e
		}
	}
	return in.AppearsWell && in.PreviouslyHealthy && in.NoFocalInfection &&
		in.WBCE9PerL >= 5.0 && in.WBCE9PerL <= 15.0 && in.BandsE9PerL <= 1.5 && in.UrineWBCPerHPF <= 10 &&
		(!in.Diarrhea || in.StoolWBCPerHPF <= 5), nil
}
