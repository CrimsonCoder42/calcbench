package trauma

// PECARNIAIFindings are the findings of the PECARN intra-abdominal injury rule besides the GCS, each true when
// PRESENT (the rule's items are worded as absences).
type PECARNIAIFindings struct {
	AbdominalWallTraumaOrSeatBeltSign bool // "no evidence of abdominal wall trauma or seat belt sign"
	AbdominalTenderness               bool // "no abdominal tenderness"
	ThoracicWallTrauma                bool // "no evidence of thoracic wall trauma"
	AbdominalPain                     bool // "no complaints of abdominal pain"
	DecreasedBreathSounds             bool // "no decreased breath sounds"
	Vomiting                          bool // "no vomiting"
}

// PECARNIntraAbdominalInjuryVeryLowRisk reports whether a child with blunt torso trauma meets every item of the
// PECARN prediction rule: no evidence of abdominal wall trauma or seat belt sign, Glasgow Coma Scale score
// greater than 13, no abdominal tenderness, no evidence of thoracic wall trauma, no complaints of abdominal pain,
// no decreased breath sounds, and no vomiting. True is the rule's "very low risk" result (the title's words).
// Errors: ErrOutOfRange when gcs is outside 3–15.
//
// ID:      MED-TRAUMA-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Holmes JF, Lillis K, Monroe D, et al. Identifying children at very low risk of clinically important blunt abdominal injuries. Ann Emerg Med. 2013;62(2):107-116.e2. doi:10.1016/j.annemergmed.2012.11.009. PMID 23375510. Abstract: "The prediction rule consisted of (in descending order of importance) no evidence of abdominal wall trauma or seat belt sign, Glasgow Coma Scale score greater than 13, no abdominal tenderness, no evidence of thoracic wall trauma, no complaints of abdominal pain, no decreased breath sounds, and no vomiting."
// Tags:    trauma, pediatrics, abdomen, pecarn, blunt trauma, decision rule
// Units:   gcs points 3–15; result boolean
// Range:   3 ≤ gcs ≤ 15 (children with blunt torso trauma, the study's population, which the caller establishes)
// Example: GCS 14, no finding → true; GCS 13, no finding → false (from "greater than 13")
func PECARNIntraAbdominalInjuryVeryLowRisk(gcs int, f PECARNIAIFindings) (bool, error) {
	if err := checkGCS(gcs); err != nil {
		return false, err
	}
	return gcs > 13 && !(f.AbdominalWallTraumaOrSeatBeltSign || f.AbdominalTenderness || f.ThoracicWallTrauma ||
		f.AbdominalPain || f.DecreasedBreathSounds || f.Vomiting), nil
}
