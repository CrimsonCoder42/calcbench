package anthropometry

import "fmt"

// MUSTScore returns the 'Malnutrition Universal Screening Tool' ('MUST') total, 0–6: the step 1 BMI score
// (> 20: 0; 18.5–20: 1; < 18.5: 2), plus the step 2 weight-loss score for unplanned loss in the past 3–6
// months (< 5 %: 0; 5–10 %: 1; > 10 %: 2), plus the step 3 acute disease effect score (2 when the person is
// acutely ill and there has been or is likely to be no nutritional intake for > 5 days).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Malnutrition Advisory Group, BAPEN. 'Malnutrition Universal Screening Tool' ('MUST'). BAPEN; 2003 (revised 2011). https://www.bapen.org.uk/pdfs/must/must_full.pdf (read 2026-11-12), p. 3: Step 1 "BMI kg/m2 >20 (>30 Obese) =0, 18.5 - 20 =1, <18.5 =2"; Step 2 "<5 =0, 5 -10 =1, >10 =2"; Step 3 "Score 2"; Step 4 "Add scores together"
// Source:  Elia M (ed). The 'MUST' Report. Nutritional screening of adults: a multidisciplinary responsibility. BAPEN; 2003. ISBN 1-899467-70-X
// Tags:    anthropometry, nutrition, malnutrition, must, bmi, weight-loss, screening, score
// Units:   bmiKgM2 kg/m²; weightLossPercent percent of body weight lost unplanned in the past 3–6 months (a gain is passed as 0); result points 0–6
// Range:   bmiKgM2 > 0; 0 ≤ weightLossPercent ≤ 100 (adults, the tool's population)
// Example: BMI 19, loss 7 %, not acutely ill → 1 + 1 + 0 = 2 (HAND-DERIVED from p. 3)
func MUSTScore(bmiKgM2, weightLossPercent float64, acuteDiseaseEffect bool) (int, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkFinite("weightLossPercent", weightLossPercent); err != nil {
		return 0, err
	}
	if err := checkPositive("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if !(weightLossPercent >= 0 && weightLossPercent <= 100) {
		return 0, fmt.Errorf("%w: weightLossPercent = %v, want 0–100", ErrOutOfRange, weightLossPercent)
	}
	points := 0
	switch {
	case bmiKgM2 < 18.5:
		points += 2
	case bmiKgM2 <= 20:
		points++
	}
	switch {
	case weightLossPercent > 10:
		points += 2
	case weightLossPercent >= 5:
		points++
	}
	if acuteDiseaseEffect {
		points += 2
	}
	return points, nil
}

// MUSTRisk is a 'MUST' overall risk of malnutrition, with the tool's label verbatim.
type MUSTRisk string

const (
	MUSTLowRisk    MUSTRisk = "Low Risk"    // score 0
	MUSTMediumRisk MUSTRisk = "Medium Risk" // score 1
	MUSTHighRisk   MUSTRisk = "High Risk"   // score 2 or more
)

// MUSTRiskCategory returns the 'MUST' step 4 overall risk of malnutrition for a total: "Score 0 Low Risk
// Score 1 Medium Risk Score 2 or more High Risk".
// Errors: ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Malnutrition Advisory Group, BAPEN. 'Malnutrition Universal Screening Tool' ('MUST'). https://www.bapen.org.uk/pdfs/must/must_full.pdf (read 2026-11-12), p. 3, Step 4 "Overall risk of malnutrition"
// Tags:    anthropometry, nutrition, malnutrition, must, category
// Units:   score points 0–6 (MUSTScore); result category label
// Range:   0 ≤ score ≤ 6
// Example: 0 → "Low Risk"; 1 → "Medium Risk"; 2 → "High Risk"
func MUSTRiskCategory(score int) (MUSTRisk, error) {
	switch {
	case score < 0 || score > 6:
		return "", fmt.Errorf("%w: score = %d, want 0–6", ErrOutOfRange, score)
	case score >= 2:
		return MUSTHighRisk, nil
	case score == 1:
		return MUSTMediumRisk, nil
	default:
		return MUSTLowRisk, nil
	}
}
