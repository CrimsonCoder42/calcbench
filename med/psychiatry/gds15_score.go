package psychiatry

// GDS15Score returns the 15-item Geriatric Depression Scale total, 0–15: one point for each "yes" answer, except
// items 1, 5, 7, 11 and 13 (answersYes[0], [4], [6], [10], [12]), which score one point for "no". Items not
// embedded.
//
// ID:      MED-PSYCH-015
// Time:    O(1)   15 items, fixed
// Space:   O(1), no allocation
// Source:  Sheikh JI, Yesavage JA. Geriatric Depression Scale (GDS): recent evidence and development of a shorter version. Clin Gerontol. 1986;5(1-2):165-173. doi:10.1300/J018v05n01_09 (not in PubMed; not open)
// Source:  Tian H, Wang Y. Behav Sci (Basel). 2023;13(2):153. doi:10.3390/bs13020153, PMC9952244 ("Each 'yes' is worth 1 point, and 'no' is worth 0 points. Items 1, 5, 7, 11, and 13 are scored in reverse")
// Source:  Islam MA, et al. Glob Ment Health. 2024, PMC11504949 ("Items 1, 5, 7, 11, and 13 … when answered negatively … ranging from 0 to 15"; ref 36 Sheikh & Yesavage 1986)
// Verified: secondary (Tian 2023 PMC9952244, ref 37 Sheikh 1986; Islam 2024 PMC11504949, ref 36 Sheikh 1986)
// Tags:    psychiatry, depression, geriatric, gds-15, questionnaire, score
// Units:   answers yes (true) / no (false); result points 0–15
// Range:   all inputs boolean
// Example: every answer "yes" → 10; every answer "no" → 5
func GDS15Score(answersYes [15]bool) int {
	n := 0
	for i, yes := range answersYes {
		reverse := i == 0 || i == 4 || i == 6 || i == 10 || i == 12
		if yes != reverse {
			n++
		}
	}
	return n
}
