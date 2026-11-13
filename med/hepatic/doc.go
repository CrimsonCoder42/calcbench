// Package hepatic holds liver calculators: the FIB-4 index and the MELD family (MELD, MELD-Na, MELD 3.0).
//
// ID short form: HEPAT (IDs are MED-HEPAT-NNN).
//
// Units: every input and result names its unit; laboratory values are in the units the sources state
// (bilirubin and creatinine mg/dL, albumin g/dL, sodium mmol/L, AST and ALT U/L, platelets 10^9/L, INR a
// ratio). The MELD sources state bounds on their inputs (values below 1 set to 1, a creatinine cap, sodium
// 125–137, albumin 1.5–3.5); those are applied as clamps, as the sources say, and are not errors. NaN or ±Inf
// in any input returns ErrNotFinite before any range check; a non-positive input returns ErrOutOfRange;
// finite in-range inputs whose result is beyond float64 return ErrDomain. Only MELD 3.0 is rounded, as its
// source says.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-HEPAT-001  FIB4       FIB-4 index from age, AST, ALT and platelets
//	MED-HEPAT-002  MELD       MELD score (not rounded)
//	MED-HEPAT-003  MELDNa     MELD-Na score (not rounded)
//	MED-HEPAT-004  MELD3      MELD 3.0 score, rounded to the nearest integer
//	MED-HEPAT-005  BARDScore  BARD score, 0–4, for nonalcoholic fatty liver disease
//	MED-HEPAT-006  GlasgowAlcoholicHepatitisScore  GAHS (5–12), Forrest 2005 Table 4, SI units
//	MED-HEPAT-007  MaddreyDiscriminantFunction     4.6 × (PT − control PT) + bilirubin mg/dL
//	MED-HEPAT-008  NAFLDActivityScore              NAS (0–8), steatosis + lobular inflammation + ballooning
package hepatic
