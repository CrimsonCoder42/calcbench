package cardio

import "fmt"

// READMITSFactors are the yes/no items of the AMI READMITS score (Nguyen 2018 Table 2; points in comments).
type READMITSFactors struct {
	CreatinineAbove2MgDL bool // "Serum creatinine >2 mg/dL": 6
	ElevatedBNP          bool // "≥50 pg/mL for BNP or ≥125 pg/mL for NT-proBNP": 8
	DiabetesMellitus     bool // 4
	Female               bool // "Not Male (ie, female sex)": 2
	NoTimelyPCI          bool // "No Intervention with Timely PCI" (within 24 hours): 1
	SystolicBelow100     bool // "Systolic blood pressure <100 mm Hg": 3
}

// READMITSScore returns the AMI READMITS score: one point per decade of age above 18 years plus the
// READMITSFactors points. The paper states "Age (per decade >18 years)" and does not say how a partial decade
// counts, so the caller passes the whole number of decades it has counted.
// Errors: ErrOutOfRange when decadesAbove18 is outside 0–9.
//
// ID:      MED-CARDIO-050
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nguyen OK, Makam AN, Clark C, Zhang S, Das SR, Halm EA. Predicting 30-day hospital readmissions in acute myocardial infarction: the AMI "READMITS" (Renal Function, Elevated Brain Natriuretic Peptide, Age, Diabetes Mellitus, Nonmale Sex, Intervention with Timely Percutaneous Coronary Intervention, and Low Systolic Blood Pressure) score. J Am Heart Assoc. 2018;7(8):e008882. doi:10.1161/JAHA.118.008882. PMID 29666065. PMC6015397, Table 2
// Tags:    cardio, myocardial-infarction, readmission, readmits, risk-score
// Units:   decadesAbove18 whole decades; result points 0–33
// Range:   0 ≤ decadesAbove18 ≤ 9 (adults hospitalised with AMI; mean age 64.8–70.5 years)
// Example: 4 decades, diabetes, female → 4 + 4 + 2 = 10 (HAND-DERIVED from Table 2)
func READMITSScore(decadesAbove18 int, f READMITSFactors) (int, error) {
	if decadesAbove18 < 0 || decadesAbove18 > 9 {
		return 0, fmt.Errorf("%w: decadesAbove18 = %d, want 0–9", ErrOutOfRange, decadesAbove18)
	}
	s := decadesAbove18
	for _, item := range [...]struct {
		present bool
		points  int
	}{
		{f.CreatinineAbove2MgDL, 6},
		{f.ElevatedBNP, 8},
		{f.DiabetesMellitus, 4},
		{f.Female, 2},
		{f.NoTimelyPCI, 1},
		{f.SystolicBelow100, 3},
	} {
		if item.present {
			s += item.points
		}
	}
	return s, nil
}

// READMITSRiskGroup is an AMI READMITS risk group, Nguyen 2018's own label.
type READMITSRiskGroup string

// The READMITS groups (Nguyen 2018 Table 4) with the observed 30-day readmission rate each reports.
const (
	READMITSExtremelyLow  READMITSRiskGroup = "Extremely low"  // ≤ 13: 2%
	READMITSLow           READMITSRiskGroup = "Low"            // 14–15: 8%
	READMITSModerate      READMITSRiskGroup = "Moderate"       // 16–17: 9%
	READMITSHigh          READMITSRiskGroup = "High"           // 18–19: 17%
	READMITSExtremelyHigh READMITSRiskGroup = "Extremely high" // ≥ 20: 34%
)

// READMITSRiskGroupOf returns the READMITS risk group for a score: ≤ 13 extremely low, 14–15 low, 16–17
// moderate, 18–19 high, ≥ 20 extremely high.
// Errors: ErrOutOfRange when score is outside 0–33.
//
// ID:      MED-CARDIO-051
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nguyen OK, et al. J Am Heart Assoc. 2018;7(8):e008882. doi:10.1161/JAHA.118.008882. PMID 29666065. PMC6015397, Table 4: extremely low (≤13) 2%, low (14–15) 8%, moderate (16–17) 9%, high (18–19) 17%, extremely high (≥20) 34%
// Tags:    cardio, myocardial-infarction, readmission, readmits, category
// Units:   score points 0–33; result group label
// Range:   0 ≤ score ≤ 33 (the sums READMITSScore allows)
// Example: 18 → "High"
func READMITSRiskGroupOf(score int) (READMITSRiskGroup, error) {
	switch {
	case score < 0 || score > 33:
		return "", fmt.Errorf("%w: score = %d, want 0–33", ErrOutOfRange, score)
	case score <= 13:
		return READMITSExtremelyLow, nil
	case score <= 15:
		return READMITSLow, nil
	case score <= 17:
		return READMITSModerate, nil
	case score <= 19:
		return READMITSHigh, nil
	}
	return READMITSExtremelyHigh, nil
}
