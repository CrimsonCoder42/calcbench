package psychiatry

// ORTOUDScore returns the revised Opioid Risk Tool for opioid use disorder (ORT-OUD) total, 0–9: one point for
// each "yes" among its nine unweighted items (family history of alcohol, illegal drug and prescription drug
// abuse; personal history of the same three; age 16–45 years; history of depression; ADD, OCD, bipolar
// disorder or schizophrenia). The source reports scores ≥ 3 as its cut-off; applying it is the caller's.
//
// ID:      MED-PSYCH-006
// Time:    O(1)   9 items, fixed
// Space:   O(1), no allocation
// Source:  Cheatle MD, Compton PA, Dhingra L, Wasser TE, O'Brien CP. Development of the revised Opioid Risk Tool to predict opioid use disorder in patients with chronic nonmalignant pain. J Pain. 2019;20(7):842-851. doi:10.1016/j.jpain.2019.01.011, PMID 30690168, PMC6768552 (nine unweighted yes = 1 items, same for men and women; "scores ≥3")
// Tags:    psychiatry, substance use, opioid, ort-oud, questionnaire, score
// Units:   answers yes (true) / no (false); result count 0–9
// Range:   all inputs boolean
// Example: age 16–45 and history of depression and one other psychiatric item → 3 (the count)
func ORTOUDScore(items [9]bool) int {
	n := 0
	for _, yes := range items {
		if yes {
			n++
		}
	}
	return n
}
