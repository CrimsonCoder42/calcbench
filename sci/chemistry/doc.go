// Package chemistry holds the stoichiometry of solutions: amount of substance from mass (n = m/M) and amount
// concentration, molarity (c = n/V).
//
// ID short form: CHEM (IDs are SCI-CHEM-NNN).
//
// Units: mass in grams, molar mass in g/mol, amount in mol, solution volume in liters, molarity in mol/L (M).
// Every parameter names its unit. NaN or ±Inf returns ErrNotFinite, a negative quantity ErrNegative, a zero
// divisor ErrZero, a result beyond float64 ErrOverflow. Results are not rounded.
//
// Functions:
//
//	SCI-CHEM-001  AmountFromMass     n = m / M
//	SCI-CHEM-002  Molarity           c = n / V
//	SCI-CHEM-003  MolarityFromMass   c = m / (M · V)
//	SCI-CHEM-004  MassForMolarity    m = c · V · M
package chemistry
