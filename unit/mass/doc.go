// Package mass converts masses between the kilogram and the pound (avoirdupois), by the exact relation of
// the NIST Guide to the SI (NIST SP 811, 2008, Appendix B.8): 1 lb = 0.453 592 37 kg exactly.
//
// ID short form: MASS (IDs are UNT-MASS-NNN).
//
// Units: every input and result names its unit (kilograms, pounds). A negative mass returns ErrNegative;
// NaN and ±Inf return ErrNotFinite; a finite input whose result is beyond float64 returns ErrOverflow.
// Results are not rounded.
//
// Functions:
//
//	UNT-MASS-001  KilogramsToPounds  m/lb = (m/kg)/0.453 592 37
//	UNT-MASS-002  PoundsToKilograms  m/kg = 0.453 592 37 · m/lb
package mass
