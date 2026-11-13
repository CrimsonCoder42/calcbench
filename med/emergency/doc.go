// Package emergency computes emergency-medicine scores from their primary sources: the TWIST score for
// testicular torsion (Barbosa 2013).
//
// ID short form: EMERG (IDs are MED-EMERG-NNN).
//
// Units: items are booleans (the finding is present); results are integer points.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-EMERG-001  TWISTScore  TWIST score (0–7) for testicular torsion in children
//	MED-EMERG-002  NEDOCSScore  NEDOCS emergency-department crowding score, Weiss 2004
//	MED-EMERG-003  NEDOCSLevel  NEDOCS level, the source's label
package emergency
