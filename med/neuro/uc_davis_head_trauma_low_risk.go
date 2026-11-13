package neuro

import "fmt"

// UCDavisHeadTraumaFindings are the five predictors of the UC Davis rule for children with blunt head trauma
// (Palchak 2003 abstract).
type UCDavisHeadTraumaFindings struct {
	AbnormalMentalStatus bool // "abnormal mental status"
	SkullFractureSigns   bool // "clinical signs of skull fracture"
	HistoryOfVomiting    bool // "a history of vomiting"
	ScalpHematoma        bool // "scalp hematoma (in children ≤2 years of age)": read only when ageYears ≤ 2
	Headache             bool // "headache"
}

// UCDavisHeadTraumaLowRisk reports whether a child with blunt head trauma has none of the five predictors of
// the UC Davis decision rule: abnormal mental status, clinical signs of skull fracture, a history of vomiting,
// scalp hematoma (counted only at ages 2 years or younger) and headache.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0 ≤ age < 18.
//
// ID:      MED-NEURO-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Palchak MJ, Holmes JF, Vance CW, et al. A decision rule for identifying children at low risk for brain injuries after blunt head trauma. Ann Emerg Med. 2003;42(4):492-506. doi:10.1067/s0196-0644(03)00425-6, PMID 14520320, abstract (absence of "abnormal mental status, clinical signs of skull fracture, a history of vomiting, scalp hematoma (in children ≤2 years of age), and headache")
// Tags:    neuro, head injury, pediatric, uc davis rule, decision rule
// Units:   ageYears years (decimal); findings yes/no; result true when no predictor is present
// Range:   0 ≤ ageYears < 18 (children; blunt head trauma, the source's population, is the caller's)
// Example: age 3 with scalp hematoma only → true (hematoma not counted after 2 years); age 2 with scalp hematoma → false (HAND-DERIVED from the abstract)
func UCDavisHeadTraumaLowRisk(ageYears float64, f UCDavisHeadTraumaFindings) (bool, error) {
	if err := checkChildAge(ageYears); err != nil {
		return false, err
	}
	hematoma := f.ScalpHematoma && ageYears <= 2
	return !(f.AbnormalMentalStatus || f.SkullFractureSigns || f.HistoryOfVomiting || hematoma || f.Headache), nil
}

// checkChildAge returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange outside 0 ≤ ageYears < 18.
func checkChildAge(ageYears float64) error {
	if err := checkRange("ageYears", ageYears, 0, 18); err != nil {
		return err
	}
	if ageYears == 18 {
		return fmt.Errorf("%w: ageYears = %v, want 0 ≤ ageYears < 18", ErrOutOfRange, ageYears)
	}
	return nil
}
