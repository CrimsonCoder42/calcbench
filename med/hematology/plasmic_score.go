package hematology

// PLASMICFindings holds the seven PLASMIC components (Singh 2026 Table 1; Lee 2022; Baysal 2023).
type PLASMICFindings struct {
	PlateletsE9PerL float64 // "Platelet count of <30 × 10⁹/L": 1
	Hemolysis       bool    // "reticulocyte count of >2.5%, or haptoglobin undetectable, or indirect bilirubin of >2.0 mg/dL": 1
	NoActiveCancer  bool    // "No active cancer": 1
	NoTransplant    bool    // "No history of solid organ or stem cell transplant": 1
	MCVfL           float64 // "MCV of <90 fL": 1
	INR             float64 // "International Normalized Ratio of <1.5": 1
	CreatinineMgDL  float64 // "Creatinine of <2.0 mg/dL": 1
}

// PLASMICScore returns the PLASMIC score, 0–7: "One point is assigned for each criterion met".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative count, MCV, INR or creatinine.
//
// ID:      MED-HEMAT-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bendapudi PK, Hurwitz S, Fry A, et al. Derivation and external validation of the PLASMIC score for rapid assessment of adults with thrombotic microangiopathies: a cohort study. Lancet Haematol. 2017;4(4):e157-e164. doi:10.1016/S2352-3026(17)30026-1, PMID 28259520 (primary; not opened)
// Source:  Singh R, Khouri A, Lee HJ, et al. A multicenter analysis revisiting the PLASMIC score: performance limitations and the impact of age, MCV, and LDH. Blood Vessel Thromb Hemost. 2026;3(2):100163. doi:10.1016/j.bvth.2026.100163, PMID 42164818, PMC13185983, Table 1 "PLASMIC score components and definitions"
// Source:  Lee CH, Huang YC, Li SS, et al. Application of PLASMIC Score in Risk Prediction of Thrombotic Thrombocytopenic Purpura: Real-World Experience From a Tertiary Medical Center in Taiwan. Front Med (Lausanne). 2022;9:893273. doi:10.3389/fmed.2022.893273, PMID 35615090, PMC9124890, Methods
// Verified: secondary (Singh 2026 PMC13185983 Table 1; Lee 2022 PMC9124890)
// Tags:    hematology, thrombotic microangiopathy, ttp, plasmic, score
// Units:   PlateletsE9PerL ×10⁹/L; MCVfL fL; INR ratio; CreatinineMgDL mg/dL; result points 0–7
// Range:   adults with thrombotic microangiopathy (the derivation population); non-negative measurements
// Example: platelets 12, hemolysis, no cancer, no transplant, MCV 88, INR 1.2, creatinine 1.1 → 7 (hand-derived)
func PLASMICScore(f PLASMICFindings) (int, error) {
	for _, in := range [...]struct {
		name string
		v    float64
	}{{"PlateletsE9PerL", f.PlateletsE9PerL}, {"MCVfL", f.MCVfL}, {"INR", f.INR}, {"CreatinineMgDL", f.CreatinineMgDL}} {
		if err := checkFloat(in.name, in.v, 0, 1e6); err != nil {
			return 0, err
		}
	}
	s := 0
	for _, met := range [...]bool{
		f.PlateletsE9PerL < 30, f.Hemolysis, f.NoActiveCancer, f.NoTransplant, f.MCVfL < 90, f.INR < 1.5, f.CreatinineMgDL < 2.0,
	} {
		if met {
			s++
		}
	}
	return s, nil
}

// PLASMICRisk is a PLASMIC risk group; values are the secondaries' words.
type PLASMICRisk string

const (
	PLASMICLow          PLASMICRisk = "low risk"          // 0–4
	PLASMICIntermediate PLASMICRisk = "intermediate risk" // 5
	PLASMICHigh         PLASMICRisk = "high risk"         // 6–7
)

// PLASMICRiskGroup returns the PLASMIC risk group: "0 to 4 (low risk), 5 (intermediate risk), and 6 to 7
// (high risk)".
// Errors: ErrOutOfRange for a score outside 0–7.
//
// ID:      MED-HEMAT-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bendapudi PK, Hurwitz S, Fry A, et al. Lancet Haematol. 2017;4(4):e157-e164. doi:10.1016/S2352-3026(17)30026-1, PMID 28259520 (primary; not opened)
// Source:  Baysal M, Hindilerden F, Ümit EG, et al. Immune Thrombotic Thrombocytopenic Purpura in Elderly Patients: The Roles of PLASMIC and French Scores. Turk J Haematol. 2023;40(4):251-257. doi:10.4274/tjh.galenos.2023.2023.0295, PMID 37791641, PMC10701326
// Source:  Lee CH, Huang YC, Li SS, et al. Front Med (Lausanne). 2022;9:893273. doi:10.3389/fmed.2022.893273, PMID 35615090, PMC9124890
// Verified: secondary (Baysal 2023 PMC10701326; Lee 2022 PMC9124890)
// Tags:    hematology, thrombotic microangiopathy, ttp, plasmic, category
// Units:   score PLASMIC points; result the secondaries' label
// Range:   score 0–7
// Example: 4 → low risk; 5 → intermediate risk; 6 → high risk
func PLASMICRiskGroup(score int) (PLASMICRisk, error) {
	switch {
	case score < 0 || score > 7:
		return "", checkFloat("score", float64(score), 0, 7)
	case score <= 4:
		return PLASMICLow, nil
	case score == 5:
		return PLASMICIntermediate, nil
	}
	return PLASMICHigh, nil
}
