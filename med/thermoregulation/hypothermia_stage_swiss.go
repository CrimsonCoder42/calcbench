package thermoregulation

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("thermoregulation: value is NaN or infinite")
	// ErrOutOfRange: an input is outside the range the source validated.
	ErrOutOfRange = errors.New("thermoregulation: input outside the range the source validated")
)

// checkFinite returns ErrNotFinite for NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// HypothermiaStage is a stage of the staging of accidental hypothermia in Paal 2016 Table 1; its value is
// the table's own stage label, verbatim.
type HypothermiaStage string

// The four stages of Paal 2016 Table 1 ("Staging of accidental hypothermia"), labels verbatim, with the
// table's core-temperature column.
const (
	HypothermiaI   HypothermiaStage = "Hypothermia I (mild)"      // 35–32 °C; "Conscious, shivering"
	HypothermiaII  HypothermiaStage = "Hypothermia II (moderate)" // <32–28 °C; "Impaired consciousness; may or may not be shivering"
	HypothermiaIII HypothermiaStage = "Hypothermia III (severe)"  // <28 °C; "Unconscious; vital signs present"
	// HypothermiaIV is never returned by the temperature functions: the table's core temperature for it is
	// "Variable" and its finding ("Apparent death; Vital signs absent") is clinical.
	HypothermiaIV HypothermiaStage = "Hypothermia IV (severe)"
)

// HypothermiaStageSwiss returns the stage of accidental hypothermia whose core-temperature band in Paal 2016
// Table 1 contains coreTemperatureC: 32 to 35 °C (both included, the table's "35–32 °C") "Hypothermia I
// (mild)"; 28 to < 32 °C ("<32–28 °C") "Hypothermia II (moderate)"; < 28 °C "Hypothermia III (severe)".
// inTable is false, with stage "", above 35 °C, where the table has no row. The table stages by clinical
// findings alongside temperature; this function reads only its core-temperature column, so it never returns
// "Hypothermia IV (severe)", whose temperature the table gives as "Variable".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange below absolute zero (−273.15 °C).
//
// ID:      MED-THERMO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Paal P, Gordon L, Strapazzon G, et al. Accidental hypothermia-an update: The content of this review is endorsed by the International Commission for Mountain Emergency Medicine (ICAR MEDCOM). Scand J Trauma Resusc Emerg Med. 2016;24(1):111. doi:10.1186/s13049-016-0303-7. PMID 27633781. PMC5025630. Table 1, "Staging of accidental hypothermia", column "Core temperature (°C)"
// Tags:    thermoregulation, temperature, hypothermia, staging, core-temperature
// Units:   coreTemperatureC °C, core (not oral or skin) temperature; result the table's stage label
// Range:   coreTemperatureC ≥ −273.15; the table's rows cover ≤ 35 °C
// Example: 35 → "Hypothermia I (mild)"; 31.9 → "Hypothermia II (moderate)"; 28 → "Hypothermia II (moderate)"; 27.9 → "Hypothermia III (severe)"; 35.1 → not in the table (Paal 2016 Table 1)
func HypothermiaStageSwiss(coreTemperatureC float64) (stage HypothermiaStage, inTable bool, err error) {
	if err := checkFinite("coreTemperatureC", coreTemperatureC); err != nil {
		return "", false, err
	}
	if coreTemperatureC < -273.15 {
		return "", false, fmt.Errorf("%w: coreTemperatureC = %v, want ≥ −273.15", ErrOutOfRange, coreTemperatureC)
	}
	return stageFromEdges(coreTemperatureC, 35, 32, 28)
}

// HypothermiaStageSwissF returns the same stage as HypothermiaStageSwiss for a core temperature in degrees
// Fahrenheit, comparing against the exact Fahrenheit images of the table's edges: 35 °C = 95 °F, 32 °C =
// 89.6 °F, 28 °C = 82.4 °F (°F = °C × 9/5 + 32, the definition of the Fahrenheit scale, NIST SP 811 (2008)
// §B.9). inTable is false above 95 °F.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange below absolute zero (−459.67 °F).
//
// ID:      MED-THERMO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Paal P, Gordon L, Strapazzon G, et al. Accidental hypothermia-an update: The content of this review is endorsed by the International Commission for Mountain Emergency Medicine (ICAR MEDCOM). Scand J Trauma Resusc Emerg Med. 2016;24(1):111. doi:10.1186/s13049-016-0303-7. PMID 27633781. PMC5025630. Table 1, "Staging of accidental hypothermia", column "Core temperature (°C)"
// Tags:    thermoregulation, temperature, hypothermia, staging, core-temperature, fahrenheit
// Units:   coreTemperatureF °F, core temperature; result the table's stage label
// Range:   coreTemperatureF ≥ −459.67; the table's rows cover ≤ 95 °F
// Example: 95 → "Hypothermia I (mild)"; 89.6 → "Hypothermia I (mild)"; 89.5 → "Hypothermia II (moderate)"; 82.3 → "Hypothermia III (severe)"
func HypothermiaStageSwissF(coreTemperatureF float64) (stage HypothermiaStage, inTable bool, err error) {
	if err := checkFinite("coreTemperatureF", coreTemperatureF); err != nil {
		return "", false, err
	}
	if coreTemperatureF < -459.67 {
		return "", false, fmt.Errorf("%w: coreTemperatureF = %v, want ≥ −459.67", ErrOutOfRange, coreTemperatureF)
	}
	return stageFromEdges(coreTemperatureF, 95, 89.6, 82.4)
}

// stageFromEdges applies Table 1's bands in one unit: t ≤ top and t ≥ mildFloor → I; t ≥ moderateFloor → II;
// below → III; above top → not in the table.
func stageFromEdges(t, top, mildFloor, moderateFloor float64) (HypothermiaStage, bool, error) {
	switch {
	case t > top:
		return "", false, nil
	case t >= mildFloor:
		return HypothermiaI, true, nil
	case t >= moderateFloor:
		return HypothermiaII, true, nil
	}
	return HypothermiaIII, true, nil
}
