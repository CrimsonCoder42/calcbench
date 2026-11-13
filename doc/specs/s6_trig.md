# S6: Trigonometry (`num/trig`)

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/` (hand-off: `notes_trig.md`)
Decisions: Devin's answers of 2026-11-12 to the catalog's Q1–Q6, as given to the build: radians only (DegToRad and
RadToDeg are the conversions); SolveTriangleSSA returns a count 0/1/2 plus two solution slots; CentralAngle stays
here (Earth distance comes later); trig owns the law of cosines (NUM-TRIG-020/021; NUM-GEOM-018/019 become REMOVED
rows); NormalizeAngle gives [0, 2π) and NormalizeAngleSigned gives (−π, π]; Acot gives (−π/2, π/2].
Assumptions: listed at the end (for the lead to number into doc/ASSUMPTIONS.md).

Detail for each function (formula, worked example, approach) is in [doc/catalog/trig.md](../catalog/trig.md)
Step 3. This spec records only the API, and where the build differs from the catalog.

Common to all 31: pure, O(1) time, O(1) space, no allocation on success. NaN or ±Inf in any input returns
`ErrNotFinite` (including ±Inf for Asec, Acsc, Acot, Sech, Csch, Coth, Acsch, Acoth, where the catalog proposed
the limit; the build rule "NaN/Inf in → error" wins). A finite input whose true result is beyond float64 returns
`ErrOverflow`, never ±Inf. Angles are radians and their names end in `Rad`. Lengths are in any one unit, and
results come back in that unit. Standard library only, no imports from other calcbench packages.

Sentinels: `ErrNotFinite`, `ErrOverflow` (deg_to_rad.go), `ErrDomain` (csc.go), `ErrNonPositive`
(law_of_sines_side.go), `ErrNotTriangle` (law_of_cosines_angle.go). Every one is wrapped with the parameter name.

`Triangle` (solve_triangle_sss.go): `{SideA, SideB, SideC, AngleARad, AngleBRad, AngleCRad float64}`, where angle A
is opposite side a. Returned by value.

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| NUM-TRIG-001 | DegToRad | `func DegToRad(angleDeg float64) (float64, error)` | degrees → radians |
| NUM-TRIG-002 | RadToDeg | `func RadToDeg(angleRad float64) (float64, error)` | radians → degrees |
| NUM-TRIG-003 | NormalizeAngle | `func NormalizeAngle(angleRad float64) (float64, error)` | radians → radians in [0, 2π) |
| NUM-TRIG-004 | NormalizeAngleSigned | `func NormalizeAngleSigned(angleRad float64) (float64, error)` | radians → radians in (−π, π] |
| NUM-TRIG-005 | AngleDifference | `func AngleDifference(fromRad, toRad float64) (float64, error)` | radians → radians in (−π, π] |
| NUM-TRIG-006 | Sec | `func Sec(xRad float64) (float64, error)` | radians → dimensionless |
| NUM-TRIG-007 | Csc | `func Csc(xRad float64) (float64, error)` | radians → dimensionless |
| NUM-TRIG-008 | Cot | `func Cot(xRad float64) (float64, error)` | radians → dimensionless |
| NUM-TRIG-009 | Asec | `func Asec(x float64) (float64, error)` | dimensionless → radians in [0, π] |
| NUM-TRIG-010 | Acsc | `func Acsc(x float64) (float64, error)` | dimensionless → radians in [−π/2, π/2] |
| NUM-TRIG-011 | Acot | `func Acot(x float64) (float64, error)` | dimensionless → radians in (−π/2, π/2] |
| NUM-TRIG-012 | Sech | `func Sech(x float64) (float64, error)` | dimensionless |
| NUM-TRIG-013 | Csch | `func Csch(x float64) (float64, error)` | dimensionless |
| NUM-TRIG-014 | Coth | `func Coth(x float64) (float64, error)` | dimensionless |
| NUM-TRIG-015 | Asech | `func Asech(x float64) (float64, error)` | dimensionless |
| NUM-TRIG-016 | Acsch | `func Acsch(x float64) (float64, error)` | dimensionless |
| NUM-TRIG-017 | Acoth | `func Acoth(x float64) (float64, error)` | dimensionless |
| NUM-TRIG-018 | LawOfSinesSide | `func LawOfSinesSide(sideA, angleARad, angleBRad float64) (float64, error)` | length, radians → length |
| NUM-TRIG-019 | LawOfSinesAngle | `func LawOfSinesAngle(sideA, sideB, angleARad float64) (angleB1Rad, angleB2Rad float64, count int, err error)` | lengths, radians → radians |
| NUM-TRIG-020 | LawOfCosinesSide | `func LawOfCosinesSide(sideA, sideB, angleCRad float64) (float64, error)` | lengths, radians → length |
| NUM-TRIG-021 | LawOfCosinesAngle | `func LawOfCosinesAngle(sideA, sideB, sideC float64) (float64, error)` | lengths → radians |
| NUM-TRIG-022 | SolveTriangleSSS | `func SolveTriangleSSS(sideA, sideB, sideC float64) (Triangle, error)` | lengths → Triangle |
| NUM-TRIG-023 | SolveTriangleSAS | `func SolveTriangleSAS(sideA, sideB, angleCRad float64) (Triangle, error)` | |
| NUM-TRIG-024 | SolveTriangleASA | `func SolveTriangleASA(angleARad, angleBRad, sideC float64) (Triangle, error)` | |
| NUM-TRIG-025 | SolveTriangleAAS | `func SolveTriangleAAS(angleARad, angleBRad, sideA float64) (Triangle, error)` | |
| NUM-TRIG-026 | SolveTriangleSSA | `func SolveTriangleSSA(sideA, sideB, angleARad float64) (solutions [2]Triangle, count int, err error)` | |
| NUM-TRIG-027 | SolveRightTriangleLegs | `func SolveRightTriangleLegs(legA, legB float64) (Triangle, error)` | C is the right angle |
| NUM-TRIG-028 | SolveRightTriangleHypotenuseLeg | `func SolveRightTriangleHypotenuseLeg(hypotenuse, legA float64) (Triangle, error)` | |
| NUM-TRIG-029 | SolveRightTriangleHypotenuseAngle | `func SolveRightTriangleHypotenuseAngle(hypotenuse, angleARad float64) (Triangle, error)` | |
| NUM-TRIG-030 | SolveRightTriangleLegAngle | `func SolveRightTriangleLegAngle(legA, angleARad float64) (Triangle, error)` | A is opposite legA |
| NUM-TRIG-031 | CentralAngle | `func CentralAngle(lat1Rad, lon1Rad, lat2Rad, lon2Rad float64) (float64, error)` | radians → radians in [0, π] |

## Per function

Source, best known, worked example and errors are as in the catalog entry of the same ID unless stated here.
Best known for all: O(1) time and space; the lower bound is reading the input, Ω(1).
Input modified: never (all parameters are scalars).

### 001 DegToRad · 002 RadToDeg
- **Approach**: (deg/180)·π and (rad/π)·180. **Lost**: deg·(π/180), which rounds π/180 first, so 180 does not give math.Pi.
- **Differs**: RadToDeg returns ErrOverflow above about 3.1e306 rad, where the product leaves float64.

### 003 NormalizeAngle · 004 NormalizeAngleSigned · 005 AngleDifference
- **Approach**: math.Mod (exact) and math.Remainder (exact) by float64(2π), with the two fix-ups in the catalog.
  **Lost**: a loop subtracting 2π (O(|θ|), and it accumulates error).
- **Differs**: −0 from NormalizeAngle is returned as +0. AngleDifference returns ErrOverflow when to − from leaves float64.

### 006–008 Sec, Csc, Cot · 012–014 Sech, Csch, Coth
- **Approach**: the reciprocal of math's function (Cot: one Sincos). **Lost**: 1/math.Tan for Cot (Tan is huge near π/2).
- **Differs**: Csc, Cot, Csch, Coth return ErrOverflow for |x| below about 5.6e-309, where 1/x leaves float64.

### 009–011 Asec, Acsc, Acot
- **Approach**: math.Acos(1/x), math.Asin(1/x), math.Atan(1/x). Acot returns π/2 for both +0 and −0 (Devin: range
  (−π/2, π/2]). **Lost**: π/2 − atan x for Acot, which is the other branch (0, π).
- **Differs**: A&S 4.4.13 is not cited (asterisked; its errata page II was not readable); the ranges come from
  DLMF 4.23.7–9 with A&S 4.4.4 (0 ≤ arccos x ≤ π, −π/2 ≤ arcsin x ≤ π/2).

### 015–017 Asech, Acsch, Acoth
- **Approach**: the logarithmic forms (A&S 4.6.23–25) arranged so no 1/x is formed: Asech = log1p(((1−x) +
  √((1−x)(1+x)))/x) for x > ½, ln(1+√((1−x)(1+x))) − ln x otherwise; Acsch = asinh(1/|x|) for |x| > 1,
  ln(1+√(1+x²)) − ln|x| otherwise, with the sign of x; Acoth = ½·log1p(2/(|x|−1)) with the sign of x.
  **Lost**: acosh(1/x), asinh(1/x), atanh(1/x): 1/x overflows for subnormal x, and its rounding is amplified at
  the ends of the domain.

### 018 LawOfSinesSide · 024 SolveTriangleASA · 025 SolveTriangleAAS
- **Approach**: a·(sin B / sin A). ErrDomain when an angle is outside (0, π) or A + B ≥ π. In ASA, AAS and
  SSA the third angle comes from `thirdAngle` (third_angle.go): A + B formed exactly as s + e (TwoSum), then
  C = ((math.Pi − s) + piLo) − e with piLo = π − math.Pi = 1.2246467991473532e-16, so C is accurate to a few
  ulps when A (or A + B) is near π and C is tiny; sin C = sin(A + B) when A + B < π/2 (C near π), sin C
  otherwise. **Lost**: π − A − B in float64 with sin(A + B) (fix 2026-11-12, reviewer): for a = 2, b = 1,
  A = π − 1e-10 it gave c = 1.0000012246 and C with relative error 1.1e-6. Also lost: the projection rule
  c = a cos B + b cos A for c, which cancels when A is obtuse and a ≈ b + c.

### 019 LawOfSinesAngle · 026 SolveTriangleSSA
- **Approach**: Kahan 2014 §8 B(a, A, b) with its 0.756 chooser among B_S, B_T, B_C. **Lost**: asin((b/a) sin A)
  alone, which loses up to half the digits near B = 90° (Kahan 2014 §9).
- **Count**: 0 when no triangle (also a = b with A ≥ π/2), 1 when a ≥ b or the two angles coincide, 2 otherwise
  (the second exists exactly when B1 > A); slot 0 holds the acute (smaller) B. Unused slots are zero.
- **Overflow**: B_T forms (a/tan A)², which overflows for tiny A; beyond 2^500 the root is taken in a factored,
  equal form (hypot, or √(q−r)·√(q+r)). Sides beyond 2^±500 are scaled by an exact power of two first.
- **SSA**: C and sin C from `thirdAngle(A, B1)` for the first triangle and C = B1 − A (exact) for the second.
- **Worked example**: Kahan 2014 Table 2 (§9, p. 8; the catalog said §8): (3, 17, 10°) → 79.73894317°.

### 020 LawOfCosinesSide · 023 SolveTriangleSAS
- **Approach**: Kahan 2014 §7, c = √((a−b)² + 4ab sin²(C/2)), evaluated as hypot(|a−b|, 2 sin(C/2)·√a·√b) so
  ab cannot overflow. **Lost**: √(a² + b² − 2ab cos C), which cancels when c ≪ a, b.
- SAS: the angle opposite the shorter side (b ≤ a) is B = atan2(b sin C, (a − b) + 2b sin²(C/2)), from
  c sin B = b sin C (DLMF 4.42.4) and c cos B = a − b cos C (DLMF 4.42.6) with 4.21.28; all terms are
  non-negative and the rounded c is not used. The other angle is `thirdAngle(C, B)`; when a = b both are B.
  **Lost** (fix 2026-11-12, reviewer): A and B by the 021 core from the rounded c, whose relative error grows
  like u/δ² as δ = π − C → 0 (100% at δ = 1e-10, A = B = 0 returned).

### 021 LawOfCosinesAngle · 022 SolveTriangleSSS
- **Approach**: Kahan 2014 §3, with exact power-of-two scaling when the longest side is beyond 2^±500, so the
  sums and products cannot overflow or underflow. A flat triangle gives 0 or π. **Lost**: acos((a²+b²−c²)/(2ab)).
- **Worked example**: Kahan 2014 Table 1 (§5, p. 4).

### 027–030 right triangles
- **Approach**: math.Hypot, math.Atan2 for both acute angles (not π/2 − A), (c−a)(c+a) for a leg from the
  hypotenuse, math.Sincos.

### 031 CentralAngle
- **Approach**: haversine form, 2·asin(min(1, √h)), with h built from sin(φ2/2 − φ1/2) and sin(λ2/2 − λ1/2) so no
  difference can overflow. **Lost**: the spherical law of cosines (acos of a value near 1 for small σ).
- **Source**: DLMF 4.42.8 and 4.21.28 (the haversine form follows); A&S 4.3.147 (haversine). Sinnott 1984 is
  mentioned in prose only (original not seen; no Source line, fix 2026-11-12).

## Out of scope

Degree twins (radians only), the law of tangents, spherical triangles beyond the central angle, Earth distance
(later, `sci/geoscience`), area and circle formulas (`num/geometry`).

## Assumptions (for doc/ASSUMPTIONS.md)

- ±Inf inputs are errors everywhere, where the catalog proposed limits for nine functions.
- A finite input whose result is beyond float64 returns ErrOverflow.
- SSA slot order: slot 0 holds the smaller B; unused slots are the zero Triangle.
- A right-triangle leg equal to the hypotenuse is ErrNotTriangle (degenerate); a flat triangle from three sides
  is allowed (angles 0 and π).
