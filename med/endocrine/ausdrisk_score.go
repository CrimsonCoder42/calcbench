package endocrine

import "fmt"

// AUSDRISKBackground is the ethnic background item of AUSDRISK, which sets both the background points and the
// waist-circumference thresholds (Chen 2010, Box 3).
type AUSDRISKBackground int

const (
	// AUSDRISKBackgroundOther: any background not listed below; 0 points; waist thresholds for "all other ethnicities".
	AUSDRISKBackgroundOther AUSDRISKBackground = iota
	// AUSDRISKBackgroundSouthernEuropeanOrPacificIslander: "Southern European ... or Pacific Islander background"; 2 points;
	// waist thresholds for "all other ethnicities".
	AUSDRISKBackgroundSouthernEuropeanOrPacificIslander
	// AUSDRISKBackgroundAsianOrAboriginalTorresStraitIslander: "Asian, Aboriginal and Torres Strait Islander"
	// background; 2 points; the lower waist thresholds the source gives for these groups.
	AUSDRISKBackgroundAsianOrAboriginalTorresStraitIslander
)

// AUSDRISKFindings holds the AUSDRISK items of Chen 2010, Box 3.
type AUSDRISKFindings struct {
	AgeYears                   float64            // completed age, years; 25–34: 0, 35–44: 2, 45–54: 4, 55–64: 6, ≥ 65: 8 points
	Male                       bool               // male sex: 3 points
	Background                 AUSDRISKBackground // 2 points unless AUSDRISKBackgroundOther
	ParentalDiabetes           bool               // parental history of diabetes: 3 points
	HighBloodGlucoseHistory    bool               // ever told of high blood glucose, including in pregnancy: 6 points
	AntihypertensiveMedication bool               // use of antihypertensive medication: 2 points
	CurrentSmoker              bool               // current smoker: 2 points
	PhysicallyInactive         bool               // < 150 min/week of walking (≥ 10 min continuous) or moderate/vigorous activity: 2 points
	WaistCm                    float64            // waist circumference, cm; category 2: 4 points, category 3: 7 points
}

// AUSDRISKScore returns the AUSDRISK total, 0–35, the points model of Chen 2010 Box 3: age band, male sex,
// background, parental diabetes, history of high blood glucose, antihypertensive medication, current smoking,
// physical inactivity and waist-circumference category (cut-offs by sex and background). The later public
// questionnaire's extra items (for example fruit and vegetable intake) are not part of this model and are not
// added.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown Background.
//
// ID:      MED-ENDO-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chen L, Magliano DJ, Balkau B, et al. AUSDRISK: an Australian Type 2 Diabetes Risk Assessment Tool based on demographic, lifestyle and simple anthropometric measures. Med J Aust. 2010;192(4):197-202. doi:10.5694/j.1326-5377.2010.tb03478.x. PMID 20170456. Box 3 (points; waist categories: Asian, Aboriginal and Torres Strait Islander men < 90 / 90–< 100 / ≥ 100 cm, women < 80 / 80–< 90 / ≥ 90 cm; all other men < 102 / 102–< 110 / ≥ 110 cm, women < 88 / 88–< 100 / ≥ 100 cm); maximum 35
// Tags:    endocrine, diabetes, type-2-diabetes, ausdrisk, risk-score, waist
// Units:   AgeYears years; WaistCm cm; result points 0–35
// Range:   AgeYears ≥ 25 (the source's population, adults aged 25 and over); WaistCm > 0
// Example: male 47, background Other, inactive, waist 105 cm, nothing else → 4 + 3 + 2 + 4 = 13 (HAND-DERIVED from Box 3)
func AUSDRISKScore(f AUSDRISKFindings) (int, error) {
	if err := checkFinite("AgeYears", f.AgeYears); err != nil {
		return 0, err
	}
	if err := checkFinite("WaistCm", f.WaistCm); err != nil {
		return 0, err
	}
	if !(f.AgeYears >= 25) {
		return 0, fmt.Errorf("%w: AgeYears = %v, want ≥ 25", ErrOutOfRange, f.AgeYears)
	}
	if err := checkPositive("WaistCm", f.WaistCm); err != nil {
		return 0, err
	}
	var lowerWaist bool
	points := 0
	switch f.Background {
	case AUSDRISKBackgroundOther:
	case AUSDRISKBackgroundSouthernEuropeanOrPacificIslander:
		points += 2
	case AUSDRISKBackgroundAsianOrAboriginalTorresStraitIslander:
		points += 2
		lowerWaist = true
	default:
		return 0, fmt.Errorf("%w: Background = %d, want 0–2", ErrOutOfRange, f.Background)
	}
	switch {
	case f.AgeYears >= 65:
		points += 8
	case f.AgeYears >= 55:
		points += 6
	case f.AgeYears >= 45:
		points += 4
	case f.AgeYears >= 35:
		points += 2
	}
	if f.Male {
		points += 3
	}
	if f.ParentalDiabetes {
		points += 3
	}
	if f.HighBloodGlucoseHistory {
		points += 6
	}
	if f.AntihypertensiveMedication {
		points += 2
	}
	if f.CurrentSmoker {
		points += 2
	}
	if f.PhysicallyInactive {
		points += 2
	}
	// Category 2 starts at cat2Cm, category 3 at cat3Cm (Box 3).
	var cat2Cm, cat3Cm float64
	switch {
	case lowerWaist && f.Male:
		cat2Cm, cat3Cm = 90, 100
	case lowerWaist:
		cat2Cm, cat3Cm = 80, 90
	case f.Male:
		cat2Cm, cat3Cm = 102, 110
	default:
		cat2Cm, cat3Cm = 88, 100
	}
	switch {
	case f.WaistCm >= cat3Cm:
		points += 7
	case f.WaistCm >= cat2Cm:
		points += 4
	}
	return points, nil
}
