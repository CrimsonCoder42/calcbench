// Package anthropometry holds body-size calculators: body mass index (metric and the CDC's imperial form)
// and the CDC adult BMI category, Haycock body surface area, adjusted body weight, the Traub–Kichen ideal
// body mass of children, the Best Guess weight of infants, and the Deurenberg body-fat percentage of adults
// and children.
//
// ID short form: ANTHRO (IDs are MED-ANTHRO-NNN).
//
// Units: every input and result names its unit, and every function takes the unit its source uses (kg and
// m for BMI, lb and in for the CDC imperial BMI, kg and cm for Haycock, months for the infant weight); no
// unit is converted here. Ages are the caller's completed age in decimal years (ageYears) or months
// (ageMonths). NaN or ±Inf in any input returns ErrNotFinite before any range check; an input outside the
// range the source validated returns ErrOutOfRange; finite in-range inputs whose result is beyond float64
// return ErrDomain. Nothing is rounded.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-ANTHRO-001  BodyMassIndex                     body mass index (kg/m²) from kg and m
//	MED-ANTHRO-002  BodyMassIndexImperial             body mass index (kg/m²) from lb and in, the CDC's × 703 form
//	MED-ANTHRO-003  BMICategoryCDCAdult               CDC adult BMI category, the CDC's own label
//	MED-ANTHRO-004  BodySurfaceAreaHaycock            body surface area (m²), Haycock 1978
//	MED-ANTHRO-005  AdjustedBodyWeight                ideal weight plus 40% of the excess (kg), Schwartz 1978
//	MED-ANTHRO-006  IdealBodyMassTraubKichen          ideal body mass of a child aged 1–17 (kg), Traub & Kichen 1983
//	MED-ANTHRO-007  EstimatedWeightInfantBestGuess    estimated weight of an infant under 12 months (kg), Best Guess
//	MED-ANTHRO-008  BodyFatPercentDeurenbergAdult     body-fat percentage from BMI, age and sex, Deurenberg adult formula
//	MED-ANTHRO-009  BodyFatPercentDeurenbergChild     body-fat percentage from BMI, age and sex, Deurenberg child formula
//	MED-ANTHRO-010  BasalMetabolicRateHarrisBenedict  predicted basal heat production
//	MED-ANTHRO-011  RelativeFatMass                   relative fat mass (RFM), an estimate of whole-body fat percentage
//	MED-ANTHRO-012  MUSTScore                         'Malnutrition Universal Screening Tool' ('MUST') total, 0–6
//	MED-ANTHRO-013  MUSTRiskCategory                  'MUST' step 4 overall risk of malnutrition for a total
//	MED-ANTHRO-014  BodyRoundnessIndex          body roundness index from waist and height (cm), Thomas 2013
//	MED-ANTHRO-015  BodyRoundnessIndexImperial  body roundness index from waist and height (in)
//	MED-ANTHRO-022  WeightAtBMI                 body weight (kg) at a BMI for a height (m), the CDC definition solved for weight
//	MED-ANTHRO-023  WeightAtBMIImperial         body weight (lb) at a BMI for a height (in), the CDC × 703 form solved for weight
package anthropometry
