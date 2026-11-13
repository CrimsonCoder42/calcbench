package cardio

import "fmt"

// TIMIUANSTEMIFactors are the seven yes/no predictors of the TIMI risk score for unstable angina / non-ST
// elevation MI, as the abstract of Antman 2000 lists them. The caller decides each one (the definition of the
// three coronary risk factors is in the paper's full text, which was not read).
type TIMIUANSTEMIFactors struct {
	ThreeOrMoreCADRiskFactors bool // "at least 3 risk factors for coronary artery disease"
	KnownCoronaryStenosis50   bool // "prior coronary stenosis of 50% or more"
	STDeviation               bool // "ST-segment deviation on electrocardiogram at presentation"
	TwoOrMoreAnginalEvents24h bool // "at least 2 anginal events in prior 24 hours"
	AspirinPrior7Days         bool // "use of aspirin in prior 7 days"
	ElevatedCardiacMarkers    bool // "elevated serum cardiac markers"
}

// TIMIScoreUANSTEMI returns the TIMI risk score for UA/NSTEMI, 0–7: one point for age ≥ 65 years and one for
// each factor present.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Antman EM, Cohen M, Bernink PJ, et al. The TIMI risk score for unstable angina/non-ST elevation MI: a method for prognostication and therapeutic decision making. JAMA. 2000;284(7):835-842. doi:10.1001/jama.284.7.835. PMID 10938172. Abstract: "age 65 years or older, at least 3 risk factors for coronary artery disease, prior coronary stenosis of 50% or more, ST-segment deviation on electrocardiogram at presentation, at least 2 anginal events in prior 24 hours, use of aspirin in prior 7 days, and elevated serum cardiac markers"; "assignment of value of 1 when a factor was present and 0 when it was absent, and summing"
// Tags:    cardio, acs, nstemi, unstable-angina, timi, risk-score
// Units:   ageYears years; result points 0–7
// Range:   ageYears 18–120 (adult trial population; the abstract states no hard limits, only the 65-year threshold is applied)
// Example: age 70, ST deviation, elevated markers → 3 (abstract: 14-day event rate 13.2% for a score of 3 in the test cohort)
func TIMIScoreUANSTEMI(ageYears float64, f TIMIUANSTEMIFactors) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	s := 0
	for _, present := range [...]bool{ageYears >= 65, f.ThreeOrMoreCADRiskFactors, f.KnownCoronaryStenosis50,
		f.STDeviation, f.TwoOrMoreAnginalEvents24h, f.AspirinPrior7Days, f.ElevatedCardiacMarkers} {
		if present {
			s++
		}
	}
	return s, nil
}

// TIMIUANSTEMIEventRatePercent returns the 14-day rate of the study's composite end point, in percent, that
// Antman 2000 reports for a TIMI UA/NSTEMI score in the test cohort: 0/1 4.7; 2 8.3; 3 13.2; 4 19.9; 5 26.2;
// 6/7 40.9.
// Errors: ErrOutOfRange when score is outside 0–7.
//
// ID:      MED-CARDIO-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Antman EM, Cohen M, Bernink PJ, et al. JAMA. 2000;284(7):835-842. doi:10.1001/jama.284.7.835. PMID 10938172. Abstract: "4.7% for a score of 0/1; 8.3% for 2; 13.2% for 3; 19.9% for 4; 26.2% for 5; and 40.9% for 6/7"
// Tags:    cardio, acs, nstemi, timi, event-rate
// Units:   score points 0–7; result percent
// Range:   0 ≤ score ≤ 7
// Example: score 3 → 13.2 (abstract)
func TIMIUANSTEMIEventRatePercent(score int) (float64, error) {
	if score < 0 || score > 7 {
		return 0, fmt.Errorf("%w: score = %d, want 0–7", ErrOutOfRange, score)
	}
	return timiUAEventRate[score], nil
}

// timiUAEventRate is the abstract's 14-day rate per score (0 and 1 share 4.7; 6 and 7 share 40.9).
var timiUAEventRate = [8]float64{4.7, 4.7, 8.3, 13.2, 19.9, 26.2, 40.9, 40.9}
