# Catalog: trigonometry

Status: DRAFT for Devin's approval (2026-11-12). Nothing here is built. Work folder: `../trig_2026-11-12/`.
Input: the names in vCalc's Trigonometry category, pasted by Devin. vCalc pages were not fetched, and nothing
but the names was taken from it. Sources were read 2026-11-12.

The rule that shapes this package: Go's `math` already has Sin, Cos, Tan, Asin, Acos, Atan, Atan2, Sinh, Cosh,
Tanh, Asinh, Acosh, Atanh, Hypot and Sincos, and we never wrap them. `num/trig` holds only what `math`
lacks: the reciprocal functions and their inverses (hyperbolic too) with their domain errors, degree and
radian handling, angle normalisation, and triangle solving.

Source keys, cited in full once here:
- **DLMF**: NIST Digital Library of Mathematical Functions, Release 1.2.x, https://dlmf.nist.gov/4, ch. 4
  "Elementary Functions" (§4.14, §4.23, §4.28, §4.37, §4.42). The equation numbers were read through a web
  summarizer; check them character by character before they go into doc comments.
- **A&S**: Abramowitz & Stegun, *Handbook of Mathematical Functions*, NBS AMS 55, 10th printing, 1972, ch. 4.
  Page images at https://personal.math.ubc.ca/~cbm/aands/. Equations 4.3.3–6, 4.3.147–149, 4.4.4–8, 4.4.13,
  4.5.3–6, 4.6.4–6 and 4.6.20–25 were seen. Note: 4.4.13 and 4.4.19 carry an errata asterisk ("See page II"),
  and page II was not opened.
- **Kahan2014**: W. Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", UC Berkeley, 2014-09-04,
  https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (§2 area, §3 angle, §7 SAS side, §8 SSA; Tables 1 and 2).
- **Kahan1987**: W. Kahan, "Branch Cuts for Complex Elementary Functions, or Much Ado About Nothing's Sign Bit",
  in Iserles & Powell (eds.), *The State of the Art in Numerical Analysis*, Clarendon Press (OUP), 1987. The
  page range 165–211 is NOT VERIFIED; the author's preprint of 1987-05-17 was read. No DOI exists.
- **SP811**: NIST Special Publication 811 (2008 ed., Thompson & Taylor), ch. 5, Table 6: "1° = (π/180) rad",
  https://www.nist.gov/pml/special-publication-811/nist-guide-si-chapter-5-units-outside-si.
- **POSIX**: IEEE Std 1003.1-2024 (POSIX), `remainder()` and `fmod()`, https://pubs.opengroup.org/onlinepubs/9799919799/.
  IEEE 754-2019 §5.3.1 itself was not readable and is NOT VERIFIED.
- **Sinnott1984**: R. W. Sinnott, "Virtues of the Haversine", *Sky and Telescope* 68(2):159, 1984. No DOI. The
  original text is NOT VERIFIED; it was seen only through a secondary FAQ.

## Step 1: Clean

42 rows: 5 collections and 37 names, in the order given. The decision rules are Devin's.
Counts: **KEEP 11 · DROP-STDLIB 13 · MOVE 15 · MERGE 1 · DROP 2.**

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| c1 | LM 25.5 Impedance | MOVE | `eng/electrical` | Impedance collection |
| c2 | LM 28.5 The principle of least time for reflection | MOVE | `sci/physics` (optics) | Reflection collection |
| c3 | LM 29.5 Summary | DROP | — | A "Summary" collection |
| c4 | LM 30.1 A real image formed by a converging mirror | MOVE | `sci/physics` (optics) | Mirror collection |
| c5 | LM 30.4 Summary | DROP | — | A "Summary" collection |
| 1 | Angle From Law of Sines | KEEP | `num/trig` NUM-TRIG-019 | Law of sines, angle; the SSA ambiguity included |
| 2 | Angle of a Circle Arc from Radius and Chord Depth | MOVE | `num/geometry` | Arc angle from radius and chord depth |
| 3 | Arc Hyperbolic Cosecant | KEEP | `num/trig` NUM-TRIG-016 | `math` has no arccsch |
| 4 | Arc Hyperbolic Cosine | DROP-STDLIB | — | `math.Acosh` |
| 5 | Arc Hyperbolic Cotangent | KEEP | `num/trig` NUM-TRIG-017 | `math` has no arccoth |
| 6 | Arc Hyperbolic Secant | KEEP | `num/trig` NUM-TRIG-015 | `math` has no arcsech |
| 7 | Arc Hyperbolic Sine | DROP-STDLIB | — | `math.Asinh` |
| 8 | Arc Hyperbolic Tangent | DROP-STDLIB | — | `math.Atanh` |
| 9 | Arccosine | DROP-STDLIB | — | `math.Acos` |
| 10 | Arcsine | DROP-STDLIB | — | `math.Asin` |
| 11 | Arctangent | DROP-STDLIB | — | `math.Atan`, `math.Atan2` |
| 12 | Area of a Circle Arc Segment | MOVE | `num/geometry` | Circle arc segment area |
| 13 | Area of Circle Around a Triangle | MOVE | `num/geometry` | Area of the circumcircle |
| 14 | Area of Circle Within a Triangle | MOVE | `num/geometry` | Area of the incircle |
| 15 | Circle - area of an arc segment | MOVE | `num/geometry` | Circle arc segment area (the same function as #12; geometry merges them) |
| 16 | Circle - radius | MOVE | `num/geometry` | Circle radius |
| 17 | Circumradius | MOVE | `num/geometry` | Circumradius |
| 18 | Cosecant | KEEP | `num/trig` NUM-TRIG-007 | `math` has no csc |
| 19 | Cosine | DROP-STDLIB | — | `math.Cos` |
| 20 | Cotangent | KEEP | `num/trig` NUM-TRIG-008 | `math` has no cot |
| 21 | Grazing Angle | MOVE | `sci/physics` | Grazing angle |
| 22 | Hyperbolic Cosecant | KEEP | `num/trig` NUM-TRIG-013 | `math` has no csch |
| 23 | Hyperbolic Cosine | DROP-STDLIB | — | `math.Cosh` |
| 24 | Hyperbolic Cotangent | KEEP | `num/trig` NUM-TRIG-014 | `math` has no coth |
| 25 | hyperbolic cotangent | MERGE | → #24 | The two hyperbolic cotangent rows are one |
| 26 | Hyperbolic Secant | KEEP | `num/trig` NUM-TRIG-012 | `math` has no sech |
| 27 | Hyperbolic Sine | DROP-STDLIB | — | `math.Sinh` |
| 28 | Hyperbolic Tangent | DROP-STDLIB | — | `math.Tanh` |
| 29 | Radius of Circle in Triangle | MOVE | `num/geometry` | Inradius |
| 30 | Radius of Displaced Circle | MOVE | `num/geometry` | Radius of displaced circle |
| 31 | Secant | KEEP | `num/trig` NUM-TRIG-006 | `math` has no sec |
| 32 | Side From Law of Sines | KEEP | `num/trig` NUM-TRIG-018 | Law of sines, side |
| 33 | Sine | DROP-STDLIB | — | `math.Sin` |
| 34 | Sine (n*Pi radian) | DROP-STDLIB | — | `math.Sin(n*math.Pi)`. Note: it is not exactly 0 for n ≠ 0, a reason for the Deg twins discussed in Q1 |
| 35 | Tangent | DROP-STDLIB | — | `math.Tan` |
| 36 | Torque at Angle | MOVE | `sci/physics` | Torque at angle |
| 37 | Transverse force on a string segment with sin | MOVE | `sci/physics` | Transverse force on a string |

(The c1–c5 and 1–37 numbering is ours. #15 duplicates #12; both are MOVEd as instructed, and geometry merges them.)

## Step 2: The standard set

What a trig package must have, independent of the list. For each: the input rows it covers, and its ID or
why it has none.

| Need | Function(s) | IDs | Covers rows |
|---|---|---|---|
| Degree ↔ radian | DegToRad, RadToDeg | 001, 002 | — |
| Normalise to [0, 2π) and to (−π, π] | NormalizeAngle, NormalizeAngleSigned | 003, 004 | — |
| Shortest signed difference | AngleDifference | 005 | — |
| Reciprocals | Sec, Csc, Cot | 006–008 | 31, 18, 20 |
| Inverse reciprocals | Asec, Acsc, Acot | 009–011 | — |
| Hyperbolic reciprocals | Sech, Csch, Coth | 012–014 | 26, 22, 24, 25 |
| Inverse hyperbolic reciprocals | Asech, Acsch, Acoth | 015–017 | 6, 3, 5 |
| Law of sines | LawOfSinesSide, LawOfSinesAngle (SSA: both solutions) | 018, 019 | 32, 1 |
| Law of cosines | LawOfCosinesSide, LawOfCosinesAngle | 020, 021 | — (overlaps NUM-GEOM-018/019: Q4) |
| Law of tangents | none | — | LATER: A&S 4.3.148 prints it, but SAS is solved stably by Kahan2014 §7 and §3 without it. No use case |
| Solve a triangle | SolveTriangleSSS, SAS, ASA, AAS, SSA | 022–026 | — |
| Right triangle (any two of hypotenuse, legs, acute angle) | SolveRightTriangleLegs, HypotenuseLeg, HypotenuseAngle, LegAngle | 027–030 | — (overlaps NUM-GEOM-020/021: Q4) |
| Great-circle distance | CentralAngle (haversine form) | 031 | — (home: Q3) |
| Spherical triangles beyond the central angle | none | — | LATER: DLMF 4.42.8–12 supply the formulas, but there is no use case yet |
| Degree twins | none proposed | — | Q1 |

## Step 3: Research, one entry per KEEP

Common to every entry:
- Angles are radians unless the name says Deg. Lengths are "same unit in and out".
- Every function is pure, with O(1) time and space and 0 allocations. Solvers return structs by value.
- NaN or ±Inf in returns `ErrNotFinite` and never a silent NaN (`.claude/rules/numeric.md`).
- The triangle solvers share one result type:
  `Triangle{SideA, SideB, SideC float64 /* one length unit */; AngleA, AngleB, AngleC float64 /* radians; A is opposite a */}`.

### DegToRad      (ID reserved: NUM-TRIG-001)
Formula: rad = (deg / 180) · π.
Inputs and units: deg float64, degrees.
Output: float64, radians.
Errors: ErrNotFinite.
Primary source: SP811, ch. 5, Table 6: "1° = (π/180) rad".
Worked example: HAND-CHECKED. 180 → math.Pi exactly; 90 → math.Pi/2 exactly; −45 → −math.Pi/4 exactly. All three are exact in the order (deg/180)·π, since deg/180 is exact for these.
Variants: none.
Best approach: (deg/180)·π, not deg·(π/180). The second rounds π/180 first, so 180·(π/180) need not equal math.Pi. O(1), 0 allocations.
Notes: degrees are never reduced here; that is NormalizeAngle's job.

### RadToDeg      (ID reserved: NUM-TRIG-002)
Formula: deg = (rad / π) · 180.
Inputs and units: rad float64, radians.
Output: float64, degrees.
Errors: ErrNotFinite.
Primary source: SP811, Table 6 (the inverse relation).
Worked example: HAND-CHECKED. math.Pi → 180 exactly; math.Pi/2 → 90 exactly.
Variants: none.
Best approach: (rad/π)·180, for the same exactness reason as 001. O(1), 0 allocations.
Notes: none.

### NormalizeAngle      (ID reserved: NUM-TRIG-003)
Formula: θ' = θ − 2π·⌊θ/(2π)⌋, so θ' ∈ [0, 2π).
Inputs and units: theta float64, radians.
Output: float64 in [0, 2π), radians.
Errors: ErrNotFinite.
Primary source: the operation is POSIX `fmod()`: "x−i*y … the result has the same sign as x and magnitude less than the magnitude of y". The interval [0, 2π) is our convention; no source states it (Q5).
Worked example: HAND-CHECKED. 7π/2 → 3π/2 (4.71238898038469); −π/2 → 3π/2; 2π → 0.
Variants: Deg form (mod 360, exact in binary for integer degrees) if Q1 adds Deg twins.
Best approach: r = math.Mod(θ, 2π); if r < 0 { r += 2π }; if r == 2π { r = 0 } (the add can round up to 2π). O(1), 0 allocations.
Notes: 2π is not representable, so for |θ| ≫ 2π the reduction is relative to float64(2π), not true 2π. The doc must say so. Must NOT loop subtracting 2π.

### NormalizeAngleSigned      (ID reserved: NUM-TRIG-004)
Formula: θ' = θ − 2π·n, with n the integer nearest θ/(2π), so θ' ∈ (−π, π].
Inputs and units: theta float64, radians.
Output: float64 in (−π, π], radians.
Errors: ErrNotFinite.
Primary source: POSIX `remainder()`: "r=x−ny … n is the integral value nearest the exact value x/y. When |n−x/y|=½, the value n is chosen to be even". Kahan1987 p. 8 uses −π ≤ arg ≤ π. The half-open (−π, π] is our convention (Q5).
Worked example: HAND-CHECKED. 7π/2 → −π/2; 3π/2 → −π/2; −π → π.
Variants: [−π, π) if Q5 prefers it.
Best approach: r = math.Remainder(θ, 2π); if r <= −π { r = π } (remainder can return −π on a tie). O(1), 0 allocations.
Notes: the same large-|θ| caveat as 003.

### AngleDifference      (ID reserved: NUM-TRIG-005)
Formula: d = NormalizeAngleSigned(to − from): the shortest signed turn from `from` to `to`, with + counter-clockwise.
Inputs and units: from, to float64, radians.
Output: float64 in (−π, π], radians.
Errors: ErrNotFinite.
Primary source: as 004 (POSIX `remainder()`); the sign convention is ours.
Worked example: HAND-CHECKED. from 350°, to 10° (in radians) → +20° (0.3490658503988659); from 10°, to 350° → −20°; a half turn → +π.
Variants: none.
Best approach: math.Remainder(to − from, 2π) with 004's tie fix. O(1), 0 allocations.
Notes: to − from can lose digits when both are huge; normalise the inputs first if they come from accumulated rotations (doc note).

### Sec      (ID reserved: NUM-TRIG-006)
Formula: sec x = 1 / cos x.
Inputs and units: x float64, radians.
Output: float64, dimensionless; |sec x| ≥ 1.
Errors: ErrNotFinite. The pole ErrDomain is unreachable: no float64 x makes cos x exactly 0, because π/2 is not representable (see the note).
Primary source: DLMF 4.14.6; A&S 4.3.5.
Worked example: HAND-CHECKED. sec 0 = 1; sec(π/3) = 1.9999999999999996 (π/3 is not exact); sec(math.Pi/2) = 1.633123935319537e16, finite and correct for the float64 input.
Variants: none.
Best approach: 1/math.Cos(x); math.Cos does exact argument reduction, so it is accurate near the poles for the given input. O(1), 0 allocations.
Notes: the function must NOT pretend math.Pi/2 is π/2 and return an error there. Its true value at that float64 is the large finite number above. A test pins it.

### Csc      (ID reserved: NUM-TRIG-007)
Formula: csc x = 1 / sin x.
Inputs and units: x float64, radians.
Output: float64; |csc x| ≥ 1.
Errors: ErrDomain at x = ±0, where sin is exactly 0 (the only float64 pole); ErrNotFinite.
Primary source: DLMF 4.14.5; A&S 4.3.4.
Worked example: HAND-CHECKED. csc(π/2) = 1; csc(π/6) = 2.0000000000000004; csc(math.Pi) = 8.165619676597685e15 (finite).
Variants: none.
Best approach: 1/math.Sin(x). O(1), 0 allocations.
Notes: must NOT return ±Inf at 0; that case is the error.

### Cot      (ID reserved: NUM-TRIG-008)
Formula: cot x = cos x / sin x.
Inputs and units: x float64, radians.
Output: float64, dimensionless.
Errors: ErrDomain at x = ±0; ErrNotFinite.
Primary source: DLMF 4.14.7; A&S 4.3.6.
Worked example: HAND-CHECKED. cot(π/4) = 1.0000000000000002; cot(π/2) = 6.123233995736766e-17 (not 0: math.Pi/2 is not exact).
Variants: none.
Best approach: s, c := math.Sincos(x); c/s. One reduction, and better than 1/math.Tan near π/2, where Tan is huge. O(1), 0 allocations.
Notes: none.

### Asec      (ID reserved: NUM-TRIG-009)
Formula: arcsec x = arccos(1/x).
Inputs and units: x float64, dimensionless, |x| ≥ 1.
Output: float64 in [0, π] excluding π/2, radians.
Errors: ErrDomain for |x| < 1; ErrNotFinite for NaN. ±Inf → π/2 (the limit in DLMF Table 4.23.1).
Primary source: DLMF 4.23.8 (principal value); A&S 4.4.7. Range: A&S 4.4.13 (0 ≤ y ≤ π/2 for x ≥ 0; π/2 < y ≤ π for x < 0).
Worked example: DLMF Table 4.23.1: arcsec 1 = 0, arcsec(−1) = π. HAND-CHECKED: arcsec 2 = π/3 (1.0471975511965976).
Variants: Deg (Q1).
Best approach: math.Acos(1/x); 1/x has one rounding and acos is well conditioned away from ±1. O(1), 0 allocations.
Notes: at ±Inf it returns the limit, not an error. Document that it differs from 1/x = ±0.

### Acsc      (ID reserved: NUM-TRIG-010)
Formula: arccsc x = arcsin(1/x).
Inputs and units: x float64, |x| ≥ 1.
Output: float64 in [−π/2, π/2] excluding 0, radians.
Errors: ErrDomain for |x| < 1; ErrNotFinite for NaN. ±Inf → ±0 (DLMF Table 4.23.1 gives 0).
Primary source: DLMF 4.23.7; A&S 4.4.6; odd symmetry, DLMF 4.23.13.
Worked example: DLMF Table 4.23.1: arccsc 1 = π/2, arccsc(−1) = −π/2. HAND-CHECKED: arccsc 2 = π/6 (0.5235987755982989).
Variants: Deg (Q1).
Best approach: math.Asin(1/x). O(1), 0 allocations.
Notes: keep the sign of zero at ±Inf (Kahan1987).

### Acot      (ID reserved: NUM-TRIG-011)
Formula: arccot x = arctan(1/x).
Inputs and units: x float64, any real.
Output: float64 in (−π/2, π/2], radians. The function is odd (DLMF 4.23.15), and arccot(±∞) = 0.
Errors: ErrNotFinite for NaN. At x = 0 the real arccot is two-valued (DLMF Table 4.23.1: "∓π/2"); the sign of zero picks: +0 → π/2, −0 → −π/2.
Primary source: DLMF 4.23.9 and Table 4.23.1; A&S 4.4.8 and 4.4.13 (the 4.4.13 row is asterisked; the errata page was not checked).
Worked example: DLMF Table 4.23.1: arccot 1 = π/4, arccot(−1) = −π/4, arccot(±∞) = 0.
Variants: the (0, π) branch, which some texts use (arccot x = π/2 − arctan x), is a different function. Q6.
Best approach: math.Atan(1/x). It gives ±π/2 at ±0 by IEEE signed infinities, which is the Kahan1987 sign-bit behaviour. O(1), 0 allocations.
Notes: must NOT silently switch to the (0, π) branch; the doc names the branch and cites DLMF.

### Sech      (ID reserved: NUM-TRIG-012)
Formula: sech x = 1 / cosh x.
Inputs and units: x float64, dimensionless.
Output: float64 in (0, 1].
Errors: ErrNotFinite (NaN). ±Inf → 0, the limit.
Primary source: DLMF 4.28.6; A&S 4.5.5.
Worked example: HAND-CHECKED. sech 0 = 1; sech 1 = 0.6480542736638855.
Variants: none.
Best approach: 1/math.Cosh(x). Cosh overflows to +Inf for |x| > ~710 and 1/Inf = 0, the correctly rounded result there. O(1), 0 allocations.
Notes: none.

### Csch      (ID reserved: NUM-TRIG-013)
Formula: csch x = 1 / sinh x.
Inputs and units: x float64.
Output: float64, nonzero.
Errors: ErrDomain at x = ±0; ErrNotFinite (NaN). ±Inf → ±0.
Primary source: DLMF 4.28.5; A&S 4.5.4.
Worked example: HAND-CHECKED. csch 1 = 0.8509181282393216; csch(−1) = −0.8509181282393216.
Variants: none.
Best approach: 1/math.Sinh(x). O(1), 0 allocations.
Notes: odd: keep the sign.

### Coth      (ID reserved: NUM-TRIG-014)
Formula: coth x = 1 / tanh x.
Inputs and units: x float64.
Output: float64 with |coth x| > 1.
Errors: ErrDomain at x = ±0; ErrNotFinite (NaN). ±Inf → ±1.
Primary source: DLMF 4.28.7; A&S 4.5.6.
Worked example: HAND-CHECKED. coth 1 = 1.3130352854993312.
Variants: none (rows #24 and #25 are this one function).
Best approach: 1/math.Tanh(x). O(1), 0 allocations.
Notes: none.

### Asech      (ID reserved: NUM-TRIG-015)
Formula: arcsech x = arccosh(1/x), for 0 < x ≤ 1.
Inputs and units: x float64 in (0, 1].
Output: float64 ≥ 0.
Errors: ErrDomain for x ≤ 0 or x > 1 (the real principal value; A&S 4.6.24 states "0<x≤1"); ErrNotFinite.
Primary source: DLMF 4.37.8; A&S 4.6.5 and log form 4.6.24, ln[1/x + (1/x² − 1)^½].
Worked example: HAND-CHECKED. arcsech 1 = 0; arcsech 0.5 = arccosh 2 = 1.3169578969248166.
Variants: none.
Best approach: math.Acosh(1/x). O(1), 0 allocations.
Notes: arcsech is ill-conditioned at x → 1, as arccosh is at 1. That is inherent, not a bug; document it.

### Acsch      (ID reserved: NUM-TRIG-016)
Formula: arccsch x = arcsinh(1/x), for x ≠ 0.
Inputs and units: x float64, nonzero.
Output: float64, nonzero, odd in x.
Errors: ErrDomain at x = ±0; ErrNotFinite (NaN). ±Inf → ±0.
Primary source: DLMF 4.37.7; A&S 4.6.4 and 4.6.23; odd symmetry, DLMF 4.37.13.
Worked example: HAND-CHECKED. arccsch 1 = arcsinh 1 = 0.881373587019543.
Variants: none.
Best approach: math.Asinh(1/x). O(1), 0 allocations.
Notes: none.

### Acoth      (ID reserved: NUM-TRIG-017)
Formula: arccoth x = arctanh(1/x) = ½·ln((x+1)/(x−1)), for |x| > 1.
Inputs and units: x float64, |x| > 1.
Output: float64, nonzero, odd in x.
Errors: ErrDomain for |x| ≤ 1 (x = ±1 is a pole; |x| < 1 is not real); ErrNotFinite (NaN). ±Inf → ±0.
Primary source: DLMF 4.37.9 (z ≠ ±1); A&S 4.6.6 and log form 4.6.25 "(x²>1)"; odd symmetry, DLMF 4.37.15.
Worked example: HAND-CHECKED. arccoth 2 = ½ ln 3 = 0.5493061443340549.
Variants: none.
Best approach: ½·math.Log1p(2/(x−1)) for x > 1, and the odd reflection for x < −1. For x ∈ (1, 2], x−1 is exact (Sterbenz), so the only roundings are 2/(x−1) and log1p. This beats atanh(1/x) near x → 1⁺, where rounding 1/x is amplified. The claim is HAND-CHECKED and must get a test with exact rational references. O(1), 0 allocations.
Notes: none.

### LawOfSinesSide      (ID reserved: NUM-TRIG-018)
Formula: b = a · sin B / sin A (from a/sin A = b/sin B).
Inputs and units: a float64 > 0 (length); angleA, angleB float64 in (0, π), with A + B < π (radians).
Output: float64 > 0, side b, in the unit of a.
Errors: ErrNonPositive (a); ErrDomain (an angle outside (0, π), or A + B ≥ π, no triangle); ErrNotFinite.
Primary source: DLMF 4.42.4; A&S 4.3.148.
Worked example: HAND-CHECKED. a = 10, A = B = 60° → b = 10; a = 3, A = atan2(3,4), B = π/2 → b = 5.
Variants: Deg (Q1).
Best approach: a·(math.Sin(B)/math.Sin(A)). O(1), 0 allocations.
Notes: none.

### LawOfSinesAngle      (ID reserved: NUM-TRIG-019)
Formula: sin B = b · sin A / a. There are 0, 1 or 2 solutions B ∈ (0, π).
Inputs and units: a, b float64 > 0 (one length unit); angleA float64 in (0, π) (radians, the angle opposite a).
Output: (B1, B2 float64, count int). count = 0 for no triangle, 1 for one (B2 unused), 2 for B1 acute and B2 = π − B1.
Errors: ErrNonPositive; ErrDomain (A ∉ (0, π)); ErrNotFinite. count = 0 is not an error (Q2).
Primary source: DLMF 4.42.4. Case analysis and stable computation: Kahan2014 §8: "this other angle could as easily be 180°−B as B. And no other non-included angle B exists at all when either a < b sin A and A ≤ 90°, or a < b and A ≥ 90°".
Worked example: Kahan2014 Table 2. (a 3, b 17, A 10°) → B = 79.73894317° (reproduced by hand), and the second triangle 100.26105683°.
  (1, 2, 30°) → 90° (one solution, the right triangle). (0.9999999998, 2, 30°) → none.
Variants: Deg (Q1).
Best approach: Kahan2014 §8's chooser among its arcsin, arctan and arccos forms (threshold (b/a)·sin A ≥ 0.756), not a bare asin of a quotient that can round above 1. O(1), 0 allocations.
Notes: a > b gives exactly one triangle (B acute). a = b·sin A gives exactly one (B = 90°). The function must NOT silently pick one of two solutions.

### LawOfCosinesSide      (ID reserved: NUM-TRIG-020; owner pending Q4: duplicates NUM-GEOM-018)
Formula: c² = a² + b² − 2ab cos C, computed as c = √((a−b)² + 4ab·sin²(C/2)).
Inputs and units: a, b float64 > 0 (one length unit); angleC float64 in (0, π) (radians, the included angle).
Output: float64 > 0, side c.
Errors: ErrNonPositive; ErrDomain (C ∉ (0, π)); ErrNotFinite; ErrOverflow.
Primary source: DLMF 4.42.5; A&S 4.3.148; the stable form is Kahan2014 §7.
Worked example: HAND-CHECKED. a = b = 10, C = 60° → 10; a = 3, b = 4, C = 90° → 5.
Variants: Deg (Q1).
Best approach: Kahan2014 §7's form, which avoids cancellation when C is small and a ≈ b. O(1), 0 allocations.
Notes: exactly one of NUM-TRIG-020 and NUM-GEOM-018 survives (Q4).

### LawOfCosinesAngle      (ID reserved: NUM-TRIG-021; owner pending Q4: duplicates NUM-GEOM-019)
Formula: C opposite c. With a ≥ b: μ = c − (a − b) if b ≥ c ≥ 0, else μ = b − (a − c) if c > b ≥ 0.
  C = 2·arctan(√(((a−b)+c)·μ / ((a+(b+c))·((a−c)+b)))). Mathematically cos C = (a² + b² − c²)/(2ab).
Inputs and units: a, b, c float64 > 0, one length unit.
Output: float64 in [0, π], radians, the angle opposite c.
Errors: ErrNonPositive; ErrNotTriangle (the triangle inequality fails); ErrNotFinite.
Primary source: DLMF 4.42.5; Kahan2014 §3.
Worked example: Kahan2014 Table 1. (10, 10, 10) → 60°; (100000, 99999.99979, 0.00029) → 1.145915591E−7°;
  (10000, 5000.000001, 15000) → 179.9985965°.
Variants: Deg (Q1).
Best approach: Kahan2014 §3, with the swap so that a ≥ b. Never acos of the quotient. O(1), 0 allocations.
Notes: a degenerate (flat) triangle gives 0 or π, not an error (Kahan2014 Table 1: area 0, C = 180).

### SolveTriangleSSS      (ID reserved: NUM-TRIG-022)
Formula: each angle from its opposite side by 021. The third angle is π − (the two others) only as a cross-check.
Inputs and units: a, b, c float64 > 0, one length unit.
Output: Triangle (all six parts).
Errors: as 021.
Primary source: DLMF 4.42.5; Kahan2014 §3.
Worked example: Kahan2014 Table 1 row 1 (10, 10, 10) → all angles 60°. HAND-CHECKED: (3, 4, 5) → C = π/2, A = 0.6435011087932844.
Variants: Deg (Q1).
Best approach: three calls of the 021 core, each stable on its own; the angles then sum to π within a few ulps. O(1), 0 allocations.
Notes: degenerate triangles are allowed, with angles 0, 0, π or similar (Kahan2014 Table 1 rows 9–10).

### SolveTriangleSAS      (ID reserved: NUM-TRIG-023)
Formula: c from 020, then A and B from 021.
Inputs and units: a, b float64 > 0; angleC in (0, π) (the included angle, radians).
Output: Triangle.
Errors: as 020.
Primary source: DLMF 4.42.5; Kahan2014 §7 (side), §3 (angles).
Worked example: HAND-CHECKED. (3, 4, C = 90°) → c = 5, A = 0.6435011087932844, B = π/2 − A.
Variants: Deg (Q1).
Best approach: Kahan2014 §7 then §3; not the law of sines for the second angle, which is ambiguous past 90°. O(1), 0 allocations.
Notes: none.

### SolveTriangleASA      (ID reserved: NUM-TRIG-024)
Formula: C = π − A − B; a = c·sin A/sin C, b = c·sin B/sin C.
Inputs and units: angleA, angleB in (0, π) with A + B < π; c float64 > 0 (the included side).
Output: Triangle.
Errors: ErrDomain (A + B ≥ π, or an angle ∉ (0, π)); ErrNonPositive; ErrNotFinite.
Primary source: DLMF 4.42.4 (the angle sum is elementary).
Worked example: HAND-CHECKED. A = B = 60°, c = 10 → equilateral, sides 10.
Variants: Deg (Q1).
Best approach: law of sines (018). O(1), 0 allocations.
Notes: none.

### SolveTriangleAAS      (ID reserved: NUM-TRIG-025)
Formula: C = π − A − B; b = a·sin B/sin A; c = a·sin C/sin A.
Inputs and units: angleA, angleB in (0, π) with A + B < π; a float64 > 0 (opposite A).
Output: Triangle.
Errors: as 024.
Primary source: DLMF 4.42.4.
Worked example: HAND-CHECKED. A = atan2(3,4), B = π/2, a = 3 → b = 5, c = 4.
Variants: Deg (Q1).
Best approach: law of sines (018). O(1), 0 allocations.
Notes: none.

### SolveTriangleSSA      (ID reserved: NUM-TRIG-026)
Formula: B from 019; C = π − A − B; c from 018.
Inputs and units: a, b float64 > 0; angleA in (0, π), the angle opposite a (not included).
Output: ([2]Triangle, count int); count ∈ {0, 1, 2} (Q2).
Errors: ErrNonPositive; ErrDomain; ErrNotFinite.
Primary source: DLMF 4.42.4; Kahan2014 §8.
Worked example: Kahan2014 Table 2 (as in 019): (3, 17, 10°) → two triangles, B = 79.73894317° or 100.26105683°.
Variants: Deg (Q1).
Best approach: 019, then 018 for c. O(1), 0 allocations (a fixed array, not a slice).
Notes: the order is fixed: Triangle[0] has the acute B.

### SolveRightTriangleLegs      (ID reserved: NUM-TRIG-027)
Formula: c = √(a² + b²); A = arctan(a/b); B = π/2 − A; C = π/2.
Inputs and units: legs a, b float64 > 0, one length unit.
Output: Triangle (C = π/2, c the hypotenuse).
Errors: ErrNonPositive; ErrNotFinite.
Primary source: DLMF 4.42.1–4.42.3 (Fig. 4.42.1: right angle C, hypotenuse c); A&S 4.3.147.
Worked example: HAND-CHECKED. (3, 4) → c = 5, A = 0.6435011087932844 (36.8698976458°).
Variants: Deg (Q1).
Best approach: math.Hypot (no overflow) and math.Atan2(a, b). O(1), 0 allocations.
Notes: none.

### SolveRightTriangleHypotenuseLeg      (ID reserved: NUM-TRIG-028)
Formula: b = √((c − a)(c + a)); A = arcsin(a/c), computed as atan2(a, b); B = π/2 − A.
Inputs and units: hypotenuse c and leg a, float64 > 0, a < c.
Output: Triangle.
Errors: ErrNonPositive; ErrNotTriangle (a ≥ c); ErrNotFinite.
Primary source: DLMF 4.42.1; A&S 4.3.147.
Worked example: HAND-CHECKED. (c 5, a 3) → b = 4, A = 0.6435011087932844.
Variants: Deg (Q1).
Best approach: (c−a)(c+a), not c² − a², which cancels when a ≈ c. O(1), 0 allocations.
Notes: a = c is degenerate; ErrNotTriangle (Q2 also decides the degenerate policy).

### SolveRightTriangleHypotenuseAngle      (ID reserved: NUM-TRIG-029)
Formula: a = c·sin A; b = c·cos A; B = π/2 − A.
Inputs and units: hypotenuse c float64 > 0; acute angle A in (0, π/2), radians.
Output: Triangle.
Errors: ErrNonPositive; ErrDomain (A ∉ (0, π/2)); ErrNotFinite.
Primary source: DLMF 4.42.1–4.42.2.
Worked example: HAND-CHECKED. c = 2, A = π/6 → a = 1 (to rounding), b = √3.
Variants: Deg (Q1).
Best approach: math.Sincos. O(1), 0 allocations.
Notes: none.

### SolveRightTriangleLegAngle      (ID reserved: NUM-TRIG-030)
Formula: given leg a and its opposite acute angle A: b = a/tan A; c = a/sin A; B = π/2 − A.
Inputs and units: leg a float64 > 0; angleA in (0, π/2), opposite a.
Output: Triangle.
Errors: ErrNonPositive; ErrDomain; ErrNotFinite.
Primary source: DLMF 4.42.1, 4.42.3.
Worked example: HAND-CHECKED. a = 3, A = atan2(3, 4) → b = 4, c = 5.
Variants: an adjacent-angle input is the caller passing π/2 − A, and the doc shows it. Deg (Q1).
Best approach: s, co := math.Sincos(A); b = a·co/s; c = a/s. O(1), 0 allocations.
Notes: none.

### CentralAngle      (ID reserved: NUM-TRIG-031; home pending Q3)
Formula: the central angle σ between two points on a sphere. Spherical law of cosines: cos σ = sin φ1 sin φ2 + cos φ1 cos φ2 cos Δλ
  (DLMF 4.42.8 with colatitudes). Haversine form: σ = 2·arcsin(√(hav Δφ + cos φ1 cos φ2 hav Δλ)), hav θ = sin²(θ/2).
  φ is latitude and λ longitude.
Inputs and units: lat1, lon1, lat2, lon2 float64, radians; |lat| ≤ π/2.
Output: float64 in [0, π], radians (multiply by a radius for a distance: that belongs to the caller or `sci/geoscience`).
Errors: ErrDomain (|lat| > π/2); ErrNotFinite.
Primary source: DLMF 4.42.8 (spherical law of cosines). Haversine definition: A&S 4.3.147 ("haversine A = hav A = ½ vers A").
  The haversine computation is from Sinnott1984 (NOT VERIFIED: original not read), so Devin should check the original before it ships.
Worked example: HAND-CHECKED. (0, 0) to (0, π/2) → π/2; pole to equator → π/2; two equator points 1° apart → π/180.
Variants: the spherical law of cosines (ill-conditioned for small σ, which is why the haversine form exists).
Best approach: the haversine form with the argument clamped to ≤ 1. It loses accuracy only near antipodes; an atan2
  form fixes that, but its source (Vincenty 1975) was not read, so it is not proposed. O(1), 0 allocations.
Notes: longitudes need no normalising (only Δλ enters, through hav). Earth distance is MOVE (Q3).

## Step 4: Build order and dependencies

Lowest first; "→" means "calls".
1. 001 DegToRad, 002 RadToDeg.
2. 003 NormalizeAngle, 004 NormalizeAngleSigned → 005 AngleDifference.
3. 006 Sec, 007 Csc, 008 Cot, 012 Sech, 013 Csch, 014 Coth (leaf functions over `math`).
4. 009 Asec, 010 Acsc, 011 Acot, 015 Asech, 016 Acsch, 017 Acoth (leaf functions over `math`).
5. 018 LawOfSinesSide; 019 LawOfSinesAngle (Kahan2014 §8 core).
6. 020 LawOfCosinesSide and 021 LawOfCosinesAngle (Kahan2014 §7 and §3 cores). Q4 decides whether these are here or in `num/geometry`.
7. 022 SSS → 021. 023 SAS → 020, 021. 024 ASA and 025 AAS → 018. 026 SSA → 019, 018.
8. 027–030 right-triangle solvers (`math.Hypot`, `math.Atan2`, `math.Sincos`).
9. 031 CentralAngle (home: Q3).
No S1 algorithm is used. Import direction recommended in Q4: `num/geometry` imports `num/trig`, never the reverse.
What geometry needs from here is offered in `../trig_2026-11-12/PENDING_SHARED.md`.
