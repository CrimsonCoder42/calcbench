// Package oncology computes oncology and immune-effector-cell-therapy gradings from their primary sources: the
// ASTCT consensus grade for cytokine release syndrome (Lee 2019).
//
// ID short form: ONCO (IDs are MED-ONCO-NNN).
//
// Units: findings are booleans or small named enums whose constants carry the source's own wording; an enum
// outside its constants returns ErrOutOfRange naming the parameter. Results are integer grades.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-ONCO-001  ASTCTCRSGrade  ASTCT consensus grade (1–4) of cytokine release syndrome, 0 when the fever condition is unmet
package oncology
