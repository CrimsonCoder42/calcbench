package respiratory

import "math"

// RESECT90Findings are the predictors of the RESECT-90 model (Taylor 2021), worded as in the paper's equation.
type RESECT90Findings struct {
	AgeYears          float64 // "age"
	Male              bool    // "male sex"
	PerformanceStatus int     // "performance status", 0–4
	DLCOPercentPred   float64 // "% predicted DLCO"
	BMIKgM2           float64 // "BMI [kg/m2]"
	CreatinineUmolL   float64 // "creatinine [μmol/l]"
	Anaemia           bool    // "anaemia" (haemoglobin < 120 g/l women, < 130 g/l men, WHO)
	Arrhythmia        bool    // "arrhythmia", pre-operative
	RightSided        bool    // "right", right-sided resection
	SegmentsResected  int     // "number of resected bronchopulmonary segments"
	Thoracotomy       bool    // "thoracotomy", open approach
	Malignant         bool    // "malignant", diagnosis confirmed or suspected
}

// RESECT90MortalityProbability returns the RESECT-90 predicted probability of death within 90 days of lung
// resection, as a fraction 0–1: Logit = −6.036 + 0.041·age + 0.493·male + 0.183·performance status −
// 0.029·%DLCO − 0.056·BMI + 0.005·creatinine + 0.242·anaemia + 0.608·arrhythmia + 0.379·right + 0.179·segments
// + 0.634·thoracotomy + 0.769·malignant; Probability = 1/(1 + exp(−Logit)).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-039
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Taylor M, Martin GP, Abah U, et al. Development and internal validation of a clinical prediction model for 90-day mortality after lung resection: the RESECT-90 score. Interact Cardiovasc Thorac Surg. 2021;33(6):921-927. doi:10.1093/icvts/ivab200. PMID 34324664, PMC8632783, Results (final model equation and probability conversion)
// Tags:    respiratory, thoracic-surgery, lung-resection, resect-90, mortality, logistic
// Units:   AgeYears years; PerformanceStatus ECOG grade 0–4; DLCOPercentPred % predicted; BMIKgM2 kg/m²; CreatinineUmolL µmol/L (mg/dL × 88.4); SegmentsResected count; the other items bool; result probability 0–1
// Range:   AgeYears 18–120; PerformanceStatus 0–4; DLCOPercentPred > 0; BMIKgM2 > 0; CreatinineUmolL > 0; SegmentsResected 0–19 (the bronchopulmonary segments of both lungs); adults undergoing lung resection
// Example: 70-year-old man, PS 1, DLCO 70 %, BMI 26, creatinine 90, right lobectomy of 3 segments by thoracotomy, malignant → Logit = −3.207, probability 0.0389 (hand-derived from the equation)
func RESECT90MortalityProbability(f RESECT90Findings) (float64, error) {
	if err := checkClosedRange("AgeYears", f.AgeYears, 18, 120); err != nil {
		return 0, err
	}
	if err := checkIntRange("PerformanceStatus", f.PerformanceStatus, 0, 4); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("DLCOPercentPred", f.DLCOPercentPred); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("BMIKgM2", f.BMIKgM2); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("CreatinineUmolL", f.CreatinineUmolL); err != nil {
		return 0, err
	}
	if err := checkIntRange("SegmentsResected", f.SegmentsResected, 0, 19); err != nil {
		return 0, err
	}
	l := -6.036 + 0.041*f.AgeYears + 0.183*float64(f.PerformanceStatus) - 0.029*f.DLCOPercentPred -
		0.056*f.BMIKgM2 + 0.005*f.CreatinineUmolL + 0.179*float64(f.SegmentsResected)
	for _, t := range [6]struct {
		on bool
		b  float64
	}{{f.Male, 0.493}, {f.Anaemia, 0.242}, {f.Arrhythmia, 0.608}, {f.RightSided, 0.379}, {f.Thoracotomy, 0.634}, {f.Malignant, 0.769}} {
		if t.on {
			l += t.b
		}
	}
	return 1 / (1 + math.Exp(-l)), nil
}
