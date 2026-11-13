// Package trig holds the trigonometry the standard library's math package lacks: degree and radian
// conversion, angle normalisation, the reciprocal circular and hyperbolic functions and their inverses,
// plane triangle solving (laws of sines and cosines, the five oblique cases, the right triangle) and the
// central angle between two points on a sphere.
//
// ID short form: TRIG (IDs are NUM-TRIG-NNN).
//
// math.Sin, Cos, Tan, Asin, Acos, Atan, Atan2, Sinh, Cosh, Tanh, Asinh, Acosh, Atanh, Hypot and Sincos are
// never wrapped here. Every angle is in radians, and its name ends in Rad; DegToRad and RadToDeg convert.
// Lengths are in any one unit, and a result comes back in that unit. NaN or ±Inf in any input returns
// ErrNotFinite, and a finite input whose true result is beyond float64 returns ErrOverflow, never ±Inf.
// Triangle parts follow DLMF §4.42: angle A is opposite side a, and in a right triangle C is the right angle.
// The triangle solvers use the cancellation-free formulas of Kahan (2014); none allocates.
//
// Functions:
//
//	NUM-TRIG-001  DegToRad                           degrees to radians
//	NUM-TRIG-002  RadToDeg                           radians to degrees
//	NUM-TRIG-003  NormalizeAngle                     angle reduced to [0, 2π)
//	NUM-TRIG-004  NormalizeAngleSigned               angle reduced to (−π, π]
//	NUM-TRIG-005  AngleDifference                    shortest signed turn from one angle to another, in (−π, π]
//	NUM-TRIG-006  Sec                                secant
//	NUM-TRIG-007  Csc                                cosecant; error at 0
//	NUM-TRIG-008  Cot                                cotangent; error at 0
//	NUM-TRIG-009  Asec                               inverse secant, principal value in [0, π]
//	NUM-TRIG-010  Acsc                               inverse cosecant, principal value in [−π/2, π/2]
//	NUM-TRIG-011  Acot                               inverse cotangent, principal value in (−π/2, π/2]
//	NUM-TRIG-012  Sech                               hyperbolic secant
//	NUM-TRIG-013  Csch                               hyperbolic cosecant; error at 0
//	NUM-TRIG-014  Coth                               hyperbolic cotangent; error at 0
//	NUM-TRIG-015  Asech                              inverse hyperbolic secant on (0, 1]
//	NUM-TRIG-016  Acsch                              inverse hyperbolic cosecant; error at 0
//	NUM-TRIG-017  Acoth                              inverse hyperbolic cotangent on |x| > 1
//	NUM-TRIG-018  LawOfSinesSide                     side from a side and two angles (law of sines)
//	NUM-TRIG-019  LawOfSinesAngle                    the 0, 1 or 2 angles from two sides and a non-included angle (SSA)
//	NUM-TRIG-020  LawOfCosinesSide                   third side from two sides and the included angle
//	NUM-TRIG-021  LawOfCosinesAngle                  angle opposite a side from three sides
//	NUM-TRIG-022  SolveTriangleSSS                   all parts of a triangle from three sides
//	NUM-TRIG-023  SolveTriangleSAS                   all parts from two sides and the included angle
//	NUM-TRIG-024  SolveTriangleASA                   all parts from two angles and the included side
//	NUM-TRIG-025  SolveTriangleAAS                   all parts from two angles and a non-included side
//	NUM-TRIG-026  SolveTriangleSSA                   the 0, 1 or 2 triangles from two sides and a non-included angle
//	NUM-TRIG-027  SolveRightTriangleLegs             right triangle from its two legs
//	NUM-TRIG-028  SolveRightTriangleHypotenuseLeg    right triangle from hypotenuse and one leg
//	NUM-TRIG-029  SolveRightTriangleHypotenuseAngle  right triangle from hypotenuse and an acute angle
//	NUM-TRIG-030  SolveRightTriangleLegAngle         right triangle from a leg and its opposite acute angle
//	NUM-TRIG-031  CentralAngle                       central angle between two latitude/longitude points (haversine form)
//	NUM-TRIG-032  UnwrapAngles                       angle sequence with its period jumps removed (phase unwrapping)
package trig
