package psychiatry

// Problem levels of the MDQ's third question, for MDQScreen.
const (
	MDQNoProblem       = 0 // "no"
	MDQMinorProblem    = 1 // "minor"
	MDQModerateProblem = 2 // "moderate"
	MDQSeriousProblem  = 3 // "serious problems"
)

// MDQResult holds the Mood Disorder Questionnaire symptom count and its screening rule.
type MDQResult struct {
	SymptomCount int  // "yes" answers to the 13 symptom questions, 0–13
	Positive     bool // 7 or more, "yes" to symptom clustering, and a moderate or serious problem
}

// MDQScreen returns the Mood Disorder Questionnaire symptom count and whether the screen is positive: "A score
// of 7 or more on the first 13 items, 'yes' to symptom clustering, and moderate or greater problems". Items not
// embedded.
// Errors: ErrOutOfRange when problemLevel is outside MDQNoProblem–MDQSeriousProblem.
//
// ID:      MED-PSYCH-014
// Time:    O(1)   13 items, fixed
// Space:   O(1), no allocation
// Source:  Hirschfeld RM, Williams JB, Spitzer RL, et al. Development and validation of a screening instrument for bipolar spectrum disorder: the Mood Disorder Questionnaire. Am J Psychiatry. 2000;157(11):1873-1875. doi:10.1176/appi.ajp.157.11.1873, PMID 11058490 (abstract: "score of 7 or more items"; full text not open)
// Source:  Nagata T, et al. BMC Psychiatry. 2013;13:69. doi:10.1186/1471-244X-13-69, PMC3599106 ("13 dichotomous (yes/no) symptom questions"; "(no, minor, moderate, or serious problems)"; the positive rule, quoted)
// Source:  Cybulski M, et al. Ann Gen Psychiatry. 2020;19:35. doi:10.1186/s12991-020-00285-4, PMC7254768 ("at least 7 'yes' answers"; "at least two of these symptoms at the same period of life"; "a moderate or serious problem")
// Verified: secondary (Nagata 2013 PMC3599106; Cybulski 2020 PMC7254768)
// Tags:    psychiatry, bipolar, mood, mdq, questionnaire, screen
// Units:   symptom answers yes/no; sameTimePeriod yes/no; problemLevel 0–3; result count and rule
// Range:   problemLevel 0–3
// Example: 7 yes, clustered, moderate problem → {7, true}; 13 yes, clustered, minor problem → {13, false}
func MDQScreen(symptomsYes [13]bool, sameTimePeriod bool, problemLevel int) (MDQResult, error) {
	if err := checkPoints("problemLevel", problemLevel, MDQNoProblem, MDQSeriousProblem); err != nil {
		return MDQResult{}, err
	}
	n := 0
	for _, yes := range symptomsYes {
		if yes {
			n++
		}
	}
	return MDQResult{SymptomCount: n, Positive: n >= 7 && sameTimePeriod && problemLevel >= MDQModerateProblem}, nil
}
