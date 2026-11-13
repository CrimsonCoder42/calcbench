package psychiatry

// CAGEScore returns the CAGE total, 0–4: the count of "yes" answers to its four questions (Cut down, Annoyed,
// Guilty, Eye-opener). The cited cut-off, two or more "yes" answers, is the caller's.
//
// ID:      MED-PSYCH-009
// Time:    O(1)   4 items, fixed
// Space:   O(1), no allocation
// Source:  Ewing JA. Detecting alcoholism. The CAGE questionnaire. JAMA. 1984;252(14):1905-1907. doi:10.1001/jama.252.14.1905, PMID 6471323 (full text not open)
// Source:  O'Rourke T, et al. Addict Behav Rep. 2022;15:100414. doi:10.1016/j.abrep.2022.100414, PMC8863361 ("four yes/no questions"; "a cut-off of two or more questions are answered with “yes”")
// Source:  Peláez-Zuberbuhler J, et al. Alcohol Alcohol. 2025;60(3):agaf016. doi:10.1093/alcalc/agaf016, PMC12005082 ("Each item has a 0 (no) and 1 (yes) score"; "≥2 indicated problem drinking")
// Verified: secondary (O'Rourke 2022 PMC8863361; Peláez-Zuberbuhler 2025 PMC12005082)
// Tags:    psychiatry, alcohol, substance use, cage, questionnaire, score
// Units:   answers yes (true) / no (false); result count 0–4
// Range:   all inputs boolean
// Example: Cut down and Guilty yes → 2 (the count)
func CAGEScore(items [4]bool) int {
	n := 0
	for _, yes := range items {
		if yes {
			n++
		}
	}
	return n
}
