package endocrine

import "fmt"

// DKAMortalityFindings holds the six predictors of the Efstathiou 2002 diabetic ketoacidosis mortality
// prediction model, each at the time point the source names.
type DKAMortalityFindings struct {
	SevereCoexistingDiseases  bool    // "severe coexisting diseases (SCD)" at presentation: 6 points
	PHAtPresentation          float64 // arterial pH at presentation; pH < 7.0: 4 points
	RegularInsulinIUAt12h     float64 // regular insulin required over the first 12 h, IU; > 50 IU: 4 points
	GlucoseMmolLAt12h         float64 // serum glucose after 12 h, mmol/L; > 16.7 mmol/L: 4 points
	DepressedMentalStateAt24h bool    // depressed mental state after 24 h: 4 points
	FeverAt24h                bool    // fever after 24 h: 3 points
}

// DKAMortalityScore returns the points, 0–25, of the Efstathiou 2002 mortality prediction model for diabetic
// ketoacidosis: 6 (SCD at presentation) + 4 (pH < 7.0 at presentation) + 4 (regular insulin required > 50 IU
// after 12 h) + 4 (serum glucose > 16.7 mmol/L after 12 h) + 4 (depressed mental state after 24 h) + 3 (fever
// after 24 h). No band is returned (the source reports risk for 0–14 and 19–25 points only).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Efstathiou SP, Tsiakou AG, Tsioulos DI, et al. A mortality prediction model in diabetic ketoacidosis. Clin Endocrinol (Oxf). 2002;57(5):595-601. doi:10.1046/j.1365-2265.2002.01636.x. PMID 12390332. Abstract: "number of points = 6 (SCD at presentation) + 4 (pH < 7.0 at presentation) + 4 (regular insulin required > 50 IU after 12 h) + 4 (serum glucose > 16.7 mmo/l after 12 h) + 4 (depressed mental state after 24 h) + 3 (fever after 24 h)"
// Tags:    endocrine, diabetes, dka, ketoacidosis, mortality, score
// Units:   PHAtPresentation pH units; RegularInsulinIUAt12h IU; GlucoseMmolLAt12h mmol/L; result points 0–25
// Range:   0 < pH ≤ 14; insulin ≥ 0 IU; glucose > 0 mmol/L (population: 154 adults with DKA, mean age 58 ± 12 years; input limits not stated in the source)
// Example: SCD, pH 6.9, 60 IU, 18 mmol/L, depressed mental state, fever → 25; none of them (pH 7.2, 30 IU, 12 mmol/L) → 0 (HAND-DERIVED from the abstract)
func DKAMortalityScore(f DKAMortalityFindings) (int, error) {
	if err := checkFinite("PHAtPresentation", f.PHAtPresentation); err != nil {
		return 0, err
	}
	if err := checkFinite("RegularInsulinIUAt12h", f.RegularInsulinIUAt12h); err != nil {
		return 0, err
	}
	if err := checkFinite("GlucoseMmolLAt12h", f.GlucoseMmolLAt12h); err != nil {
		return 0, err
	}
	if !(f.PHAtPresentation > 0 && f.PHAtPresentation <= 14) {
		return 0, fmt.Errorf("%w: PHAtPresentation = %v, want 0 < pH ≤ 14", ErrOutOfRange, f.PHAtPresentation)
	}
	if !(f.RegularInsulinIUAt12h >= 0) {
		return 0, fmt.Errorf("%w: RegularInsulinIUAt12h = %v, want ≥ 0", ErrOutOfRange, f.RegularInsulinIUAt12h)
	}
	if err := checkPositive("GlucoseMmolLAt12h", f.GlucoseMmolLAt12h); err != nil {
		return 0, err
	}
	points := 0
	if f.SevereCoexistingDiseases {
		points += 6
	}
	if f.PHAtPresentation < 7.0 {
		points += 4
	}
	if f.RegularInsulinIUAt12h > 50 {
		points += 4
	}
	if f.GlucoseMmolLAt12h > 16.7 {
		points += 4
	}
	if f.DepressedMentalStateAt24h {
		points += 4
	}
	if f.FeverAt24h {
		points += 3
	}
	return points, nil
}
