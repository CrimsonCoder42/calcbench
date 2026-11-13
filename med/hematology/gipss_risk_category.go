package hematology

import "fmt"

// GIPSSCategory is a GIPSS risk category; its value is the source's label, verbatim.
type GIPSSCategory string

// Tefferi 2018 GIPSS risk categories.
const (
	GIPSSLow           GIPSSCategory = "Low"            // 0
	GIPSSIntermediate1 GIPSSCategory = "Intermediate-1" // 1
	GIPSSIntermediate2 GIPSSCategory = "Intermediate-2" // 2
	GIPSSHigh          GIPSSCategory = "High"           // ≥ 3
)

// GIPSSRiskCategory returns the GIPSS risk category for a score: 0 "Low"; 1 "Intermediate-1"; 2
// "Intermediate-2"; ≥ 3 "High".
// Errors: ErrOutOfRange when score is outside 0–6.
//
// ID:      MED-HEMAT-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tefferi A, Guglielmelli P, Nicolosi M, et al. GIPSS: genetically inspired prognostic scoring system for primary myelofibrosis. Leukemia. 2018;32(7):1631-1642. doi:10.1038/s41375-018-0107-z, PMID 29654267, Results and Figure 1
// Tags:    hematology, myelofibrosis, gipss, category
// Units:   score points 0–6; result the source's label
// Range:   0 ≤ score ≤ 6
// Example: 0 → "Low"; 1 → "Intermediate-1"; 2 → "Intermediate-2"; 3 → "High"
func GIPSSRiskCategory(score int) (GIPSSCategory, error) {
	if score < 0 || score > 6 {
		return "", fmt.Errorf("%w: score = %d, want 0–6", ErrOutOfRange, score)
	}
	switch score {
	case 0:
		return GIPSSLow, nil
	case 1:
		return GIPSSIntermediate1, nil
	case 2:
		return GIPSSIntermediate2, nil
	}
	return GIPSSHigh, nil
}
