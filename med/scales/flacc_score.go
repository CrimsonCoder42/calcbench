package scales

import "fmt"

// FLACCScore returns the FLACC behavioural pain score, 0–10: the sum of the five categories Face, Legs,
// Activity, Cry and Consolability, each rated 0, 1 or 2. Items not embedded: the behavioural descriptors of
// each rating are the instrument's copyrighted text (Merkel 1997); only the arithmetic is here. No bands are
// returned: the bands found (0, 1–3, 4–6, 7–10) are not cited to Merkel 1997 by either secondary.
// Errors: ErrOutOfRange when a category is outside 0–2, naming the category.
//
// ID:      MED-SCALE-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Merkel SI, Voepel-Lewis T, Shayevitz JR, Malviya S. The FLACC: a behavioral scale for scoring postoperative pain in young children. Pediatr Nurs. 1997;23(3):293-297. PMID 9220806 (primary; abstract names the five categories; full text not open)
// Source:  Uitti JM, et al. BMC Pediatr. 2018;18:392. PMC6302518, Table 1 (credit line "From Merkel SI, Voepel-Lewis T, Shayevitz JR, Malviya S. The FLACC ... Pediatr Nurs. 1997;23:293–297") and Methods ("Each of the five categories Face; Legs; Activity; Cry; Consolability is scored from 0 to 2, which results in a total score between 0 and 10")
// Source:  Curtis SJ, et al. BMC Pediatr. 2007;7:27. PMC1950500, FLACC table (credit line "Adapted from The FLACC: A behavioral scale for scoring postoperative pain in young children, by S Merkel and others, 1997, Pediatr Nurse 23(3), p. 293–297"; each category 0–2; "Ten is the maximum score")
// Verified: secondary (Uitti 2018 PMC6302518 Table 1; Curtis 2007 PMC1950500)
// Tags:    scales, pain, flacc, pediatric, behavioural, observational, score
// Units:   each category points 0–2; result points 0–10
// Range:   face, legs, activity, cry, consolability each 0–2
// Example: face 1, legs 0, activity 1, cry 2, consolability 1 → 5; all 2 → 10 (HAND-DERIVED from the sum)
func FLACCScore(face, legs, activity, cry, consolability int) (int, error) {
	for _, e := range [...]error{
		itemZeroToTwo("face", face), itemZeroToTwo("legs", legs), itemZeroToTwo("activity", activity),
		itemZeroToTwo("cry", cry), itemZeroToTwo("consolability", consolability),
	} {
		if e != nil {
			return 0, e
		}
	}
	return face + legs + activity + cry + consolability, nil
}

// itemZeroToTwo returns ErrOutOfRange, wrapped with name, unless 0 ≤ v ≤ 2 (the 0–2 item scale of FLACC and
// PAINAD).
func itemZeroToTwo(name string, v int) error {
	if v < 0 || v > 2 {
		return fmt.Errorf("%w: %s = %d, want 0–2", ErrOutOfRange, name, v)
	}
	return nil
}
