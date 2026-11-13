package psychiatry

// DAST10Score returns the Drug Abuse Screening Test (DAST-10) total, 0–10: "Each question answered with a 'Yes'
// scores 1 point, except for Question 3, which scores 1 point for 'No'." answersYes[2] is question 3. Items not
// embedded.
//
// ID:      MED-PSYCH-011
// Time:    O(1)   10 items, fixed
// Space:   O(1), no allocation
// Source:  Skinner HA. The drug abuse screening test. Addict Behav. 1982;7(4):363-371. doi:10.1016/0306-4603(82)90005-3, PMID 7183189 (the 28-item DAST; full text not open)
// Source:  Lam LP, Leung WC, Ip P, et al. Validation of the Drug Abuse Screening Test (DAST-10): a study on illicit drug use among Chinese pregnant women. Sci Rep. 2015;5:11420. doi:10.1038/srep11420, PMID 26091290, PMC4473689, Methods (scoring, quoted)
// Source:  Magableh RM, et al. Sci World J. 2025;2025:5597067. doi:10.1155/tswj/5597067, PMC12149511, Methods ("Each yes answer receives 1 point, except for the third question (i.e., no = 1 point), with a possible score range of 0–10")
// Verified: secondary (Lam 2015 PMC4473689; Magableh 2025 PMC12149511)
// Tags:    psychiatry, substance use, drug use, dast-10, questionnaire, score
// Units:   answers yes (true) / no (false); result points 0–10
// Range:   all inputs boolean
// Example: all "yes" → 9 (question 3 scores for "no"); all "no" → 1
func DAST10Score(answersYes [10]bool) int {
	n := 0
	for i, yes := range answersYes {
		if yes != (i == 2) {
			n++
		}
	}
	return n
}

// DAST10Level is a DAST-10 level; its value is the source's words, verbatim.
type DAST10Level string

// Lam 2015, citing Bohn 1991: "0: No problem; 1–2 Low level; 3–5 Moderate level; 6–8 Substantial level; 9–10
// Severe level".
const (
	DAST10NoProblem   DAST10Level = "No problem"        // 0
	DAST10Low         DAST10Level = "Low level"         // 1–2
	DAST10Moderate    DAST10Level = "Moderate level"    // 3–5
	DAST10Substantial DAST10Level = "Substantial level" // 6–8
	DAST10Severe      DAST10Level = "Severe level"      // 9–10
)

// DAST10LevelOf returns the DAST-10 level for a total (0–10): 0, 1–2, 3–5, 6–8, 9–10.
// Errors: ErrOutOfRange when total is outside 0–10.
//
// ID:      MED-PSYCH-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bohn MJ, Babor TF, Kranzler HR. Validity of the Drug Abuse Screening Test (DAST-10) in inpatient substance abusers. NIDA Res Monogr. 1991;119:233 (not open)
// Source:  Lam LP, Leung WC, Ip P, et al. Sci Rep. 2015;5:11420. doi:10.1038/srep11420, PMC4473689, Methods (the five levels with their ranges, citing Bohn 1991)
// Source:  Magableh RM, et al. Sci World J. 2025;2025:5597067. doi:10.1155/tswj/5597067, PMC12149511, Results (low 1–2, moderate 3–5, substantial 6–8, severe 9–10)
// Verified: secondary (Lam 2015 PMC4473689; Magableh 2025 PMC12149511)
// Tags:    psychiatry, substance use, drug use, dast-10, category
// Units:   total points 0–10; result the source's words
// Range:   total 0–10
// Example: 0 → "No problem"; 2 → "Low level"; 3 → "Moderate level"; 9 → "Severe level"
func DAST10LevelOf(total int) (DAST10Level, error) {
	if err := checkPoints("total", total, 0, 10); err != nil {
		return "", err
	}
	switch {
	case total == 0:
		return DAST10NoProblem, nil
	case total <= 2:
		return DAST10Low, nil
	case total <= 5:
		return DAST10Moderate, nil
	case total <= 8:
		return DAST10Substantial, nil
	}
	return DAST10Severe, nil
}
