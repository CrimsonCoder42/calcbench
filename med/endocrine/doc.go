// Package endocrine holds diabetes and endocrine calculators: the glucose management indicator, the BeAM value,
// the C-peptide to glucose ratio, the IADPSG gestational-diabetes thresholds, the Efstathiou diabetic
// ketoacidosis mortality points, the AUSDRISK type 2 diabetes points and the Diabetes Distress Scale score and
// its Fisher 2012 categories.
//
// ID short form: ENDO (IDs are MED-ENDO-NNN).
//
// Units: every input and result names its unit, and every function takes the unit its source uses (glucose
// mg/dL or mmol/L as the source states it, C-peptide ng/mL, waist cm, pH units, insulin IU); where a source
// prints its thresholds in both systems (IADPSG) each system has its own function with the source's own
// numbers, and no conversion is applied. Ages are the caller's completed age in decimal years. NaN or ±Inf in
// any input returns ErrNotFinite before any range check; an input outside the range the source validated, or
// an unknown level of a category input, returns ErrOutOfRange; finite in-range inputs whose result is beyond
// float64 return ErrDomain. Nothing is rounded. Questionnaire item wording is not embedded.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-ENDO-001  GlucoseManagementIndicatorPercent  GMI (%) from mean CGM glucose in mg/dL, Bergenstal 2018
//	MED-ENDO-002  GlucoseManagementIndicatorMmolMol  GMI (mmol/mol) from mean CGM glucose in mmol/L, Bergenstal 2018
//	MED-ENDO-003  BeAMValue                          bedtime minus pre-breakfast glucose (mg/dL), Zisman 2016
//	MED-ENDO-004  CPeptideGlucoseRatio               C-peptide (ng/mL) / glucose (mg/dL) × 100, Saisho 2016
//	MED-ENDO-005  GDMThresholdIADPSGMgDL             whether a 75-g OGTT value meets an IADPSG threshold (mg/dL)
//	MED-ENDO-006  GDMThresholdIADPSGMmolL            whether a 75-g OGTT value meets an IADPSG threshold (mmol/L)
//	MED-ENDO-007  DKAMortalityScore                  DKA mortality prediction model points (0–25), Efstathiou 2002
//	MED-ENDO-008  AUSDRISKScore                      AUSDRISK points (0–35), Chen 2010 Box 3
//	MED-ENDO-009  DiabetesDistressScore17            DDS17 mean-item score (1–6) from the 17 ratings
//	MED-ENDO-010  DiabetesDistressCategoryFisher     Fisher 2012 DDS17 category, the source's own label
//	MED-ENDO-011  TimeInRangesMgDL                   percent of glucose readings (mg/dL) in each consensus range
//	MED-ENDO-012  TimeInRangesMmolL                  percent of glucose readings (mmol/L) in each consensus range
//	MED-ENDO-013  A1CCategoryADAPercent              ADA category of a hemoglobin A1C in percent (NGSP)
//	MED-ENDO-014  A1CCategoryADAMmolMol              ADA category of a hemoglobin A1C in mmol/mol
//	MED-ENDO-015  FastingGlucoseCategoryADAMgDL      ADA category of a fasting plasma glucose (FPG) in mg/dL
//	MED-ENDO-016  FastingGlucoseCategoryADAMmolL     ADA category of a fasting plasma glucose in mmol/L
//	MED-ENDO-017  OGTTTwoHourCategoryADAMgDL         ADA category of the 2-h plasma glucose
//	MED-ENDO-018  OGTTTwoHourCategoryADAMmolL        ADA category of the 2-h plasma glucose during a 75-g OGTT in mmol/L
//	MED-ENDO-019  RandomGlucoseCriterionADAMgDL      whether a random plasma glucose meets the fourth ADA Table 2.1 criterion
//	MED-ENDO-020  RandomGlucoseCriterionADAMmolL     whether a random plasma glucose in mmol/L meets the fourth ADA Table 2.1
//	MED-ENDO-021  HypoglycemiaLevelADAMgDL           ADA hypoglycemia level of a glucose in mg/dL
//	MED-ENDO-022  HypoglycemiaLevelADAMmolL          ADA hypoglycemia level of a glucose in mmol/L
//	MED-ENDO-023  CGMTargetsMetBattelino             which of the 2019 International Consensus on Time in Range targets for people
//	MED-ENDO-024  HOMAIR             HOMA-IR, insulin (µU/mL) × glucose (mmol/L) / 22.5
//	MED-ENDO-025  HOMABeta           HOMA-%B, 20 × insulin / (glucose mmol/L − 3.5), percent
//	MED-ENDO-026  QUICKI             1/(log10 insulin µU/mL + log10 glucose mg/dL)
//	MED-ENDO-027  METSIR             METS-IR, ln(2G + TG) × BMI / ln(HDL), mg/dL
//	MED-ENDO-028  MyxedemaComaScore  myxedema coma diagnostic score (0–230), Popoveniuc 2014
//	MED-ENDO-029  OSTAIndex          OSTA index, trunc(0.2 × (kg − years)), Koh 2001
//	MED-ENDO-030  OSTAIndexImperial  OSTA index from pounds
//	MED-ENDO-031  OSTARiskCategory   OSTA risk category, the source's label
package endocrine
