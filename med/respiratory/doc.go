// Package respiratory computes respiratory-physiology quantities from their primary sources.
//
// ID short form: RESP (IDs are MED-RESP-NNN).
//
// Units: every input and result names its unit. Every function enforces the input range its source
// validated and returns ErrOutOfRange outside it; NaN and ±Inf return ErrNotFinite.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-RESP-001  DiffusingCapacityRoughtonForster      lung diffusing capacity from DM, θ and VC (Roughton–Forster)
//	MED-RESP-002  SixMinuteWalkDistanceEnright          reference six-minute walk distance
//	MED-RESP-003  SixMinuteWalkDistanceEnrightImperial  reference six-minute walk distance
//	MED-RESP-004  AIRQScore                             Asthma Impairment and Risk Questionnaire total (0–10)
//	MED-RESP-005  AIRQControlLevel                      AIRQ control level for a total score
//	MED-RESP-006  DeltaPScore                           Dutch-English LEMS Tumor Association Prediction (DELTA-P) score (0–6)
//	MED-RESP-007  PackYears                             cumulative smoking exposure in pack-years
//	MED-RESP-008  PackYearsFromCigarettes               pack-years from cigarettes per day
//	MED-RESP-009  NEWS2SpO2Scale1Points                 NEWS2 points (0–3) for an oxygen saturation on SpO₂ Scale 1
//	MED-RESP-010  NEWS2SpO2Scale2Points                 NEWS2 points (0–3) for an oxygen saturation on SpO₂ Scale 2
//	MED-RESP-011  CURB65Score                           CURB-65 score (0–5) for community-acquired pneumonia
//	MED-RESP-012  ROXIndex                              ROX index: (SpO₂ / FiO₂) / respiratory rate
//	MED-RESP-013  PaO2FiO2Ratio                         PaO₂/FiO₂ ratio, in mmHg: arterial oxygen tension divided
//	MED-RESP-014  BerlinARDSOxygenationCategory         Berlin Definition's oxygenation category for a PaO₂/FiO₂ ratio
//	MED-RESP-015  SpO2FiO2Ratio                         SpO₂/FiO₂ (S/F) ratio: oxygen saturation by pulse oximetry
//	MED-RESP-016  SpO2FiO2FromPaO2FiO2Rice              S/F ratio that the Rice 2007 regression S/F = 64 + 0.84 ×
//	MED-RESP-017  RapidShallowBreathingIndex            rapid shallow breathing index f/VT, in breaths/min/L
//	MED-RESP-018  OxygenationIndex                      oxygenation index OI = mean airway pressure × FiO₂ / PaO₂
//	MED-RESP-019  OxygenSaturationIndex                 oxygen saturation index OSI = mean airway pressure × FiO₂ / SpO₂
//	MED-RESP-020  MuLBSTAScore                          MuLBSTA score (0–20) for viral pneumonia
//	MED-RESP-021  PESARDScore                           PE-SARD bleeding score (0–5 points, in steps of 0.5)
//	MED-RESP-022  PESARDRiskGroup                       PE-SARD bleeding-risk group for a score: 0 low, 1–2.5 intermediate, > 2.5 high
//	MED-RESP-023  PIROScore                             PIRO score for community-acquired pneumonia (0–8)
//	MED-RESP-024  PIRORiskLevel                         PIRO risk level for a score: 0–2 Low, 3 Mild, 4 high, 5–8 Very high
//	MED-RESP-025  SCARFScore                            Sequential Clinical Assessment of Respiratory Function score
//	MED-RESP-026  PeakFlowZoneNHLBI                     NHLBI Asthma Action Plan zone for a peak expiratory flow measured against the
//	MED-RESP-027  CHOSENScore                         CHOSEN COVID-19 risk score (0–55) from age, SpO₂ and albumin
//	MED-RESP-028  PESIScore                           Pulmonary Embolism Severity Index total (age + item points)
//	MED-RESP-029  PESIRiskClass                       PESI risk class I–V for a total
//	MED-RESP-030  OttawaCOPDRiskScale                 Ottawa COPD Risk Scale total (0–16)
//	MED-RESP-031  OttawaCOPDRiskCategory              Ottawa COPD risk category: 0 Low, 1–2 Medium, 3–4 High, ≥ 5 Very high
//	MED-RESP-032  DRIPScore                           Drug Resistance in Pneumonia score (0–14)
//	MED-RESP-033  FeverPAINScore                      FeverPAIN score (0–5)
//	MED-RESP-034  PRESETScore                         PRESET ECMO score (0–15)
//	MED-RESP-035  PRESETRiskClass                     PRESET risk class: 0–5 I, 6–9 II, 10–15 III
//	MED-RESP-036  QCSIScore                           quick COVID-19 Severity Index (0–12)
//	MED-RESP-037  POMPECMortalityPercent              POMPE-C 30-day mortality, percent (weight in lb)
//	MED-RESP-038  POMPECMortalityPercentKg            POMPE-C 30-day mortality, percent (weight in kg)
//	MED-RESP-039  RESECT90MortalityProbability        RESECT-90 probability of death within 90 days of lung resection
//	MED-RESP-040  STUMBLScore                         blunt chest wall trauma (STUMBL) risk score
//	MED-RESP-041  UtahCOVID19Score                    Utah COVID-19 risk score (0.5–24.5)
//	MED-RESP-042  PARSScore                           Pediatric Asthma Risk Score (0–14)
//	MED-RESP-043  LIPSScore                           Lung Injury Prediction Score
//	MED-RESP-045  CRB65Score                          CRB-65 score (0–4) for community-acquired pneumonia
//	MED-RESP-046  BODEIndex                           BODE index (0–10) for COPD
//	MED-RESP-047  LightsCriteria                      Light's criteria met (0–3) and whether a pleural effusion is an exudate
//	MED-RESP-048  COPDAssessmentTestScore             COPD Assessment Test total (0–40)
//	MED-RESP-049  SilvermanAndersenScore              Silverman–Andersen respiratory severity score (0–10)
//	MED-RESP-050  PEFAmplitudePercentMean             one day's PEF amplitude percent mean (highest − lowest) / mean × 100, GINA
//	MED-RESP-051  PEFDiurnalVariabilityPercent        average daily diurnal PEF variability over the days given, GINA
//	MED-RESP-052  PEFDiurnalVariabilityExcessiveGINA  whether diurnal PEF variability exceeds GINA Box 1-2 (> 10 % adults, > 13 % children 6–11)
//	MED-RESP-053  PEFChangeConsistentWithAsthmaGINA   whether a PEF change is at least 20 % (GINA, adults)
//	MED-RESP-054  AcuteAsthmaPEFBandSIGN              SIGN 158 Table 15 PEF cell (< 33 / 33–50 / > 50–75 % best or predicted), adults
package respiratory
