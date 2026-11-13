package psychiatry

// PCPTSD5Score returns the Primary Care PTSD Screen for DSM-5 total, 0–5: 0 when the trauma-exposure question
// is answered "No"; otherwise the count of "yes" answers to its five symptom items. The source reports 3 as
// "The optimally sensitive score"; applying it is the caller's.
//
// ID:      MED-PSYCH-005
// Time:    O(1)   5 items, fixed
// Space:   O(1), no allocation
// Source:  Prins A, Bovin MJ, Smolenski DJ, et al. The Primary Care PTSD Screen for DSM-5 (PC-PTSD-5): development and evaluation within a veteran primary care sample. J Gen Intern Med. 2016;31(10):1206-1211. doi:10.1007/s11606-016-3703-5, PMID 27170304, PMC5023594 ("If 'No,' screen total = 0; if 'Yes,' continue with screening"; five yes/no items; cut score 3)
// Tags:    psychiatry, ptsd, trauma, pc-ptsd-5, questionnaire, score
// Units:   answers yes (true) / no (false); result count 0–5
// Range:   all inputs boolean
// Example: exposure yes, four symptom items yes → 4; exposure no → 0 whatever the items (source)
func PCPTSD5Score(traumaExposure bool, items [5]bool) int {
	if !traumaExposure {
		return 0
	}
	n := 0
	for _, yes := range items {
		if yes {
			n++
		}
	}
	return n
}
