package gastro

import "fmt"

// AcutePancreatitisFindings holds the eight items of the Jin acute pancreatitis prediction model, assessed in
// the emergency department in a person with serum lipase ≥ 3 times the upper limit of normal (Jin 2024 Table 1).
type AcutePancreatitisFindings struct {
	PriorEpisodes                int     // prior episodes of acute pancreatitis, ≥ 0; 0–3 score their count, ≥ 4 score 4
	Cholelithiasis               bool    // history of cholelithiasis: 2 points
	AbdominalSurgeryPrior2Months bool    // abdominal surgery in the prior 2 months: −2 points
	EpigastricPain               bool    // epigastric pain: 2 points
	WorseningSeverity            bool    // worsening severity of pain: 1 point
	PainDurationAtLeast5Days     bool    // pain duration ≥ 5 days: −1 point (< 5 days: 0)
	PainSeverity                 int     // pain severity on a 0–10 scale: 0–3 mild 0, 4–6 moderate 2, 7–10 severe 3
	LipaseTimesULN               float64 // serum lipase as a multiple of the upper limit of normal, ≥ 3: < 10 0, ≥ 10 to < 20 1, ≥ 20 2
}

// AcutePancreatitisDiagnosisScore returns the points of the Jin clinical model for the early diagnosis of
// acute pancreatitis in the emergency department, −3 to 14, by Jin 2024 Table 1 (the 2018 derivation model,
// prospectively validated). No band is returned; the source's reported thresholds are ≥ 6 and ≥ 8 points.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-GASTRO-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Jin DX, Lacson R, Cochon LR, et al. A clinical model for the early diagnosis of acute pancreatitis in the emergency department. Pancreas. 2018;47(7):871-879. doi:10.1097/MPA.0000000000001102. PMID 29975351 (derivation; abstract: "diagnostic risk threshold of 8 points or higher")
// Source:  Jin DX, Lacson R, Eskian M, et al. Prospective validation of a prediction model for the diagnosis of acute pancreatitis. JAMA Netw Open. 2024;7(6):e2419014. doi:10.1001/jamanetworkopen.2024.19014. PMID 38941094. PMC11214112, Table 1 "Prediction Model for the Diagnosis of AP: Categories and Points" (range −3 to 14)
// Tags:    gastro, pancreas, acute-pancreatitis, lipase, emergency, score
// Units:   counts and 0–10 pain severity as integers; LipaseTimesULN multiples of the upper limit of normal; result points −3 to 14
// Range:   PriorEpisodes ≥ 0; PainSeverity 0–10; LipaseTimesULN ≥ 3 (the model's population: emergency department visits with lipase ≥ 3 × ULN)
// Example: 1 prior episode, cholelithiasis, epigastric pain, worsening, pain 8/10, lipase 25 × ULN → 1 + 2 + 2 + 1 + 3 + 2 = 11 (Table 1)
func AcutePancreatitisDiagnosisScore(f AcutePancreatitisFindings) (int, error) {
	if err := checkFinite("LipaseTimesULN", f.LipaseTimesULN); err != nil {
		return 0, err
	}
	if f.PriorEpisodes < 0 {
		return 0, fmt.Errorf("%w: PriorEpisodes = %d, want ≥ 0", ErrOutOfRange, f.PriorEpisodes)
	}
	if f.PainSeverity < 0 || f.PainSeverity > 10 {
		return 0, fmt.Errorf("%w: PainSeverity = %d, want 0–10", ErrOutOfRange, f.PainSeverity)
	}
	if !(f.LipaseTimesULN >= 3) {
		return 0, fmt.Errorf("%w: LipaseTimesULN = %v, want ≥ 3", ErrOutOfRange, f.LipaseTimesULN)
	}
	points := min(f.PriorEpisodes, 4)
	if f.Cholelithiasis {
		points += 2
	}
	if f.AbdominalSurgeryPrior2Months {
		points -= 2
	}
	if f.EpigastricPain {
		points += 2
	}
	if f.WorseningSeverity {
		points++
	}
	if f.PainDurationAtLeast5Days {
		points--
	}
	switch {
	case f.PainSeverity >= 7:
		points += 3
	case f.PainSeverity >= 4:
		points += 2
	}
	switch {
	case f.LipaseTimesULN >= 20:
		points += 2
	case f.LipaseTimesULN >= 10:
		points++
	}
	return points, nil
}
