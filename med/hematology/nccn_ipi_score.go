package hematology

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when a float input is NaN or ±Inf.
	ErrNotFinite = errors.New("hematology: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside its physical range or the range the source allows.
	ErrOutOfRange = errors.New("hematology: input outside the range the source validated")
)

// checkFloat returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange outside lo–hi, wrapped with name.
func checkFloat(name string, v, lo, hi float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v–%v", ErrOutOfRange, name, v, lo, hi)
	}
	return nil
}

// pts returns w when present is true and 0 otherwise.
func pts(present bool, w int) int {
	if present {
		return w
	}
	return 0
}

// NCCNIPIScore returns the NCCN-IPI score, 0–8 (Table 3): age > 40 to ≤ 60 years 1, > 60 to ≤ 75 2, > 75 3;
// LDH normalized to the upper limit of normal > 1 to ≤ 3 1, > 3 2; Ann Arbor stage III–IV 1; extranodal
// disease in bone marrow, CNS, liver/GI tract or lung 1; performance status (ECOG) ≥ 2 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEMAT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zhou Z, Sehn LH, Rademaker AW, et al. An enhanced International Prognostic Index (NCCN-IPI) for patients with diffuse large B-cell lymphoma treated in the rituximab era. Blood. 2014;123(6):837-842. doi:10.1182/blood-2013-09-524108, PMID 24264230, Table 3
// Tags:    hematology, lymphoma, dlbcl, nccn-ipi, prognostic score
// Units:   ageYears years; ldhRatioToULN serum LDH ÷ upper limit of normal (dimensionless); ecog 0–4; result points 0–8
// Range:   0 ≤ ageYears ≤ 130; 0 ≤ ldhRatioToULN ≤ 100; 0 ≤ ecog ≤ 4 (physical limits; Table 3 bands are open-ended)
// Example: age 65 (2), LDH ratio 2 (1), stage IV (1), marrow involved (1), ECOG 1 (0) → 5 (Table 3)
func NCCNIPIScore(ageYears, ldhRatioToULN float64, stageIIIorIV, extranodalMajorOrgan bool, ecog int) (int, error) {
	if err := checkFloat("ageYears", ageYears, 0, 130); err != nil {
		return 0, err
	}
	if err := checkFloat("ldhRatioToULN", ldhRatioToULN, 0, 100); err != nil {
		return 0, err
	}
	if ecog < 0 || ecog > 4 {
		return 0, fmt.Errorf("%w: ecog = %d, want 0–4", ErrOutOfRange, ecog)
	}
	s := 0
	switch {
	case ageYears > 75:
		s += 3
	case ageYears > 60:
		s += 2
	case ageYears > 40:
		s++
	}
	switch {
	case ldhRatioToULN > 3:
		s += 2
	case ldhRatioToULN > 1:
		s++
	}
	return s + pts(stageIIIorIV, 1) + pts(extranodalMajorOrgan, 1) + pts(ecog >= 2, 1), nil
}
