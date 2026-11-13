package cardio

// WellsPEFindings are the seven items of Wells' PE score (Wells 2000 abstract; points in the comments).
type WellsPEFindings struct {
	ClinicalSignsOfDVT      bool // "Clinical symptoms of DVT": 3.0
	NoAlternativeDiagnosis  bool // "no alternative diagnosis": 3.0
	HeartRateAbove100       bool // "heart rate >100": 1.5
	ImmobilizationOrSurgery bool // "immobilization or surgery in the previous four weeks": 1.5
	PreviousDVTOrPE         bool // "previous DVT/PE": 1.5
	Hemoptysis              bool // 1.0
	Malignancy              bool // 1.0
}

// WellsPEScore returns Wells' clinical score for pulmonary embolism, 0–12.5 points. The paper's groups:
// three-level < 2.0 low, 2.0–6.0 moderate, > 6.0 high; two-level ≤ 4.0 "PE unlikely", > 4.0 "PE likely".
// Errors: none.
//
// ID:      MED-CARDIO-040
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wells PS, Anderson DR, Rodger M, et al. Derivation of a simple clinical model to categorize patients probability of pulmonary embolism. Thromb Haemost. 2000;83(3):416-420. doi:10.1055/s-0037-1613830. PMID 10744147. Abstract: "Clinical symptoms of DVT (3.0), no alternative diagnosis (3.0), heart rate >100 (1.5), immobilization or surgery in the previous four weeks (1.5), previous DVT/PE (1.5), hemoptysis (1.0) and malignancy (1.0)"
// Tags:    cardio, pulmonary-embolism, vte, wells, risk-score
// Units:   result points 0–12.5 (multiples of 0.5)
// Range:   hospital and ambulatory persons with suspected PE (Canadian cohort); no numeric inputs
// Example: heart rate > 100 and hemoptysis → 2.5 (HAND-DERIVED)
func WellsPEScore(f WellsPEFindings) float64 {
	s := 0.0
	add := func(present bool, pts float64) {
		if present {
			s += pts
		}
	}
	add(f.ClinicalSignsOfDVT, 3)
	add(f.NoAlternativeDiagnosis, 3)
	add(f.HeartRateAbove100, 1.5)
	add(f.ImmobilizationOrSurgery, 1.5)
	add(f.PreviousDVTOrPE, 1.5)
	add(f.Hemoptysis, 1)
	add(f.Malignancy, 1)
	return s
}
