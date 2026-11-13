// Package electrical holds circuit relations: Ohm's law with power (any two of V, I, R, P give the other two),
// resistors in series and parallel, the resistance of a wire, and the resistive voltage drop along a feeder.
//
// ID short form: ELEC (IDs are ENG-ELEC-NNN).
//
// Units: volts, amperes, ohms, watts, meters and ohms per meter, named in every parameter. NaN or ±Inf returns
// ErrNotFinite, a negative input ErrNegative, a zero divisor ErrZero, a result beyond float64 ErrOverflow.
// Quantities are magnitudes (DC, or AC RMS with a purely resistive load). Results are not rounded.
//
// Functions:
//
//	ENG-ELEC-001  OhmsLawFromVoltageCurrent      R = V/I, P = V·I
//	ENG-ELEC-002  OhmsLawFromVoltageResistance   I = V/R, P = V²/R
//	ENG-ELEC-003  OhmsLawFromVoltagePower        I = P/V, R = V²/P
//	ENG-ELEC-004  OhmsLawFromCurrentResistance   V = I·R, P = I²·R
//	ENG-ELEC-005  OhmsLawFromCurrentPower        V = P/I, R = P/I²
//	ENG-ELEC-006  OhmsLawFromResistancePower     V = √(P·R), I = √(P/R)
//	ENG-ELEC-007  VoltageDropTwoWire             ΔV = I · R' · 2L (DC or single-phase, resistive)
//	ENG-ELEC-008  VoltageDropThreePhase          ΔV(line-to-line) = √3 · I · R' · L (balanced, resistive)
//	ENG-ELEC-009  SeriesResistance               R = R1 + … + Rn (compensated sum)
//	ENG-ELEC-010  ParallelResistance             1/R = 1/R1 + … + 1/Rn (compensated sum)
//	ENG-ELEC-011  WireResistance                 R = L / (σ · πd²/4)
package electrical
