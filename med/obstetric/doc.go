// Package obstetric holds pregnancy calculators: the ACOG Committee Opinion No. 579 term categories, and the due date and
// gestational age from the last menstrual period (ACOG Committee Opinion No. 700).
//
// ID short form: OB (IDs are MED-OB-NNN).
//
// Units: gestational age is completed weeks and days past the completed weeks (w d/7), as int, the way the
// source states it. An input outside the range the source defines returns ErrOutOfRange. How a gestational age
// was estimated (last menstrual period, ultrasound) is the caller's; nothing here reads a clock.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-OB-001  TermCategoryACOG  ACOG CO 579 category: "early term", "full term", "late term" or "postterm"
//	MED-OB-002  IOTASimpleRules  IOTA simple rules: "malignant", "benign" or "inconclusive" from M and B feature counts
//	MED-OB-003  DueDateFromLMP         estimated due date, LMP + 280 days (ACOG CO 700)
//	MED-OB-004  GestationalAgeFromLMP  gestational age (completed weeks and days) on a date, counted from the LMP
//	MED-OB-005  GestationalWeightGainIOM       IOM 2009 Table S-1 total and weekly weight gain (lb) for a prepregnancy BMI, singleton
//	MED-OB-006  GestationalWeightGainTwinsIOM  IOM 2009 provisional total weight gain (lb) for twins by prepregnancy BMI
package obstetric
