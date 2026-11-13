package trauma

// PedSRCFindings are the PedSRC rule's yes/no findings besides AST, each true when the finding is PRESENT.
type PedSRCFindings struct {
	AbnormalAbdominalExam     bool // "abnormal abdominal examination"
	AbnormalChestXRay         bool // "abnormal chest x-ray"
	AbdominalPain             bool // "report of abdominal pain"
	AbnormalPancreaticEnzymes bool // "abnormal pancreatic enzymes"
}

// PedSRCVeryLowRisk returns the number of the PedSRC prediction rule's five variables present (0–5): aspartate
// aminotransferase above 200 U/L, abnormal abdominal examination, abnormal chest x-ray, report of abdominal pain
// and abnormal pancreatic enzymes; and whether none is present, the rule's "very low risk" result (the title's
// words).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when astUL is negative.
//
// ID:      MED-TRAUMA-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Streck CJ, Vogel AM, Zhang J, et al. Identifying children at very low risk for blunt intra-abdominal injury in whom CT of the abdomen can be avoided safely. J Am Coll Surg. 2017;224(4):449-458.e3. doi:10.1016/j.jamcollsurg.2016.12.041. PMID 28130170. Abstract: "The prediction rule consisted of (in descending order of significance): aspartate aminotransferase >200 U/L, abnormal abdominal examination, abnormal chest x-ray, report of abdominal pain, and abnormal pancreatic enzymes."
// Tags:    trauma, pediatrics, abdomen, pedsrc, blunt trauma, decision rule
// Units:   astUL U/L; result count 0–5 and boolean
// Range:   astUL ≥ 0 (physical limit; children with blunt torso trauma, the study's population, which the caller establishes)
// Example: AST 200, none → (0, true); AST 201 → (1, false) (from ">200 U/L")
func PedSRCVeryLowRisk(astUL float64, f PedSRCFindings) (int, bool, error) {
	if err := finite("astUL", astUL); err != nil {
		return 0, false, err
	}
	if err := atLeastZero("astUL", astUL); err != nil {
		return 0, false, err
	}
	n := count(astUL > 200, f.AbnormalAbdominalExam, f.AbnormalChestXRay, f.AbdominalPain, f.AbnormalPancreaticEnzymes)
	return n, n == 0, nil
}
