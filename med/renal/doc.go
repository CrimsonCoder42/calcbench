// Package renal holds kidney-function calculators: estimated GFR (CKD-EPI 2021 creatinine and
// creatinine–cystatin C, bedside Schwartz for children), Cockcroft–Gault creatinine clearance, the
// fractional excretions of sodium and urea, and the KDIGO GFR category.
//
// ID short form: RENAL (IDs are MED-RENAL-NNN).
//
// Units: every input and result names its unit. Creatinine has one function per system: the mg/dL form
// and an SI twin taking µmol/L, which divides by 88.4 µmol/L per mg/dL (KDIGO 2024, Kidney Int.
// 2024;105(4S):S117–S314, p. S126) once at entry and then applies the same equation. Ages are the
// caller's completed age in decimal years. NaN or ±Inf in any input returns ErrNotFinite before any range
// check; an input outside the range the source validated returns ErrOutOfRange; finite in-range inputs whose
// result is beyond float64 return ErrDomain. Nothing is rounded.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-RENAL-001  EGFRCKDEPI2021                       eGFR, 2021 CKD-EPI creatinine, creatinine in mg/dL
//	MED-RENAL-002  EGFRCKDEPI2021SI                     eGFR, 2021 CKD-EPI creatinine, creatinine in µmol/L
//	MED-RENAL-003  EGFRCKDEPI2021CystatinC              eGFR, 2021 CKD-EPI creatinine–cystatin C, creatinine in mg/dL
//	MED-RENAL-004  EGFRCKDEPI2021CystatinCSI            eGFR, 2021 CKD-EPI creatinine–cystatin C, creatinine in µmol/L
//	MED-RENAL-005  CreatinineClearanceCockcroftGault    creatinine clearance (mL/min), Cockcroft–Gault, creatinine in mg/dL
//	MED-RENAL-006  CreatinineClearanceCockcroftGaultSI  creatinine clearance (mL/min), Cockcroft–Gault, creatinine in µmol/L
//	MED-RENAL-007  EGFRSchwartzBedside                  eGFR of a child, bedside Schwartz, creatinine in mg/dL
//	MED-RENAL-008  EGFRSchwartzBedsideSI                eGFR of a child, bedside Schwartz, creatinine in µmol/L
//	MED-RENAL-009  FractionalExcretionSodium            fractional excretion of sodium (FENa), percent
//	MED-RENAL-010  FractionalExcretionUrea              fractional excretion of urea (FEUrea), percent
//	MED-RENAL-011  GFRCategoryKDIGO                     KDIGO GFR category G1–G5 with the source's term
//	MED-RENAL-012  AKINStage                            Acute Kidney Injury Network stage, 0–3 (0
//	MED-RENAL-013  STONEScore                           STONE score, 0–13, from Moore 2014 Table 3
//	MED-RENAL-014  STONEProbabilityGroup                STONE-score group of Moore 2014
//	MED-RENAL-015  UrineProteinCreatinineRatio          protein/creatinine ratio of a single voided urine sample
//	MED-RENAL-016  UrineAnionGap                        urine anion gap, Na⁺ + K⁺ − Cl⁻, in mmol/L,
//	MED-RENAL-017  EGFRCKiDU25Creatinine    eGFR of a child or young adult, CKiD U25 height/creatinine, creatinine in mg/dL
//	MED-RENAL-018  EGFRCKiDU25CreatinineSI  eGFR, CKiD U25 height/creatinine, creatinine in µmol/L
//	MED-RENAL-019  EGFRCKiDU25CystatinC     eGFR of a child or young adult, CKiD U25 cystatin C
//	MED-RENAL-020  McMahonScore             McMahon rhabdomyolysis risk score, 0–19
//	MED-RENAL-021  HIVCKDRiskScore          Scherzer 2014 5-year CKD risk score in HIV, 0–15
//	MED-RENAL-022  EGFRMDRD                 eGFR, 4-variable MDRD (IDMS, 175), creatinine in mg/dL
//	MED-RENAL-023  EGFRMDRDSI               eGFR, 4-variable MDRD (IDMS, 175), creatinine in µmol/L
//	MED-RENAL-024  KtV                      dialysis dose Kt/V from clearance, time and urea volume
//	MED-RENAL-025  RIFLEClassify            ADQI RIFLE class from creatinine and urine output, creatinine in mg/dL
//	MED-RENAL-026  RIFLEClassifySI          ADQI RIFLE class, creatinine in µmol/L
//	MED-RENAL-027  ThakarScore              Thakar (Cleveland Clinic) acute renal failure after cardiac surgery score (0–17)
package renal
