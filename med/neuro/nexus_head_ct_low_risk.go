package neuro

// NEXUSHeadCTFindings are the seven findings besides age of the NEXUS Head CT instrument (Mower 2005 abstract).
type NEXUSHeadCTFindings struct {
	SkullFractureEvidence bool // "evidence of skull fracture"
	ScalpHematoma         bool // "scalp hematoma"
	NeurologicDeficit     bool // "neurologic deficit"
	AbnormalAlertness     bool // "abnormal alertness"
	AbnormalBehavior      bool // "abnormal behavior"
	Coagulopathy          bool // "coagulopathy"
	PersistentVomiting    bool // "persistent vomiting"
}

// NEXUSHeadCTLowRisk reports whether a person meets the NEXUS Head CT instrument's low-risk criteria: age under
// 65 years and none of the seven findings (skull fracture, scalp hematoma, neurologic deficit, abnormal
// alertness, abnormal behavior, coagulopathy, persistent vomiting).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-NEURO-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Mower WR, Hoffman JR, Herbert M, Wolfson AB, Pollack CV Jr, Zucker MI; NEXUS II Investigators. Developing a decision instrument to guide computed tomographic imaging of blunt head injury patients. J Trauma. 2005;59(4):954-959. doi:10.1097/01.ta.0000187813.79047.42, PMID 16374287, abstract ("patients under age 65 who had no evidence of skull fracture, scalp hematoma, neurologic deficit, abnormal alertness, abnormal behavior, coagulopathy, or persistent vomiting")
// Tags:    neuro, head injury, ct, nexus, decision rule
// Units:   ageYears years; findings yes/no; result true when every low-risk criterion is met
// Range:   0 ≤ ageYears ≤ 130 (physical limits; blunt head injury, the source's population, is the caller's)
// Example: age 64.9, no finding → true; age 65 → false; age 30 with scalp hematoma → false (HAND-DERIVED from the abstract)
func NEXUSHeadCTLowRisk(ageYears float64, f NEXUSHeadCTFindings) (bool, error) {
	if err := checkRange("ageYears", ageYears, 0, 130); err != nil {
		return false, err
	}
	found := f.SkullFractureEvidence || f.ScalpHematoma || f.NeurologicDeficit || f.AbnormalAlertness ||
		f.AbnormalBehavior || f.Coagulopathy || f.PersistentVomiting
	return ageYears < 65 && !found, nil
}
