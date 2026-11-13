package hepatic

import "math"

// MELD returns the MELD score, dimensionless and not rounded: 9.57 × ln(creatinine) + 3.78 × ln(bilirubin) +
// 11.20 × ln(INR) + 6.43. Values below 1 are set to 1; creatinine is capped at 4.0 mg/dL, and set to 4.0 when
// dialysisTwiceInPastWeek. The quoted form states no rounding, so none is applied.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEPAT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kim WR, Mannalithara A, Heimbach JK, et al. MELD 3.0: the Model for End-Stage Liver Disease updated for the modern era. Gastroenterology. 2021;161(6):1887-1895.e4. doi:10.1053/j.gastro.2021.08.050. PMID 34481845. Methods (full text PMC8608337), the MELD formula as quoted: "MELD = 9.57 × loge(creatinine) + 3.78 × loge(bilirubin) + 11.20 × loge(INR) + 6.43"
// Source:  Kamath PS, Wiesner RH, Malinchoc M, et al. A model to predict survival in patients with end-stage liver disease. Hepatology. 2001;33(2):464-470. doi:10.1053/jhep.2001.22172. PMID 11172350 (the original model; its own coefficient form not read)
// Tags:    hepatic, liver, meld, bilirubin, inr, creatinine
// Units:   bilirubinMgDL mg/dL; inrRatio dimensionless; creatinineMgDL mg/dL; dialysisTwiceInPastWeek true when dialysed at least twice in the past week; result dimensionless score
// Range:   every numeric input > 0 (limits not stated by the source; the bounds above are its rules, applied as clamps)
// Example: bilirubin 1, INR 1, creatinine 1 → 6.43 (hand-derived: every log is 0)
func MELD(bilirubinMgDL, inrRatio, creatinineMgDL float64, dialysisTwiceInPastWeek bool) (float64, error) {
	if err := checkMELDInputs(bilirubinMgDL, inrRatio, creatinineMgDL); err != nil {
		return 0, err
	}
	return meld(bilirubinMgDL, inrRatio, creatinineMgDL, dialysisTwiceInPastWeek), nil
}

// checkMELDInputs validates the three laboratory inputs of MED-HEPAT-002 and -003: finiteness first, then ranges.
func checkMELDInputs(bilirubinMgDL, inrRatio, creatinineMgDL float64) error {
	if err := checkFinite("bilirubinMgDL", bilirubinMgDL); err != nil {
		return err
	}
	if err := checkFinite("inrRatio", inrRatio); err != nil {
		return err
	}
	if err := checkFinite("creatinineMgDL", creatinineMgDL); err != nil {
		return err
	}
	if err := checkPositive("bilirubinMgDL", bilirubinMgDL); err != nil {
		return err
	}
	if err := checkPositive("inrRatio", inrRatio); err != nil {
		return err
	}
	return checkPositive("creatinineMgDL", creatinineMgDL)
}

// meld is the quoted MELD formula, with its bounds, on validated inputs.
func meld(bilirubinMgDL, inrRatio, creatinineMgDL float64, dialysisTwiceInPastWeek bool) float64 {
	bilirubin := math.Max(bilirubinMgDL, 1)
	inr := math.Max(inrRatio, 1)
	creatinine := math.Min(math.Max(creatinineMgDL, 1), 4.0)
	if dialysisTwiceInPastWeek {
		creatinine = 4.0
	}
	return 9.57*math.Log(creatinine) + 3.78*math.Log(bilirubin) + 11.20*math.Log(inr) + 6.43
}
