package scales

// DangerAssessment5Count returns the number of "yes" answers, 0–5, to the five questions of the brief
// emergency-department intimate-partner-violence risk assessment Snider et al. derived from "20 items of
// Campbell's Danger Assessment" (the DA-5). The study reports that "A positive answer to any three questions
// has a sensitivity of 83%"; applying that threshold is the caller's. No question wording is embedded.
//
// ID:      MED-SCALE-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Snider C, Webster D, O'Sullivan CS, Campbell J. Intimate partner violence: development of a brief risk assessment for the emergency department. Acad Emerg Med. 2009;16(11):1208-1216. doi:10.1111/j.1553-2712.2009.00457.x, PMID 20053241, abstract
// Tags:    scales, danger assessment, intimate partner violence, screening, questionnaire
// Units:   each answer yes (true) or no; result count 0–5
// Range:   any combination of the five answers
// Example: three yes → 3; none → 0; all → 5 (abstract)
func DangerAssessment5Count(q1, q2, q3, q4, q5 bool) int {
	n := 0
	for _, yes := range [...]bool{q1, q2, q3, q4, q5} {
		if yes {
			n++
		}
	}
	return n
}
