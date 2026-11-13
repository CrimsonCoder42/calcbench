package cardio

import "fmt"

// HEARTScore returns the HEART score for chest pain, 0–10, as the HEART Pathway trial defines it: History,
// ECG, Risk factors and Troponin are each rated 0–1–2 by the caller (history slightly / moderately / highly
// suspicious; ECG normal / nonspecific changes / changes consistent with ACS; 0 / 1–2 / ≥ 3 risk factors;
// troponin negative / 1–3× / > 3× the normal limit), and Age is scored from ageYears: < 45 → 0, 45–65 → 1,
// > 65 → 2. Six 2008's original table prints "45-65" and "≥65", overlapping at 65; Mahler 2015's "> 65" is
// used here.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-057
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Mahler SA, Riley RF, Hiestand BC, et al. The HEART Pathway randomized trial: identifying emergency department patients with acute chest pain for early discharge. Circ Cardiovasc Qual Outcomes. 2015;8(2):195-203. doi:10.1161/CIRCOUTCOMES.114.001384. PMID 25737484. PMC4413911, Table 2 (HEART score elements)
// Source:  Six AJ, Backus BE, Kelder JC. Chest pain in the emergency room: value of the HEART score. Neth Heart J. 2008;16(6):191-196. doi:10.1007/BF03086144. PMID 18665203. PMC2442661, Table 1 (the original five elements, 0–2 each)
// Tags:    cardio, chest-pain, acs, heart-score, risk-score
// Units:   history, ecg, riskFactors, troponin points 0–2; ageYears years; result points 0–10
// Range:   each rated item 0–2; ageYears 18–110 (adults with acute chest pain; no upper limit stated)
// Example: history 1, ECG 0, age 50, risk factors 1, troponin 0 → 1 + 0 + 1 + 1 + 0 = 3 (HAND-DERIVED)
func HEARTScore(history, ecg int, ageYears float64, riskFactors, troponin int) (int, error) {
	s, err := hearPoints(history, ecg, ageYears, riskFactors)
	if err != nil {
		return 0, err
	}
	if troponin < 0 || troponin > 2 {
		return 0, fmt.Errorf("%w: troponin = %d, want 0–2", ErrOutOfRange, troponin)
	}
	return s + troponin, nil
}

// hearPoints sums History, ECG, Age and Risk factors after checking them.
func hearPoints(history, ecg int, ageYears float64, riskFactors int) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	for _, it := range [...]struct {
		name string
		v    int
	}{{"history", history}, {"ecg", ecg}, {"riskFactors", riskFactors}} {
		if it.v < 0 || it.v > 2 {
			return 0, fmt.Errorf("%w: %s = %d, want 0–2", ErrOutOfRange, it.name, it.v)
		}
	}
	if err := checkRange("ageYears", ageYears, 18, 110); err != nil {
		return 0, err
	}
	age := 0
	switch {
	case ageYears > 65:
		age = 2
	case ageYears >= 45:
		age = 1
	}
	return history + ecg + age + riskFactors, nil
}
