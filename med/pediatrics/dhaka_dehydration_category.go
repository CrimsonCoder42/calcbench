package pediatrics

import "fmt"

// DHAKACategory is a DHAKA dehydration category; its value is the source's label, verbatim.
type DHAKACategory string

// Levine 2015 Table 3, "Assigned Dehydration Categories".
const (
	DHAKANoDehydration     DHAKACategory = "No dehydration"     // 0–1
	DHAKASomeDehydration   DHAKACategory = "Some dehydration"   // 2–3
	DHAKASevereDehydration DHAKACategory = "Severe dehydration" // ≥ 4
)

// DHAKADehydrationCategory returns the DHAKA category for a score: 0–1 "No dehydration"; 2–3 "Some
// dehydration"; ≥ 4 "Severe dehydration".
// Errors: ErrOutOfRange when score is outside 0–12.
//
// ID:      MED-PEDS-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Levine AC, Glavis-Bloom J, Modi P, et al. Empirically derived dehydration scoring and decision tree models for children with diarrhea. Glob Health Sci Pract. 2015;3(3):405-418. doi:10.9745/GHSP-D-15-00097, PMID 26374802, Table 3
// Tags:    pediatrics, dehydration, diarrhea, dhaka, category
// Units:   score points 0–12; result the source's label
// Range:   0 ≤ score ≤ 12
// Example: 1 → "No dehydration"; 2 → "Some dehydration"; 4 → "Severe dehydration" (Table 3)
func DHAKADehydrationCategory(score int) (DHAKACategory, error) {
	if score < 0 || score > 12 {
		return "", fmt.Errorf("%w: score = %d, want 0–12", ErrOutOfRange, score)
	}
	switch {
	case score <= 1:
		return DHAKANoDehydration, nil
	case score <= 3:
		return DHAKASomeDehydration, nil
	}
	return DHAKASevereDehydration, nil
}
