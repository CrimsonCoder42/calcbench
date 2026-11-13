package respiratory

import (
	"fmt"
	"math"
)

// OxygenationIndex returns the oxygenation index OI = mean airway pressure × FiO₂ / PaO₂ (PALICC-2, Table 4
// footnote), with FiO₂ entered in percent: the reading under which the guideline's thresholds (OI ≥ 4 for
// PARDS on invasive ventilation; "Mild/moderate: OI < 16", "Severe: OI ≥ 16") apply. In fraction form this is
// MAP × FiO₂ × 100 / PaO₂.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the product overflows.
//
// ID:      MED-RESP-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Emeriaud G, López-Fernández YM, et al. Executive Summary of the Second International Guidelines for the Diagnosis and Management of Pediatric Acute Respiratory Distress Syndrome (PALICC-2). Pediatr Crit Care Med. 2023;24(2):143-168. doi:10.1097/PCC.0000000000003147. PMID 36661420, PMC9848214, Table 4 footnote: "OI = mean airway pressure (MAP) (cm H2O) × Fio2/Pao2 (mm Hg)"
// Tags:    respiratory, oxygenation-index, oi, pards, palicc, pediatric, ventilation
// Units:   meanAirwayPressureCmH2O cmH₂O; fio2Percent % (21–100); paO2MmHg mmHg; result unitless index
// Range:   meanAirwayPressureCmH2O > 0; 21 ≤ fio2Percent ≤ 100; paO2MmHg > 0
// Example: MAP 15 cmH₂O, FiO₂ 60 %, PaO₂ 60 mmHg → 15 × 60 / 60 = 15 (hand-derived)
func OxygenationIndex(meanAirwayPressureCmH2O, fio2Percent, paO2MmHg float64) (float64, error) {
	if err := checkPositiveFinite("meanAirwayPressureCmH2O", meanAirwayPressureCmH2O); err != nil {
		return 0, err
	}
	if err := checkClosedRange("fio2Percent", fio2Percent, 21, 100); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("paO2MmHg", paO2MmHg); err != nil {
		return 0, err
	}
	return pressureTimesFiO2Over(meanAirwayPressureCmH2O, fio2Percent, paO2MmHg)
}

// OxygenSaturationIndex returns the oxygen saturation index OSI = mean airway pressure × FiO₂ / SpO₂ (PALICC-2,
// Table 4 footnote), with FiO₂ entered in percent as for OxygenationIndex. The guideline uses SpO₂ only when it is
// ≤ 97 %, and the function enforces that limit.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the product overflows.
//
// ID:      MED-RESP-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Emeriaud G, López-Fernández YM, et al. Executive Summary of the Second International Guidelines for the Diagnosis and Management of Pediatric Acute Respiratory Distress Syndrome (PALICC-2). Pediatr Crit Care Med. 2023;24(2):143-168. doi:10.1097/PCC.0000000000003147. PMID 36661420, PMC9848214, Table 4 footnote: "OSI = MAP (cm H2O) × Fio2/Spo2"; "When Spo2 is used, ensure that Spo2 is ≤ 97%"
// Tags:    respiratory, oxygen-saturation-index, osi, pards, palicc, pediatric, ventilation
// Units:   meanAirwayPressureCmH2O cmH₂O; fio2Percent % (21–100); spo2Percent %; result unitless index
// Range:   meanAirwayPressureCmH2O > 0; 21 ≤ fio2Percent ≤ 100; 0 < spo2Percent ≤ 97
// Example: MAP 12 cmH₂O, FiO₂ 50 %, SpO₂ 92 % → 600 / 92 = 6.5217… (hand-derived)
func OxygenSaturationIndex(meanAirwayPressureCmH2O, fio2Percent, spo2Percent float64) (float64, error) {
	if err := checkPositiveFinite("meanAirwayPressureCmH2O", meanAirwayPressureCmH2O); err != nil {
		return 0, err
	}
	if err := checkClosedRange("fio2Percent", fio2Percent, 21, 100); err != nil {
		return 0, err
	}
	if err := checkPositiveAtMost("spo2Percent", spo2Percent, 97); err != nil {
		return 0, err
	}
	return pressureTimesFiO2Over(meanAirwayPressureCmH2O, fio2Percent, spo2Percent)
}

// pressureTimesFiO2Over returns p·f/d, with ErrDomain if the result is not finite.
func pressureTimesFiO2Over(p, f, d float64) (float64, error) {
	r := p * f / d
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: MAP × FiO₂ / denominator overflows", ErrDomain)
	}
	return r, nil
}
