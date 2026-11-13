package obstetric

import (
	"fmt"
	"math"
)

// GWGCategory is a prepregnancy BMI category of IOM 2009 Table S-1, with the source's label.
type GWGCategory string

// The four prepregnancy BMI categories of IOM 2009 Table S-1, labels verbatim.
const (
	GWGUnderweight  GWGCategory = "Underweight"   // < 18.5 kg/m²
	GWGNormalWeight GWGCategory = "Normal weight" // 18.5–24.9 kg/m²
	GWGOverweight   GWGCategory = "Overweight"    // 25.0–29.9 kg/m²
	GWGObese        GWGCategory = "Obese"         // ≥ 30.0 kg/m²
)

// GWGRecommendation is one row of IOM 2009 Table S-1: the total gain range at term and the rate of gain in the
// second and third trimesters, mean and range, all in pounds.
type GWGRecommendation struct {
	Category          GWGCategory
	TotalLowLb        float64 // total weight gain at term, low end (lb)
	TotalHighLb       float64 // total weight gain at term, high end (lb)
	RateMeanLbPerWeek float64 // 2nd and 3rd trimester rate, mean (lb/week)
	RateLowLbPerWeek  float64 // 2nd and 3rd trimester rate, range low (lb/week)
	RateHighLbPerWeek float64 // 2nd and 3rd trimester rate, range high (lb/week)
}

// iomTableS1 holds Table S-1 in category order.
var iomTableS1 = [4]GWGRecommendation{
	{GWGUnderweight, 28, 40, 1, 1, 1.3},
	{GWGNormalWeight, 25, 35, 1, 0.8, 1},
	{GWGOverweight, 15, 25, 0.6, 0.5, 0.7},
	{GWGObese, 11, 20, 0.5, 0.4, 0.6},
}

// gwgIndex returns the Table S-1 row index for a prepregnancy BMI.
func gwgIndex(prePregnancyBMIKgM2 float64) (int, error) {
	if math.IsNaN(prePregnancyBMIKgM2) || math.IsInf(prePregnancyBMIKgM2, 0) || !(prePregnancyBMIKgM2 > 0) {
		return 0, fmt.Errorf("%w: prePregnancyBMIKgM2 = %v, want finite and > 0", ErrOutOfRange, prePregnancyBMIKgM2)
	}
	switch {
	case prePregnancyBMIKgM2 < 18.5:
		return 0, nil
	case prePregnancyBMIKgM2 < 25:
		return 1, nil
	case prePregnancyBMIKgM2 < 30:
		return 2, nil
	}
	return 3, nil
}

// GestationalWeightGainIOM returns the IOM 2009 Table S-1 row for a singleton pregnancy and a prepregnancy BMI:
// the category label, the total weight gain range at term, and the second- and third-trimester rate of gain
// (mean and range). The table's ranges 18.5–24.9 and 25.0–29.9 are read as 18.5 to < 25 and 25 to < 30. The
// table's rates assume a 0.5–2 kg (1.1–4.4 lb) gain in the first trimester (its footnote).
// Errors: ErrOutOfRange when prePregnancyBMIKgM2 is NaN, ±Inf or ≤ 0.
//
// ID:      MED-OB-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Institute of Medicine and National Research Council. Weight Gain During Pregnancy: Reexamining the Guidelines. Washington, DC: The National Academies Press; 2009. doi:10.17226/12584. Summary, Table S-1 "New Recommendations for Total and Rate of Weight Gain During Pregnancy, by Prepregnancy BMI" (read 2026-11-13)
// Tags:    obstetric, pregnancy, weight-gain, gestational-weight-gain, bmi, iom
// Units:   prePregnancyBMIKgM2 kg/m²; result totals lb, rates lb/week
// Range:   prePregnancyBMIKgM2 > 0; singleton pregnancy (twins: MED-OB-006)
// Example: BMI 22 → "Normal weight", 25–35 lb, 1 (0.8–1) lb/week; BMI 31 → "Obese", 11–20 lb, 0.5 (0.4–0.6) lb/week (Table S-1)
func GestationalWeightGainIOM(prePregnancyBMIKgM2 float64) (GWGRecommendation, error) {
	i, err := gwgIndex(prePregnancyBMIKgM2)
	if err != nil {
		return GWGRecommendation{}, err
	}
	return iomTableS1[i], nil
}

// GestationalWeightGainTwinsIOM returns the IOM 2009 provisional total weight gain range at term for a twin
// pregnancy: "Normal weight" 37–54 lb, "Overweight" 31–50 lb, "Obese" 25–42 lb. The source gives no twin range
// for "Underweight", so a BMI below 18.5 returns ErrOutOfRange.
// Errors: ErrOutOfRange when prePregnancyBMIKgM2 is NaN, ±Inf, ≤ 0 or < 18.5.
//
// ID:      MED-OB-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Institute of Medicine and National Research Council. Weight Gain During Pregnancy: Reexamining the Guidelines. Washington, DC: The National Academies Press; 2009. doi:10.17226/12584. Summary, provisional guidelines for twins below Table S-1 (read 2026-11-13)
// Tags:    obstetric, pregnancy, weight-gain, gestational-weight-gain, twins, bmi, iom
// Units:   prePregnancyBMIKgM2 kg/m²; lowLb, highLb lb
// Range:   prePregnancyBMIKgM2 ≥ 18.5 (no twin range for "Underweight")
// Example: BMI 22 → "Normal weight", 37–54 lb; BMI 27 → "Overweight", 31–50 lb (IOM 2009 Summary)
func GestationalWeightGainTwinsIOM(prePregnancyBMIKgM2 float64) (category GWGCategory, lowLb, highLb float64, err error) {
	i, err := gwgIndex(prePregnancyBMIKgM2)
	if err != nil {
		return "", 0, 0, err
	}
	switch i {
	case 1:
		return GWGNormalWeight, 37, 54, nil
	case 2:
		return GWGOverweight, 31, 50, nil
	case 3:
		return GWGObese, 25, 42, nil
	}
	return "", 0, 0, fmt.Errorf("%w: prePregnancyBMIKgM2 = %v: the source gives no twin range for %q", ErrOutOfRange, prePregnancyBMIKgM2, GWGUnderweight)
}
