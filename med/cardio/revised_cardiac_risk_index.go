package cardio

// RCRIFactors are the six yes/no predictors of the Revised Cardiac Risk Index (Lee 1999 abstract, verbatim in
// the comments).
type RCRIFactors struct {
	HighRiskSurgery        bool // "high-risk type of surgery"
	IschemicHeartDisease   bool // "history of ischemic heart disease"
	CongestiveHeartFailure bool // "history of congestive heart failure"
	CerebrovascularDisease bool // "history of cerebrovascular disease"
	InsulinTreatment       bool // "preoperative treatment with insulin"
	CreatinineAbove2MgDL   bool // "preoperative serum creatinine >2.0 mg/dL"
}

// RevisedCardiacRiskIndex returns the Revised Cardiac Risk Index, 0–6: one point per factor present.
// The abstract's major cardiac complication rates for 0, 1, 2, ≥3 factors: 0.5, 1.3, 4, 9% (derivation) and
// 0.4, 0.9, 7, 11% (validation). The I–IV class labels were not read and are not returned.
// Errors: none.
//
// ID:      MED-CARDIO-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lee TH, Marcantonio ER, Mangione CM, et al. Derivation and prospective validation of a simple index for prediction of cardiac risk of major noncardiac surgery. Circulation. 1999;100(10):1043-1049. doi:10.1161/01.cir.100.10.1043. PMID 10477528. Abstract: "high-risk type of surgery, history of ischemic heart disease, history of congestive heart failure, history of cerebrovascular disease, preoperative treatment with insulin, and preoperative serum creatinine >2.0 mg/dL"
// Tags:    cardio, perioperative, noncardiac-surgery, rcri, lee, risk-score
// Units:   result points 0–6
// Range:   persons aged ≥ 50 having elective major noncardiac surgery (4,315 persons, Boston); age is not an input
// Example: insulin-treated diabetes and creatinine 2.4 mg/dL → 2 (HAND-DERIVED)
func RevisedCardiacRiskIndex(f RCRIFactors) int {
	s := 0
	for _, present := range [...]bool{f.HighRiskSurgery, f.IschemicHeartDisease, f.CongestiveHeartFailure,
		f.CerebrovascularDisease, f.InsulinTreatment, f.CreatinineAbove2MgDL} {
		if present {
			s++
		}
	}
	return s
}
