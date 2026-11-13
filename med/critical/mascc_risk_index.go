package critical

import "fmt"

// MASCCBurden is the burden-of-illness item of the MASCC risk index (Klastersky 2000 abstract).
type MASCCBurden int

const (
	// MASCCBurdenNoOrMildSymptoms is "a burden of illness indicating absence of symptoms or mild symptoms": weight 5.
	MASCCBurdenNoOrMildSymptoms MASCCBurden = iota
	// MASCCBurdenModerateSymptoms is "moderate symptoms": weight 3.
	MASCCBurdenModerateSymptoms
	// MASCCBurdenNeither is a burden of illness that is neither of the two weighted levels: weight 0.
	MASCCBurdenNeither
)

// MASCCFactors are the yes/no predictive factors of the MASCC risk index besides burden and age.
type MASCCFactors struct {
	NoHypotension             bool // "absence of hypotension": 5
	NoCOPD                    bool // "absence of chronic obstructive pulmonary disease": 4
	SolidTumorOrNoPriorFungal bool // "presence of solid tumor or absence of previous fungal infection in patients with hematologic malignancies": 4
	Outpatient                bool // "outpatient status": 3
	NoDehydration             bool // "absence of dehydration": 3
}

// MASCCResult is the MASCC risk-index score and whether it reaches the abstract's low-risk threshold.
type MASCCResult struct {
	Score   int  // 0–26
	LowRisk bool // score ≥ 21: "a Multinational Association for Supportive Care in Cancer risk-index score >/= 21 identified low-risk patients"
}

// MASCCRiskIndex returns the MASCC risk-index score for febrile neutropenia, 0–26, as the sum of the abstract's
// weights: burden of illness (no or mild symptoms 5, moderate 3), absence of hypotension 5, absence of COPD 4,
// solid tumor or no previous fungal infection in hematologic malignancy 4, outpatient status 3, absence of
// dehydration 3, age less than 60 years 2; and whether the score is 21 or more.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an age outside 0–130 or an unknown burden.
//
// ID:      MED-CRIT-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Klastersky J, Paesmans M, Rubenstein EB, et al. The Multinational Association for Supportive Care in Cancer risk index: a multinational scoring system for identifying low-risk febrile neutropenic cancer patients. J Clin Oncol. 2000;18(16):3038-3051. doi:10.1200/JCO.2000.18.16.3038. PMID 10944139. Abstract: "burden of illness indicating absence of symptoms or mild symptoms (weight, 5 ...) or moderate symptoms (weight, 3 ...); absence of hypotension (weight, 5 ...); absence of chronic obstructive pulmonary disease (weight, 4 ...); presence of solid tumor or absence of previous fungal infection in patients with hematologic malignancies (weight, 4 ...); outpatient status (weight, 3 ...); absence of dehydration (weight, 3 ...); and age less than 60 years (weight, 2 ...)"
// Tags:    critical, oncology, febrile neutropenia, mascc, score
// Units:   ageYears years; result points 0–26
// Range:   0 ≤ ageYears ≤ 130 (physical limits; febrile neutropenia in cancer, the study's population)
// Example: no or mild symptoms, every favourable factor, age 45 → 26, LowRisk; moderate symptoms, inpatient, age 70, all others favourable → 3+5+4+4+3 = 19 (HAND-DERIVED from the weights)
func MASCCRiskIndex(burden MASCCBurden, f MASCCFactors, ageYears float64) (MASCCResult, error) {
	if err := finite("ageYears", ageYears); err != nil {
		return MASCCResult{}, err
	}
	if err := within("ageYears", ageYears, 0, 130); err != nil {
		return MASCCResult{}, err
	}
	s := 0
	switch burden {
	case MASCCBurdenNoOrMildSymptoms:
		s = 5
	case MASCCBurdenModerateSymptoms:
		s = 3
	case MASCCBurdenNeither:
	default:
		return MASCCResult{}, fmt.Errorf("%w: burden = %d, want 0–2", ErrOutOfRange, burden)
	}
	if f.NoHypotension {
		s += 5
	}
	if f.NoCOPD {
		s += 4
	}
	if f.SolidTumorOrNoPriorFungal {
		s += 4
	}
	if f.Outpatient {
		s += 3
	}
	if f.NoDehydration {
		s += 3
	}
	if ageYears < 60 {
		s += 2
	}
	return MASCCResult{Score: s, LowRisk: s >= 21}, nil
}
