package cardio

// ATRIAStrokeFactors are the ATRIA stroke items besides age and sex (Singer 2013, Table 3); each scores 1.
type ATRIAStrokeFactors struct {
	PriorStroke       bool // selects the "with prior stroke" age column
	Diabetes          bool
	HeartFailure      bool // "CHF"
	Hypertension      bool
	Proteinuria       bool
	EGFRBelow45OrESRD bool // "eGFR<45 or ESRD"
}

// ATRIAStrokeScore returns the ATRIA stroke risk score for atrial fibrillation, 0–12 without a prior stroke
// and 7–15 with one. Age points (no prior stroke / prior stroke): ≥ 85 6/9; 75–84 5/7; 65–74 3/7; < 65 0/8.
// Female, diabetes, CHF, hypertension, proteinuria and eGFR < 45 or ESRD score 1 each.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Singer DE, Chang Y, Borowsky LH, et al. A new risk scheme to predict ischemic stroke and other thromboembolism in atrial fibrillation: the ATRIA study stroke risk score. J Am Heart Assoc. 2013;2(3):e000250. doi:10.1161/JAHA.113.000250. PMID 23782923. PMC3698792, Table 3 "ATRIA Stroke Risk Model Point Scoring System"; "Possible point scores range from 0 to 12 for those without a prior stroke and from 7 to 15 for those with a prior stroke"
// Tags:    cardio, atrial-fibrillation, stroke, thromboembolism, atria, risk-score
// Units:   ageYears years (completed years: 74.9 counts as 65–74); result points
// Range:   ageYears 18–120 (all aged ≥ 18; 10,927 adults with nonvalvular AF off warfarin)
// Example: woman aged 80, hypertension, no prior stroke → 5 + 1 + 1 = 7 (HAND-DERIVED from Table 3)
func ATRIAStrokeScore(ageYears float64, female bool, f ATRIAStrokeFactors) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	col := 0
	if f.PriorStroke {
		col = 1
	}
	var band int
	switch {
	case ageYears >= 85:
		band = 0
	case ageYears >= 75:
		band = 1
	case ageYears >= 65:
		band = 2
	default:
		band = 3
	}
	s := atriaStrokeAgePoints[band][col]
	for _, present := range [...]bool{female, f.Diabetes, f.HeartFailure, f.Hypertension, f.Proteinuria, f.EGFRBelow45OrESRD} {
		if present {
			s++
		}
	}
	return s, nil
}

// atriaStrokeAgePoints is Table 3's age rows (≥85, 75–84, 65–74, <65) × (no prior stroke, prior stroke).
var atriaStrokeAgePoints = [4][2]int{{6, 9}, {5, 7}, {3, 7}, {0, 8}}

// ATRIAStrokeCategory is a risk group of Singer 2013; its value is the paper's wording.
type ATRIAStrokeCategory string

// Singer 2013: "low (0 to 5 points), moderate (6 points), and high (7 to 15 points)".
const (
	ATRIAStrokeLow      ATRIAStrokeCategory = "low"
	ATRIAStrokeModerate ATRIAStrokeCategory = "moderate"
	ATRIAStrokeHigh     ATRIAStrokeCategory = "high"
)

// ATRIAStrokeGroup returns the Singer 2013 group for an ATRIA stroke score: 0–5 "low", 6 "moderate",
// 7–15 "high".
// Errors: ErrOutOfRange when score is outside 0–15.
//
// ID:      MED-CARDIO-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Singer DE, et al. J Am Heart Assoc. 2013;2(3):e000250. doi:10.1161/JAHA.113.000250. PMID 23782923. PMC3698792: "collapsed into low (0 to 5 points), moderate (6 points), and high (7 to 15 points) risk categories"
// Tags:    cardio, atrial-fibrillation, stroke, atria, category
// Units:   score points 0–15; result the paper's group word
// Range:   0 ≤ score ≤ 15
// Example: 5 → "low"; 6 → "moderate"; 7 → "high"
func ATRIAStrokeGroup(score int) (ATRIAStrokeCategory, error) {
	switch {
	case score < 0 || score > 15:
		return "", errScore(score, 0, 15)
	case score <= 5:
		return ATRIAStrokeLow, nil
	case score == 6:
		return ATRIAStrokeModerate, nil
	default:
		return ATRIAStrokeHigh, nil
	}
}
