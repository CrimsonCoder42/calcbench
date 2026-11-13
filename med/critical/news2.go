package critical

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when a float input is NaN or ±Inf.
	ErrNotFinite = errors.New("critical: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside its physical range or the range the source validated.
	ErrOutOfRange = errors.New("critical: input outside the range the source validated")
)

// finite returns ErrNotFinite, wrapped with name, when v is NaN or ±Inf.
func finite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// within returns ErrOutOfRange, wrapped with name, unless lo ≤ v ≤ hi. NaN passes; check finite first.
func within(name string, v, lo, hi float64) error {
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v–%v", ErrOutOfRange, name, v, lo, hi)
	}
	return nil
}

// atLeastZero returns ErrOutOfRange, wrapped with name, when v < 0. NaN passes; check finite first.
func atLeastZero(name string, v float64) error {
	if v < 0 {
		return fmt.Errorf("%w: %s = %v, want ≥ 0", ErrOutOfRange, name, v)
	}
	return nil
}

// NEWS2Scores holds the NEWS2 score of each parameter, as printed in RCP Chart 1, and their aggregate.
type NEWS2Scores struct {
	RespirationRate int  // 0–3
	SpO2            int  // 0–3, from SpO2 Scale 1 or Scale 2
	AirOrOxygen     int  // 0 (Air) or 2 (Oxygen)
	SystolicBP      int  // 0–3
	Pulse           int  // 0–3
	Consciousness   int  // 0 (Alert) or 3 (CVPU)
	Temperature     int  // 0–3
	Total           int  // aggregate 0–20
	RedScore        bool // true when any single parameter scores 3 (RCP Chart 2: "Score of 3 in any individual parameter")
}

// NEWS2 returns the National Early Warning Score 2 for one set of observations: each parameter's score from
// RCP Chart 1 and their aggregate (0–20). SpO2 is scored on Scale 1 unless scale2 is true; on Scale 2 a
// saturation of 93% or more scores 0 on air and 1, 2 or 3 on oxygen (Chart 1: "≥93 on air", "93–94 on oxygen",
// "95–96 on oxygen", "≥97 on oxygen"). alert is the chart's "Alert"; false means new confusion, voice, pain or
// unresponsive ("CVPU"). The chart prints temperature at 0.1 °C; a reading between two printed values is scored
// by the band whose upper edge it does not exceed (e.g. 36.05 °C scores as 36.1–38.0, i.e. 0).
// Errors: ErrNotFinite for a NaN or ±Inf temperature; ErrOutOfRange for a negative rate, pressure or pulse, an
// SpO2 outside 0–100, or a temperature outside 10–50 °C.
//
// ID:      MED-CRIT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Royal College of Physicians. National Early Warning Score (NEWS) 2: Standardising the assessment of acute-illness severity in the NHS. London: RCP, 2017. Chart 1 "The NEWS scoring system"
// Tags:    critical, news2, early warning, vital signs, score
// Units:   respirationRatePerMin breaths/min; spo2Percent %; systolicMmHg mmHg; pulsePerMin beats/min; temperatureC °C; result points
// Range:   respirationRatePerMin ≥ 0; 0 ≤ spo2Percent ≤ 100; systolicMmHg ≥ 0; pulsePerMin ≥ 0; 10 ≤ temperatureC ≤ 50 (physical limits; the chart's outer bands are open-ended)
// Example: RR 16, SpO2 98 (Scale 1, air), SBP 120, pulse 70, alert, 37.0 °C → 0; RR 26, SpO2 90 on oxygen (Scale 1), SBP 88, pulse 135, CVPU, 35.0 °C → 3+3+2+3+3+3+3 = 20 (HAND-DERIVED from Chart 1)
func NEWS2(respirationRatePerMin, spo2Percent int, scale2, onOxygen bool, systolicMmHg, pulsePerMin int, alert bool, temperatureC float64) (NEWS2Scores, error) {
	if math.IsNaN(temperatureC) || math.IsInf(temperatureC, 0) {
		return NEWS2Scores{}, fmt.Errorf("%w: temperatureC = %v", ErrNotFinite, temperatureC)
	}
	switch {
	case respirationRatePerMin < 0:
		return NEWS2Scores{}, fmt.Errorf("%w: respirationRatePerMin = %d, want ≥ 0", ErrOutOfRange, respirationRatePerMin)
	case spo2Percent < 0 || spo2Percent > 100:
		return NEWS2Scores{}, fmt.Errorf("%w: spo2Percent = %d, want 0–100", ErrOutOfRange, spo2Percent)
	case systolicMmHg < 0:
		return NEWS2Scores{}, fmt.Errorf("%w: systolicMmHg = %d, want ≥ 0", ErrOutOfRange, systolicMmHg)
	case pulsePerMin < 0:
		return NEWS2Scores{}, fmt.Errorf("%w: pulsePerMin = %d, want ≥ 0", ErrOutOfRange, pulsePerMin)
	case temperatureC < 10 || temperatureC > 50:
		return NEWS2Scores{}, fmt.Errorf("%w: temperatureC = %v, want 10–50", ErrOutOfRange, temperatureC)
	}
	var s NEWS2Scores

	// Respiration rate (per minute): ≤8 → 3; 9–11 → 1; 12–20 → 0; 21–24 → 2; ≥25 → 3.
	switch r := respirationRatePerMin; {
	case r <= 8:
		s.RespirationRate = 3
	case r <= 11:
		s.RespirationRate = 1
	case r <= 20:
		s.RespirationRate = 0
	case r <= 24:
		s.RespirationRate = 2
	default:
		s.RespirationRate = 3
	}

	if scale2 {
		// SpO2 Scale 2 (%): ≤83 → 3; 84–85 → 2; 86–87 → 1; 88–92 or ≥93 on air → 0;
		// 93–94 on oxygen → 1; 95–96 on oxygen → 2; ≥97 on oxygen → 3.
		switch p := spo2Percent; {
		case p <= 83:
			s.SpO2 = 3
		case p <= 85:
			s.SpO2 = 2
		case p <= 87:
			s.SpO2 = 1
		case p <= 92 || !onOxygen:
			s.SpO2 = 0
		case p <= 94:
			s.SpO2 = 1
		case p <= 96:
			s.SpO2 = 2
		default:
			s.SpO2 = 3
		}
	} else {
		// SpO2 Scale 1 (%): ≤91 → 3; 92–93 → 2; 94–95 → 1; ≥96 → 0.
		switch p := spo2Percent; {
		case p <= 91:
			s.SpO2 = 3
		case p <= 93:
			s.SpO2 = 2
		case p <= 95:
			s.SpO2 = 1
		default:
			s.SpO2 = 0
		}
	}

	// Air or oxygen?: Oxygen → 2; Air → 0.
	if onOxygen {
		s.AirOrOxygen = 2
	}

	// Systolic blood pressure (mmHg): ≤90 → 3; 91–100 → 2; 101–110 → 1; 111–219 → 0; ≥220 → 3.
	switch b := systolicMmHg; {
	case b <= 90:
		s.SystolicBP = 3
	case b <= 100:
		s.SystolicBP = 2
	case b <= 110:
		s.SystolicBP = 1
	case b <= 219:
		s.SystolicBP = 0
	default:
		s.SystolicBP = 3
	}

	// Pulse (per minute): ≤40 → 3; 41–50 → 1; 51–90 → 0; 91–110 → 1; 111–130 → 2; ≥131 → 3.
	switch p := pulsePerMin; {
	case p <= 40:
		s.Pulse = 3
	case p <= 50:
		s.Pulse = 1
	case p <= 90:
		s.Pulse = 0
	case p <= 110:
		s.Pulse = 1
	case p <= 130:
		s.Pulse = 2
	default:
		s.Pulse = 3
	}

	// Consciousness: Alert → 0; CVPU → 3.
	if !alert {
		s.Consciousness = 3
	}

	// Temperature (°C): ≤35.0 → 3; 35.1–36.0 → 1; 36.1–38.0 → 0; 38.1–39.0 → 1; ≥39.1 → 2.
	switch c := temperatureC; {
	case c <= 35.0:
		s.Temperature = 3
	case c <= 36.0:
		s.Temperature = 1
	case c <= 38.0:
		s.Temperature = 0
	case c <= 39.0:
		s.Temperature = 1
	default:
		s.Temperature = 2
	}

	s.Total = s.RespirationRate + s.SpO2 + s.AirOrOxygen + s.SystolicBP + s.Pulse + s.Consciousness + s.Temperature
	s.RedScore = s.RespirationRate == 3 || s.SpO2 == 3 || s.SystolicBP == 3 || s.Pulse == 3 ||
		s.Consciousness == 3 || s.Temperature == 3
	return s, nil
}
