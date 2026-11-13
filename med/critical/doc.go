// Package critical computes acute-illness scores from their primary sources: the National Early Warning
// Score 2 (Royal College of Physicians 2017) and SCORTEN (Bastuji-Garin 2000).
//
// ID short form: CRIT (IDs are MED-CRIT-NNN).
//
// Units: every input names its unit (respirationRatePerMin, spo2Percent, systolicMmHg, pulsePerMin,
// temperatureC, ureaMmolL, glucoseMmolL, bicarbonateMmolL, ageYears, detachmentPercentBSA). An input outside
// its physical range returns ErrOutOfRange; NaN and ±Inf return ErrNotFinite. Category functions return the
// source's own label, verbatim, as a named string type.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-CRIT-001  NEWS2                       NEWS2 parameter scores and aggregate
//	MED-CRIT-002  NEWS2ClinicalRisk           NEWS2 clinical-risk band for an aggregate score and the red-score flag
//	MED-CRIT-003  SCORTEN                     SCORTEN (0–7), the count of its seven risk factors present
//	MED-CRIT-004  ATLASScore                  five-component ATLAS score for Clostridium difficile infection, 0–10
//	MED-CRIT-005  MASCCRiskIndex              MASCC risk-index score for febrile neutropenia
//	MED-CRIT-006  NeutrophilLymphocyteRatio   ratio of the neutrophil count to the lymphocyte count
//	MED-CRIT-007  PECARNFebrileInfantLowRisk  whether a febrile infant 60 days and younger meets the PECARN rule's low-risk
//	MED-CRIT-008  QSOFAScore                  quick SOFA (qSOFA) score, 0–3
//	MED-CRIT-009  RochesterLowRisk            whether a febrile infant 60 days of age or younger meets every Rochester
//	MED-CRIT-010  ShapiroBloodCultureRule     ShapiroBloodCultureRule applies the Shapiro decision rule for blood cultures in
//	MED-CRIT-011  SICScore                    sepsis-induced coagulopathy (SIC) score of Iba 2017 Table 3
//	MED-CRIT-012  CPOTScore           Critical-Care Pain Observation Tool score (0–8)
//	MED-CRIT-013  SIRSCriteria        the four SIRS criteria, their count and "two or more"
//	MED-CRIT-014  SOFAScore           SOFA organ scores and total, 0–24
//	MED-CRIT-015  MSOFAScore          modified SOFA (Grissom 2010) components and total, 0–19
//	MED-CRIT-016  SOFA2Score          SOFA-2 (Ranzani 2025) system scores and total, 0–24
//	MED-CRIT-017  PhoenixSepsisScore  Phoenix Sepsis Score of a child and the consensus criteria
//	MED-CRIT-018  REMSScore           Rapid Emergency Medicine Score, 0–26
//	MED-CRIT-019  LRINECScore         LRINEC score for necrotizing fasciitis, 0–13
//	MED-CRIT-020  SAPS3Score          SAPS 3 admission score and hospital mortality, bilirubin and creatinine in mg/dL
//	MED-CRIT-021  SAPS3ScoreSI        SAPS 3 admission score and hospital mortality, bilirubin and creatinine in µmol/L
package critical
