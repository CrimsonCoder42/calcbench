# Hand-off: num/trig (builder notes, 2026-11-12)

Built: NUM-TRIG-001…031, all 31 from doc/catalog/trig.md, in catalog order. Held back: none.
Package: `num/trig/` (doc.go + 31 source files + 31 test files). Spec: `doc/specs/s6_trig.md` (Status: BUILDING).
Checks run while building: `gofmt -l num/trig` (clean), `go build ./...` (whole module OK), `go vet ./num/trig/` (clean).
NOT run: `go test`, benchmarks (the lead verifies). A scratch program outside the repo exercised the functions
against Kahan's Tables 1 and 2 and a big.Float oracle while building; it found and fixed the bugs listed in (e).
Imports: standard library only (math, errors, fmt; tests also math/big, math/rand/v2). Nothing from num/* or algo/*.

## (a) INDEX.md rows (final form, sorted by ID)

> Superseded by fix_trig.md (2026-11-12): the Source cells of NUM-TRIG-023 and -031 changed; use the lines there.

Does = the doc.go line; Time, Space, Source (several joined by "; ") and Tags are word for word from the doc comment.
These replace the rows in doc/trig_2026-11-12/PENDING_SHARED.md §1 (which used short citations and "(Q3)/(Q6)" notes).
Also for the lead: NUM-GEOM-018 and NUM-GEOM-019 become REMOVED rows pointing at NUM-TRIG-020 and NUM-TRIG-021 (Devin's decision).

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-TRIG-001 | DegToRad | num/trig | degrees to radians | O(1) | O(1), no allocation | Thompson and Taylor, Guide for the Use of the International System of Units (SI), NIST Special Publication 811, 2008 ed., §5.1.1, Table 6 ("1° = (π/180) rad"), https://www.nist.gov/pml/special-publication-811 (read 2026-11-12) | trig, angle, unit, degree, radian |
| NUM-TRIG-002 | RadToDeg | num/trig | radians to degrees | O(1) | O(1), no allocation | Thompson and Taylor, Guide for the Use of the International System of Units (SI), NIST Special Publication 811, 2008 ed., §5.1.1, Table 6 ("1° = (π/180) rad"), https://www.nist.gov/pml/special-publication-811 (read 2026-11-12) | trig, angle, unit, degree, radian |
| NUM-TRIG-003 | NormalizeAngle | num/trig | angle reduced to [0, 2π) | O(1) | O(1), no allocation | IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), fmod(), DESCRIPTION ("x−i*y … the result has the same sign as x and magnitude less than the magnitude of y"), https://pubs.opengroup.org/onlinepubs/9799919799/functions/fmod.html (read 2026-11-12) | trig, angle, normalize, wrap |
| NUM-TRIG-004 | NormalizeAngleSigned | num/trig | angle reduced to (−π, π] | O(1) | O(1), no allocation | IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), remainder(), DESCRIPTION ("r=x−ny … n is the integral value nearest the exact value x/y"), https://pubs.opengroup.org/onlinepubs/9799919799/functions/remainder.html (read 2026-11-12) | trig, angle, normalize, wrap |
| NUM-TRIG-005 | AngleDifference | num/trig | shortest signed turn from one angle to another, in (−π, π] | O(1) | O(1), no allocation | IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), remainder(), DESCRIPTION ("r=x−ny … n is the integral value nearest the exact value x/y"), https://pubs.opengroup.org/onlinepubs/9799919799/functions/remainder.html (read 2026-11-12) | trig, angle, difference, heading |
| NUM-TRIG-006 | Sec | num/trig | secant | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.14.6, https://dlmf.nist.gov/4.14 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.5, p. 72 | trig, secant, reciprocal |
| NUM-TRIG-007 | Csc | num/trig | cosecant; error at 0 | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.14.5, https://dlmf.nist.gov/4.14 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.4, p. 72 | trig, cosecant, reciprocal |
| NUM-TRIG-008 | Cot | num/trig | cotangent; error at 0 | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.14.7, https://dlmf.nist.gov/4.14 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.6, p. 72 | trig, cotangent, reciprocal |
| NUM-TRIG-009 | Asec | num/trig | inverse secant, principal value in [0, π] | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.23.8 and Table 4.23.1, https://dlmf.nist.gov/4.23 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.4.7 and 4.4.4 (0 ≤ arccos x ≤ π), p. 79 | trig, inverse, arcsecant |
| NUM-TRIG-010 | Acsc | num/trig | inverse cosecant, principal value in [−π/2, π/2] | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.23.7, 4.23.13 and Table 4.23.1, https://dlmf.nist.gov/4.23 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.4.6 and 4.4.4 (−π/2 ≤ arcsin x ≤ π/2), p. 79 | trig, inverse, arccosecant |
| NUM-TRIG-011 | Acot | num/trig | inverse cotangent, principal value in (−π/2, π/2] | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.23.9, 4.23.15 and Table 4.23.1, https://dlmf.nist.gov/4.23 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.4.8, p. 79 | trig, inverse, arccotangent |
| NUM-TRIG-012 | Sech | num/trig | hyperbolic secant | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.28.6, https://dlmf.nist.gov/4.28 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.5.5, p. 83 | trig, hyperbolic, secant |
| NUM-TRIG-013 | Csch | num/trig | hyperbolic cosecant; error at 0 | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.28.5, https://dlmf.nist.gov/4.28 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.5.4, p. 83 | trig, hyperbolic, cosecant |
| NUM-TRIG-014 | Coth | num/trig | hyperbolic cotangent; error at 0 | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.28.7, https://dlmf.nist.gov/4.28 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.5.6, p. 83 | trig, hyperbolic, cotangent |
| NUM-TRIG-015 | Asech | num/trig | inverse hyperbolic secant on (0, 1] | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.37.8, https://dlmf.nist.gov/4.37 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.6.5, p. 86, and 4.6.24 (the logarithmic form, 0 < x ≤ 1), p. 87 | trig, hyperbolic, inverse |
| NUM-TRIG-016 | Acsch | num/trig | inverse hyperbolic cosecant; error at 0 | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.37.7 and 4.37.13, https://dlmf.nist.gov/4.37 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.6.4, p. 86, and 4.6.23 (the logarithmic form, x ≠ 0), p. 87 | trig, hyperbolic, inverse |
| NUM-TRIG-017 | Acoth | num/trig | inverse hyperbolic cotangent on \|x\| > 1 | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.37.9 and 4.37.15, https://dlmf.nist.gov/4.37 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.6.6, p. 86, and 4.6.25 (the logarithmic form, x² > 1), p. 87 | trig, hyperbolic, inverse |
| NUM-TRIG-018 | LawOfSinesSide | num/trig | side from a side and two angles (law of sines) | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.148, p. 79 | trig, triangle, law-of-sines |
| NUM-TRIG-019 | LawOfSinesAngle | num/trig | the 0, 1 or 2 angles from two sides and a non-included angle (SSA) | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12); Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §8 (formulas B_S, B_T, B_C and the selection B(a, A, b)), p. 6, and §9 Table 2, p. 8, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12) | trig, triangle, law-of-sines, ssa, ambiguous-case |
| NUM-TRIG-020 | LawOfCosinesSide | num/trig | third side from two sides and the included angle | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12); Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §7 ("c := √((a−b)² + 4 a b sin²(C/2))"), p. 6, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12) | trig, triangle, law-of-cosines |
| NUM-TRIG-021 | LawOfCosinesAngle | num/trig | angle opposite a side from three sides | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12); Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §3 (the formula for C with μ), p. 3, and §5 Table 1, p. 4, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12) | trig, triangle, law-of-cosines, needle |
| NUM-TRIG-022 | SolveTriangleSSS | num/trig | all parts of a triangle from three sides | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12); Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §3, p. 3, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12) | trig, triangle, solve, sss |
| NUM-TRIG-023 | SolveTriangleSAS | num/trig | all parts from two sides and the included angle | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12); Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §7 (side c), p. 6, and §3 (angles), p. 3, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12) | trig, triangle, solve, sas |
| NUM-TRIG-024 | SolveTriangleASA | num/trig | all parts from two angles and the included side | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.148, p. 79 | trig, triangle, solve, asa |
| NUM-TRIG-025 | SolveTriangleAAS | num/trig | all parts from two angles and a non-included side | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.148, p. 79 | trig, triangle, solve, aas |
| NUM-TRIG-026 | SolveTriangleSSA | num/trig | the 0, 1 or 2 triangles from two sides and a non-included angle | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12); Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §8 (B(a, A, b)), p. 6, and §9 Table 2, p. 8, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12) | trig, triangle, solve, ssa, ambiguous-case |
| NUM-TRIG-027 | SolveRightTriangleLegs | num/trig | right triangle from its two legs | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1–4.42.3 (Fig. 4.42.1, C the right angle) and 4.42.5 with C = π/2, https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78 | trig, triangle, right-triangle, solve |
| NUM-TRIG-028 | SolveRightTriangleHypotenuseLeg | num/trig | right triangle from hypotenuse and one leg | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1 (Fig. 4.42.1, C the right angle) and 4.42.5 with C = π/2, https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78 | trig, triangle, right-triangle, solve |
| NUM-TRIG-029 | SolveRightTriangleHypotenuseAngle | num/trig | right triangle from hypotenuse and an acute angle | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1–4.42.2 (Fig. 4.42.1, C the right angle), https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78 | trig, triangle, right-triangle, solve |
| NUM-TRIG-030 | SolveRightTriangleLegAngle | num/trig | right triangle from a leg and its opposite acute angle | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1 and 4.42.3 (Fig. 4.42.1, C the right angle), https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78 | trig, triangle, right-triangle, solve |
| NUM-TRIG-031 | CentralAngle | num/trig | central angle between two latitude/longitude points (haversine form) | O(1) | O(1), no allocation | NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.8 (with b, c the colatitudes and A = Δλ) and 4.21.28 (cos 2z = 1 − 2 sin² z), which give the haversine form, https://dlmf.nist.gov/4.42 (read 2026-11-12); Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147 ("haversine A = hav A = ½ vers A"), p. 78; Sinnott, "Virtues of the Haversine", Sky and Telescope, 1984, 68(2):159 (original not seen; formula and locator as reproduced in the comp.infosystems.gis FAQ Q5.1, R. G. Chamberlain, 1996, and in Wikipedia "Haversine formula", both read 2026-11-12) | trig, spherical, great-circle, haversine |

## (b) doc/sources/ records (for the lead to fold in)

### dlmf_ch4.md
- **Citation**: NIST Digital Library of Mathematical Functions, F. W. J. Olver et al. (eds.), Release 1.2.8 (2026-09-15),
  ch. 4 "Elementary Functions", https://dlmf.nist.gov/4 (read 2026-11-12).
- **Taken**: 4.14.5–4.14.7 (csc, sec, cot); 4.21.28 (cos 2z = 1 − 2 sin² z); 4.23.7–4.23.9, 4.23.13, 4.23.15 and Table 4.23.1
  (arccsc, arcsec, arccot, oddness, special values); 4.28.5–4.28.7 (csch, sech, coth); 4.37.7–4.37.9, 4.37.13, 4.37.15
  (arccsch, arcsech, arccoth, oddness); 4.42.1–4.42.5 and Figs. 4.42.1–4.42.2 (right and plane triangles: definitions,
  law of sines, law of cosines); 4.42.8 (spherical law of cosines).
- **Used by**: NUM-TRIG-006…031.
- **Checked**: each section page fetched on 2026-11-12 and every cited equation number matched to its formula
  (the catalog's warning about summarizer-read numbers is cleared). Table 4.23.1 values are test cases (Asec, Acsc, Acot).

### abramowitz_stegun_1972.md
- **Citation**: Abramowitz and Stegun (eds.), *Handbook of Mathematical Functions*, NBS Applied Mathematics Series 55,
  10th printing, 1972. Page images at https://personal.math.ubc.ca/~cbm/aands/ (read 2026-11-12).
- **Taken**: p. 72: 4.3.4–4.3.6 (csc, sec, cot) and 4.3.11–4.3.12 (test identities); p. 78: 4.3.147 (right triangle, haversine);
  p. 79: 4.3.148 (plane triangle: law of sines, law of cosines), 4.4.4 (arcsin, arccos ranges), 4.4.6–4.4.8; p. 83: 4.5.4–4.5.6,
  and 4.5.17–4.5.18 (test identities); p. 86: 4.6.4–4.6.6; p. 87: 4.6.23–4.6.25 (logarithmic forms with their domains).
- **Used by**: NUM-TRIG-006…018, 024, 025, 027…031.
- **Checked**: each page image read on 2026-11-12. NOT used: 4.4.13 (ranges of principal values), whose arccos/arcsec and
  arccot/arccsc rows carry the errata asterisk "See page II"; page II is not on the UBC site (404), so the ranges are cited
  from DLMF 4.23 and A&S 4.4.4 instead.

### kahan_2014_triangle.md

> Superseded by fix_trig.md (2026-11-12): NUM-TRIG-023 no longer uses Kahan §3 (angles now from DLMF 4.42.4, 4.42.6, 4.21.28); §7 is still used for its side c.
- **Citation**: W. Kahan, "Miscalculating Area and Angles of a Needle-like Triangle (from Lecture Notes for Introductory
  Numerical Analysis Classes)", UC Berkeley, 2014-09-04, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12).
- **Taken**: §3, p. 3 (angle C with μ, a ≥ b swap; "do not remove parentheses"); §5 Table 1, p. 4 (worked examples for C);
  §7, p. 6 (c := √((a−b)² + 4ab sin²(C/2))); §8, p. 6 (B_S, B_T, B_C and the selection B(a, A, b) with the 0.756 threshold);
  §9 Table 2, p. 8 (worked examples for B).
- **Used by**: NUM-TRIG-019, 020, 021, 022, 023, 026.
- **Checked**: PDF read page by page. Table 1 reproduced in law_of_cosines_angle_test.go (all 10 numeric rows plus the
  error rows) and Table 2 in law_of_sines_angle_test.go (15 rows). Where Kahan's decimal data round to float64 the angle
  moves (rows 4, 5, 7 of Table 1); those rows carry a wider stated tolerance and are also checked against a 512-bit oracle.
  Correction to the catalog: Table 2 is in §9 (p. 8), not §8.

### nist_sp811_2008.md
- **Citation**: A. Thompson and B. N. Taylor, *Guide for the Use of the International System of Units (SI)*, NIST Special
  Publication 811, 2008 ed., §5.1.1, Table 6 ("1° = (π/180) rad"), https://www.nist.gov/pml/special-publication-811 (read 2026-11-12).
- **Used by**: NUM-TRIG-001, 002. **Checked**: chapter 5 page fetched; table number, section and row text confirmed.

### posix_2024.md
- **Citation**: IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), fmod() and remainder(), DESCRIPTION,
  https://pubs.opengroup.org/onlinepubs/9799919799/functions/fmod.html and .../remainder.html (read 2026-11-12).
- **Used by**: NUM-TRIG-003 (fmod), 004, 005 (remainder). **Checked**: both pages fetched; the quoted definitions match.

### sinnott_1984.md

> Superseded by fix_trig.md (2026-11-12): NUM-TRIG-031 no longer cites Sinnott 1984 (prose mention only); the record was retired to `_delete/doc/sources/sinnott_1984.md`.
- **Citation**: R. W. Sinnott, "Virtues of the Haversine", *Sky and Telescope*, 1984, 68(2):159. No DOI.
- **Taken**: the haversine form of the central angle, σ = 2 arcsin(min(1, √(hav Δφ + cos φ1 cos φ2 hav Δλ))).
- **Used by**: NUM-TRIG-031 (third Source line).
- **Checked**: original NOT seen. Formula and locator reproduced by two secondary sources, read 2026-11-12:
  (1) comp.infosystems.gis FAQ Q5.1, answer by R. G. Chamberlain (Caltech/JPL), 1996, https://movable-type.co.uk/scripts/gis-faq-5.1.html
  ("a = sin²(dlat/2) + cos(lat1)·cos(lat2)·sin²(dlon/2); c = 2·arcsin(min(1, sqrt(a)))", citing vol. 68, no. 2, p. 159);
  (2) Wikipedia, "Haversine formula", https://en.wikipedia.org/wiki/Haversine_formula (hav θ = hav Δφ + cos φ1 cos φ2 hav Δλ,
  citing Sky and Telescope 68(2), 159 (1984), and deriving it from the spherical law of cosines).
  The formula is also derived from confirmed primaries: DLMF 4.42.8 with colatitudes, and 4.21.28.

## (c) Citations

> Superseded by fix_trig.md (2026-11-12): Sinnott 1984 is no longer a Source of NUM-TRIG-031; NUM-TRIG-023 cites DLMF 4.42.4, 4.42.6 and 4.21.28 instead of Kahan §3.

Confirmed (how):
- DLMF 4.14, 4.21, 4.23, 4.28, 4.37, 4.42: section pages fetched, equation numbers matched.
- A&S 4.3.4–6, 4.3.11–12, 4.3.147–148, 4.4.4, 4.4.6–8, 4.5.4–6, 4.5.17–18, 4.6.4–6, 4.6.23–25: page images read (pp. 72, 78, 79, 83, 86, 87).
- Kahan 2014 §3, §5 Table 1, §7, §8, §9 Table 2: PDF read.
- NIST SP 811 §5.1.1 Table 6; POSIX fmod() and remainder(): pages fetched.
- Sinnott 1984: two secondary sources with the locator (above), plus derivation from DLMF 4.42.8 and 4.21.28.

NOT confirmed, and therefore not cited:
- A&S 4.4.13 (errata page II unreadable): dropped from NUM-TRIG-009 and 011; DLMF and A&S 4.4.4 cover the ranges.
- Kahan 1987 ("Branch Cuts…", in Iserles & Powell (eds.), Clarendon Press): the page range 165–211 appeared only in a
  search summary. Not needed: POSIX remainder() covers NUM-TRIG-004, and the (−π, π] interval is Devin's decision. Dropped.
- IEEE 754-2019 §5.3.1: not cited (as in the catalog).

Functions held back: none.

## (d) Assumptions (for doc/ASSUMPTIONS.md) and decisions applied

- Devin's decisions applied: radians only; SSA returns ([2]Triangle, count); CentralAngle here; trig owns the law of cosines
  (020, 021); NormalizeAngle [0, 2π), NormalizeAngleSigned (−π, π] (a −π tie → +π); Acot (−π/2, π/2] (±0 → +π/2).
- A: ±Inf input returns ErrNotFinite everywhere. The catalog proposed limits for Asec, Acsc, Acot, Sech, Csch, Coth, Acsch,
  Acoth; the build rule "NaN/Inf in → error" was applied instead.
- A: a finite input whose true result is beyond float64 returns ErrOverflow (RadToDeg, AngleDifference, Csc, Cot, Csch, Coth,
  LawOfSinesSide, LawOfCosinesSide, the solvers). Not in the catalog.
- A: SSA slot 0 holds the smaller B; unused slots are zero. A right-triangle leg equal to the hypotenuse is ErrNotTriangle;
  a flat triangle from three sides is allowed (angles 0 and π).
- A: sentinels are ErrNotFinite, ErrOverflow, ErrDomain, ErrNonPositive, ErrNotTriangle, each declared beside its first user
  (deg_to_rad.go, csc.go, law_of_sines_side.go, law_of_cosines_angle.go) and wrapped with the parameter name.
- A: Triangle {SideA, SideB, SideC, AngleARad, AngleBRad, AngleCRad} lives in solve_triangle_sss.go; geometry can import it.

## (e) Bugs found and fixed while building (all inside num/trig)

- Subnormal sides: power-of-two scaling built 2^1063 (= +Inf), so SSS on 1e-320 sides gave NaN. Fixed: exact math.Ldexp per side.
- ASA, AAS and SSA took sin C from C = π − A − B, which rounds to math.Pi when A + B is tiny (sin = 1.2e-16 instead of ≈ A + B).
  Fixed: sin C = sin(A + B); the second SSA triangle uses C = B1 − A. Tests pin both.
- Kahan's B_T formed (a/tan A)², which overflows for tiny A and silently returned B = A. Fixed with an equal factored root beyond
  2^500 (Kahan's literal form is kept everywhere else). Test pins it.
- The second-SSA-solution guard rejected real, very thin triangles. Fixed: it exists exactly when B1 > A.

## (f) FOUND, NOT FIXED (outside owned paths)

- doc/catalog/trig.md: cites Kahan 2014 Table 2 as §8 (it is §9, p. 8); proposes ±Inf limits (superseded, see (d)); keeps
  A&S 4.4.13 and Kahan 1987 citations (dropped here); still marks Q1–Q6 open. The catalog should be updated to match.
- doc/trig_2026-11-12/PENDING_SHARED.md §1: rows are superseded by (a) above.
- No test or benchmark was run (by instruction). Two expected constants in tests come from 50-digit Python decimal
  (Asech(0.9), Acoth(1 + 2^-52)); the test comments say so.
