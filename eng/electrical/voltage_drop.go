package electrical

import (
	"fmt"
	"math"
)

// VoltageDropTwoWire returns the voltage lost along a two-wire run (DC, or single-phase AC feeding a resistive
// load), ΔV = I · R' · 2L: the current flows out and back, so the conductor path is twice the one-way length.
// Ohm's law on that loop resistance. Reactance and power factor are not modelled.
// Errors: ErrNotFinite; ErrNegative for any negative input; ErrOverflow.
//
// ID:      ENG-ELEC-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.4 "Ohm's Law" (V = IR) and §9.3 "Resistivity and Resistance" (R proportional to length), applied to the out-and-back conductor length 2L
// Tags:    voltage-drop, wire, conductor, feeder, single-phase, dc, electrical
// Units:   currentAmperes A; resistanceOhmsPerMeter Ω/m (one conductor); oneWayLengthMeters m; result V
// Example: 10 A, 0.0012 Ω/m, 152.4 m (500 ft) → 3.65760 V
func VoltageDropTwoWire(currentAmperes, resistanceOhmsPerMeter, oneWayLengthMeters float64) (float64, error) {
	return drop(2, currentAmperes, resistanceOhmsPerMeter, oneWayLengthMeters)
}

// VoltageDropThreePhase returns the line-to-line voltage lost along a balanced three-phase run feeding a
// resistive load, ΔV = √3 · I · R' · L: each conductor drops I·R'·L in phase with its line current, and the
// line-to-line magnitude of a balanced set is √3 times the line-to-neutral one.
// Errors: ErrNotFinite; ErrNegative for any negative input; ErrOverflow.
//
// ID:      ENG-ELEC-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Glover JD, Sarma MS, Overbye TJ. Power System Analysis and Design, 5th ed., Cengage, 2012, §2.5 "Balanced Three-Phase Circuits" (V_LL = √3 V_LN for a balanced set); Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.4 (V = IR per conductor)
// Tags:    voltage-drop, wire, conductor, feeder, three-phase, electrical
// Units:   currentAmperes A (line current); resistanceOhmsPerMeter Ω/m (one conductor); oneWayLengthMeters m; result V (line-to-line)
// Example: 10 A, 0.0012 Ω/m, 152.4 m → 3.16757 V
func VoltageDropThreePhase(currentAmperes, resistanceOhmsPerMeter, oneWayLengthMeters float64) (float64, error) {
	return drop(math.Sqrt(3), currentAmperes, resistanceOhmsPerMeter, oneWayLengthMeters)
}

func drop(factor, currentAmperes, resistanceOhmsPerMeter, oneWayLengthMeters float64) (float64, error) {
	if err := checkInput("currentAmperes", currentAmperes); err != nil {
		return 0, err
	}
	if err := checkInput("resistanceOhmsPerMeter", resistanceOhmsPerMeter); err != nil {
		return 0, err
	}
	if err := checkInput("oneWayLengthMeters", oneWayLengthMeters); err != nil {
		return 0, err
	}
	v := factor * currentAmperes * resistanceOhmsPerMeter * oneWayLengthMeters
	if math.IsInf(v, 0) {
		return 0, fmt.Errorf("%w: voltage drop", ErrOverflow)
	}
	return v, nil
}
