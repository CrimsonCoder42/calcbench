package electrical

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or infinite input.
	ErrNotFinite = errors.New("electrical: value is NaN or infinite")
	// ErrNegative reports a negative input.
	ErrNegative = errors.New("electrical: value is negative")
	// ErrZero reports a zero divisor.
	ErrZero = errors.New("electrical: divisor is zero")
	// ErrOverflow reports a result beyond float64.
	ErrOverflow = errors.New("electrical: result overflows float64")
)

// Circuit is the four quantities of a resistive element: volts across it, amperes through it, its resistance
// in ohms and the power it dissipates in watts.
type Circuit struct {
	Volts, Amperes, Ohms, Watts float64
}

func checkInput(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, v)
	}
	return nil
}

func checkPair(an string, a float64, bn string, b float64) error {
	if err := checkInput(an, a); err != nil {
		return err
	}
	return checkInput(bn, b)
}

func finish(c Circuit) (Circuit, error) {
	for _, v := range [4]float64{c.Volts, c.Amperes, c.Ohms, c.Watts} {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return Circuit{}, ErrOverflow
		}
	}
	return c, nil
}

// OhmsLawFromVoltageCurrent returns the circuit with the given voltage and current: R = V/I, P = V·I.
// Errors: ErrNotFinite; ErrNegative; ErrZero for currentAmperes = 0; ErrOverflow.
//
// ID:      ENG-ELEC-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.4 "Ohm's Law" (V = IR) and §9.5 "Electrical Energy and Power" (P = IV = I²R = V²/R)
// Tags:    ohms-law, voltage, current, resistance, power, circuit, electrical
// Units:   voltageVolts V; currentAmperes A; result V, A, Ω, W
// Example: 12 V, 2 A → 6 Ω, 24 W
func OhmsLawFromVoltageCurrent(voltageVolts, currentAmperes float64) (Circuit, error) {
	if err := checkPair("voltageVolts", voltageVolts, "currentAmperes", currentAmperes); err != nil {
		return Circuit{}, err
	}
	if currentAmperes == 0 {
		return Circuit{}, fmt.Errorf("%w: currentAmperes = 0", ErrZero)
	}
	return finish(Circuit{voltageVolts, currentAmperes, voltageVolts / currentAmperes, voltageVolts * currentAmperes})
}

// OhmsLawFromVoltageResistance returns the circuit with the given voltage and resistance: I = V/R, P = V²/R.
// Errors: ErrNotFinite; ErrNegative; ErrZero for resistanceOhms = 0; ErrOverflow.
//
// ID:      ENG-ELEC-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.4 "Ohm's Law" (V = IR) and §9.5 "Electrical Energy and Power" (P = V²/R)
// Tags:    ohms-law, voltage, current, resistance, power, circuit, electrical
// Units:   voltageVolts V; resistanceOhms Ω; result V, A, Ω, W
// Example: 12 V, 6 Ω → 2 A, 24 W
func OhmsLawFromVoltageResistance(voltageVolts, resistanceOhms float64) (Circuit, error) {
	if err := checkPair("voltageVolts", voltageVolts, "resistanceOhms", resistanceOhms); err != nil {
		return Circuit{}, err
	}
	if resistanceOhms == 0 {
		return Circuit{}, fmt.Errorf("%w: resistanceOhms = 0", ErrZero)
	}
	i := voltageVolts / resistanceOhms
	return finish(Circuit{voltageVolts, i, resistanceOhms, voltageVolts * i})
}

// OhmsLawFromVoltagePower returns the circuit with the given voltage and power: I = P/V, R = V²/P.
// Errors: ErrNotFinite; ErrNegative; ErrZero for voltageVolts = 0 or powerWatts = 0; ErrOverflow.
//
// ID:      ENG-ELEC-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.5 "Electrical Energy and Power" (P = IV = V²/R)
// Tags:    ohms-law, voltage, current, resistance, power, circuit, electrical
// Units:   voltageVolts V; powerWatts W; result V, A, Ω, W
// Example: 120 V, 60 W → 0.5 A, 240 Ω
func OhmsLawFromVoltagePower(voltageVolts, powerWatts float64) (Circuit, error) {
	if err := checkPair("voltageVolts", voltageVolts, "powerWatts", powerWatts); err != nil {
		return Circuit{}, err
	}
	if voltageVolts == 0 {
		return Circuit{}, fmt.Errorf("%w: voltageVolts = 0", ErrZero)
	}
	if powerWatts == 0 {
		return Circuit{}, fmt.Errorf("%w: powerWatts = 0 (resistance unbounded)", ErrZero)
	}
	i := powerWatts / voltageVolts
	return finish(Circuit{voltageVolts, i, voltageVolts / i, powerWatts})
}

// OhmsLawFromCurrentResistance returns the circuit with the given current and resistance: V = I·R, P = I²·R.
// Errors: ErrNotFinite; ErrNegative; ErrOverflow.
//
// ID:      ENG-ELEC-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.4 "Ohm's Law" (V = IR) and §9.5 "Electrical Energy and Power" (P = I²R)
// Tags:    ohms-law, voltage, current, resistance, power, circuit, electrical
// Units:   currentAmperes A; resistanceOhms Ω; result V, A, Ω, W
// Example: 2 A, 6 Ω → 12 V, 24 W
func OhmsLawFromCurrentResistance(currentAmperes, resistanceOhms float64) (Circuit, error) {
	if err := checkPair("currentAmperes", currentAmperes, "resistanceOhms", resistanceOhms); err != nil {
		return Circuit{}, err
	}
	v := currentAmperes * resistanceOhms
	return finish(Circuit{v, currentAmperes, resistanceOhms, v * currentAmperes})
}

// OhmsLawFromCurrentPower returns the circuit with the given current and power: V = P/I, R = P/I².
// Errors: ErrNotFinite; ErrNegative; ErrZero for currentAmperes = 0; ErrOverflow.
//
// ID:      ENG-ELEC-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.5 "Electrical Energy and Power" (P = IV = I²R)
// Tags:    ohms-law, voltage, current, resistance, power, circuit, electrical
// Units:   currentAmperes A; powerWatts W; result V, A, Ω, W
// Example: 2 A, 24 W → 12 V, 6 Ω
func OhmsLawFromCurrentPower(currentAmperes, powerWatts float64) (Circuit, error) {
	if err := checkPair("currentAmperes", currentAmperes, "powerWatts", powerWatts); err != nil {
		return Circuit{}, err
	}
	if currentAmperes == 0 {
		return Circuit{}, fmt.Errorf("%w: currentAmperes = 0", ErrZero)
	}
	v := powerWatts / currentAmperes
	return finish(Circuit{v, currentAmperes, v / currentAmperes, powerWatts})
}

// OhmsLawFromResistancePower returns the circuit with the given resistance and power: V = √(P·R), I = √(P/R).
// Errors: ErrNotFinite; ErrNegative; ErrZero for resistanceOhms = 0; ErrOverflow.
//
// ID:      ENG-ELEC-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.5 "Electrical Energy and Power" (P = I²R = V²/R)
// Tags:    ohms-law, voltage, current, resistance, power, circuit, electrical
// Units:   resistanceOhms Ω; powerWatts W; result V, A, Ω, W
// Example: 6 Ω, 24 W → 12 V, 2 A
func OhmsLawFromResistancePower(resistanceOhms, powerWatts float64) (Circuit, error) {
	if err := checkPair("resistanceOhms", resistanceOhms, "powerWatts", powerWatts); err != nil {
		return Circuit{}, err
	}
	if resistanceOhms == 0 {
		return Circuit{}, fmt.Errorf("%w: resistanceOhms = 0", ErrZero)
	}
	// √P·√R avoids overflow of P·R for large inputs.
	return finish(Circuit{math.Sqrt(powerWatts) * math.Sqrt(resistanceOhms), math.Sqrt(powerWatts / resistanceOhms), resistanceOhms, powerWatts})
}
