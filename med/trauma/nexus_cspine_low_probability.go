package trauma

// NEXUSCSpineFindings are the five NEXUS findings, each true when the finding is PRESENT.
type NEXUSCSpineFindings struct {
	MidlineCervicalTenderness bool // the low-probability criterion is "no midline cervical tenderness"
	FocalNeurologicDeficit    bool // "no focal neurologic deficit"
	AbnormalAlertness         bool // "normal alertness"
	Intoxication              bool // "no intoxication"
	PainfulDistractingInjury  bool // "no painful, distracting injury"
}

// NEXUSCSpineLowProbability reports whether the NEXUS decision instrument classifies a blunt-trauma subject as
// "having a low probability of injury": true only when all five criteria are met, i.e. none of the five
// findings is present.
//
// ID:      MED-TRAUMA-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hoffman JR, Mower WR, Wolfson AB, Todd KH, Zucker MI. Validity of a set of clinical criteria to rule out injury to the cervical spine in patients with blunt trauma. N Engl J Med. 2000;343(2):94-99. doi:10.1056/NEJM200007133430203. PMID 10891516. Abstract: "The decision instrument required patients to meet five criteria in order to be classified as having a low probability of injury: no midline cervical tenderness, no focal neurologic deficit, normal alertness, no intoxication, and no painful, distracting injury."
// Tags:    trauma, cervical spine, nexus, blunt trauma, decision rule
// Units:   result boolean
// Range:   any combination of findings (blunt trauma, the instrument's population)
// Example: no finding → true; intoxication only → false (from the abstract's five criteria)
func NEXUSCSpineLowProbability(f NEXUSCSpineFindings) bool {
	return !(f.MidlineCervicalTenderness || f.FocalNeurologicDeficit || f.AbnormalAlertness ||
		f.Intoxication || f.PainfulDistractingInjury)
}
