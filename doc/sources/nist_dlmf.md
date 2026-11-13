# NIST Digital Library of Mathematical Functions

**Citation**: NIST Digital Library of Mathematical Functions, F. W. J. Olver, A. B. Olde Daalhuis, D. W. Lozier, B. I. Schneider, R. F. Boisvert, C. W. Clark, B. R. Miller, B. V. Saunders, H. S. Cohl, M. A. McClain (eds.), https://dlmf.nist.gov. Release 1.2.8 (2026-09-15) where a doc comment names it; read 2026-11-12. Chapter authors: ch. 1 Algebraic and Analytic Methods; ch. 4 Elementary Functions; ch. 5 Gamma Function (Askey, Roy); ch. 7 Error Functions; ch. 19 Elliptic Integrals (Carlson); ch. 26 Combinatorial Analysis (Bressoud).

§1.5 (coordinate systems, NUM-LINALG-052…055) is recorded separately in [nist_dlmf_1_5.md](nist_dlmf_1_5.md).

## Ch. 1 §1.11: zeros of polynomials (num/algebra)

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Horner scheme | §1.11(i) | NUM-ALGEB-001 | Hand-checked example [−1, 2, −6, 2] at 3 → 5; oracle vs math/big |
| Extended Horner scheme (the derivative) | §1.11(i) | NUM-ALGEB-002 | Oracle vs math/big |
| Quadratic formula | §1.11(iii) | NUM-ALGEB-008 | Vieta and known-root oracles |
| Cubic (Cardano) and its Example z³ − 6z² + 6z − 2 | §1.11(iii), Example | NUM-ALGEB-009 | `TestCubicRoots_workedExample` |
| Quartic Example z⁴ − 4z³ + 5z + 2 | §1.11(iii), Example | NUM-ALGEB-010 test | `TestPolynomialRoots_workedExample` |

Status: page read 2026-11-12. Equation labels 1.11.12–1.11.15 NOT confirmed, so the section is cited, not the equations.

## Ch. 4: elementary functions (num/trig, num/geometry)

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| csc, sec, cot | eqs. 4.14.5–4.14.7 | NUM-TRIG-006, -007, -008 | Hand-derived examples; reciprocal oracles |
| cos 2z = 1 − 2 sin² z (haversine form of 4.42.8) | eq. 4.21.28 | NUM-TRIG-031 | Oracle vs the spherical law of cosines |
| 1 − cos C = 2 sin²(C/2), so the angle denominator a − b cos C = (a − b) + 2b sin²(C/2) has no cancelling term | eq. 4.21.28 | NUM-TRIG-023 | `TestSolveTriangleSAS_nearPi` (solve_triangle_sas_test.go:84), 256-bit oracle |
| arccsc, arcsec, arccot; oddness; special values | eqs. 4.23.7–4.23.9, 4.23.13, 4.23.15, Table 4.23.1 | NUM-TRIG-009, -010, -011 | Table 4.23.1 values are test cases |
| csch, sech, coth | eqs. 4.28.5–4.28.7 | NUM-TRIG-012, -013, -014 | Hand-derived examples |
| arccsch, arcsech, arccoth; oddness | eqs. 4.37.7–4.37.9, 4.37.13, 4.37.15 | NUM-TRIG-015, -016, -017 | Examples; 50-digit decimal constants for Asech(0.9), Acoth(1 + 2⁻⁵²) |
| Right and plane triangles: definitions, law of sines, law of cosines | eqs. 4.42.1–4.42.5, Figs. 4.42.1–4.42.2 | NUM-TRIG-018…030 | Kahan 2014 Tables 1 and 2; 512-bit oracle |
| Law of sines c sin B = b sin C and projection formula c cos B = a − b cos C, giving tan B = b sin C / (a − b cos C) (angle opposite the shorter side) | eqs. 4.42.4, 4.42.6 | NUM-TRIG-023 | `TestSolveTriangleSAS_nearPi`; 256-bit oracle (sin/cos/atan by the series 4.19.1, 4.19.2, 4.24.3) |
| Third angle C = π − A − B and sin C, formed by the unexported helper `thirdAngle` (third_angle.go:14; π split as math.Pi + piLo, TwoSum, Sterbenz). The side then follows by the law of sines c = a sin C / sin A | eq. 4.42.4 (the Source line each of 024–026 carries; the angle sum A + B + C = π itself has no separate equation cited) | NUM-TRIG-024, -025, -026 (and NUM-TRIG-023's second angle) | `TestThirdAngle` (third_angle_test.go:107); `…_nearPi` tests in solve_triangle_{asa,aas,ssa}_test.go |
| Spherical law of cosines | eq. 4.42.8 | NUM-TRIG-031 | Oracle |
| Sine series sin z = z − z³/3! + … | eq. 4.19.1 (Release 1.2.8, 2026-09-15) | NUM-GEOM-040 | Series cross-check at small angles |

Status: every section page fetched 2026-11-12 and each cited equation number matched to its formula (the catalog's warning about summarizer-read numbers is cleared). Eqs. 4.42.6 and 4.21.28 for NUM-TRIG-023 confirmed on the page in the trig fix pass of 2026-11-12 (fix_trig.md). The trig Source lines name Release 1.2.8 (2026-09-15).

## Ch. 5: gamma function (num/combinatorics)

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| n! = Γ(n+1) | eq. 5.4.1 | NUM-COMBI-001, -002, -004 | Table constants vs math/big |
| 1 < (2π)^(−½) x^(½−x) eˣ Γ(x) < e^(1/(12x)) | eq. 5.6.1 | NUM-COMBI-016 | Bound tested n = 1…170 vs exact n! |
| Stirling series for ln Γ (coefficients 1/12, −1/360, 1/1260, −1/1680) | eq. 5.11.1 | NUM-COMBI-002, -004, -016; the δ series inside NUM-PROB-026, -032, -035, -041 | Oracle vs math/big; Loader δ table recomputed at 60 digits |

Status: fetched from dlmf.nist.gov 2026-11-12.

## Ch. 7: error functions (num/prob)

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| erf, erfc | §7.2, eqs. 7.2.1–7.2.2 | NUM-PROB-048, -049, -050 | NIST e-Handbook normal table (0.47500) and critical values |
| Φ(z) = ½ erfc(−z/√2), Q(z) = ½ erfc(z/√2) | §7.1 (notation; no equation number) | NUM-PROB-048, -049, -050 | As above |

Status: dlmf.nist.gov/7.1 and /7.2 read 2026-11-12.

## Ch. 19: elliptic integrals (num/geometry)

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Ellipse perimeter L(a, b) = 4aE(k), k² = 1 − b²/a² | eq. 19.9.9 (Release 1.2.8) | NUM-GEOM-046 | `TestEllipsePerimeterRamanujan_vsExact`; Carlson check values |
| Ramanujan's approximation in λ, leading error 3λ¹⁰/2¹⁷ | eq. 19.9.10 (Release 1.2.8) | NUM-GEOM-047 | Error behaviour vs 19.9.9 |
| E(k) as k′ → 0 (used for b/a < 2⁻³⁰) | eq. 19.12.2 (Release 1.2.8) | NUM-GEOM-046 | Threshold test |
| 6R_G from three R_D | eq. 19.21.11 (Release 1.2.8) | NUM-GEOM-070 | Carlson check values |
| Ellipsoid S = 3V R_G(a⁻², b⁻², c⁻²), V = 4πabc/3 | eq. 19.33.1 and §19.33 (Release 1.2.8) | NUM-GEOM-069, -070 | Sphere and spheroid special cases |

Status: pages read 2026-11-12 (equations quoted from the page source). NUM-GEOM-040, -046, -047, -069, -070 cite Release 1.2.8 (2026-09-15) in their Source lines. 19.8 read but not cited (Carlson 1995 eqs. (42)–(46) used instead). 19.21.10 read and rejected: its terms cancel for flat ellipsoids.

## Ch. 26: combinatorial analysis (num/combinatorics)

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Falling factorial | §26.1 | NUM-COMBI-005 | Oracle vs math/big |
| Binomial coefficient, symmetry, recurrence; Table 26.3.1 | eqs. 26.3.1, 26.3.2, 26.3.6, Table 26.3.1 | NUM-COMBI-003, -004 | C(6,3) = 20, C(8,2) = 28, C(10,5) = 252 |
| Multinomial coefficient | eq. 26.4.2 | NUM-COMBI-008, -017 | Oracle vs math/big |
| Catalan numbers | eqs. 26.5.1, 26.5.4, Table 26.5.1 | NUM-COMBI-012, -019 | C(10) = 16 796; 26.5.4 coded as a math/big oracle |
| Bell numbers | eqs. 26.7.2, 26.7.6, Table 26.7.1 | NUM-COMBI-015, -022 | B(5) = 52, B(10) = 115 975; 26.7.6 as oracle |
| Stirling numbers of both kinds | eqs. 26.8.3, 26.8.6, 26.8.18, 26.8.22, Tables 26.8.1, 26.8.2 | NUM-COMBI-013, -014, -020, -021 | s(6,3) = −225, s(10,5) = −269 325; S(6,3) = 90, S(10,5) = 42 525 |
| Derangements | eq. 26.13.4 | NUM-COMBI-011, -018 | Coded as a math/big oracle |
| Twelvefold way (kⁿ, one-to-one, C(k+n−1, n)) | §26.17, Table 26.17.1 | NUM-COMBI-005, -006, -007 | Oracles |

Status: each equation re-read on dlmf.nist.gov 2026-11-12. NOT used: §26.1 for n! (it defines the falling factorial only; NUM-COMBI-001 cites 5.4.1); Table 26.4.1 (layout not read reliably; NUM-COMBI-008 cites eq. 26.4.2 only).
