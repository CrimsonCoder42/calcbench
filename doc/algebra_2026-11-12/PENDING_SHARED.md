# Pending shared edits: algebra catalog pass, 2026-11-12

A subject agent never edits shared files. Each block below is in its final form for Devin or a merge pass to
paste. Nothing here is applied yet. Paste only after the catalog (`doc/catalog/algebra.md`) is approved; rows
that depend on an open question say which.

## 1. INDEX.md rows (paste at BUILD time, after the functions exist; sorted by ID)

The Source and Time/Space cells must match the doc comments word for word; the build pass regenerates them
from the code. These are the intended text.

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-ALGEB-001 | PolynomialEvaluate | num/algebra | value of a polynomial at x by Horner's rule | O(n)   n = degree = len(coeffs)−1 | O(1), no allocation | Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.4 (Horner's rule); Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §5.1 | polynomial, horner, evaluate, algebra |
| NUM-ALGEB-002 | PolynomialValueAndDerivative | num/algebra | value and first derivative of a polynomial at x in one pass | O(n)   n = degree | O(1), no allocation | Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §5.2; NIST DLMF §1.11(i) "Extended Horner Scheme", https://dlmf.nist.gov/1.11 (read 2026-11-12) | polynomial, horner, derivative, newton, algebra |
| NUM-ALGEB-003 | PolynomialMultiply | num/algebra | coefficients of the product of two polynomials | O(n·m)   n = len(a), m = len(b) | O(n+m) for the result (one allocation) | Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6 | polynomial, multiply, binomial, foil, convolution, algebra |
| NUM-ALGEB-004 | LinearEvaluate | num/algebra | y = m·x + b with one rounding | O(1) | O(1) | Stewart, Calculus: Early Transcendentals, 8th ed., Cengage, 2016, Appendix B (slope-intercept form); IEEE Std 754-2019, §5.4.1 (fusedMultiplyAdd) | linear, line, slope-intercept, algebra |
| NUM-ALGEB-005 | LinearSolve | num/algebra | x from y = m·x + b | O(1) | O(1) | Stewart, Calculus: Early Transcendentals, 8th ed., Cengage, 2016, Appendix B (slope-intercept form) | linear, solve, inverse, algebra |
| NUM-ALGEB-006 | QuadraticDiscriminant | num/algebra | b² − 4ac to within 2 ulps | O(1) | O(1) | Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8; Boldo, "Kahan's algorithm for a correct discriminant computation at last formally proven", IEEE Transactions on Computers, 2009, 58(2):220–225, doi:10.1109/TC.2008.200 | quadratic, discriminant, fma, kahan, algebra |
| NUM-ALGEB-007 | QuadraticRealRoots | num/algebra | real roots of a·x² + b·x + c without cancellation | O(1) | O(1) | Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge University Press, 2007, §5.6, pp. 227–229; Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §1.8 | quadratic, roots, quadratic-formula, algebra |
| NUM-ALGEB-008 | QuadraticRoots | num/algebra | both complex roots of a·x² + b·x + c | O(1) | O(1) | Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge University Press, 2007, §5.6, pp. 227–229; NIST DLMF §1.11(iii), https://dlmf.nist.gov/1.11 (read 2026-11-12) | quadratic, roots, complex, quadratic-formula, algebra |
| NUM-ALGEB-009 | CubicRoots | num/algebra | the three complex roots of a cubic, Newton-polished | O(1) | O(1) | Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge University Press, 2007, §5.6, pp. 227–229; NIST DLMF §1.11(iii), eqs. 1.11.12–1.11.15, https://dlmf.nist.gov/1.11 (read 2026-11-12) | cubic, roots, cardano, polynomial, algebra |
| NUM-ALGEB-010 | PolynomialRoots | num/algebra | all complex roots of a real polynomial (Aberth–Ehrlich) | O(n²·k)   n = degree, k = sweeps ≤ 50+n | O(n) for the result plus O(n) flags | Aberth, "Iteration methods for finding all zeros of a polynomial simultaneously", Mathematics of Computation, 1973, 27(122):339–344, doi:10.1090/S0025-5718-1973-0329236-7; Ehrlich, "A modified Newton method for polynomials", Communications of the ACM, 1967, 10(2):107–108, doi:10.1145/363067.363115; Bini, "Numerical computation of polynomial zeros by means of Aberth's method", Numerical Algorithms, 1996, 13:179–200, doi:10.1007/BF02207694 | polynomial, roots, aberth, ehrlich, quartic, algebra |

Only if Q1 is approved (`num/integer`), at that package's build time:

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-INTEG-001 | IsPrime | num/integer | whether n is prime, exact for every uint64 | O(log n)   12 bases | O(1) | Sorenson, Webster, "Strong pseudoprimes to twelve prime bases", Mathematics of Computation, 2017, 86(304):985–1003, doi:10.1090/mcom/3134; Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.5.4, Algorithm P | prime, primality, miller-rabin, number-theory |
| NUM-INTEG-002 | PrimeFactors | num/integer | prime factors of n with multiplicity, ascending | O(n^(1/4)) multiplications expected per split (heuristic) | O(log n) | Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15:331–334, doi:10.1007/BF01933667; Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980, 20:176–184, doi:10.1007/BF01933190 | prime, factorization, pollard-rho, number-theory |
| NUM-INTEG-003 | SmallestPrimeFactor | num/integer | least prime factor of n | as NUM-INTEG-002 | O(log n) | Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15:331–334, doi:10.1007/BF01933667; Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980, 20:176–184, doi:10.1007/BF01933190 | prime, factor, smallest-prime-factor, number-theory |

## 2. doc/scope/SCOPE.md row (slice number = next free at merge time)

| S? | NUM | Algebra: polynomials (evaluate, derivative, multiply, roots to any degree), linear and quadratic (10 functions) | [s?_algebra.md](../specs/s?_algebra.md) — catalog: [../catalog/algebra.md](../catalog/algebra.md) | NOT STARTED (catalog awaiting approval) |

## 3. doc/taxonomy.md (only if Q1 is approved)

Add to "Groups and planned packages", after `num/geometry`:

| NUM | `num/integer` | INTEG | primes, factorisation, gcd, lcm, modular arithmetic |

Add to "Where an outside catalog's categories land":

| Number Theory (primes, factors, divisibility) | `num/integer` | |

## 4. Requests to other subjects (the rows MOVED out of algebra)

- **num/geometry**: Distance between two points (use `math.Hypot`), Distance between a point and a line,
  Mid-point between two points, Slope between two points, Line equation from two points, 2nd-degree
  polynomial tangent line. The tangent line needs NUM-ALGEB-002 PolynomialValueAndDerivative: it may import
  `num/algebra` if Devin allows imports between `num` packages (Q6); otherwise the tangent line moves to
  `num/algebra` as a NUM-ALGEB ID.
- **num/logs**: e-to-power-minus-one and Scalb, as routed. The standard library already has `math.Expm1` and
  `math.Ldexp` (scalb for radix 2); recommendation to that subject: DROP both (Q4).
- **num/stats**: Sum of squares; Weighted geometric mean (the three vCalc variants are one function with
  weights); Weighted grade average (a weighted arithmetic mean).
- **num/linalg**: vector addition (from vCalc collection "LM 7_3 Techniques for adding vectors").
- **eng/mechanical**: Conveyor capacity, Conveyor mass capacity, Conveyor volume output.
- **Basic Operations catalog pass (also `num/algebra`)**: percent change and ratio helpers, not in this list
  (Q5). Its IDs continue from NUM-ALGEB-011.

## 5. doc/sources/ files to create at build time (one per source)

knuth_taocp_vol2.md · higham_2002.md (exists; add §1.8, §5.1, §5.2 and IDs NUM-ALGEB-001/002/007) ·
nist_dlmf.md · numerical_recipes_3e.md · jeannerod_louvet_muller_2013.md · boldo_2009.md · goldberg_1991.md ·
stewart_calculus_8e.md · ieee_754_2019.md · aberth_1973.md · ehrlich_1967.md · bini_1996.md ·
wilkinson_1959.md · (num/integer, if approved) sorenson_webster_2017.md · pomerance_selfridge_wagstaff_1980.md ·
pollard_1975.md · brent_1980.md · miller_1976.md · rabin_1980.md.

## 6. doc/DECISIONS.md rows (after Devin answers Q1–Q6; D-numbers assigned at merge)

One row per answered question, in the form: decision · the alternative · why it lost (texts in
`../catalog/algebra.md` and this work folder's README, "Questions").

APPLIED 2026-11-12 (lockdown pass; see doc/lockdown_2026-11-12/APPLY_LOG.md)
