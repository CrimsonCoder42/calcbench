# Pending shared changes: trigonometry (for Devin or a merge pass to apply)

These files are not edited by this subject: INDEX.md, doc/scope/SCOPE.md, doc/taxonomy.md (.claude/CLAUDE.md,
"Several subjects at once"). Every row is in final form. The rows depend on the answers to Q1–Q6 in the
report; where an answer changes a row, the row says so.

## 1. INDEX.md rows (apply when each function is BUILT, not before; sorted by ID)

Time and Space are the same for all 31: `O(1)` and `O(1), no allocation`. Source citations are short here;
the full keys are in doc/catalog/trig.md, and each doc comment carries the full citation.

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-TRIG-001 | DegToRad | num/trig | degrees to radians | O(1) | O(1), no allocation | NIST SP 811 (2008), ch. 5, Table 6 | trig, angle, unit, degree, radian |
| NUM-TRIG-002 | RadToDeg | num/trig | radians to degrees | O(1) | O(1), no allocation | NIST SP 811 (2008), ch. 5, Table 6 | trig, angle, unit, degree, radian |
| NUM-TRIG-003 | NormalizeAngle | num/trig | angle reduced to [0, 2π) | O(1) | O(1), no allocation | IEEE Std 1003.1-2024, fmod() | trig, angle, normalize, wrap |
| NUM-TRIG-004 | NormalizeAngleSigned | num/trig | angle reduced to (−π, π] | O(1) | O(1), no allocation | IEEE Std 1003.1-2024, remainder(); Kahan 1987 | trig, angle, normalize, wrap |
| NUM-TRIG-005 | AngleDifference | num/trig | shortest signed turn from one angle to another, in (−π, π] | O(1) | O(1), no allocation | IEEE Std 1003.1-2024, remainder() | trig, angle, difference, heading |
| NUM-TRIG-006 | Sec | num/trig | secant | O(1) | O(1), no allocation | DLMF 4.14.6; A&S 4.3.5 | trig, secant, reciprocal |
| NUM-TRIG-007 | Csc | num/trig | cosecant; error at 0 | O(1) | O(1), no allocation | DLMF 4.14.5; A&S 4.3.4 | trig, cosecant, reciprocal |
| NUM-TRIG-008 | Cot | num/trig | cotangent; error at 0 | O(1) | O(1), no allocation | DLMF 4.14.7; A&S 4.3.6 | trig, cotangent, reciprocal |
| NUM-TRIG-009 | Asec | num/trig | inverse secant, principal value in [0, π] | O(1) | O(1), no allocation | DLMF 4.23.8; A&S 4.4.7, 4.4.13 | trig, inverse, arcsecant |
| NUM-TRIG-010 | Acsc | num/trig | inverse cosecant, principal value in [−π/2, π/2] | O(1) | O(1), no allocation | DLMF 4.23.7; A&S 4.4.6 | trig, inverse, arccosecant |
| NUM-TRIG-011 | Acot | num/trig | inverse cotangent, principal value in (−π/2, π/2] (DLMF branch; Q6) | O(1) | O(1), no allocation | DLMF 4.23.9, Table 4.23.1; A&S 4.4.8 | trig, inverse, arccotangent |
| NUM-TRIG-012 | Sech | num/trig | hyperbolic secant | O(1) | O(1), no allocation | DLMF 4.28.6; A&S 4.5.5 | trig, hyperbolic, secant |
| NUM-TRIG-013 | Csch | num/trig | hyperbolic cosecant; error at 0 | O(1) | O(1), no allocation | DLMF 4.28.5; A&S 4.5.4 | trig, hyperbolic, cosecant |
| NUM-TRIG-014 | Coth | num/trig | hyperbolic cotangent; error at 0 | O(1) | O(1), no allocation | DLMF 4.28.7; A&S 4.5.6 | trig, hyperbolic, cotangent |
| NUM-TRIG-015 | Asech | num/trig | inverse hyperbolic secant on (0, 1] | O(1) | O(1), no allocation | DLMF 4.37.8; A&S 4.6.5, 4.6.24 | trig, hyperbolic, inverse |
| NUM-TRIG-016 | Acsch | num/trig | inverse hyperbolic cosecant; error at 0 | O(1) | O(1), no allocation | DLMF 4.37.7; A&S 4.6.4, 4.6.23 | trig, hyperbolic, inverse |
| NUM-TRIG-017 | Acoth | num/trig | inverse hyperbolic cotangent on \|x\| > 1 | O(1) | O(1), no allocation | DLMF 4.37.9; A&S 4.6.6, 4.6.25 | trig, hyperbolic, inverse |
| NUM-TRIG-018 | LawOfSinesSide | num/trig | side from a side and two angles (law of sines) | O(1) | O(1), no allocation | DLMF 4.42.4; A&S 4.3.148 | trig, triangle, law-of-sines |
| NUM-TRIG-019 | LawOfSinesAngle | num/trig | the 0, 1 or 2 angles from two sides and a non-included angle (SSA) | O(1) | O(1), no allocation | DLMF 4.42.4; Kahan 2014 §8 | trig, triangle, law-of-sines, ssa, ambiguous-case |
| NUM-TRIG-020 | LawOfCosinesSide | num/trig | third side from two sides and the included angle (Q4) | O(1) | O(1), no allocation | DLMF 4.42.5; Kahan 2014 §7 | trig, triangle, law-of-cosines |
| NUM-TRIG-021 | LawOfCosinesAngle | num/trig | angle opposite a side from three sides (Q4) | O(1) | O(1), no allocation | DLMF 4.42.5; Kahan 2014 §3 | trig, triangle, law-of-cosines, needle |
| NUM-TRIG-022 | SolveTriangleSSS | num/trig | all parts of a triangle from three sides | O(1) | O(1), no allocation | DLMF 4.42.5; Kahan 2014 §3 | trig, triangle, solve, sss |
| NUM-TRIG-023 | SolveTriangleSAS | num/trig | all parts from two sides and the included angle | O(1) | O(1), no allocation | DLMF 4.42.5; Kahan 2014 §3, §7 | trig, triangle, solve, sas |
| NUM-TRIG-024 | SolveTriangleASA | num/trig | all parts from two angles and the included side | O(1) | O(1), no allocation | DLMF 4.42.4 | trig, triangle, solve, asa |
| NUM-TRIG-025 | SolveTriangleAAS | num/trig | all parts from two angles and a non-included side | O(1) | O(1), no allocation | DLMF 4.42.4 | trig, triangle, solve, aas |
| NUM-TRIG-026 | SolveTriangleSSA | num/trig | the 0, 1 or 2 triangles from two sides and a non-included angle | O(1) | O(1), no allocation | DLMF 4.42.4; Kahan 2014 §8 | trig, triangle, solve, ssa, ambiguous-case |
| NUM-TRIG-027 | SolveRightTriangleLegs | num/trig | right triangle from its two legs | O(1) | O(1), no allocation | DLMF 4.42.1–4.42.3; A&S 4.3.147 | trig, triangle, right-triangle, solve |
| NUM-TRIG-028 | SolveRightTriangleHypotenuseLeg | num/trig | right triangle from hypotenuse and one leg | O(1) | O(1), no allocation | DLMF 4.42.1; A&S 4.3.147 | trig, triangle, right-triangle, solve |
| NUM-TRIG-029 | SolveRightTriangleHypotenuseAngle | num/trig | right triangle from hypotenuse and an acute angle | O(1) | O(1), no allocation | DLMF 4.42.1–4.42.2 | trig, triangle, right-triangle, solve |
| NUM-TRIG-030 | SolveRightTriangleLegAngle | num/trig | right triangle from a leg and its opposite acute angle | O(1) | O(1), no allocation | DLMF 4.42.1, 4.42.3 | trig, triangle, right-triangle, solve |
| NUM-TRIG-031 | CentralAngle | num/trig (Q3) | central angle between two latitude/longitude points (haversine form) | O(1) | O(1), no allocation | DLMF 4.42.8; A&S 4.3.147; Sinnott 1984 (NOT VERIFIED) | trig, spherical, great-circle, haversine |

## 2. SCOPE.md row (number the slice when the catalog is approved)

| S<n> | NUM | Trigonometry: `num/trig`, 31 functions (angle utilities, reciprocals and inverses, hyperbolic reciprocals, triangle solvers, central angle) | [catalog/trig.md](../catalog/trig.md) | SPEC (catalog awaiting approval) |

## 3. Taxonomy changes proposed (doc/taxonomy.md)

- The `num/trig` row's "Holds" column, now "trigonometry": replace with "reciprocal and inverse functions math lacks,
  degree/radian, angle normalisation, triangle solving (plane, right), central angle; never wraps math".
- Add one sentence under the table: "Import direction inside NUM: `num/geometry` may import `num/trig`; `num/trig` imports nothing from `num/`." (Q4).
- Add `sci/geoscience` to the Earth great-circle distance row of the vCalc mapping, if Q3 is answered as recommended.

## 4. Offers to geometry (doc/catalog/geometry.md). Geometry may mark these rows BLOCKED ON trig

| Geometry needs | Offered by trig | Note |
|---|---|---|
| Circumradius from a side and its opposite angle (R = a / (2 sin A), the extended law of sines) | NUM-TRIG-018/019 primitives; or geometry calls `math.Sin` directly | Geometry owns the circumradius function itself |
| Angles of a triangle from three sides (for incircle and circumcircle rows) | NUM-TRIG-021, 022 | Q4: if accepted, NUM-GEOM-019 becomes MOVE to NUM-TRIG-021 |
| Third side from SAS | NUM-TRIG-020 | Q4: if accepted, NUM-GEOM-018 becomes MOVE to NUM-TRIG-020 |
| Arc angles in a fixed interval | NUM-TRIG-003/004 | for arc-segment and arc-from-chord rows (trig rows #2, #12, #15) |
| Degree inputs at the public edge | NUM-TRIG-001/002 | only if geometry offers Deg twins |
| Right-triangle lengths only (NUM-GEOM-020 hypotenuse, 021 leg) | none needed | Geometry keeps them; trig's 027/028 compute the same lengths internally as part of a full solve. If Devin prefers no internal re-derivation, 027/028 call geometry, which reverses the import direction: Q4 |

## 5. Requests to other subjects

- `sci/geoscience` (no catalog yet): an Earth great-circle distance = radius × NUM-TRIG-031, with the Earth
  radius taken from its own primary source (for example the IUGG/GRS80 mean radius; not researched here).
  Depends on trig. If Q3 moves CentralAngle there, this row and NUM-TRIG-031 merge.
- `sci/physics`: trig rows #21 (grazing angle), #36 (torque at angle), #37 (transverse force on a string) and the
  optics collections c2 and c4. `eng/electrical`: collection c1 (impedance). Routed here; not built by trig.

APPLIED 2026-11-12 (lockdown pass; see doc/lockdown_2026-11-12/APPLY_LOG.md)
