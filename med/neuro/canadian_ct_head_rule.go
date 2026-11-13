package neuro

// CanadianCTHeadFindings are the seven yes/no factors of the Canadian CT Head Rule (Stiell 2001 abstract).
type CanadianCTHeadFindings struct {
	GCSBelow15At2Hours        bool // high-risk: "failure to reach GCS of 15 within 2 h"
	SuspectedOpenSkullFx      bool // high-risk: "suspected open skull fracture"
	BasalSkullFractureSign    bool // high-risk: "any sign of basal skull fracture"
	VomitingMoreThan2         bool // high-risk: "vomiting >2 episodes"
	AgeOver65                 bool // high-risk: "age >65 years"
	AmnesiaBeforeImpactOver30 bool // medium-risk: "amnesia before impact >30 min"
	DangerousMechanism        bool // medium-risk: "dangerous mechanism of injury"
}

// CanadianCTHeadResult says which of the rule's two factor groups has at least one factor present.
type CanadianCTHeadResult struct {
	HighRisk   bool // any of the five "high-risk factors" (the source: 100% sensitive for "need for neurological intervention")
	MediumRisk bool // any of the two "additional medium-risk factors" (98.4% sensitive for "clinically important brain injury" with the high-risk factors)
}

// CanadianCTHeadRule applies the Canadian CT Head Rule for minor head injury: whether any of its five
// high-risk factors, and whether any of its two additional medium-risk factors, is present. The rule was
// derived in adults with a Glasgow Coma Scale score of 13–15 after minor head injury; that the person belongs to
// this population is the caller's.
//
// ID:      MED-NEURO-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stiell IG, Wells GA, Vandemheen K, et al. The Canadian CT Head Rule for patients with minor head injury. Lancet. 2001;357(9266):1391-1396. doi:10.1016/s0140-6736(00)04561-x, PMID 11356436, abstract ("five high-risk factors (failure to reach GCS of 15 within 2 h, suspected open skull fracture, any sign of basal skull fracture, vomiting >2 episodes, or age >65 years) and two additional medium-risk factors (amnesia before impact >30 min and dangerous mechanism of injury)")
// Tags:    neuro, head injury, ct, canadian ct head rule, decision rule
// Units:   yes/no findings; result booleans
// Range:   adults, GCS 13–15 after minor head injury (the source's population; the caller's)
// Example: age over 65 alone → HighRisk; dangerous mechanism alone → MediumRisk only; none → neither (HAND-DERIVED from the abstract)
func CanadianCTHeadRule(f CanadianCTHeadFindings) CanadianCTHeadResult {
	return CanadianCTHeadResult{
		HighRisk: f.GCSBelow15At2Hours || f.SuspectedOpenSkullFx || f.BasalSkullFractureSign ||
			f.VomitingMoreThan2 || f.AgeOver65,
		MediumRisk: f.AmnesiaBeforeImpactOver30 || f.DangerousMechanism,
	}
}
