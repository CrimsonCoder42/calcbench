package pediatrics

import "fmt"

// DHAKADehydrationScore returns the DHAKA Dehydration Score, 0–12 (Table 3), from the finding level of each
// sign: generalAppearance 0 normal, 1 restless/irritable, 2 lethargic/unconscious → 0, 2, 4 points;
// respirations 0 normal, 1 deep → 0, 2; skinPinch 0 normal, 1 slow, 2 very slow → 0, 2, 4; tears 0 normal,
// 1 decreased, 2 absent → 0, 1, 2.
// Errors: ErrOutOfRange when a level is outside its range, naming the sign.
//
// ID:      MED-PEDS-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Levine AC, Glavis-Bloom J, Modi P, et al. Empirically derived dehydration scoring and decision tree models for children with diarrhea. Glob Health Sci Pract. 2015;3(3):405-418. doi:10.9745/GHSP-D-15-00097, PMID 26374802, Table 3
// Tags:    pediatrics, dehydration, diarrhea, dhaka, score
// Units:   each sign's level (category code); result points 0–12
// Range:   generalAppearance 0–2; respirations 0–1; skinPinch 0–2; tears 0–2
// Example: restless (2), normal breathing (0), slow pinch (2), decreased tears (1) → 5; all worst → 12 (Table 3)
func DHAKADehydrationScore(generalAppearance, respirations, skinPinch, tears int) (int, error) {
	switch {
	case generalAppearance < 0 || generalAppearance > 2:
		return 0, fmt.Errorf("%w: generalAppearance = %d, want 0–2", ErrOutOfRange, generalAppearance)
	case respirations < 0 || respirations > 1:
		return 0, fmt.Errorf("%w: respirations = %d, want 0–1", ErrOutOfRange, respirations)
	case skinPinch < 0 || skinPinch > 2:
		return 0, fmt.Errorf("%w: skinPinch = %d, want 0–2", ErrOutOfRange, skinPinch)
	case tears < 0 || tears > 2:
		return 0, fmt.Errorf("%w: tears = %d, want 0–2", ErrOutOfRange, tears)
	}
	return 2*generalAppearance + 2*respirations + 2*skinPinch + tears, nil
}
