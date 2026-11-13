package respiratory

import "fmt"

// SilvermanAndersenScore returns the Silverman–Andersen respiratory severity score (0–10), the sum of five
// criteria each graded 0, 1 or 2 by the observer: "five clinical categories, i.e., sternal and intercostal
// retractions, nasal flaring, expiratory grunting, and upper chest movements, with 0–2 points for each category.
// Zero points indicate normal respiration and ten points severe respiratory distress". grades[i] is the caller's
// grade for criterion i+1, in the order upper chest movement, lower chest (intercostal) retraction, xiphoid
// retraction, nares dilatation, expiratory grunt; the order does not change the sum.
// Errors: ErrOutOfRange for a grade outside 0–2.
//
// ID:      MED-RESP-049
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Silverman WA, Andersen DH. A controlled clinical trial of effects of water mist on obstructive respiratory signs, death rate and necropsy findings among premature infants. Pediatrics. 1956;17(1):1-10. PMID 13353856 (primary; no abstract, not opened)
// Source:  Brenne H, Follestad T, Bergseng H, Eriksen BH, Søraunet K, Grunewaldt KH. PLoS One. 2023, PMC10313036 ("from zero to two points for each item … The sum score is therefore in the range from zero to 10"; ref 9 Silverman & Andersen 1956)
// Source:  Nussbaum C, Lengauer M, Puchwein-Schwepcke AF, Weiss VBN, et al. Noninvasive ventilation in preterm infants: factors influencing weaning decisions and the role of the Silverman-Andersen score. Children (Basel). 2022;9(9):1292. doi:10.3390/children9091292. PMID 36138599, PMC9497739, Methods (five categories, 0–2 points each)
// Verified: secondary (Brenne 2023 PMC10313036, ref 9; Nussbaum 2022 Children PMC9497739, ref 10)
// Tags:    respiratory, neonatal, respiratory-distress, silverman-andersen, score
// Units:   grades five integers 0–2; result points 0–10
// Range:   each grade 0–2; newborn infants
// Example: grades 1, 1, 0, 2, 1 → 5 (hand-derived: the sum)
func SilvermanAndersenScore(grades [5]int) (int, error) {
	s := 0
	for i, g := range grades {
		if g < 0 || g > 2 {
			return 0, fmt.Errorf("%w: grades[%d] = %d, want 0 ≤ grade ≤ 2", ErrOutOfRange, i, g)
		}
		s += g
	}
	return s, nil
}
