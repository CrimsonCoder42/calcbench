package neuro

// CATCHFindings are the seven yes/no factors of the CATCH rule (Osmond 2010 abstract).
type CATCHFindings struct {
	GCSBelow15At2Hours      bool // high-risk: "failure to reach score of 15 on the Glasgow coma scale within two hours"
	SuspectedOpenSkullFx    bool // high-risk: "suspicion of open skull fracture"
	WorseningHeadache       bool // high-risk: "worsening headache"
	Irritability            bool // high-risk: "irritability"
	LargeBoggyScalpHematoma bool // medium-risk: "large, boggy hematoma of the scalp"
	BasalSkullFractureSign  bool // medium-risk: "signs of basal skull fracture"
	DangerousMechanism      bool // medium-risk: "dangerous mechanism of injury"
}

// CATCHResult says which of the rule's two factor groups has at least one factor present.
type CATCHResult struct {
	HighRisk   bool // any of the four "high-risk factors" (the source: 100.0% sensitive for "the need for neurologic intervention")
	MediumRisk bool // any of the three "additional medium-risk factors" (98.1% sensitivity "for the prediction of brain injury")
}

// CATCHRule applies the CATCH rule (Canadian Assessment of Tomography for Childhood Head injury): whether any of
// its four high-risk factors, and whether any of its three additional medium-risk factors, is present. The rule
// was derived in children with blunt head trauma and a Glasgow Coma Scale score of 13–15; that the child belongs
// to this population is the caller's.
//
// ID:      MED-NEURO-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Osmond MH, Klassen TP, Wells GA, et al. CATCH: a clinical decision rule for the use of computed tomography in children with minor head injury. CMAJ. 2010;182(4):341-348. doi:10.1503/cmaj.091421, PMID 20142371, abstract ("four high-risk factors (failure to reach score of 15 on the Glasgow coma scale within two hours, suspicion of open skull fracture, worsening headache and irritability) and three additional medium-risk factors (large, boggy hematoma of the scalp; signs of basal skull fracture; dangerous mechanism of injury)")
// Tags:    neuro, head injury, pediatric, ct, catch, decision rule
// Units:   yes/no findings; result booleans
// Range:   children with blunt head trauma, GCS 13–15 (the source's population; the caller's)
// Example: irritability alone → HighRisk; boggy scalp hematoma alone → MediumRisk only (HAND-DERIVED from the abstract)
func CATCHRule(f CATCHFindings) CATCHResult {
	return CATCHResult{
		HighRisk:   f.GCSBelow15At2Hours || f.SuspectedOpenSkullFx || f.WorseningHeadache || f.Irritability,
		MediumRisk: f.LargeBoggyScalpHematoma || f.BasalSkullFractureSign || f.DangerousMechanism,
	}
}
