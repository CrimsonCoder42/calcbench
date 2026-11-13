// Package chemistry holds clinical-chemistry calculators: albumin-adjusted calcium (Payne), calculated serum
// osmolality (Dorwart–Chalmers), estimated average glucose from HbA1c (ADAG), the albumin-adjusted anion gap
// (Figge) and LDL cholesterol by the Martin–Hopkins method.
//
// ID short form: CHEM (IDs are MED-CHEM-NNN).
//
// Units: every input and result names its unit. Where clinicians use both systems there is one function
// per system; the SI twin applies its documented factor once at entry (calcium 0.2495 mmol/L per mg/dL,
// KDIGO 2024, Kidney Int. 2024;105(4S):S117–S314, p. S126; albumin 10 g/L per g/dL by definition; lipids the
// Martin et al. 2013 Figure 2 caption factors 0.0259 and 0.0113) and then applies the same equation. NaN or
// ±Inf in any input returns ErrNotFinite before any range check; an input outside the range the source
// validated returns ErrOutOfRange; finite in-range inputs whose result is beyond float64 return ErrDomain.
// Nothing is rounded.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-CHEM-001  CorrectedCalciumPayne          albumin-adjusted calcium, mg/dL (Payne 1973)
//	MED-CHEM-002  CorrectedCalciumPayneSI        albumin-adjusted calcium, mmol/L from mmol/L and g/L
//	MED-CHEM-003  SerumOsmolalityDorwart         calculated serum osmolality, mOsm/kg, glucose and BUN in mg/dL
//	MED-CHEM-004  SerumOsmolalityDorwartSI       calculated serum osmolality, mOsm/kg, glucose and urea in mmol/L
//	MED-CHEM-005  EstimatedAverageGlucose        estimated average glucose from HbA1c, mg/dL (ADAG)
//	MED-CHEM-006  EstimatedAverageGlucoseMmolL   estimated average glucose from HbA1c, mmol/L (ADAG)
//	MED-CHEM-007  AdjustedAnionGapFigge          albumin-adjusted anion gap, albumin in g/L (Figge 1998)
//	MED-CHEM-008  AdjustedAnionGapFiggeGDL       albumin-adjusted anion gap, albumin in g/dL (Figge 1998)
//	MED-CHEM-009  LDLCholesterolMartinHopkins    LDL cholesterol, mg/dL, Martin–Hopkins 180-cell factor
//	MED-CHEM-010  LDLCholesterolMartinHopkinsSI  LDL cholesterol, mmol/L, Martin–Hopkins 180-cell factor
//	MED-CHEM-011  AnionGap                      serum anion gap Na − (Cl + HCO3), mEq/L, Kitabchi 2009
//	MED-CHEM-012  EffectiveSerumOsmolality      effective osmolality 2Na + glucose/18, mOsm/kg, Kitabchi 2009
//	MED-CHEM-013  CorrectedSodiumHyperglycemia  Na + 1.6 per 100 mg/dL glucose above 100, mEq/L, Kitabchi 2009
package chemistry
