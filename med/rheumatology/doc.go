// Package rheumatology computes rheumatology classification-criteria scores from their primary sources: the
// 2016 ACR/EULAR primary Sjögren's syndrome criteria, the 2012 EULAR/ACR polymyalgia rheumatica scoring
// algorithm, the 2015 ACR/EULAR gout criteria, and the 2012 SLICC systemic lupus erythematosus rule.
//
// ID short form: RHEUM (IDs are MED-RHEUM-NNN).
//
// Units: items are booleans (the finding is present) or integers with the range the source's table allows;
// serum urate is in mg/dL. An item outside its range returns ErrOutOfRange naming the item; NaN and ±Inf return
// ErrNotFinite. These are classification criteria for research cohorts, as their sources state; the functions
// return the source's score or rule result and nothing else.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-RHEUM-001  SjogrenACREULAR2016Score   2016 ACR/EULAR primary Sjögren's syndrome weighted score (0–9)
//	MED-RHEUM-002  PMRScoreWithoutUltrasound  2012 EULAR/ACR PMR algorithm score without ultrasound (0–6)
//	MED-RHEUM-003  PMRScoreWithUltrasound     2012 EULAR/ACR PMR algorithm score with ultrasound (0–8)
//	MED-RHEUM-004  GoutACREULAR2015Score      2015 ACR/EULAR gout criteria domain score (−6 to 23)
//	MED-RHEUM-005  SLICC2012Classified        whether the SLICC 2012 SLE classification rule is met
//	MED-RHEUM-006  AcuteGoutDiagnosticRuleScore  Janssens 2010 acute gout rule score (0–13)
//	MED-RHEUM-009  PASScore                      PAS / PAS-II (0–10), Wolfe 2005
//	MED-RHEUM-010  PASActivityCategory           PAS / PAS-II activity category, the source's label
package rheumatology
