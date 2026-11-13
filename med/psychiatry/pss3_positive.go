package psychiatry

// PSS3Positive applies the positive-screen rule of the three-item Patient Safety Screener (ED-SAFE): positive
// for "individuals who either endorse active ideation or report a suicide attempt within the past 6 months". The
// depressed-mood item does not enter the rule. Items not embedded.
//
// ID:      MED-PSYCH-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Boudreaux ED, Miller I, Goldstein AB, et al. The Emergency Department Safety Assessment and Follow-up Evaluation (ED-SAFE): method and design considerations. Contemp Clin Trials. 2013;36(1):14-24. doi:10.1016/j.cct.2013.05.008, PMID 23707435, PMC3979300, Methods ("The three-item Patient Safety Screener (PSS) ... begins with the depressed mood item"; "The other two items assess active ideation ... and lifetime suicide attempt"; the positive rule, quoted)
// Tags:    psychiatry, suicide, screening, pss-3, ed-safe, decision rule
// Units:   yes/no answers; result true when the screen is positive
// Range:   all inputs boolean
// Example: active ideation → true; lifetime attempt 2 years ago, no ideation → false; attempt 3 months ago → true
func PSS3Positive(activeIdeation, attemptWithinPast6Months bool) bool {
	return activeIdeation || attemptWithinPast6Months
}
