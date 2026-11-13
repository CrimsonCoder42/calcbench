# Notes: num/algebra (builder hand-off, 2026-11-12)

Built: NUM-ALGEB-001…010 (10 functions), `num/algebra/`. Held back: none. Spec: `doc/specs/s5_algebra.md` (Status BUILDING).
Checks run while building: `gofmt -l num/algebra` (clean), `go build ./...` (ok), `go vet ./num/algebra/` (ok).
Tests and benchmarks NOT run (lead verifies). Benchmarks: `Benchmark<Function>/<size>`; O(n) polynomial
functions at degrees 1e3/1e4/1e5, PolynomialMultiply and PolynomialRoots at degrees 1e1/1e2/1e3, scalars as `/scalar`.

## (a) INDEX.md rows (generated from the doc comments; Does = doc.go line)

> Superseded by fix_algebra.md (2026-11-12): NUM-ALGEB-010 Source (Bini narrowed, Graillat & Ménissier-Morain 2012 added) and Space changed; use the lines there.

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-ALGEB-001 | PolynomialEvaluate | num/algebra | value of a polynomial at x by Horner's rule | O(n)   n = degree = len(coeffs)−1; n fused multiply-adds | O(1), no allocation | Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.4, pp. 486–488 (Horner's rule); Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §5.1 "Horner's Method", p. 94; NIST Digital Library of Mathematical Functions, §1.11(i) (Horner scheme), https://dlmf.nist.gov/1.11 (read 2026-11-12) | polynomial, horner, evaluate, algebra |
| NUM-ALGEB-002 | PolynomialValueAndDerivative | num/algebra | value and first derivative of a polynomial at x in one pass | O(n)   n = degree = len(coeffs)−1; 2n fused multiply-adds | O(1), no allocation | NIST Digital Library of Mathematical Functions, §1.11(i) (extended Horner scheme, the derivative), https://dlmf.nist.gov/1.11 (read 2026-11-12); Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §5.2 "Evaluating Derivatives" | polynomial, horner, derivative, newton, algebra |
| NUM-ALGEB-003 | PolynomialMultiply | num/algebra | coefficients of the product of two polynomials | O(n·m)   n = len(a), m = len(b); n·m fused multiply-adds | O(n+m) for the result (one allocation), O(1) extra | Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §30.1, eqs. (30.1)–(30.2) (product in coefficient form) | polynomial, multiply, binomial, foil, convolution, algebra |
| NUM-ALGEB-004 | LinearEvaluate | num/algebra | y = m·x + b with one rounding | O(1) | O(1) | Abramson, College Algebra, 2nd ed., OpenStax (Rice University), 2021, §4.1 "Linear Functions" (slope-intercept form), https://openstax.org/books/college-algebra-2e/pages/4-1-linear-functions (read 2026-11-12); IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, clause 5.4.1 (fusedMultiplyAdd, one rounding) | linear, line, slope-intercept, algebra |
| NUM-ALGEB-005 | LinearSolve | num/algebra | x from y = m·x + b | O(1) | O(1) | Abramson, College Algebra, 2nd ed., OpenStax (Rice University), 2021, §4.1 "Linear Functions" (slope-intercept form), https://openstax.org/books/college-algebra-2e/pages/4-1-linear-functions (read 2026-11-12) | linear, solve, inverse, algebra |
| NUM-ALGEB-006 | QuadraticDiscriminant | num/algebra | b² − 4ac to within 2 ulps | O(1) | O(1) | Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8; Boldo, "Kahan's algorithm for a correct discriminant computation at last formally proven", IEEE Transactions on Computers, 2009, 58(2):220–225, doi:10.1109/TC.2008.200; Goldberg, "What every computer scientist should know about floating-point arithmetic", ACM Computing Surveys, 1991, 23(1):5–48, doi:10.1145/103162.103163, §1.4 "Cancellation", p. 10 (the example) | quadratic, discriminant, fma, kahan, algebra |
| NUM-ALGEB-007 | QuadraticRealRoots | num/algebra | real roots of a·x² + b·x + c without cancellation | O(1) | O(1) | Press, Teukolsky, Vetterling, Flannery, Numerical Recipes in C: The Art of Scientific Computing, 2nd ed., Cambridge University Press, 1992, §5.6, eqs. (5.6.4)–(5.6.5); Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §1.8 "Solving a Quadratic Equation" | quadratic, roots, quadratic-formula, algebra |
| NUM-ALGEB-008 | QuadraticRoots | num/algebra | both complex roots of a·x² + b·x + c | O(1) | O(1) | Press, Teukolsky, Vetterling, Flannery, Numerical Recipes in C: The Art of Scientific Computing, 2nd ed., Cambridge University Press, 1992, §5.6, eqs. (5.6.4)–(5.6.5); NIST Digital Library of Mathematical Functions, §1.11(iii) (quadratic equations), https://dlmf.nist.gov/1.11 (read 2026-11-12) | quadratic, roots, complex, quadratic-formula, algebra |
| NUM-ALGEB-009 | CubicRoots | num/algebra | the three complex roots of a cubic, Newton-polished | O(1) | O(1), no allocation | Press, Teukolsky, Vetterling, Flannery, Numerical Recipes in C: The Art of Scientific Computing, 2nd ed., Cambridge University Press, 1992, §5.6, eqs. (5.6.10)–(5.6.12) and (5.6.15)–(5.6.17); NIST Digital Library of Mathematical Functions, §1.11(iii) (cubic equations, Example), https://dlmf.nist.gov/1.11 (read 2026-11-12); Kahan, "To Solve a Real Cubic Equation", lecture notes, University of California, Berkeley, 1986, https://people.eecs.berkeley.edu/~wkahan/Math128/Cubic.pdf (read 2026-11-12) (Newton polishing, deflation to a quadratic) | cubic, roots, cardano, viete, polynomial, algebra |
| NUM-ALGEB-010 | PolynomialRoots | num/algebra | all complex roots of a real polynomial (Aberth–Ehrlich) | O(n²·k)   n = degree = len(coeffs)−1, k = sweeps ≤ 100 + n | O(n) for the result (one allocation), O(n) extra (hull indices and converged flags, at most two allocations) | Aberth, "Iteration methods for finding all zeros of a polynomial simultaneously", Mathematics of Computation, 1973, 27(122):339–344, doi:10.1090/S0025-5718-1973-0329236-7; Ehrlich, "A modified Newton method for polynomials", Communications of the ACM, 1967, 10(2):107–108, doi:10.1145/363067.363115; Bini, "Numerical computation of polynomial zeros by means of Aberth's method", Numerical Algorithms, 1996, 13(2):179–200, doi:10.1007/BF02207694 (start points, stopping rule) | polynomial, roots, aberth, ehrlich, quartic, quintic, algebra |

## (b) doc/sources/ records (for the lead to fold in)

> Superseded by fix_algebra.md (2026-11-12): bini_1996.md now takes only the Newton-polygon start points (3.8 weight and stopping rule removed; 0.7 offset is our own choice); graillat_menissier_morain_2012.md added.

| File | Citation | Taken | IDs | How checked |
|---|---|---|---|---|
| knuth_taocp_vol2.md (new) | Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.4, pp. 486–488 | Horner's rule; n multiplications and n additions are optimal for a general polynomial (Theorems A, C there) | NUM-ALGEB-001 | Hand-checked example [−1,2,−6,2] at 3 → 5; oracle vs math/big within Higham's γ₂ₙ bound. Locator via Wikipedia "Horner's method" (cites Knuth 1997 §4.6.4 pp. 486–488) and Erickson's notes (jeffe.cs.illinois.edu/teaching/497/08-polynomials.pdf); equation number not confirmed, so not cited |
| higham_2002.md (exists; add rows) | Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002 | §5.1 "Horner's Method", p. 94 (backward error γ₂ₙ); §5.2 "Evaluating Derivatives"; §1.8 "Solving a Quadratic Equation"; §3.1 γₖ (test tolerances) | NUM-ALGEB-001, -002, -007; tolerances in -001/-002/-003 tests | Section titles and p. 94 confirmed by search snippets of the TOC (SIAM page returned 403); §1.8 page not confirmed, cited by section only |
| nist_dlmf.md (new) | NIST Digital Library of Mathematical Functions, §1.11, https://dlmf.nist.gov/1.11 (read 2026-11-12) | §1.11(i) Horner and extended Horner (derivative); §1.11(iii) quadratic formula, cubic (Cardano) and the two Examples (cubic z³−6z²+6z−2, quartic z⁴−4z³+5z+2) | NUM-ALGEB-001, -002, -008, -009, -010 | Page read 2026-11-12; Examples reproduced in TestCubicRoots_workedExample and TestPolynomialRoots_workedExample. Equation labels 1.11.12–15 not confirmed, so not cited |
| cormen_introduction_to_algorithms.md (exists; add row) | Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §30.1, eqs. (30.1)–(30.2) | Product of polynomials in coefficient form, c_j = Σ a_k b_{j−k}, Θ(n²) | NUM-ALGEB-003 | Locator via walkccc.me/CLRS/Chap30/30.1 (Exercise 30.1-1 "using equations (30.1) and (30.2)") and a search summary of the CLRS ch. 30 PDF; hand-checked example; oracle vs exact convolution |
| openstax_college_algebra_2e.md (new) | Abramson, College Algebra, 2nd ed., OpenStax (Rice University), 2021, §4.1 "Linear Functions", https://openstax.org/books/college-algebra-2e/pages/4-1-linear-functions (read 2026-11-12) | Slope-intercept form f(x) = mx + b | NUM-ALGEB-004, -005 | Definition box read on the page; hand-checked examples; oracles vs math/big |
| ieee_754_2019.md (new) | IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, clause 5.4.1 | fusedMultiplyAdd computes (x×y)+z with one rounding | NUM-ALGEB-004 (and every FMA use) | Clause text read in a copy of the standard (www-users.cse.umn.edu/~vinals/.../IEEE Standard 754-2019.pdf); TestLinearEvaluate_oracle checks bit-exact single rounding |
| jeannerod_louvet_muller_2013.md (new) | Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8 | Kahan's FMA determinant, relative error ≤ 2u | NUM-ALGEB-006 (and 007/008/009 through it) | AMS landing page and abstract (search); oracle asserts ≤ 2u vs math/big on 600 inputs incl. near-double roots |
| boldo_2009.md (new) | Boldo, "Kahan's algorithm for a correct discriminant computation at last formally proven", IEEE Transactions on Computers, 2009, 58(2):220–225, doi:10.1109/TC.2008.200 | Second source: the discriminant algorithm, 2-ulp bound, formally proven | NUM-ALGEB-006 | Crossref record |
| goldberg_1991.md (new) | Goldberg, "What every computer scientist should know about floating-point arithmetic", ACM Computing Surveys, 1991, 23(1):5–48, doi:10.1145/103162.103163, §1.4 "Cancellation", p. 10 | Worked example a = 1.22, b = 3.34, c = 2.28 (exact 0.0292, 3-digit 0.1) | NUM-ALGEB-006 test | Read in the original typeset PDF; reproduced in TestQuadraticDiscriminant_workedExample. Note: the cancellation-free quadratic formula is eq. (5) there, not (4) as the catalog says |
| numerical_recipes_c_2e.md (new) | Press, Teukolsky, Vetterling, Flannery, Numerical Recipes in C: The Art of Scientific Computing, 2nd ed., Cambridge University Press, 1992, §5.6 (pp. 183–185) | (5.6.4)–(5.6.5) q form of the quadratic; (5.6.10) Q, R; (5.6.11)–(5.6.12) trigonometric three real roots; (5.6.15)–(5.6.17) one real root | NUM-ALGEB-007, -008, -009 | Read in nrbook.com chap05c.pdf; hand-checked examples; oracles (Vieta, exact discriminant sign, known roots) |
| kahan_1986_cubic.md (new) | Kahan, "To Solve a Real Cubic Equation", lecture notes, University of California, Berkeley, 1986, https://people.eecs.berkeley.edu/~wkahan/Math128/Cubic.pdf (read 2026-11-12) | Polish one real root by Newton, then deflate to a quadratic solved with an accurate discriminant | NUM-ALGEB-009 | PDF read; oracle vs known roots |
| aberth_1973.md (new) | Aberth, "Iteration methods for finding all zeros of a polynomial simultaneously", Mathematics of Computation, 1973, 27(122):339–344, doi:10.1090/S0025-5718-1973-0329236-7 | The simultaneous iteration | NUM-ALGEB-010 | Crossref; DLMF quartic Example; oracle vs known roots (300 polynomials, degree 2–14) |
| ehrlich_1967.md (new) | Ehrlich, "A modified Newton method for polynomials", Communications of the ACM, 1967, 10(2):107–108, doi:10.1145/363067.363115 | The same iteration, independently | NUM-ALGEB-010 | Search records (scilit, arXiv 1508.03359 citing it); ACM page 403 |
| bini_1996.md (new) | Bini, "Numerical computation of polynomial zeros by means of Aberth's method", Numerical Algorithms, 1996, 13(2):179–200, doi:10.1007/BF02207694 | Start points on circles from the upper convex hull of (i, log\|aᵢ\|), angle offset 0.7; stopping rule \|p(z)\| < ε·s(\|z\|), s weights (1 + 3.8i) | NUM-ALGEB-010 | Crossref record. The method details were read in the comments of Bini's own reference code (pzeros, netlib na10; mirror jblevins.org/mirror/amiller/pzeros.f90), NOT in the paper text (Springer not readable) |
| wilkinson_1959.md (new) | Wilkinson, "The evaluation of the zeros of ill-conditioned polynomials. Part I", Numerische Mathematik, 1959, 1:150–166, doi:10.1007/BF01386381 | Test case (x−1)…(x−10) and its conditioning | NUM-ALGEB-010 test | EuDML / Springer records |

## (c) Citations

> Superseded by fix_algebra.md (2026-11-12): Bini's 3.8 weight and stopping rule are no longer used, and the 0.7 offset is no longer attributed to Bini.

Confirmed (how): DLMF §1.11(i),(iii) and both Examples (page read); IEEE 754-2019 cl. 5.4.1 (standard text read);
NR in C 2nd ed. §5.6 eq. numbers (book PDF read); Goldberg §1.4 p. 10 and the example (original PDF read);
Kahan 1986 (PDF read); OpenStax §4.1 (page read); Jeannerod et al., Boldo, Aberth, Bini, Wilkinson (Crossref / publisher
records); Knuth §4.6.4 pp. 486–488 (two secondary sources); Higham §5.1 p. 94, §5.2, §1.8 (TOC search snippets);
CLRS §30.1 eqs. (30.1)–(30.2) (two secondary sources); Ehrlich (search records).

NOT confirmed, so dropped or narrowed: Knuth §4.6.4 eq. (2) (cite section and pages only); Knuth §4.6 introduction
for the product (replaced by CLRS §30.1); Stewart, Calculus 8e, Appendix A and B (no preview; replaced by OpenStax
§4.1); NR 3rd ed. equation numbers (cite the 2nd ed., read directly; the 3rd ed. §5.6 pp. 227–229 is in the catalog
only); DLMF equation labels 1.11.12–1.11.15 (cite §1.11(iii)); Higham §1.8 page; Edelman–Murakami DOI (not cited);
Bini's 3.8 weight and 0.7 offset confirmed in Bini's reference code, not in the paper text.
Corrections for the catalog / PENDING_SHARED (not edited by me): Goldberg's alternative formula is eq. (5), not (4);
PENDING_SHARED Source cells for 001–010 differ from the code; use the rows in (a).

Functions held back: none.

## (d) FOUND, NOT FIXED

> Superseded by fix_algebra.md (2026-11-12): the (1 + 3.8n)·Σ|coeffs| overflow guard is now n·Σ|coeffs|, and the Space line now says at most three extra allocations.

- Tests were written but not run (per the build rules). Tolerances that depend on numerics most: CubicRoots
  worked example (1e-14) and oracle; PolynomialRoots worked example (1e-13), Wilkinson (1e-7 abs), fourfold root
  (1e-3), the real-snap of the DLMF quartic (asserts imag == 0 exactly), and the oracles' requirement that every
  non-real root have an exact conjugate.
- CubicRoots: coefficient ratios beyond about 2^±1000 can underflow a scaled coefficient to 0 (the smallest root
  then comes back 0 instead of a tiny value). Not documented in the doc comment; candidate for DECISIONS/ASSUMPTIONS.
- PolynomialRoots returns ErrOverflow when (1 + 3.8n)·Σ|coeffs[i]| exceeds MaxFloat64 (coefficients near 1e305),
  even if the roots are representable; documented in the doc comment. Scaling would fix it at the cost of underflow risk.
- PolynomialRoots Space says "at most two allocations" extra: Go may stack-allocate the small flags slice.
- doc/catalog/algebra.md and doc/algebra_2026-11-12/PENDING_SHARED.md still carry the pre-build citations
  (Stewart, NR 3rd ed., Knuth eq. (2), Goldberg eq. (4)); the INDEX rows in (a) supersede them.
- Pending shared edits for the lead: SCOPE row for S5 (spec path doc/specs/s5_algebra.md); doc/sources files above;
  DECISIONS rows for Q2–Q6 as Devin decided them.
