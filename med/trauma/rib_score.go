package trauma

// RibScoreFindings are the six radiographic variables of RibScore (Chapman 2016 abstract).
type RibScoreFindings struct {
	SixOrMoreRibFractures           bool // "(1) six or more rib fractures"
	BilateralFractures              bool // "(2) bilateral fractures"
	FlailChest                      bool // "(3) flail chest"
	ThreeOrMoreSeverelyDisplaced    bool // "(4) three or more severely (bicortical) displaced fractures"
	FirstRibFracture                bool // "(5) first rib fracture"
	FractureInAllThreeAnatomicAreas bool // "(6) at least one fracture in all three anatomic areas (anterior, lateral, and posterior)"
}

// RibScore returns RibScore, 0–6: one point for each of its six radiographic variables present.
//
// ID:      MED-TRAUMA-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chapman BC, Herbert B, Rodil M, et al. RibScore: a novel radiographic score based on fracture pattern that predicts pneumonia, respiratory failure, and tracheostomy. J Trauma Acute Care Surg. 2016;80(1):95-101. doi:10.1097/TA.0000000000000867. PMID 26683395. Abstract: "assigning 1 point for each variable": "(1) six or more rib fractures, (2) bilateral fractures, (3) flail chest, (4) three or more severely (bicortical) displaced fractures, (5) first rib fracture, and (6) at least one fracture in all three anatomic areas (anterior, lateral, and posterior)"
// Tags:    trauma, rib fractures, ribscore, chest, score
// Units:   result points 0–6
// Range:   any combination of findings
// Example: bilateral and flail chest → 2; all six → 6 (from the abstract's one point per variable)
func RibScore(f RibScoreFindings) int {
	return count(f.SixOrMoreRibFractures, f.BilateralFractures, f.FlailChest, f.ThreeOrMoreSeverelyDisplaced,
		f.FirstRibFracture, f.FractureInAllThreeAnatomicAreas)
}
