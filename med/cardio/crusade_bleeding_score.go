package cardio

import "fmt"

// CRUSADEFactors are the yes/no items of the CRUSADE bleeding score (Subherwal 2009 Table 4; points in comments).
type CRUSADEFactors struct {
	Female               bool // "Female": 8
	SignsOfCHF           bool // "Signs of CHF at presentation": 7
	PriorVascularDisease bool // "Prior vascular disease": 6
	DiabetesMellitus     bool // "Diabetes mellitus": 6
}

// CRUSADEBleedingScore returns the CRUSADE bleeding score, 1–96 points, from Subherwal 2009 Table 4: baseline
// hematocrit < 31 → 9, 31–33.9 → 7, 34–36.9 → 3, 37–39.9 → 2, ≥ 40 → 0; creatinine clearance (Cockcroft–Gault)
// ≤ 15 → 39, > 15–30 → 35, > 30–60 → 28, > 60–90 → 17, > 90–120 → 7, > 120 → 0; heart rate ≤ 70 → 0, 71–80 → 1,
// 81–90 → 3, 91–100 → 6, 101–110 → 8, 111–120 → 10, ≥ 121 → 11; systolic pressure ≤ 90 → 10, 91–100 → 8,
// 101–120 → 5, 121–180 → 1, 181–200 → 3, ≥ 201 → 5; plus the CRUSADEFactors. Heart rate and systolic pressure
// between the printed integers fall in the lower band (a heart rate of 70.5 scores 0).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-048
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Subherwal S, Bach RG, Chen AY, et al. Baseline risk of major bleeding in non-ST-segment-elevation myocardial infarction: the CRUSADE Bleeding Score. Circulation. 2009;119(14):1873-1882. doi:10.1161/CIRCULATIONAHA.108.828541. PMID 19332461. PMC3767035, Table 4
// Tags:    cardio, nstemi, bleeding, crusade, risk-score
// Units:   hematocritPercent %; creatinineClearanceMLMin mL/min (Cockcroft–Gault); heartRateBpm beats/min; systolicMmHg mmHg; result points
// Range:   hematocritPercent 10–70; creatinineClearanceMLMin 0–300; heartRateBpm 20–250; systolicMmHg 40–300 (NSTEMI admissions in the CRUSADE registry; no hard limits stated)
// Example: Hct 32, CrCl 45, HR 95, SBP 110, female, diabetes → 7 + 28 + 6 + 5 + 8 + 6 = 60 (HAND-DERIVED from Table 4)
func CRUSADEBleedingScore(hematocritPercent, creatinineClearanceMLMin, heartRateBpm, systolicMmHg float64, f CRUSADEFactors) (int, error) {
	if err := checkBounded([]boundedInput{
		{"hematocritPercent", hematocritPercent, 10, 70},
		{"creatinineClearanceMLMin", creatinineClearanceMLMin, 0, 300},
		{"heartRateBpm", heartRateBpm, 20, 250},
		{"systolicMmHg", systolicMmHg, 40, 300},
	}); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case hematocritPercent < 31:
		s += 9
	case hematocritPercent < 34:
		s += 7
	case hematocritPercent < 37:
		s += 3
	case hematocritPercent < 40:
		s += 2
	}
	switch {
	case creatinineClearanceMLMin <= 15:
		s += 39
	case creatinineClearanceMLMin <= 30:
		s += 35
	case creatinineClearanceMLMin <= 60:
		s += 28
	case creatinineClearanceMLMin <= 90:
		s += 17
	case creatinineClearanceMLMin <= 120:
		s += 7
	}
	switch {
	case heartRateBpm <= 70:
	case heartRateBpm <= 80:
		s += 1
	case heartRateBpm <= 90:
		s += 3
	case heartRateBpm <= 100:
		s += 6
	case heartRateBpm <= 110:
		s += 8
	case heartRateBpm <= 120:
		s += 10
	default:
		s += 11
	}
	switch {
	case systolicMmHg <= 90:
		s += 10
	case systolicMmHg <= 100:
		s += 8
	case systolicMmHg <= 120:
		s += 5
	case systolicMmHg <= 180:
		s += 1
	case systolicMmHg <= 200:
		s += 3
	default:
		s += 5
	}
	if f.Female {
		s += 8
	}
	if f.SignsOfCHF {
		s += 7
	}
	if f.PriorVascularDisease {
		s += 6
	}
	if f.DiabetesMellitus {
		s += 6
	}
	return s, nil
}

// CRUSADEBleedingGroup is a CRUSADE bleeding-risk group, Subherwal 2009's own label.
type CRUSADEBleedingGroup string

// The CRUSADE groups (Subherwal 2009, Results, risk quintiles), with the in-hospital major bleeding rate each reports.
const (
	CRUSADEVeryLow  CRUSADEBleedingGroup = "Very Low"  // ≤ 20: 3.1%
	CRUSADELow      CRUSADEBleedingGroup = "Low"       // 21–30: 5.5%
	CRUSADEModerate CRUSADEBleedingGroup = "Moderate"  // 31–40: 8.6%
	CRUSADEHigh     CRUSADEBleedingGroup = "High"      // 41–50: 11.9%
	CRUSADEVeryHigh CRUSADEBleedingGroup = "Very High" // > 50: 19.5%
)

// CRUSADEBleedingGroupOf returns the CRUSADE bleeding-risk group for a score: ≤ 20 very low, 21–30 low,
// 31–40 moderate, 41–50 high, > 50 very high.
// Errors: ErrOutOfRange when score is outside 1–96.
//
// ID:      MED-CARDIO-049
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Subherwal S, et al. Circulation. 2009;119(14):1873-1882. doi:10.1161/CIRCULATIONAHA.108.828541. PMID 19332461. PMC3767035, Results: quintiles very low (≤20) 3.1%, low (21–30) 5.5%, moderate (31–40) 8.6%, high (41–50) 11.9%, very high (>50) 19.5%
// Tags:    cardio, nstemi, bleeding, crusade, category
// Units:   score points 1–96; result group label
// Range:   1 ≤ score ≤ 96 (the sums Table 4 allows)
// Example: 60 → "Very High"
func CRUSADEBleedingGroupOf(score int) (CRUSADEBleedingGroup, error) {
	switch {
	case score < 1 || score > 96:
		return "", fmt.Errorf("%w: score = %d, want 1–96", ErrOutOfRange, score)
	case score <= 20:
		return CRUSADEVeryLow, nil
	case score <= 30:
		return CRUSADELow, nil
	case score <= 40:
		return CRUSADEModerate, nil
	case score <= 50:
		return CRUSADEHigh, nil
	}
	return CRUSADEVeryHigh, nil
}
