// Package pediatrics computes paediatric scores and reference values from their primary sources: the Kocher
// criteria for septic arthritis of the hip (Kocher 1999), the AASM recommended sleep duration (Paruthi 2016)
// and the DHAKA dehydration score (Levine 2015).
//
// ID short form: PEDS (IDs are MED-PEDS-NNN).
//
// Units: every input names its unit (ageMonths, esrMmPerHour, wbcPerMm3); item levels are integers with the
// range the source's table allows. An input outside its range returns ErrOutOfRange naming it; NaN and ±Inf
// return ErrNotFinite. Category functions return the source's own label, verbatim, as a named string type.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-PEDS-001  KocherCriteriaCount        number of Kocher's four predictors present (0–4)
//	MED-PEDS-002  RecommendedSleepHoursAASM  AASM 2016 recommended sleep per 24 h for an age, as [min, max] hours
//	MED-PEDS-003  DHAKADehydrationScore      DHAKA dehydration score (0–12) from four clinical signs
//	MED-PEDS-004  DHAKADehydrationCategory   DHAKA dehydration category for a score
//	MED-PEDS-005  WAT1Score  Withdrawal Assessment Tool-1 (0–12), Franck 2008
//	MED-PEDS-006  PEWSScore  Brighton Paediatric Early Warning Score (0–13), Monaghan 2005
//	MED-PEDS-007  BacterialMeningitisScore             Bacterial Meningitis Score, 0–6 (Nigrovic 2002; validated 2007)
//	MED-PEDS-008  WHOWeightIncrement1MonthPercentiles  WHO 1-month weight-increment percentiles, birth–12 months, g (WHO 2009 data files)
package pediatrics
