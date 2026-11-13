package scales

// HARKScore returns the HARK score, 0–4: one point for each "yes" to the four HARK questions (Humiliation,
// Afraid, Rape, Kick): "one point is given for every yes answer" (Table 1). The study reported "a HARK cut off
// score of ≥1"; applying it is the caller's. No question wording is embedded.
//
// ID:      MED-SCALE-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Sohal H, Eldridge S, Feder G. The sensitivity and specificity of four questions (HARK) to identify intimate partner violence: a diagnostic accuracy study in general practice. BMC Fam Pract. 2007;8:49. doi:10.1186/1471-2296-8-49, PMID 17727730, Table 1 and Results
// Tags:    scales, hark, intimate partner violence, screening, questionnaire
// Units:   each answer yes (true) or no; result points 0–4
// Range:   any combination of the four answers
// Example: humiliation only → 1; all four → 4; none → 0 (Table 1)
func HARKScore(humiliation, afraid, rape, kick bool) int {
	n := 0
	for _, yes := range [...]bool{humiliation, afraid, rape, kick} {
		if yes {
			n++
		}
	}
	return n
}
