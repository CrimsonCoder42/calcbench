package cardio

import "fmt"

// DLCNLDLBand is the LDL-cholesterol group of the Dutch Lipid Clinic Network criteria, chosen by the caller
// from Table 1's printed bands (the printed bands leave gaps such as 8.4–8.5 mmol/L, so the band is an input,
// not computed from a number).
type DLCNLDLBand int

// LDL bands (Nordestgaard 2013, Table 1, group 4).
const (
	DLCNLDLBelow4  DLCNLDLBand = iota // below 4.0 mmol/L: 0
	DLCNLDL4To49                      // "4.0–4.9 mmol/L (155–190 mg/dL)": 1
	DLCNLDL5To64                      // "5.0–6.4 mmol/L (191–250 mg/dL)": 3
	DLCNLDL65To84                     // "6.5–8.4 mmol/L (251–325 mg/dL)": 5
	DLCNLDLAbove85                    // ">8.5 mmol/L (>325 mg/dL)": 8
)

// DLCNFindings are the yes/no items of the Dutch Lipid Clinic Network criteria (Table 1). Within each group
// only the highest applicable item counts.
type DLCNFindings struct {
	RelativePrematureCHD          bool // group 1 (i): 1
	RelativeLDLAbove95th          bool // group 1 (ii): 1
	RelativeXanthomaOrArcus       bool // group 1 (iii): 2
	ChildLDLAbove95th             bool // group 1 (iv): 2
	PrematureCHD                  bool // group 2 (i): 2
	PrematureCerebralOrPeripheral bool // group 2 (ii): 1
	TendonXanthoma                bool // group 3 (i): 6
	CornealArcusUnder45           bool // group 3 (ii): 4
	CausativeMutation             bool // group 5: 8 (LDLR, APOB or PCSK9)
}

// DutchLipidClinicFHScore returns the Dutch Lipid Clinic Network score for heterozygous familial
// hypercholesterolaemia in adults, 0–26: the highest applicable item in each of the five groups, summed.
// Errors: ErrOutOfRange for an LDL band other than the five constants.
//
// ID:      MED-CARDIO-038
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nordestgaard BG, Chapman MJ, Humphries SE, et al. Familial hypercholesterolaemia is underdiagnosed and undertreated in the general population. Eur Heart J. 2013;34(45):3478-3490. doi:10.1093/eurheartj/eht273. PMID 23956253. PMC3844152, Table 1; footnote "per group only one score, the highest applicable, can be chosen"
// Tags:    cardio, lipids, familial-hypercholesterolemia, dlcn, score
// Units:   ldl the caller's band; result points 0–26
// Range:   adults (Table 1 is for adults); ldl one of the five bands
// Example: tendon xanthoma, LDL 6.5–8.4 mmol/L, relative with premature CHD → 6 + 5 + 1 = 12 (HAND-DERIVED from Table 1)
func DutchLipidClinicFHScore(ldl DLCNLDLBand, f DLCNFindings) (int, error) {
	if ldl < DLCNLDLBelow4 || ldl > DLCNLDLAbove85 {
		return 0, fmt.Errorf("%w: ldl = %d, want one of the five LDL bands", ErrOutOfRange, ldl)
	}
	s := dlcnLDLPoints[ldl]
	switch {
	case f.RelativeXanthomaOrArcus, f.ChildLDLAbove95th:
		s += 2
	case f.RelativePrematureCHD, f.RelativeLDLAbove95th:
		s++
	}
	switch {
	case f.PrematureCHD:
		s += 2
	case f.PrematureCerebralOrPeripheral:
		s++
	}
	switch {
	case f.TendonXanthoma:
		s += 6
	case f.CornealArcusUnder45:
		s += 4
	}
	if f.CausativeMutation {
		s += 8
	}
	return s, nil
}

var dlcnLDLPoints = [5]int{0, 1, 3, 5, 8}

// DLCNDiagnosis is a Dutch Lipid Clinic Network group; its value is Table 1's wording.
type DLCNDiagnosis string

// Nordestgaard 2013, Table 1 footnote.
const (
	DLCNDefinite DLCNDiagnosis = "definite FH" // > 8
	DLCNProbable DLCNDiagnosis = "probable FH" // 6–8
	DLCNPossible DLCNDiagnosis = "possible FH" // 3–5
	DLCNUnlikely DLCNDiagnosis = "unlikely FH" // 0–2
)

// DutchLipidClinicFHGroup returns the Table 1 group for a DLCN score: > 8 "definite FH", 6–8 "probable FH",
// 3–5 "possible FH", 0–2 "unlikely FH".
// Errors: ErrOutOfRange when score is outside 0–26.
//
// ID:      MED-CARDIO-039
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nordestgaard BG, et al. Eur Heart J. 2013;34(45):3478-3490. doi:10.1093/eurheartj/eht273. PMID 23956253. PMC3844152, Table 1 footnote: "A 'definite FH' diagnosis can be made if the subject scores >8 points. A 'probable FH' ... 6 to 8 points. A 'possible FH' ... 3 to 5 points. An 'unlikely FH' ... 0 to 2 points."
// Tags:    cardio, lipids, familial-hypercholesterolemia, dlcn, category
// Units:   score points 0–26; result Table 1's group words
// Range:   0 ≤ score ≤ 26
// Example: 2 → "unlikely FH"; 3 → "possible FH"; 6 → "probable FH"; 9 → "definite FH"
func DutchLipidClinicFHGroup(score int) (DLCNDiagnosis, error) {
	switch {
	case score < 0 || score > 26:
		return "", errScore(score, 0, 26)
	case score > 8:
		return DLCNDefinite, nil
	case score >= 6:
		return DLCNProbable, nil
	case score >= 3:
		return DLCNPossible, nil
	default:
		return DLCNUnlikely, nil
	}
}
