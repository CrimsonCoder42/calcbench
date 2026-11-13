// Package sports holds sports scoring and rating calculators: the World Handicap System Handicap Index from a
// scoring record.
//
// ID short form: SPORT (IDs are LIF-SPORT-NNN).
//
// Units: golf Score Differentials and Handicap Indexes are strokes, as float64 in tenths, the way the Rules of
// Handicapping state them. Every function checks NaN and ±Inf (ErrNotFinite) and its range (ErrOutOfRange).
//
// Results are numbers and the source's own labels; nothing here advises.
//
// Functions:
//
//	LIF-SPORT-001  HandicapIndexWHS  WHS Handicap Index from up to the 20 most recent Score Differentials (Rule 5.2, Table 5.2a)
package sports
