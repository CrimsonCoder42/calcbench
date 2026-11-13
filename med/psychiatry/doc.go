// Package psychiatry computes psychiatry and substance-use instrument totals and published bands from their
// primary sources: the Current Opioid Misuse Measure total (Butler 2007) and the Opioid Risk Tool bands
// (Webster 2005).
//
// ID short form: PSYCH (IDs are MED-PSYCH-NNN).
//
// Units: item ratings and totals are integer points with the range the source allows; an item outside it
// returns ErrOutOfRange naming the item. No item wording is embedded (the instruments' text is the
// publishers'); only the arithmetic and the published bands are here.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-PSYCH-001  COMMScore               Current Opioid Misuse Measure total (0–68) from its 17 item ratings
//	MED-PSYCH-002  OpioidRiskToolCategory  Opioid Risk Tool band for a total: low, moderate or high risk
//	MED-PSYCH-003  GAD7Score               GAD-7 anxiety severity score, 0–21
//	MED-PSYCH-004  GAD7SeverityOf          GAD-7 severity band for a total (0–21)
//	MED-PSYCH-005  PCPTSD5Score            Primary Care PTSD Screen for DSM-5 total, 0–5
//	MED-PSYCH-006  ORTOUDScore             revised Opioid Risk Tool for opioid use disorder (ORT-OUD) total, 0–9
//	MED-PSYCH-007  CRAFFTScore             CRAFFT total, 0–6: the count of "yes" answers to its six items
//	MED-PSYCH-008  AUDITCScore                   AUDIT-C total, 0–12
//	MED-PSYCH-009  CAGEScore                     CAGE total, 0–4: the count of "yes" answers
//	MED-PSYCH-010  CoronavirusAnxietyScaleScore  Coronavirus Anxiety Scale total, 0–20
//	MED-PSYCH-011  DAST10Score                   DAST-10 total, 0–10 (question 3 scores for "no")
//	MED-PSYCH-012  DAST10LevelOf                 DAST-10 level for a total (0–10)
//	MED-PSYCH-013  ASRSScreenerCount             ASRS v1.1 screener items at threshold, 0–6
//	MED-PSYCH-014  MDQScreen                     Mood Disorder Questionnaire symptom count and screening rule
//	MED-PSYCH-015  GDS15Score                    15-item Geriatric Depression Scale total, 0–15
//	MED-PSYCH-016  COWSScore                     Clinical Opiate Withdrawal Scale total, 0–48
//	MED-PSYCH-017  COWSSeverityOf                COWS withdrawal band for a total (0–48)
//	MED-PSYCH-018  PSS3Positive                  Patient Safety Screener (ED-SAFE) positive-screen rule
package psychiatry
