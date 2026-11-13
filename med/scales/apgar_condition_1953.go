package scales

import "fmt"

// ApgarCondition is one of Apgar's 1953 groups by score; its value is the source's own words, verbatim.
type ApgarCondition string

// Apgar 1953 groups, verbatim from p. 267: "groups 8, 9, and 10 indicate infants in good condition, 0, 1,
// and 2, poor condition, and the remaining scores, fair condition."
const (
	ApgarGoodCondition ApgarCondition = "good condition" // 8, 9, 10
	ApgarFairCondition ApgarCondition = "fair condition" // 3–7, "the remaining scores"
	ApgarPoorCondition ApgarCondition = "poor condition" // 0, 1, 2
)

// ApgarCondition1953 returns Apgar's 1953 group for an Apgar score, with the source's own words: 8–10
// "good condition"; 3–7 "fair condition"; 0–2 "poor condition". Later bands (AAP/ACOG) are not built here.
// Errors: ErrOutOfRange when score is outside 0–10.
//
// ID:      MED-SCALE-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Apgar V. A proposal for a new method of evaluation of the newborn infant. Curr Res Anesth Analg. 1953;32(4):260-267. PMID 13083014. p. 267: "groups 8, 9, and 10 indicate infants in good condition, 0, 1, and 2, poor condition, and the remaining scores, fair condition."
// Tags:    scales, apgar, newborn, neonatal, category
// Units:   score points 0–10; result the source's group words
// Range:   0 ≤ score ≤ 10
// Example: 2 → "poor condition"; 3 → "fair condition"; 8 → "good condition" (Apgar 1953 p. 267)
func ApgarCondition1953(score int) (ApgarCondition, error) {
	if score < 0 || score > 10 {
		return "", fmt.Errorf("%w: score = %d, want 0–10", ErrOutOfRange, score)
	}
	switch {
	case score <= 2:
		return ApgarPoorCondition, nil
	case score <= 7:
		return ApgarFairCondition, nil
	}
	return ApgarGoodCondition, nil
}
