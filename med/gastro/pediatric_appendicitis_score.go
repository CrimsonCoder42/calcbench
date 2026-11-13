package gastro

// PediatricAppendicitisFindings holds the eight variables of Samuel's Pediatric Appendicitis Score, in the
// source's order. Each is judged by the caller (the source's abstract gives no laboratory thresholds).
type PediatricAppendicitisFindings struct {
	CoughPercussionHoppingTendernessRLQ bool // (1) cough/percussion/hopping tenderness in the right lower quadrant: 2 points
	Anorexia                            bool // (2) anorexia: 1 point
	Pyrexia                             bool // (3) pyrexia: 1 point
	NauseaEmesis                        bool // (4) nausea/emesis: 1 point
	TendernessRightIliacFossa           bool // (5) tenderness over the right iliac fossa: 2 points
	Leukocytosis                        bool // (6) leukocytosis: 1 point
	PolymorphonuclearNeutrophilia       bool // (7) polymorphonuclear neutrophilia: 1 point
	MigrationOfPain                     bool // (8) migration of pain: 1 point
}

// PediatricAppendicitisScore returns Samuel's Pediatric Appendicitis Score (PAS), 0–10: one point for each
// variable, except the two physical signs (1) and (5), which score 2. No band is returned (the abstract
// states none).
// Errors: none; every combination of findings is valid.
//
// ID:      MED-GASTRO-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Samuel M. Pediatric appendicitis score. J Pediatr Surg. 2002;37(6):877-881. doi:10.1053/jpsu.2002.32893. PMID 12037754. Abstract: "Each of these variables was assigned a score of 1, except for physical signs (1 and 5), which were scored 2 to obtain a total of 10."
// Tags:    gastro, appendicitis, pediatrics, pas, score
// Units:   result points 0–10
// Range:   population: 1,170 children aged 4 to 15 years with abdominal pain suggestive of acute appendicitis (no age input)
// Example: RLQ cough/hop tenderness, anorexia, leukocytosis → 2 + 1 + 1 = 4 (HAND-DERIVED from the abstract)
func PediatricAppendicitisScore(f PediatricAppendicitisFindings) int {
	points := 0
	add := func(on bool, p int) {
		if on {
			points += p
		}
	}
	add(f.CoughPercussionHoppingTendernessRLQ, 2)
	add(f.Anorexia, 1)
	add(f.Pyrexia, 1)
	add(f.NauseaEmesis, 1)
	add(f.TendernessRightIliacFossa, 2)
	add(f.Leukocytosis, 1)
	add(f.PolymorphonuclearNeutrophilia, 1)
	add(f.MigrationOfPain, 1)
	return points
}
