package neuro

import "fmt"

// SICHNCCT is the non-contrast CT categorization of the secondary intracerebral hemorrhage score; its value
// is the source's points.
type SICHNCCT int

// Delgado Almandoz 2010, Table 2: "High probability: 2", "Indeterminate: 1", "Low probability: 0".
const (
	SICHNCCTLowProbability  SICHNCCT = 0
	SICHNCCTIndeterminate   SICHNCCT = 1
	SICHNCCTHighProbability SICHNCCT = 2
)

// SICHScore returns the secondary intracerebral hemorrhage (sICH) score (0–6): NCCT categorization (low
// probability 0, indeterminate 1, high probability 2) + age 18–45 years 2, 46–70 years 1, ≥ 71 years 0 +
// female sex 1 + neither known hypertension nor impaired coagulation 1.
// Errors: ErrOutOfRange outside Range.
//
// ID:      MED-NEURO-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Delgado Almandoz JE, Schaefer PW, Goldstein JN, et al. Practical scoring system for the identification of patients with intracerebral hemorrhage at highest risk of harboring an underlying vascular etiology: the Secondary Intracerebral Hemorrhage Score. AJNR Am J Neuroradiol. 2010;31(9):1653-1660. doi:10.3174/ajnr.A2156, PMID 20581068, PMC3682824, Table 2 ("The SICH score is calculated by adding the total number of points")
// Tags:    neuro, stroke, intracerebral hemorrhage, sich, vascular etiology, score
// Units:   ageYears completed years; result points 0–6
// Range:   ncct 0–2; ageYears ≥ 18 (the table's youngest band)
// Example: high probability CT, age 40, female, no hypertension or coagulopathy → 2 + 2 + 1 + 1 = 6 (hand-derived from Table 2)
func SICHScore(ncct SICHNCCT, ageYears int, female, neitherHypertensionNorImpairedCoagulation bool) (int, error) {
	if ncct < SICHNCCTLowProbability || ncct > SICHNCCTHighProbability {
		return 0, fmt.Errorf("%w: ncct = %d, want 0–2", ErrOutOfRange, ncct)
	}
	if ageYears < 18 {
		return 0, fmt.Errorf("%w: ageYears = %d, want ≥ 18", ErrOutOfRange, ageYears)
	}
	n := int(ncct)
	switch {
	case ageYears <= 45:
		n += 2
	case ageYears <= 70:
		n++
	}
	if female {
		n++
	}
	if neitherHypertensionNorImpairedCoagulation {
		n++
	}
	return n, nil
}
