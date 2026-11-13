package trauma

// CanadianCSpineFindings are the yes/no items of the Canadian C-Spine Rule besides age (Stiell 2001 abstract).
type CanadianCSpineFindings struct {
	DangerousMechanism      bool // high-risk factor: "dangerous mechanism"
	ParesthesiasExtremities bool // high-risk factor: "paresthesias in extremities"
	SimpleRearEndCollision  bool // low-risk factor: "simple rear-end motor vehicle collision"
	SittingPositionInED     bool // low-risk factor: "sitting position in ED"
	AmbulatorySinceInjury   bool // low-risk factor: "ambulatory at any time since injury"
	DelayedOnsetNeckPain    bool // low-risk factor: "delayed onset of neck pain"
	NoMidlineTenderness     bool // low-risk factor: "absence of midline C-spine tenderness"
	AbleToRotateNeck45      bool // "able to actively rotate neck 45 degrees to the left and right"
}

// CanadianCSpineResult holds the answers to the rule's three questions and the rule's outcome.
type CanadianCSpineResult struct {
	HighRiskFactor bool // question 1: a high-risk factor that "mandates radiography" is present
	LowRiskFactor  bool // question 2: a low-risk factor "that allows safe assessment of range of motion" is present
	Radiography    bool // the rule's outcome: true when question 1 is yes, or question 2 is no, or the neck cannot rotate 45° both ways
}

// CanadianCSpineRule applies the three questions of the Canadian C-Spine Rule: (1) any high-risk factor (age
// 65 years or older, dangerous mechanism, paresthesias in extremities) "that mandates radiography"; (2) any
// low-risk factor "that allows safe assessment of range of motion"; (3) able to actively rotate the neck 45° to
// the left and right. Radiography is the outcome when (1) is yes, otherwise when (2) is no, otherwise when (3) is
// no. AbleToRotateNeck45 is read only when range of motion is assessed (no high-risk and some low-risk factor).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-TRAUMA-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stiell IG, Wells GA, Vandemheen KL, et al. The Canadian C-spine rule for radiography in alert and stable trauma patients. JAMA. 2001;286(15):1841-1848. doi:10.1001/jama.286.15.1841. PMID 11597285. Abstract: "3 main questions: (1) is there any high-risk factor present that mandates radiography (ie, age >/=65 years, dangerous mechanism, or paresthesias in extremities)? (2) is there any low-risk factor present that allows safe assessment of range of motion (ie, simple rear-end motor vehicle collision, sitting position in ED, ambulatory at any time since injury, delayed onset of neck pain, or absence of midline C-spine tenderness)? and (3) is the patient able to actively rotate neck 45 degrees to the left and right?"
// Tags:    trauma, cervical spine, canadian c-spine rule, decision rule, radiography
// Units:   ageYears years; result booleans
// Range:   0 ≤ ageYears ≤ 130 (physical limits; alert and stable adults after trauma, the title's population, which the caller establishes)
// Example: age 70 → Radiography (question 1); age 30, sitting in ED, rotates 45° → no radiography; age 30, no low-risk factor → Radiography (HAND-DERIVED from the abstract's three questions)
func CanadianCSpineRule(ageYears float64, f CanadianCSpineFindings) (CanadianCSpineResult, error) {
	if err := checkAge(ageYears); err != nil {
		return CanadianCSpineResult{}, err
	}
	var r CanadianCSpineResult
	r.HighRiskFactor = ageYears >= 65 || f.DangerousMechanism || f.ParesthesiasExtremities
	r.LowRiskFactor = f.SimpleRearEndCollision || f.SittingPositionInED || f.AmbulatorySinceInjury ||
		f.DelayedOnsetNeckPain || f.NoMidlineTenderness
	r.Radiography = r.HighRiskFactor || !r.LowRiskFactor || !f.AbleToRotateNeck45
	return r, nil
}
