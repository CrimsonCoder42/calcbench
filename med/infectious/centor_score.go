package infectious

// CentorScore returns the Centor score, 0–4: the number of the four clinical variables present — "tonsillar
// exudates, swollen tender anterior cervical nodes, lack of a cough, and history of fever" (adults in the
// emergency room). The paper reports the probability of a positive culture by count ("all 4 variables had a
// 56% probability ...; 3 variables, 32%; 2 variables, 15%; 1 variable, 6.5%; and 0 variables, 2.5%"); those
// are its population's figures and are not returned. The McIsaac age adjustment is a different score and is
// not applied here.
//
// ID:      MED-INFECT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Centor RM, Witherspoon JM, Dalton HP, Brody CE, Link K. The diagnosis of strep throat in adults in the emergency room. Med Decis Making. 1981;1(3):239-246. doi:10.1177/0272989X8100100304, PMID 6763125, abstract
// Tags:    infectious, pharyngitis, streptococcus, centor, score
// Units:   each variable present or absent; result count 0–4
// Range:   any combination of the four variables
// Example: exudates, nodes, no cough, fever → 4; fever alone → 1 (abstract)
func CentorScore(tonsillarExudates, tenderAnteriorCervicalNodes, coughAbsent, historyOfFever bool) int {
	n := 0
	for _, v := range [...]bool{tonsillarExudates, tenderAnteriorCervicalNodes, coughAbsent, historyOfFever} {
		if v {
			n++
		}
	}
	return n
}
