package neuro

// CanadianTIAFindings are the thirteen items of the Canadian TIA Score (Perry 2021, Table 1).
type CanadianTIAFindings struct {
	FirstTIAInLifetime          bool // +2 "First transient ischaemic attack (in lifetime)"
	SymptomsTenMinutesOrMore    bool // +2 "Symptoms ≥10 minutes"
	HistoryOfCarotidStenosis    bool // +2 "Past history of carotid stenosis"
	OnAntiplateletTherapy       bool // +3 "Already on antiplatelet therapy"
	HistoryOfGaitDisturbance    bool // +1 "History of gait disturbance"
	HistoryOfUnilateralWeakness bool // +1 "History of unilateral weakness"
	HistoryOfVertigo            bool // −3 "History of vertigo"
	TriageDiastolicAtLeast110   bool // +3 "Initial triage diastolic blood pressure ≥110 mm Hg"
	DysarthriaOrAphasia         bool // +1 "Dysarthria or aphasia (history or examination)"
	AtrialFibrillationOnECG     bool // +2 "Atrial fibrillation on electrocardiogram"
	InfarctionOnCT              bool // +1 "Infarction (new or old) on computed tomography"
	PlateletsAtLeast400         bool // +2 "Platelet count ≥400×10⁹/L"
	GlucoseAtLeast15MmolL       bool // +3 "Glucose ≥15 mmol/L"
}

// CanadianTIAScore returns the Canadian TIA Score, −3 to 23: the sum of the points of the thirteen items
// present (Perry 2021, Table 1). The thresholds inside the items (≥10 minutes, diastolic ≥110 mm Hg, platelets
// ≥400×10⁹/L, glucose ≥15 mmol/L) are applied by the caller. The risk groups are not returned.
//
// ID:      MED-NEURO-030
// Time:    O(1)   13 items, fixed
// Space:   O(1), no allocation
// Source:  Perry JJ, Sivilotti MLA, Émond M, et al. Prospective validation of Canadian TIA Score and comparison with ABCD2 and ABCD2i for subsequent stroke risk after transient ischaemic attack: multicentre prospective cohort study. BMJ. 2021;372:n49. doi:10.1136/bmj.n49, PMID 33541890, PMC7859838, Table 1 ("Canadian TIA Score"; "The assigned score, which ranges from −3 to 23")
// Tags:    neuro, tia, stroke, canadian tia score, score
// Units:   yes/no items; result points −3 to 23
// Range:   all inputs boolean (adults with TIA or minor stroke in the emergency department, the source's population; the caller's)
// Example: every item but vertigo → 23; vertigo alone → −3; first TIA and symptoms ≥10 min → 4
func CanadianTIAScore(f CanadianTIAFindings) int {
	items := [...]struct {
		yes bool
		pts int
	}{
		{f.FirstTIAInLifetime, 2}, {f.SymptomsTenMinutesOrMore, 2}, {f.HistoryOfCarotidStenosis, 2},
		{f.OnAntiplateletTherapy, 3}, {f.HistoryOfGaitDisturbance, 1}, {f.HistoryOfUnilateralWeakness, 1},
		{f.HistoryOfVertigo, -3}, {f.TriageDiastolicAtLeast110, 3}, {f.DysarthriaOrAphasia, 1},
		{f.AtrialFibrillationOnECG, 2}, {f.InfarctionOnCT, 1}, {f.PlateletsAtLeast400, 2},
		{f.GlucoseAtLeast15MmolL, 3},
	}
	s := 0
	for _, it := range items {
		if it.yes {
			s += it.pts
		}
	}
	return s
}
