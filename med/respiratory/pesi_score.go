package respiratory

import "fmt"

// PESIFindings are the eleven predictors of the Pulmonary Embolism Severity Index (Aujesky 2005, Table 2),
// worded as in the table. The measured items are scored against the table's thresholds.
type PESIFindings struct {
	AgeYears              int     // "Age, per yr": points = age in years
	Male                  bool    // "Male sex" +10
	Cancer                bool    // "Cancer" +30
	HeartFailure          bool    // "Heart failure" +10
	ChronicLungDisease    bool    // "Chronic lung disease" +10
	PulsePerMin           float64 // "Pulse ⩾ 110/min" +20
	SystolicMmHg          float64 // "Systolic blood pressure < 100 mm Hg" +30
	RespiratoryRatePerMin float64 // "Respiratory rate ⩾ 30/min" +20
	TemperatureC          float64 // "Temperature < 36°C" +20
	AlteredMentalStatus   bool    // "Altered mental status" +60
	SaO2Percent           float64 // "Arterial oxygen saturation < 90%" +20
}

// PESIScore returns the Pulmonary Embolism Severity Index total, in points: age in years, plus 10 for male sex,
// 30 cancer, 10 heart failure, 10 chronic lung disease, 20 pulse ≥ 110/min, 30 systolic < 100 mmHg, 20
// respiratory rate ≥ 30/min, 20 temperature < 36 °C, 60 altered mental status, 20 arterial oxygen saturation
// < 90%.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Aujesky D, Obrosky DS, Stone RA, et al. Derivation and validation of a prognostic model for pulmonary embolism. Am J Respir Crit Care Med. 2005;172(8):1041-1046. doi:10.1164/rccm.200506-862OC. PMID 16020800, PMC2718410, Table 2 (predictors and points)
// Tags:    respiratory, pulmonary-embolism, pesi, prognosis, score
// Units:   AgeYears completed years; PulsePerMin beats/min; SystolicMmHg mmHg; RespiratoryRatePerMin breaths/min; TemperatureC °C; SaO2Percent %; result points
// Range:   AgeYears 18–120 ("18 years or older"); PulsePerMin, RespiratoryRatePerMin ≥ 0; SystolicMmHg > 0; TemperatureC 20–45; SaO2Percent 0–100
// Example: 70-year-old man with cancer, pulse 115, otherwise none → 70 + 10 + 30 + 20 = 130 (hand-derived from Table 2)
func PESIScore(f PESIFindings) (int, error) {
	if err := checkIntRange("AgeYears", f.AgeYears, 18, 120); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("PulsePerMin", f.PulsePerMin); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("SystolicMmHg", f.SystolicMmHg); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("RespiratoryRatePerMin", f.RespiratoryRatePerMin); err != nil {
		return 0, err
	}
	if err := checkClosedRange("TemperatureC", f.TemperatureC, 20, 45); err != nil {
		return 0, err
	}
	if err := checkClosedRange("SaO2Percent", f.SaO2Percent, 0, 100); err != nil {
		return 0, err
	}
	s := f.AgeYears
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(f.Male, 10)
	add(f.Cancer, 30)
	add(f.HeartFailure, 10)
	add(f.ChronicLungDisease, 10)
	add(f.PulsePerMin >= 110, 20)
	add(f.SystolicMmHg < 100, 30)
	add(f.RespiratoryRatePerMin >= 30, 20)
	add(f.TemperatureC < 36, 20)
	add(f.AlteredMentalStatus, 60)
	add(f.SaO2Percent < 90, 20)
	return s, nil
}

// PESIClass is a PESI risk class; its value is the source's label, verbatim.
type PESIClass string

// The five classes of Aujesky 2005, Table 2 footnote / Results, verbatim: "⩽ 65 class I, very low risk; 66–85
// class II, low risk; 86–105 class III, intermediate risk; 106–125 class IV, high risk; > 125 class V, very
// high risk".
const (
	PESIClassI   PESIClass = "class I, very low risk"
	PESIClassII  PESIClass = "class II, low risk"
	PESIClassIII PESIClass = "class III, intermediate risk"
	PESIClassIV  PESIClass = "class IV, high risk"
	PESIClassV   PESIClass = "class V, very high risk"
)

// PESIRiskClass returns the PESI risk class for a total: ≤ 65 class I, 66–85 class II, 86–105 class III,
// 106–125 class IV, > 125 class V.
// Errors: ErrOutOfRange for a total below 18 (the least possible: age 18 and nothing else) or above 350 (age 120 plus all 230 item points).
//
// ID:      MED-RESP-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Aujesky D, Obrosky DS, Stone RA, et al. Derivation and validation of a prognostic model for pulmonary embolism. Am J Respir Crit Care Med. 2005;172(8):1041-1046. doi:10.1164/rccm.200506-862OC. PMID 16020800, PMC2718410, Table 2 (risk classes)
// Tags:    respiratory, pulmonary-embolism, pesi, category
// Units:   total points; result the source's class label
// Range:   total 18–350 (age 18–120 plus at most 230 item points)
// Example: 65 → "class I, very low risk"; 66 → "class II, low risk"; 126 → "class V, very high risk"
func PESIRiskClass(total int) (PESIClass, error) {
	if total < 18 || total > 350 {
		return "", fmt.Errorf("%w: total = %d, want 18 ≤ total ≤ 350", ErrOutOfRange, total)
	}
	switch {
	case total <= 65:
		return PESIClassI, nil
	case total <= 85:
		return PESIClassII, nil
	case total <= 105:
		return PESIClassIII, nil
	case total <= 125:
		return PESIClassIV, nil
	}
	return PESIClassV, nil
}
