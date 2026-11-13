package neuro

import "fmt"

// MTOQ4Score returns the 4-item Migraine Treatment Optimization Questionnaire (mTOQ-4) sum, 0–8. Each item's
// answer is scored "never (0), rarely (0), <half the time (1), ≥half the time (2)"; the caller passes the points.
// Errors: ErrOutOfRange when an item is outside 0–2, naming its index.
//
// ID:      MED-NEURO-033
// Time:    O(1)   4 items, fixed
// Space:   O(1), no allocation
// Source:  Lipton RB, Fanning KM, Serrano D, Reed ML, Cady R, Buse DC. Ineffective acute treatment of episodic migraine is associated with new-onset chronic migraine. Neurology. 2015;84(7):688-695. doi:10.1212/WNL.0000000000001256, PMID 25609757, PMC4336107, Methods ("Response options are scored as never (0), rarely (0), <half the time (1), ≥half the time (2)"; "Sum scores range from 0 to 8")
// Tags:    neuro, migraine, headache, mtoq-4, questionnaire, score
// Units:   item points 0–2; result points 0–8
// Range:   each item 0–2
// Example: all ≥ half the time → 8; two items at 1, two at 2 → 6
func MTOQ4Score(items [4]int) (int, error) {
	total := 0
	for i, v := range items {
		if v < 0 || v > 2 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 0–2", ErrOutOfRange, i, v)
		}
		total += v
	}
	return total, nil
}

// MTOQ4Category is an mTOQ-4 treatment-efficacy category; its value is the source's words, verbatim.
type MTOQ4Category string

// Lipton 2015: "very poor treatment efficacy (0), poor treatment efficacy (1–5), moderate treatment efficacy
// (6–7), and maximum treatment efficacy (8)".
const (
	MTOQ4VeryPoor MTOQ4Category = "very poor treatment efficacy" // 0
	MTOQ4Poor     MTOQ4Category = "poor treatment efficacy"      // 1–5
	MTOQ4Moderate MTOQ4Category = "moderate treatment efficacy"  // 6–7
	MTOQ4Maximum  MTOQ4Category = "maximum treatment efficacy"   // 8
)

// MTOQ4CategoryOf returns the mTOQ-4 treatment-efficacy category for a sum (0–8): 0, 1–5, 6–7, 8.
// Errors: ErrOutOfRange when total is outside 0–8.
//
// ID:      MED-NEURO-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lipton RB, Fanning KM, Serrano D, Reed ML, Cady R, Buse DC. Ineffective acute treatment of episodic migraine is associated with new-onset chronic migraine. Neurology. 2015;84(7):688-695. doi:10.1212/WNL.0000000000001256, PMID 25609757, PMC4336107, Methods (the four categories, quoted)
// Tags:    neuro, migraine, headache, mtoq-4, category
// Units:   total points 0–8; result the source's words
// Range:   total 0–8
// Example: 0 → very poor; 5 → poor; 6 → moderate; 8 → maximum
func MTOQ4CategoryOf(total int) (MTOQ4Category, error) {
	if err := checkInt("total", total, 0, 8); err != nil {
		return "", err
	}
	switch {
	case total == 0:
		return MTOQ4VeryPoor, nil
	case total <= 5:
		return MTOQ4Poor, nil
	case total <= 7:
		return MTOQ4Moderate, nil
	}
	return MTOQ4Maximum, nil
}
