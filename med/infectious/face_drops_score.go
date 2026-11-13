package infectious

// FACEDROPSScore returns the FACE DROPS score, −1 to 27, for acute facial palsy in a Lyme-endemic area: Fever
// (temperature ≥ 100.4 °F) +8; Aches (arthralgia/myalgia) +6; Cephalalgia (headache other than isolated otalgia
// or postauricular pain) +3; Exhaustion (unusual fatigue around palsy onset) +4; Dermatomal or radicular signs
// +3; Otalgia or postauricular pain ipsilateral to the droop −1; Stiff neck (nuchal rigidity) +3. The paper
// reports "scores ≤4 predicted BP with ≥93.5% accuracy while scores of ≥7 predicted LDFP"; those accuracies are
// its population's figures and are not returned.
// Errors: none; every combination of the items is defined.
//
// ID:      MED-INFECT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  McEntire CRS, Chung SY, Chang B, et al. FACE DROPS: A clinical risk assessment tool for differentiation of acute Lyme disease-associated facial palsy from Bell palsy. Neurol Clin Pract. 2025;15(3):e200476. doi:10.1212/CPJ.0000000000200476, PMID 40290705, PMC12021021, Table 3 (scoring) and Table 4 (cut-offs ≤4, 5–6, ≥7)
// Tags:    infectious, lyme disease, facial palsy, bell palsy, face drops, score
// Units:   each item present or absent; result points −1 to 27
// Range:   any combination of the items (ages 9–85, Lyme-endemic northeastern United States, the study's population)
// Example: fever, aches, stiff neck and otalgia → 16; otalgia alone → −1; all positive items, no otalgia → 27 (Table 3)
func FACEDROPSScore(fever, aches, cephalalgia, exhaustion, dermatomalOrRadicular, otalgiaOrPostauricularPain, stiffNeck bool) int {
	s := 0
	if fever {
		s += 8
	}
	if aches {
		s += 6
	}
	if cephalalgia {
		s += 3
	}
	if exhaustion {
		s += 4
	}
	if dermatomalOrRadicular {
		s += 3
	}
	if otalgiaOrPostauricularPain {
		s--
	}
	if stiffNeck {
		s += 3
	}
	return s
}
