package respiratory

// DeltaPFindings are the six DELTA-P items, each true when present, worded as in the source's abstract.
type DeltaPFindings struct {
	AgeAtOnsetAtLeast50   bool // "age at onset ≥ 50 years"
	SmokingAtDiagnosis    bool // "smoking at diagnosis"
	WeightLossAtLeast5Pct bool // "weight loss ≥ 5%"
	BulbarInvolvement     bool // "bulbar involvement"
	ErectileDysfunction   bool // "erectile dysfunction"
	KarnofskyBelow70      bool // "Karnofsky performance status lower than 70"
}

// DeltaPScore returns the Dutch-English LEMS Tumor Association Prediction (DELTA-P) score (0–6): one point for each
// item present.
// Every input is valid, so there is no error result.
//
// ID:      MED-RESP-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Titulaer MJ, Maddison P, Sont JK, et al. Clinical Dutch-English Lambert-Eaton Myasthenic Syndrome (LEMS) Tumor Association Prediction Score Accurately Predicts Small-Cell Lung Cancer in the LEMS. J Clin Oncol. 2011;29(7):902-908. doi:10.1200/JCO.2010.32.0440. PMID 21245427, abstract: the six items and "A DELTA-P score of 0 or 1 ... scores of 4, 5, and 6"
// Tags:    respiratory, lems, sclc, delta-p, score
// Units:   one bool per item; result points 0–6
// Range:   persons with Lambert–Eaton myasthenic syndrome (the source's population); not enforced
// Example: bulbar involvement and erectile dysfunction only → 2 (hand-derived: one point per item)
func DeltaPScore(f DeltaPFindings) int {
	n := 0
	for _, present := range [6]bool{f.AgeAtOnsetAtLeast50, f.SmokingAtDiagnosis, f.WeightLossAtLeast5Pct,
		f.BulbarInvolvement, f.ErectileDysfunction, f.KarnofskyBelow70} {
		if present {
			n++
		}
	}
	return n
}
