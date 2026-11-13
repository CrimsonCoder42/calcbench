package psychiatry

import "fmt"

// checkPoints returns ErrOutOfRange when an item is outside [lo, hi], naming it.
func checkPoints(name string, v, lo, hi int) error {
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %d, want %d–%d", ErrOutOfRange, name, v, lo, hi)
	}
	return nil
}

// AUDITCScore returns the AUDIT-C total, 0–12: the sum of the first three AUDIT consumption questions
// (frequency of drinking, typical quantity, frequency of heavy drinking), each answered on 0–4 points. The cited
// cut-offs, ≥ 3 "typically for women" and ≥ 4 "typically for men", are the caller's.
// Errors: ErrOutOfRange when an item is outside 0–4, naming it.
//
// ID:      MED-PSYCH-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bush K, Kivlahan DR, McDonell MB, Fihn SD, Bradley KA. The AUDIT alcohol consumption questions (AUDIT-C): an effective brief screening test for problem drinking. Arch Intern Med. 1998;158(16):1789-1795. doi:10.1001/archinte.158.16.1789, PMID 9738608 (full text not open)
// Source:  Flentje A, Barger BT, Capriotti MR, et al. Screening gender minority people for harmful alcohol use. PLoS One. 2020;15(4):e0231022. doi:10.1371/journal.pone.0231022, PMC7138294, Methods ("summing AUDIT consumption items 1, 2, and 3, for a total score of 0–12"; "items answered on a scale of 0–4"; "cutoff scores of ≥ 3 (typically for women) and ≥ 4 (typically for men)")
// Source:  Bottenfield M, et al. Cureus. 2026;18(2):e103164. doi:10.7759/cureus.103164, PMC12970570 (AUDIT-C "greater than or equal to three (female) and four (male)")
// Verified: secondary (Flentje 2020 PMC7138294; Bottenfield 2026 PMC12970570)
// Tags:    psychiatry, alcohol, substance use, audit-c, questionnaire, score
// Units:   item points 0–4; result points 0–12
// Range:   each item 0–4
// Example: 2, 1, 1 → 4; all 4 → 12 (the sum)
func AUDITCScore(frequency, typicalQuantity, heavyDrinkingFrequency int) (int, error) {
	if err := checkPoints("frequency", frequency, 0, 4); err != nil {
		return 0, err
	}
	if err := checkPoints("typicalQuantity", typicalQuantity, 0, 4); err != nil {
		return 0, err
	}
	if err := checkPoints("heavyDrinkingFrequency", heavyDrinkingFrequency, 0, 4); err != nil {
		return 0, err
	}
	return frequency + typicalQuantity + heavyDrinkingFrequency, nil
}
