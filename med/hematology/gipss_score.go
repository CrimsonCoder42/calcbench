package hematology

import "fmt"

// GIPSSKaryotype is the karyotype class GIPSS scores.
type GIPSSKaryotype int

// Karyotype classes of Tefferi 2018 (the definitions of each class are the source's).
const (
	GIPSSKaryotypeFavorable    GIPSSKaryotype = iota // 0 points
	GIPSSKaryotypeUnfavorable                        // 1 point
	GIPSSKaryotypeVeryHighRisk                       // "VHR", 2 points
)

// GIPSSScore returns the GIPSS (genetically inspired prognostic scoring system) score for primary
// myelofibrosis, 0–6: very-high-risk karyotype 2, unfavorable karyotype 1; absence of a type 1/like CALR
// mutation 1; ASXL1 mutation 1; SRSF2 mutation 1; U2AF1 Q157 mutation 1.
// Errors: ErrOutOfRange for a karyotype other than the three constants.
//
// ID:      MED-HEMAT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tefferi A, Guglielmelli P, Nicolosi M, et al. GIPSS: genetically inspired prognostic scoring system for primary myelofibrosis. Leukemia. 2018;32(7):1631-1642. doi:10.1038/s41375-018-0107-z, PMID 29654267, Results "Development of a new risk model (GIPSS)"
// Tags:    hematology, myelofibrosis, gipss, prognostic score
// Units:   karyotype class; mutations present/absent; result points 0–6
// Range:   karyotype one of GIPSSKaryotypeFavorable, GIPSSKaryotypeUnfavorable, GIPSSKaryotypeVeryHighRisk
// Example: unfavorable karyotype, no type 1/like CALR, ASXL1 → 3; favorable, type 1/like CALR, no mutations → 0
func GIPSSScore(karyotype GIPSSKaryotype, calrType1LikeAbsent, asxl1, srsf2, u2af1Q157 bool) (int, error) {
	if karyotype < GIPSSKaryotypeFavorable || karyotype > GIPSSKaryotypeVeryHighRisk {
		return 0, fmt.Errorf("%w: karyotype = %d, want 0–2", ErrOutOfRange, karyotype)
	}
	return int(karyotype) + pts(calrType1LikeAbsent, 1) + pts(asxl1, 1) + pts(srsf2, 1) + pts(u2af1Q157, 1), nil
}
