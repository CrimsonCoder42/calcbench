// Package infectious computes infectious-disease scores from their primary sources: the Centor score for
// streptococcal pharyngitis (Centor 1981) and the Menza score for HIV acquisition (Menza 2009).
//
// ID short form: INFECT (IDs are MED-INFECT-NNN).
//
// Units: items are booleans (the finding is present); results are integer points.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-INFECT-001  CentorScore  Centor 1981 count (0–4) of its four clinical variables
//	MED-INFECT-002  MenzaScore   Menza 2009 simple-model score (0–19) for HIV acquisition
//	MED-INFECT-003  FACEDROPSScore    FACE DROPS score (−1 to 27) for Lyme facial palsy
//	MED-INFECT-004  HIVRASPTotalRisk  RASP total risk of HIV transmission for one exposure (probability)
package infectious
