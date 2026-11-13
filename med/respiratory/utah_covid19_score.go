package respiratory

// UtahCOVID19Findings are the items of the Utah COVID-19 risk score (Webb 2022, Table 6), worded as in the
// table.
type UtahCOVID19Findings struct {
	AgeYears int  // "0.5 for every decade": 0–10 → 0.5, 11–20 → 1, …, 91–100 → 5, > 100 → 5.5
	Male     bool // "Male" 1
	// CommunityOfColor is the table's "Communities of color" item, 2 points.
	CommunityOfColor bool
	// High risk comorbidities, 2 points each.
	DiabetesMellitus          bool // "Diabetes Mellitus"
	SeverelyImmunocompromised bool // "Severely Immunocompromised"
	ObesityBMIOver30          bool // "Obesity (BMI>30)"
	// Other comorbidities, 1 point each.
	Hypertension             bool // "Hypertension"
	CoronaryArteryDisease    bool // "Coronary Artery Disease"
	CardiacArrhythmia        bool // "Cardiac Arrhythmia"
	CongestiveHeartFailure   bool // "Congestive Heart Failure"
	ChronicKidneyDisease     bool // "Chronic Kidney Disease"
	ChronicPulmonaryDisease  bool // "Chronic Pulmonary Disease"
	ChronicLiverDisease      bool // "Chronic Liver Disease"
	CerebrovascularDisease   bool // "Cerebrovascular Disease"
	ChronicNeurologicDisease bool // "Chronic Neurologic Disease"
	Dyspnea                  bool // Symptom: "Dyspnea" 1
}

// UtahCOVID19Score returns the Utah COVID-19 risk score of Webb 2022, Table 6, in points (multiples of 0.5,
// 0.5–24.5): age 0.5 per decade (0–10 → 0.5, 11–20 → 1, 21–30 → 1.5, …, 91–100 → 5, > 100 → 5.5), male 1,
// communities of color 2, diabetes mellitus 2, severely immunocompromised 2, obesity (BMI > 30) 2, each of nine
// other comorbidities 1, dyspnea 1. The paper reports that "A score of 6" optimized sensitivity and specificity
// for hospitalization; that comparison is the caller's.
// Errors: ErrOutOfRange for an age outside Range.
//
// ID:      MED-RESP-041
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Webb BJ, Levin NM, Grisel N, et al. Simple scoring tool to estimate risk of hospitalization and mortality in ambulatory and emergency department patients with COVID-19. PLoS One. 2022;17(3):e0261508. doi:10.1371/journal.pone.0261508. PMID 35239664, PMC8893609, Table 6
// Tags:    respiratory, covid-19, utah, hospitalization, score
// Units:   AgeYears completed years; the other items bool; result points (exact multiples of 0.5)
// Range:   AgeYears 0–120 (Table 6 starts at the 0–10 band)
// Example: 65-year-old man with diabetes and hypertension → 3.5 + 1 + 2 + 1 = 7.5 (hand-derived from Table 6)
func UtahCOVID19Score(f UtahCOVID19Findings) (float64, error) {
	if err := checkIntRange("AgeYears", f.AgeYears, 0, 120); err != nil {
		return 0, err
	}
	// Age band index: 1 for 0–10, ⌈age/10⌉ for 11–100, 11 above 100.
	decades := (f.AgeYears + 9) / 10
	if decades < 1 {
		decades = 1
	}
	if decades > 11 {
		decades = 11
	}
	halves := decades // the score in half points
	if f.Male {
		halves += 2
	}
	for _, two := range [4]bool{f.CommunityOfColor, f.DiabetesMellitus, f.SeverelyImmunocompromised, f.ObesityBMIOver30} {
		if two {
			halves += 4
		}
	}
	for _, one := range [10]bool{f.Hypertension, f.CoronaryArteryDisease, f.CardiacArrhythmia, f.CongestiveHeartFailure,
		f.ChronicKidneyDisease, f.ChronicPulmonaryDisease, f.ChronicLiverDisease, f.CerebrovascularDisease,
		f.ChronicNeurologicDisease, f.Dyspnea} {
		if one {
			halves += 2
		}
	}
	return float64(halves) / 2, nil
}
