// Package area converts areas between the square foot and the square meter, by the exact factor of the NIST
// Guide to the SI (NIST SP 811, 2008, Appendix B.8): 1 ft² = 9.290304 × 10⁻² m² exactly.
//
// ID short form: AREA (IDs are UNT-AREA-NNN).
//
// Units: every input and result names its unit (squareFeet, squareMeters). A negative area returns ErrNegative;
// NaN and ±Inf return ErrNotFinite; a result beyond float64 returns ErrOverflow. Results are not rounded.
//
// Functions:
//
//	UNT-AREA-001  SquareFeetToSquareMeters  m² = ft² × 0.09290304
//	UNT-AREA-002  SquareMetersToSquareFeet  ft² = m² / 0.09290304
package area
