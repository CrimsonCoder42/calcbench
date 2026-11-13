package cardio

// BostonSyncopeFindings are the eight risk-factor groups of the Boston Syncope Rule (Grossman 2007, verbatim names).
type BostonSyncopeFindings struct {
	AcuteCoronarySyndrome    bool // "Signs and symptoms of Acute Coronary Syndrome"
	ConductionDisease        bool // "Signs of conduction disease"
	WorrisomeCardiacHistory  bool // "Worrisome cardiac history"
	ValvularHeartDisease     bool // "Valvular heart disease by history or physical examination"
	FamilyHistorySuddenDeath bool // "Family history of sudden death"
	PersistentAbnormalVitals bool // "Persistent abnormal vital signs in the ED"
	VolumeDepletion          bool // "Volume depletion"
	PrimaryCNSEvent          bool // "Primary central nervous system event"
}

// BostonSyncopeRuleCount returns how many of the Boston Syncope Rule's eight risk-factor groups are present,
// 0–8, and whether any is present (the rule is positive when any one risk factor is present).
// Errors: none.
//
// ID:      MED-CARDIO-055
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Grossman SA, Fischer C, Lipsitz LA, et al. Predicting adverse outcomes in syncope. J Emerg Med. 2007;33(3):233-239. doi:10.1016/j.jemermed.2007.04.001. PMID 17976548. PMC2276584, Table 1 (risk factors) and Methods
// Tags:    cardio, syncope, emergency, boston, decision-rule
// Units:   result count 0–8 and any-present flag
// Range:   adults 18 years and older presenting to the ED with syncope; no numeric inputs
// Example: volume depletion only → (1, true)
func BostonSyncopeRuleCount(f BostonSyncopeFindings) (count int, anyPresent bool) {
	for _, p := range [...]bool{f.AcuteCoronarySyndrome, f.ConductionDisease, f.WorrisomeCardiacHistory,
		f.ValvularHeartDisease, f.FamilyHistorySuddenDeath, f.PersistentAbnormalVitals, f.VolumeDepletion,
		f.PrimaryCNSEvent} {
		if p {
			count++
		}
	}
	return count, count > 0
}
