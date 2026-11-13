package cardio

// GenevaRevisedFindings are the yes/no items of the revised Geneva score besides age and heart rate
// (Le Gal 2006 abstract; points in the comments).
type GenevaRevisedFindings struct {
	PreviousDVTOrPE           bool // 3
	SurgeryOrFracture1Month   bool // "surgery or fracture within 1 month": 2
	ActiveMalignancy          bool // 2
	UnilateralLowerLimbPain   bool // 3
	Hemoptysis                bool // 2
	DeepVeinPalpationAndEdema bool // "pain on lower-limb deep venous palpation and unilateral edema": 4
}

// GenevaScoreRevised returns the revised Geneva score for pulmonary embolism, 0–22: age > 65 years 1; heart
// rate 75–94 beats/min 3 or ≥ 95 beats/min 5; and the item points. The paper's groups: low 0–3, intermediate
// 4–10, high ≥ 11.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120 or heartRateBpm ≤ 0.
//
// ID:      MED-CARDIO-041
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Le Gal G, Righini M, Roy PM, et al. Prediction of pulmonary embolism in the emergency department: the revised Geneva score. Ann Intern Med. 2006;144(3):165-171. doi:10.7326/0003-4819-144-3-200602070-00004. PMID 16461960. Abstract: "age older than 65 years (1 point), previous deep venous thrombosis or pulmonary embolism (3 points), surgery or fracture within 1 month (2 points), active malignant condition (2 points), unilateral lower limb pain (3 points), hemoptysis (2 points), heart rate of 75 to 94 beats/min (3 points) or 95 beats/min or more (5 points), and pain on lower-limb deep venous palpation and unilateral edema (4 points)"
// Tags:    cardio, pulmonary-embolism, vte, geneva, risk-score
// Units:   ageYears years; heartRateBpm beats/min (whole beats: 94.5 counts as 75–94); result points 0–22
// Range:   ageYears 18–120; heartRateBpm > 0 (ED cohorts of 3 European university hospitals; no limits stated)
// Example: age 70, heart rate 100, hemoptysis → 1 + 5 + 2 = 8 (HAND-DERIVED)
func GenevaScoreRevised(ageYears, heartRateBpm float64, f GenevaRevisedFindings) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkFinite("heartRateBpm", heartRateBpm); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	if err := checkPositive("heartRateBpm", heartRateBpm); err != nil {
		return 0, err
	}
	s := 0
	if ageYears > 65 {
		s++
	}
	switch {
	case heartRateBpm >= 95:
		s += 5
	case heartRateBpm >= 75:
		s += 3
	}
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(f.PreviousDVTOrPE, 3)
	add(f.SurgeryOrFracture1Month, 2)
	add(f.ActiveMalignancy, 2)
	add(f.UnilateralLowerLimbPain, 3)
	add(f.Hemoptysis, 2)
	add(f.DeepVeinPalpationAndEdema, 4)
	return s, nil
}
