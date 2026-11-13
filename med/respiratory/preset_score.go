package respiratory

import (
	"fmt"
	"math"
)

// PRESETScore returns the PRESET (PREdiction of Survival on ECMO Therapy) score (0–15) of Hilder 2017,
// Table 4: mean arterial pressure > 100 → 0, 91–100 → 1, 81–90 → 2, 71–80 → 3, ≤ 70 → 4 mmHg; lactate
// ≤ 1.50 → 0, 1.51–3.00 → 1, 3.01–6.00 → 2, 6.01–10.00 → 3, > 10.00 → 4 mmol/L; arterial pH > 7.300 → 0,
// 7.201–7.300 → 1, 7.101–7.200 → 2, ≤ 7.100 → 3; platelets > 200 → 0, 101–200 → 1, ≤ 100 → 2 ×1000/µL;
// hospital days before ECMO ≤ 2 → 0, 3–7 → 1, > 7 → 2. A value between two printed bands (MAP 90.5, pH
// 7.2005) is compared with the upper edge of the lower band: MAP > 90 scores 1, pH > 7.200 scores 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hilder M, Herbstreit F, Adamzik M, et al. Comparison of mortality prediction models in acute respiratory distress syndrome undergoing extracorporeal membrane oxygenation and development of a novel prediction score: the PREdiction of Survival on ECMO Therapy-Score (PRESET-Score). Crit Care. 2017;21(1):301. doi:10.1186/s13054-017-1888-6. PMID 29233160, PMC5728043, Table 4
// Tags:    respiratory, ards, ecmo, preset, mortality, score
// Units:   meanArterialMmHg mmHg; lactateMmolL mmol/L; arterialPH pH units; plateletsE3PerUL ×1000/µL (= ×10⁹/L); hospitalDaysBeforeECMO whole days; result points 0–15
// Range:   meanArterialMmHg > 0; lactateMmolL ≥ 0; arterialPH 6.5–8.0 (a plausibility bound; the paper states no limits); plateletsE3PerUL ≥ 0; hospitalDaysBeforeECMO ≥ 0; adults with ARDS on veno-venous ECMO
// Example: MAP 75, lactate 4, pH 7.15, platelets 150, 5 days → 3 + 2 + 2 + 1 + 1 = 9 (hand-derived from Table 4)
func PRESETScore(meanArterialMmHg, lactateMmolL, arterialPH, plateletsE3PerUL float64, hospitalDaysBeforeECMO int) (int, error) {
	if err := checkPositiveFinite("meanArterialMmHg", meanArterialMmHg); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("lactateMmolL", lactateMmolL); err != nil {
		return 0, err
	}
	if err := checkClosedRange("arterialPH", arterialPH, 6.5, 8.0); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("plateletsE3PerUL", plateletsE3PerUL); err != nil {
		return 0, err
	}
	if err := checkIntRange("hospitalDaysBeforeECMO", hospitalDaysBeforeECMO, 0, math.MaxInt32); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case meanArterialMmHg > 100:
	case meanArterialMmHg > 90:
		s++
	case meanArterialMmHg > 80:
		s += 2
	case meanArterialMmHg > 70:
		s += 3
	default:
		s += 4
	}
	switch {
	case lactateMmolL <= 1.5:
	case lactateMmolL <= 3:
		s++
	case lactateMmolL <= 6:
		s += 2
	case lactateMmolL <= 10:
		s += 3
	default:
		s += 4
	}
	switch {
	case arterialPH > 7.3:
	case arterialPH > 7.2:
		s++
	case arterialPH > 7.1:
		s += 2
	default:
		s += 3
	}
	switch {
	case plateletsE3PerUL > 200:
	case plateletsE3PerUL > 100:
		s++
	default:
		s += 2
	}
	switch {
	case hospitalDaysBeforeECMO <= 2:
	case hospitalDaysBeforeECMO <= 7:
		s++
	default:
		s += 2
	}
	return s, nil
}

// PRESETRisk is a PRESET risk class; its value is the source's label, verbatim.
type PRESETRisk string

// The classes of Hilder 2017 (Results / Fig. 3), verbatim: "Class I (0–5)", "Class II (6–9)", "Class III (10–15)".
const (
	PRESETClassI   PRESETRisk = "Class I"
	PRESETClassII  PRESETRisk = "Class II"
	PRESETClassIII PRESETRisk = "Class III"
)

// PRESETRiskClass returns the PRESET risk class for a score: 0–5 Class I, 6–9 Class II, 10–15 Class III.
// Errors: ErrOutOfRange for a score outside 0–15.
//
// ID:      MED-RESP-035
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hilder M, Herbstreit F, Adamzik M, et al. Comparison of mortality prediction models in acute respiratory distress syndrome undergoing extracorporeal membrane oxygenation and development of a novel prediction score: the PREdiction of Survival on ECMO Therapy-Score (PRESET-Score). Crit Care. 2017;21(1):301. doi:10.1186/s13054-017-1888-6. PMID 29233160, PMC5728043, Results (risk classes I–III)
// Tags:    respiratory, ecmo, preset, category
// Units:   score points 0–15; result the source's class label
// Range:   score 0–15
// Example: 5 → "Class I"; 6 → "Class II"; 10 → "Class III"
func PRESETRiskClass(score int) (PRESETRisk, error) {
	if score < 0 || score > 15 {
		return "", fmt.Errorf("%w: score = %d, want 0 ≤ score ≤ 15", ErrOutOfRange, score)
	}
	switch {
	case score <= 5:
		return PRESETClassI, nil
	case score <= 9:
		return PRESETClassII, nil
	}
	return PRESETClassIII, nil
}
