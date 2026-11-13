// Package cardio computes cardiovascular quantities from their primary sources.
//
// ID short form: CARDIO (IDs are MED-CARDIO-NNN).
//
// Units: every input and result names its unit (ageYears, qtMs, rrSeconds, SI units for the physics of
// MED-CARDIO-003). Every function enforces the input range its source validated and returns ErrOutOfRange
// outside it; NaN and ±Inf return ErrNotFinite.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-CARDIO-001  MaxHeartRateTanaka                 age-predicted maximal heart rate (Tanaka 2001)
//	MED-CARDIO-002  QTcFridericia                      QT interval corrected to RR = 1 s by the cube root (Fridericia 1920)
//	MED-CARDIO-003  PulseWaveVelocityMoensKorteweg     pressure-wave speed in an elastic tube (Korteweg 1878, eq. 10)
//	MED-CARDIO-004  TIMIScoreUANSTEMI                  TIMI risk score for UA/NSTEMI, 0–7 (Antman 2000)
//	MED-CARDIO-005  TIMIUANSTEMIEventRatePercent       14-day event rate per TIMI UA/NSTEMI score (Antman 2000)
//	MED-CARDIO-006  TIMIRiskIndex                      HR × (age/10)² / SBP for STEMI (Morrow 2001)
//	MED-CARDIO-007  MarburgHeartScore                  chest-pain CAD score in primary care, 0–5 (Bösner 2010)
//	MED-CARDIO-008  INTERCHESTScore                    chest-pain CAD score in primary care, −1 to 5 (Aerts 2017)
//	MED-CARDIO-009  (withdrawn before release: Canadian Syncope Risk Score point table not confirmed; ID not reused)
//	MED-CARDIO-010  OESILScore                         syncope score, 0–4 (Colivicchi 2003)
//	MED-CARDIO-011  OESILMortalityPercent              12-month mortality per OESIL score (Colivicchi 2003)
//	MED-CARDIO-012  RevisedCardiacRiskIndex            perioperative cardiac risk factors, 0–6 (Lee 1999)
//	MED-CARDIO-013  SubtleAnteriorSTEMI4Variable       4-variable LAD-occlusion formula (Driver 2017)
//	MED-CARDIO-014  SubtleAnteriorSTEMI3Variable       3-variable LAD-occlusion formula (Smith 2012)
//	MED-CARDIO-015  ShanghaiBrugadaScore               Shanghai score for Brugada syndrome (Antzelevitch 2016)
//	MED-CARDIO-016  TisdaleQTRiskScore                 QTc-prolongation risk score, 0–21 (Tisdale 2013)
//	MED-CARDIO-017  TisdaleQTRiskGroup                 Tisdale Low / Moderate / High
//	MED-CARDIO-018  LACEIndex                          30-day death or readmission index, 0–19 (van Walraven 2010)
//	MED-CARDIO-019  CHADS2Score                        CHADS₂, 0–6 (Gage 2001)
//	MED-CARDIO-020  CHADS2StrokeRatePer100PersonYears  stroke rate per CHADS₂ score (Gage 2001)
//	MED-CARDIO-021  ATRIABleedingScore                 warfarin hemorrhage score, 0–10 (Fang 2011)
//	MED-CARDIO-022  ATRIABleedingGroup                 ATRIA bleeding low / intermediate / high risk
//	MED-CARDIO-023  ATRIAStrokeScore                   AF stroke score, 0–15 (Singer 2013)
//	MED-CARDIO-024  ATRIAStrokeGroup                   ATRIA stroke low / moderate / high
//	MED-CARDIO-025  ORBITBleedingScore                 AF bleeding score on OAC, 0–7 (O'Brien 2015)
//	MED-CARDIO-026  ORBITBleedingGroup                 ORBIT low / medium / high
//	MED-CARDIO-027  DAPTScore                          DAPT score, −2 to 10 (Yeh 2016)
//	MED-CARDIO-028  HeartRateRecovery1Min              peak minus 1-minute heart rate (Cole 1999)
//	MED-CARDIO-029  MaxHeartRateGulati                 206 − 0.88 × age, women (Gulati 2010)
//	MED-CARDIO-030  MaxHeartRateNes                    211 − 0.64 × age (Nes 2013, HUNT)
//	MED-CARDIO-031  CardiacPowerOutputW                MAP × CO / 451, watts (Fincke 2004)
//	MED-CARDIO-032  DiastolicShockIndex                HR / diastolic pressure (Ospina-Tascón 2020)
//	MED-CARDIO-033  RespiratoryAdjustedShockIndex      HR / SBP × RR / 10 (Caputo 2018)
//	MED-CARDIO-034  ACEFII                             ACEF II cardiac-surgery score (Ranucci 2018)
//	MED-CARDIO-035  IEMortalityScore                   6-month mortality score in infective endocarditis (Park 2016)
//	MED-CARDIO-036  GillmoreATTRStage                  ATTR-CM stage I–III (Gillmore 2018)
//	MED-CARDIO-037  SIPACutoff                         age-adjusted pediatric shock-index cut-off (Acker 2015)
//	MED-CARDIO-038  DutchLipidClinicFHScore            Dutch Lipid Clinic Network FH score, 0–26 (Nordestgaard 2013)
//	MED-CARDIO-039  DutchLipidClinicFHGroup            definite / probable / possible / unlikely FH
//	MED-CARDIO-040  WellsPEScore                       Wells PE score, 0–12.5 (Wells 2000)
//	MED-CARDIO-041  GenevaScoreRevised                 revised Geneva PE score, 0–22 (Le Gal 2006)
//	MED-CARDIO-042  DDimerAgeAdjustedCutoffUgL         age × 10 µg/L above 50 years (Douma 2010)
//	MED-CARDIO-043  HERDOO2Count                       HERDOO2 characteristics, 0–4 (Rodger 2008)
//	MED-CARDIO-044  IMPEDEVTEScore                     IMPEDE VTE score in myeloma (Sanfilippo 2019)
//	MED-CARDIO-045  TRiPCastScore                      TRiP(cast) VTE score after leg immobilisation (Nemeth 2020)
//	MED-CARDIO-046  PulsePressureVariationPercent      respiratory variation of the arterial pulse pressure
//	MED-CARDIO-047  CARTScore                       Cardiac Arrest Risk Triage score, 0–57 (Churpek 2012)
//	MED-CARDIO-048  CRUSADEBleedingScore            CRUSADE bleeding score, 1–96 (Subherwal 2009)
//	MED-CARDIO-049  CRUSADEBleedingGroupOf          CRUSADE Very Low / Low / Moderate / High / Very High
//	MED-CARDIO-050  READMITSScore                   AMI READMITS 30-day readmission score, 0–33 (Nguyen 2018)
//	MED-CARDIO-051  READMITSRiskGroupOf             READMITS extremely low … extremely high
//	MED-CARDIO-052  RESCUEIHCAScore                 RESCUE-IHCA score for ECPR in-hospital arrest (Tonna 2022)
//	MED-CARDIO-053  ARCHBRHighBleedingRisk          ARC-HBR: ≥1 major or ≥2 minor criteria (Urban 2019)
//	MED-CARDIO-054  DukeISCVIDCategoryOf            2023 Duke-ISCVID Definite / Possible / Rejected IE (Fowler 2023)
//	MED-CARDIO-055  BostonSyncopeRuleCount          Boston Syncope Rule groups present, 0–8 (Grossman 2007)
//	MED-CARDIO-056  ModifiedShockIndex              heart rate / mean arterial pressure (Liu 2012)
//	MED-CARDIO-057  HEARTScore                      HEART score, 0–10, HEART Pathway bands (Mahler 2015)
//	MED-CARDIO-058  HEARTPathwayLowRisk             HEART 0–3 and negative 0 h and 3 h troponin (Mahler 2015)
//	MED-CARDIO-059  HEARScore                       HEART without troponin, 0–8 (Smith 2020)
//	MED-CARDIO-060  PCPHF10YearRiskPercent          10-year heart-failure risk, PCP-HF equations (Khan 2019)
//	MED-CARDIO-062  YEARSPulmonaryEmbolismExcluded  YEARS items + D-dimer 1000/500 ng/mL rule (van der Hulle 2017)
//	MED-CARDIO-063  SanFranciscoSyncopeRuleCount    San Francisco Syncope Rule items present, 0–5 (Quinn 2004)
//	MED-CARDIO-064  FourPEPSScore                   4PEPS score (−4 to 22), Roy 2021 Table 3
//	MED-CARDIO-065  FourPEPSLevel                   4PEPS clinical pretest probability level, the source's label
package cardio
