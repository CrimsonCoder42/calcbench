package psychiatry

// CRAFFTScore returns the CRAFFT total, 0–6: the count of "yes" answers to its six items (the mnemonic's
// Car, Relax, Alone, Forget, Friends, Trouble). The source reports "Two or more yes answers" as its cut-off;
// applying it is the caller's. The later CRAFFT 2.x gating by prior-use questions is not in this source and is
// not applied.
//
// ID:      MED-PSYCH-007
// Time:    O(1)   6 items, fixed
// Space:   O(1), no allocation
// Source:  Knight JR, Shrier LA, Bravender TD, Farrell M, Vander Bilt J, Shaffer HJ. A new brief screen for adolescent substance abuse. Arch Pediatr Adolesc Med. 1999;153(6):591-596. doi:10.1001/archpedi.153.6.591, PMID 10357299, abstract ("6 items were arranged into a mnemonic (CRAFFT)"; "Two or more yes answers")
// Tags:    psychiatry, substance use, adolescent, crafft, questionnaire, score
// Units:   answers yes (true) / no (false); result count 0–6
// Range:   all inputs boolean
// Example: Car and Relax yes, the rest no → 2 (the count)
func CRAFFTScore(items [6]bool) int {
	n := 0
	for _, yes := range items {
		if yes {
			n++
		}
	}
	return n
}
