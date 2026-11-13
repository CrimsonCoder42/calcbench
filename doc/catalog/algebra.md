# Catalog: algebra (`num/algebra`, SHORT `ALGEB`)

Catalog pass, 2026-11-12. Docs only: nothing is built until Devin approves this file.
Input: the vCalc "Algebra" list, NAMES ONLY (no vCalc page was fetched; nothing else was taken from it).
Work folder: [../algebra_2026-11-12/](../algebra_2026-11-12/README.md). Shared-file rows: `PENDING_SHARED.md` there.

Citation status, used in every entry: **CHECKED** = the locator was confirmed on the web on 2026-11-12
(DLMF page read; DOIs and tables of contents confirmed by search). **CONFIRM IN PRINT** = the book and section are
right but the equation or page number must be checked against the printed copy before the build cites it.

## Step 1: clean

The input's one "Collections" line holds two collections, so it is two rows here: 37 rows from 36 lines.

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| 1 | 10 to the x | DROP | — | standard library: `math.Pow10` (integer n), `math.Pow(10, x)`. Listed both as MOVE num/logs and as DROP in the request; DROP wins (Q4) |
| 2 | 2nd Degree Polynomial | MERGE | num/algebra → NUM-ALGEB-001 | one evaluator for every degree |
| 3 | 2nd Degree Polynomial Tangent Line | MOVE | num/geometry | tangent line is geometry per the request; it needs the derivative from NUM-ALGEB-002 (Q6) |
| 4 | 3rd Degree Polynomial | MERGE | num/algebra → NUM-ALGEB-001 | as row 2 |
| 5 | 4th Degree Polynomial | MERGE | num/algebra → NUM-ALGEB-001 | as row 2 |
| 6 | 5th Degree Polynomial | MERGE | num/algebra → NUM-ALGEB-001 | as row 2 |
| 7 | 6th Degree Polynomial | MERGE | num/algebra → NUM-ALGEB-001 | as row 2 |
| 8 | Competition Analysis | DROP | — | not a formula (a vCalc site feature) |
| 9 | Complex number square root | DROP | — | standard library: `math/cmplx.Sqrt` (principal root, branch cut on the negative real axis) |
| 10 | Conveyor Capacity | MOVE | eng/mechanical | conveyor engineering, not algebra |
| 11 | Conveyor Mass Capacity | MOVE | eng/mechanical | as row 10 |
| 12 | Conveyor Volume Output | MOVE | eng/mechanical | as row 10 |
| 13 | Distance between a point and line | MOVE | num/geometry | geometry |
| 14 | Distance between two points | MOVE | num/geometry | geometry (`math.Hypot` avoids overflow) |
| 15 | e-to-power minus one | MOVE | num/logs | as the request says; note `math.Expm1` already covers it (Q4) |
| 16 | Floor | DROP | — | standard library: `math.Floor` |
| 17 | Line Equation from 2 Points | MOVE | num/geometry | geometry (slope and intercept from two points) |
| 18 | Linear Equation | KEEP | num/algebra → NUM-ALGEB-004, 005 | y = m·x + b and its inverse |
| 19 | mid-point between two points | MOVE | num/geometry | geometry |
| 20 | Multiplication of Two Binomials | KEEP | num/algebra → NUM-ALGEB-003 | kept as the general `PolynomialMultiply`; a binomial is its smallest case, so no binomial-only function |
| 21 | Quadratic Equation | KEEP | num/algebra → NUM-ALGEB-006, 007, 008 | discriminant, real roots, complex roots |
| 22 | Quadratic Formula | MERGE | num/algebra → NUM-ALGEB-007, 008 | same function as row 21 |
| 23 | Quartic Function | MERGE | num/algebra → NUM-ALGEB-001, 010 | evaluator plus the general root finder; no separate quartic closed form (see 010, Notes) |
| 24 | regression test -- data set access (Weighted Geometric Mean) | DROP | — | a vCalc test page, not a formula |
| 25 | Scalb | MOVE | num/logs | as the request says; note `math.Ldexp` is exactly scalb for base 2 (Q4) |
| 26 | Slope Between Two Points | MOVE | num/geometry | geometry |
| 27 | Smallest Prime Factor | MOVE | num/integer (PROPOSED) | number theory, not algebra; researched below under "Routed" so the owner can build it (Q1) |
| 28 | Square Root of Complex Number | DROP | — | duplicate of row 9; `math/cmplx.Sqrt` |
| 29 | Sum of Squares | MOVE | num/stats | descriptive statistic |
| 30 | Weighted Geometric Mean (bounded criteria) | MOVE | num/stats | one weighted geometric mean; the three vCalc variants are input shapes of one function |
| 31 | Weighted Geometric Mean (DataSet criteria) | MOVE | num/stats | as row 30 |
| 32 | Weighted Geometric Mean (ten input criteria) | MOVE | num/stats | as row 30 |
| 33 | Weighted Grade Average | MOVE | num/stats | a weighted arithmetic mean |
| 34 | eX | DROP | — | standard library: `math.Exp` |
| 35 | Collections: LM 7_3 Techniques for adding vectors | DROP | — | a collection, not a formula; vector addition belongs to num/linalg (PENDING_SHARED.md) |
| 36 | Collections: LM 7_6 Summary and Problems | DROP | — | a collection of textbook problems, not a formula |
| 37 | Dataset: vCalc Competition on ddMMyyyy | DROP | — | a dataset, not a formula |

Counts: KEEP 3 · MERGE 7 · MOVE 17 · DROP 10 · LATER 0 (no abstract-algebra row is in the list; vCalc's
category text mentions group, ring and field theory, which stay LATER and are not catalogued).

## Step 2: research (KEEP rows become 10 functions)

Conventions for every entry below (an API question, Q2): a polynomial is `coeffs []float64` in ASCENDING
order, `coeffs[i]` multiplies x^i, so degree = len(coeffs) − 1. All values are dimensionless (pure maths;
a caller that attaches units keeps them consistent itself). Sentinel errors proposed for the package:
`ErrNotFinite` (a NaN or ±Inf input), `ErrOverflow` (finite inputs, non-finite result), `ErrEmpty`,
`ErrDegenerate` (zero leading coefficient, zero slope), `ErrNoConvergence`. Each is wrapped with the parameter name.

### PolynomialEvaluate      (ID reserved: NUM-ALGEB-001)
Formula: p(x) = Σ_{i=0..n} a_i·x^i, evaluated by Horner's rule: s ← a_n; for i = n−1 down to 0: s ← s·x + a_i.
  a_i = coeffs[i], n = degree, x = the point.
Inputs and units: `coeffs []float64` dimensionless, len ≥ 1, all finite; `x float64` dimensionless, finite.
Output: `float64`, p(x).
Errors: len(coeffs) = 0 → ErrEmpty; any NaN/Inf coefficient or x → ErrNotFinite; finite inputs with an Inf
  result → ErrOverflow. Trailing zero coefficients are allowed (evaluation does not need the true degree).
Primary source: Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed.,
  Addison-Wesley, 1997, §4.6.4 "Evaluation of Polynomials" (Horner's rule, eq. (2)) [section CHECKED,
  equation CONFIRM IN PRINT]; error analysis: Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed.,
  SIAM, 2002, §5.1 "Horner's Method", p. 94 [CHECKED]. Second source: NIST DLMF §1.11(i) "Horner's Scheme",
  https://dlmf.nist.gov/1.11 (read 2026-11-12) [CHECKED].
Worked example: none numeric in the sources. HAND-CHECKED: p(x) = −1 + 2x − 6x² + 2x³ (coeffs [−1, 2, −6, 2]),
  x = 3: s = 2 → 2·3 − 6 = 0 → 0·3 + 2 = 2 → 2·3 − 1 = 5. p(3) = 5.
Variants: complex x (internal helper for NUM-ALGEB-010, unexported until asked for); a running error bound
  (Higham §5.1) as a later variant `PolynomialEvaluateWithBound`, not now.
Best approach: Horner with `math.FMA(s, x, a_i)` (one rounding per step). n multiply-adds, the minimum for a
  general polynomial (Knuth §4.6.4). Time O(n), n = len(coeffs) − 1; Space O(1), no allocation.
Notes: Horner is backward stable: the computed value is the exact value of a polynomial whose coefficients
  are each perturbed by at most γ_2n relative (Higham §5.1); near a root the RELATIVE error can be large and
  that is the conditioning of the problem, not a defect. Must NOT expand powers with `math.Pow`, must NOT
  reorder coefficients, must NOT allocate.

### PolynomialValueAndDerivative      (ID reserved: NUM-ALGEB-002)
Formula: Horner's rule carried twice in one pass: s ← a_n, d ← 0; for i = n−1 down to 0: d ← d·x + s;
  s ← s·x + a_i. Then s = p(x), d = p′(x).
Inputs and units: as NUM-ALGEB-001.
Output: `(value, derivative float64)`, p(x) and dp/dx at x.
Errors: as NUM-ALGEB-001; overflow of either result → ErrOverflow.
Primary source: Higham 2002 (as above), §5.2 "Evaluating Derivatives" [section CONFIRM IN PRINT]; Knuth
  TAOCP Vol. 2, 3rd ed., 1997, §4.6.4 [CHECKED]; DLMF §1.11(i) "Extended Horner Scheme",
  https://dlmf.nist.gov/1.11 (read 2026-11-12) [CHECKED].
Worked example: HAND-CHECKED with NUM-ALGEB-001's p: p′(x) = 2 − 12x + 6x², p′(3) = 2 − 36 + 54 = 20.
  Expected (5, 20).
Variants: k-th derivative (extended Horner, DLMF §1.11(i)) LATER if a caller needs it.
Best approach: one pass, 2n FMAs. Time O(n); Space O(1), no allocation.
Notes: this is the Newton step that NUM-ALGEB-009 and 010 use to polish roots, and the slope the tangent
  line in num/geometry needs (Q6). Must NOT compute the derivative coefficients into a new slice.

### PolynomialMultiply      (ID reserved: NUM-ALGEB-003)
Formula: (Σ a_i x^i)(Σ b_j x^j) = Σ_k c_k x^k, c_k = Σ_{i+j=k} a_i·b_j, k = 0..n+m.
  The binomial product (a₁x + a₀)(b₁x + b₀) = a₁b₁x² + (a₁b₀ + a₀b₁)x + a₀b₀ is the case n = m = 1.
Inputs and units: `a, b []float64` ascending coefficients, dimensionless, len ≥ 1, finite.
Output: `[]float64`, len(a) + len(b) − 1 ascending coefficients.
Errors: either empty → ErrEmpty; NaN/Inf → ErrNotFinite; an Inf coefficient from finite inputs → ErrOverflow.
Primary source: Knuth TAOCP Vol. 2, 3rd ed., 1997, §4.6 "Polynomial Arithmetic" (introduction:
  sum and product of polynomials) [section CHECKED, page CONFIRM IN PRINT]. Elementary identity: Stewart,
  Calculus: Early Transcendentals, 8th ed., Cengage, 2016, Appendix A "Numbers, Inequalities, and Absolute
  Values" — Reference Page 1 "Algebra" (FOIL / special products) [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: (2x + 3)(x − 4): a = [3, 2], b = [−4, 1] → c = [−12, −5, 2], i.e. 2x² − 5x − 12.
Variants: FFT/Karatsuba multiplication is faster for degree in the hundreds; LATER, only if a caller needs it
  (and FFT loses accuracy on coefficients of mixed magnitude).
Best approach: schoolbook convolution, each c_k accumulated with `math.FMA`. Time O(n·m), n = len(a),
  m = len(b); Space O(n + m) for the result, one allocation, nothing else.
Notes: each c_k is a dot product; FMA accumulation keeps one rounding per term. Compensated dot product
  (Ogita, Rump, Oishi, "Accurate sum and dot product", SIAM J. Sci. Comput., 2005, 26(6):1955–1988,
  doi:10.1137/030601818) is a later variant if a caller needs it. Must NOT trim trailing zeros silently.

### LinearEvaluate      (ID reserved: NUM-ALGEB-004)
Formula: y = m·x + b. m = slope, b = y-intercept, x = the point.
Inputs and units: `slope, intercept, x float64`, dimensionless (the caller's y-unit per x-unit for slope), finite.
Output: `float64`, y.
Errors: NaN/Inf → ErrNotFinite; Inf result → ErrOverflow.
Primary source: Stewart, Calculus: Early Transcendentals, 8th ed., Cengage, 2016, Appendix B "Coordinate
  Geometry and Lines", slope-intercept form [equation number CONFIRM IN PRINT]. Rounding: IEEE Std 754-2019,
  §5.4.1 fusedMultiplyAdd [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: m = 2, b = 3, x = 4 → 11.
Variants: none (a polynomial of degree 1 through NUM-ALGEB-001 gives the same number; this exists because it
  is the name people search for).
Best approach: `math.FMA(slope, x, intercept)`, one rounding. Time O(1); Space O(1).
Notes: must NOT round twice (m*x + b in two operations).

### LinearSolve      (ID reserved: NUM-ALGEB-005)
Formula: x = (y − b)/m, the solution of m·x + b = y.
Inputs and units: `slope, intercept, y float64`, finite; slope ≠ 0.
Output: `float64`, x.
Errors: slope = 0 → ErrDegenerate (no solution or every x, both reported as degenerate); NaN/Inf → ErrNotFinite;
  Inf result (tiny slope) → ErrOverflow.
Primary source: as NUM-ALGEB-004 (Stewart, Appendix B).
Worked example: HAND-CHECKED: m = 2, b = 3, y = 11 → 4 (the inverse of 004's example).
Variants: none.
Best approach: one subtraction and one division. Time O(1); Space O(1).
Notes: y − b can cancel; that is the conditioning of the data. Must NOT compare slope with a tolerance
  (no source defines one; `.claude/rules/numeric.md`).

### QuadraticDiscriminant      (ID reserved: NUM-ALGEB-006)
Formula: D = b² − 4ac, for a·x² + b·x + c. Computed as a 2×2 determinant with FMA (Kahan's algorithm):
  w = 4a·c; e = FMA(−4a, c, w); f = FMA(b, b, −w); D = f + e. (4a is exact barring overflow.)
Inputs and units: `a, b, c float64`, dimensionless, finite.
Output: `float64`, D.
Errors: NaN/Inf → ErrNotFinite; b² or 4ac outside float64 range → handled by power-of-two scaling (Best
  approach); if D itself overflows → ErrOverflow. a = 0 is allowed here (D is defined).
Primary source: Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate
  computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264,
  doi:10.1090/S0025-5718-2013-02679-8 (relative error ≤ 2u with FMA) [DOI CHECKED, pages CONFIRM]. Second:
  Boldo, "Kahan's algorithm for a correct discriminant computation at last formally proven", IEEE
  Transactions on Computers, 2009, 58(2):220–225, doi:10.1109/TC.2008.200 [CHECKED]. Why it matters: Goldberg,
  "What Every Computer Scientist Should Know About Floating-Point Arithmetic", ACM Computing Surveys, 1991,
  23(1):5–48, doi:10.1145/103162.103163, §1.4 "Cancellation" [section CONFIRM IN PRINT].
Worked example: Goldberg 1991 §1.4: a = 1.22, b = 3.34, c = 2.28; exact D = 11.1556 − 11.1264 = 0.0292, where
  3-digit decimal arithmetic gives 0.1. The test checks the float64 result against the exact discriminant of
  the float64-rounded inputs computed with `math/big`, to ≤ 2 ulps.
Variants: none.
Best approach: as Formula; 3 FMAs. Time O(1); Space O(1).
Notes: the naive b*b − 4*a*c loses every digit when b² ≈ 4ac (a near-double root). Must NOT use the naive form.

### QuadraticRealRoots      (ID reserved: NUM-ALGEB-007)
Formula: D from NUM-ALGEB-006. If D < 0 no real roots. Else q = −½·(b + sgn(b)·√D) (sgn(0) = +1),
  x₁ = q/a, x₂ = c/q; if q = 0 (b = 0 and D = 0, so c = 0) both roots are 0.
Inputs and units: `a, b, c float64`, dimensionless, finite, a ≠ 0.
Output: `(lo, hi float64, n int, err error)`: n ∈ {0, 1, 2} distinct real roots, lo ≤ hi; for n = 1 lo = hi
  (double root); for n = 0 lo = hi = 0.
Errors: a = 0 → ErrDegenerate (caller uses LinearSolve; Q2); NaN/Inf → ErrNotFinite; a root that is ±Inf from
  finite inputs → ErrOverflow. D < 0 is NOT an error: n = 0, err = nil.
Primary source: Press, Teukolsky, Vetterling, Flannery, Numerical Recipes: The Art of Scientific Computing,
  3rd ed., Cambridge University Press, 2007, §5.6 "Quadratic and Cubic Equations", pp. 227–229 (the q form)
  [section and pages CHECKED, equation numbers CONFIRM IN PRINT]. Second: Higham 2002, §1.8 "Solving a
  Quadratic Equation", p. 10 [CHECKED]; Goldberg 1991 §1.4, eq. (4) [CONFIRM IN PRINT]; DLMF §1.11(iii)
  "Quadratic Equations", https://dlmf.nist.gov/1.11 (read 2026-11-12) [CHECKED].
Worked example: HAND-CHECKED: x² − 3x + 2 → (1, 2, 2). Cancellation case: x² − 1e8·x + 1 → hi = 1e8,
  lo = 1e-8 to full precision (the textbook formula gives lo with most digits wrong; the test asserts
  |lo − 1e-8| ≤ 2 ulps). Double root: x² − 2x + 1 → (1, 1, 1).
Variants: complex roots are NUM-ALGEB-008.
Best approach: scale a, b, c by the same power of two (exact, from `math.Frexp` of the largest magnitude) so
  b² and 4ac cannot overflow, then the q form; never subtract two nearly equal numbers. Time O(1); Space O(1).
Notes: "n = 1" is decided by D == 0 exactly as computed (accurate to 2 ulps by 006); no tolerance invented.
  Must NOT fall back silently to a linear solve when a = 0.

### QuadraticRoots      (ID reserved: NUM-ALGEB-008)
Formula: as NUM-ALGEB-007 when D ≥ 0. When D < 0: re = −b/(2a), im = √(−D)/(2|a|), roots re ± i·im
  (no cancellation: re and im are computed separately).
Inputs and units: `a, b, c float64`, dimensionless, finite, a ≠ 0.
Output: `(r1, r2 complex128, err error)`, r1 has the smaller real part (ties: smaller imaginary part);
  real roots have imaginary part exactly 0.
Errors: as NUM-ALGEB-007.
Primary source: Numerical Recipes 3rd ed., 2007, §5.6 (complex case) [CONFIRM IN PRINT]; DLMF §1.11(iii)
  [CHECKED].
Worked example: HAND-CHECKED: x² + 2x + 5 → −1 − 2i, −1 + 2i. x² − 3x + 2 → 1 + 0i, 2 + 0i.
Variants: complex coefficients LATER.
Best approach: same scaling and q form; complex branch has no cancellation (re and im computed separately).
  Time O(1); Space O(1).
Notes: must NOT use `cmplx.Sqrt` of a negative real D (works, but hides the exact-zero imaginary part rule).

### CubicRoots      (ID reserved: NUM-ALGEB-009)
Formula: a·x³ + b·x² + c·x + d, a ≠ 0. Divide by a: x³ + A·x² + B·x + C. Q = (A² − 3B)/9,
  R = (2A³ − 9AB + 27C)/54. If R² < Q³ (three real roots): θ = arccos(R/√Q³),
  x_k = −2√Q·cos((θ + 2πk)/3) − A/3, k = 0, 1, 2. Else: S = −sgn(R)·(|R| + √(R² − Q³))^(1/3),
  T = Q/S (T = 0 if S = 0), real root x₁ = (S + T) − A/3; then the other two come from deflating to a quadratic.
  Equivalent Cardano form: DLMF eqs. 1.11.12–1.11.15.
Inputs and units: `a, b, c, d float64`, dimensionless, finite, a ≠ 0.
Output: `([3]complex128, error)`, sorted by real part then imaginary part; real roots have imaginary part 0.
Errors: a = 0 → ErrDegenerate; NaN/Inf → ErrNotFinite; Q³ or R² overflow → removed by power-of-two scaling of
  x (substitute x = 2^k·y, exact), else ErrOverflow.
Primary source: Numerical Recipes 3rd ed., 2007, §5.6 "Quadratic and Cubic Equations", pp. 227–229
  [section CHECKED, equation numbers CONFIRM IN PRINT]. Second: DLMF §1.11(iii) "Cubic Equations", eqs.
  1.11.12–1.11.15, https://dlmf.nist.gov/1.11 (read 2026-11-12) [CHECKED]. Polishing and deflation: Kahan,
  "To Solve a Real Cubic Equation", lecture notes, UC Berkeley, 1986,
  https://people.eecs.berkeley.edu/~wkahan/Math128/Cubic.pdf [NOT YET READ — read before build].
Worked example: DLMF §1.11(iii) Example: z³ − 6z² + 6z − 2 → 2 + ∛4 + ∛2 ≈ 4.84732210186307 and
  2 + ∛4·ρ + ∛2·ρ², its conjugate, ≈ 0.57633894906846 ± 0.28360600102688i, ρ = e^{2πi/3}. The test builds the
  expected values from `math.Cbrt` exactly as DLMF writes them. HAND-CHECKED: x³ − 6x² + 11x − 6 → 1, 2, 3.
Variants: real roots only (`CubicRealRoots`) LATER if a caller asks.
Best approach: closed form for one real root, then at most 3 Newton steps with NUM-ALGEB-002, each kept only
  if |p| decreases; deflate by synthetic division to a quadratic; solve that with NUM-ALGEB-008. Time O(1);
  Space O(1).
Notes: the textbook Cardano form loses accuracy near multiple roots and when terms cancel; Newton polishing
  restores the simple roots, multiple roots stay at ~ε^(1/m) (that is the problem's conditioning). Must NOT
  call `math.Pow(x, 1.0/3)` on a negative x (NaN): use `math.Cbrt`.

### PolynomialRoots      (ID reserved: NUM-ALGEB-010)
Formula: all n complex zeros of p(x) = Σ a_i x^i by the Aberth–Ehrlich simultaneous iteration:
  for each approximation z_k: N_k = p(z_k)/p′(z_k); z_k ← z_k − N_k / (1 − N_k·Σ_{j≠k} 1/(z_k − z_j)).
  Start points: Bini's points on circles from the Newton polygon of |a_i|; stop for z_k when
  |p(z_k)| ≤ 4u·Σ|a_i||z_k|^i (the value is below its own rounding error, Higham §5.1 bound).
Inputs and units: `coeffs []float64`, ascending, dimensionless, finite, len ≥ 2, coeffs[n] ≠ 0.
Output: `([]complex128, error)`, len n, sorted by real part then imaginary part; exact zeros for each leading
  zero coefficient a_0 = … = a_{k−1} = 0 (factored out exactly before iterating).
Errors: len < 2 → ErrEmpty; coeffs[n] = 0 → ErrDegenerate; NaN/Inf → ErrNotFinite; not converged within
  the iteration cap → ErrNoConvergence (the best approximations are still returned).
Primary source: Aberth, "Iteration methods for finding all zeros of a polynomial simultaneously",
  Mathematics of Computation, 1973, 27(122):339–344, doi:10.1090/S0025-5718-1973-0329236-7 [CHECKED];
  Ehrlich, "A modified Newton method for polynomials", Communications of the ACM, 1967, 10(2):107–108,
  doi:10.1145/363067.363115 [CHECKED]; Bini, "Numerical computation of polynomial zeros by means of Aberth's
  method", Numerical Algorithms, 1996, 13:179–200, doi:10.1007/BF02207694 (start points, stopping rule)
  [CHECKED].
Worked example: DLMF §1.11(iii) quartic Example, https://dlmf.nist.gov/1.11 (read 2026-11-12):
  z⁴ − 4z³ + 5z + 2 → (3 ± √17)/2 and (1 ± √5)/2, i.e. −0.6180339887498949, −0.5615528128088303,
  1.618033988749895, 3.5615528128088303. HAND-CHECKED: x⁴ − 1 → −1, −i, i, 1; (x−1)(x−2)…(x−10) → 1…10
  within the bound set by the conditioning (Wilkinson, "The evaluation of the zeros of ill-conditioned
  polynomials. Part I", Numerische Mathematik, 1959, 1:150–166, doi:10.1007/BF01386381 [DOI CONFIRM]).
Variants: Laguerre with deflation (Numerical Recipes 3rd ed. §9.5) and companion-matrix eigenvalues
  (Edelman, Murakami, "Polynomial roots from companion matrix eigenvalues", Mathematics of Computation,
  1995, 64(210):763–776, doi:10.1090/S0025-5718-1995-1262279-2 [DOI CONFIRM]) were weighed; see Q3.
Best approach: Aberth–Ehrlich. Convergence: cubic near simple roots, linear near multiple roots; no proof of
  global convergence, but with Bini's start points it is the method of MPSolve and fails in practice only on
  adversarial inputs, hence the cap and ErrNoConvergence. No deflation, so no accumulated deflation error.
  Time O(n²·k), n = degree, k = sweeps (cap proposed: 50 + n, reported in the doc comment); Space O(n), the
  result plus one []bool of converged flags.
Notes: real coefficients only; roots of real polynomials come in conjugate pairs, which the output
  enforces by symmetrising pairs whose imaginary parts are within the stopping bound (that bound is the
  source's, not ours). Must NOT promise more accuracy than the conditioning allows; must NOT use random start
  points (purity: no randomness unless passed in).

## Routed: researched here for the proposed `num/integer` (SHORT `INTEG`), not built here

Row 27 is MOVED. If Devin approves `num/integer` (Q1), the owner takes these three entries as written.
IDs are PROVISIONAL until the package exists in `doc/taxonomy.md`.

### IsPrime      (provisional ID: NUM-INTEG-001)
Formula: n < 2 → false; trial division by primes ≤ 37; then strong-probable-prime (Miller–Rabin) tests to
  the 12 bases 2, 3, 5, …, 37: write n − 1 = 2^s·d, d odd; n passes base a if a^d ≡ 1 or a^(2^r·d) ≡ −1
  (mod n) for some 0 ≤ r < s. Prime iff it passes all 12.
Inputs and units: `n uint64`, dimensionless.
Output: `bool`. Errors: none (every uint64 has an answer).
Primary source: deterministic for n < ψ₁₂ = 318665857834031151167461 > 2⁶⁴: Sorenson, Webster, "Strong
  pseudoprimes to twelve prime bases", Mathematics of Computation, 2017, 86(304):985–1003,
  doi:10.1090/mcom/3134 [CHECKED]. Test: Miller, "Riemann's hypothesis and tests for primality", J. Computer
  and System Sciences, 1976, 13(3):300–317, doi:10.1016/S0022-0000(76)80043-8; Rabin, "Probabilistic
  algorithm for testing primality", J. Number Theory, 1980, 12(1):128–138, doi:10.1016/0022-314X(80)90084-0
  [both DOIs CONFIRM]; Knuth TAOCP Vol. 2, 3rd ed., 1997, §4.5.4, Algorithm P [CHECKED].
Worked example: Pomerance, Selfridge, Wagstaff, "The pseudoprimes to 25·10⁹", Mathematics of Computation,
  1980, 35(151):1003–1026, doi:10.1090/S0025-5718-1980-0572872-7 [CHECKED]: 3215031751 = 151·751·28351 is
  a strong pseudoprime to bases 2, 3, 5, 7 → IsPrime = false. HAND-CHECKED: 2⁶⁴ − 59 = 18446744073709551557
  (largest prime below 2⁶⁴) → true; the test cross-checks with `math/big.Int.ProbablyPrime(0)`, which Go
  documents as exact below 2⁶⁴.
Best approach: modular multiply with `math/bits.Mul64` + `bits.Div64` (128-bit product, no overflow).
  Time O(log n) multiplications per base, 12 bases; Space O(1).
Notes: `math/big.ProbablyPrime(0)` exists in the standard library but allocates; ours does not. Must NOT use
  random bases.

### PrimeFactors      (provisional ID: NUM-INTEG-002)
Formula: trial division by small primes (to 2¹⁰, say), then for each remaining composite cofactor m (found
  by IsPrime) a factor by Pollard's rho, f(x) = x² + c mod m, with Brent's cycle detection and batched gcds;
  recurse on both parts. Sort ascending.
Inputs and units: `n uint64`, n ≥ 2. Output: `([]uint64, error)` prime factors with multiplicity, ascending.
Errors: n < 2 → ErrOutOfRange.
Primary source: Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15:331–334,
  doi:10.1007/BF01933667 [CHECKED]; Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980,
  20:176–184, doi:10.1007/BF01933190 [CHECKED]; Knuth TAOCP Vol. 2, 3rd ed., 1997, §4.5.4, Algorithms A and B
  [CHECKED].
Worked example: HAND-CHECKED: 2⁶⁴ − 1 = (2³² − 1)(2³² + 1) = 3·5·17·257·65537 · 641·6700417
  → [3, 5, 17, 257, 641, 65537, 6700417] (641·6700417 = 4294967297 checked by hand).
Best approach: as Formula. Time O(n^(1/4)) multiplications expected per split (heuristic, Pollard 1975 and
  Brent 1980; not proven), O(log n) splits; Space O(log n), ≤ 64 factors.
Notes: deterministic: c runs 1, 2, 3, … instead of random (a fixed sequence is allowed; randomness is not).

### SmallestPrimeFactor      (provisional ID: NUM-INTEG-003)
Formula: the least element of PrimeFactors(n) (trial division returns it at once when it is small).
Inputs and units: `n uint64`, n ≥ 2. Output: `(uint64, error)`. Errors: n < 2 → ErrOutOfRange.
Primary source: as NUM-INTEG-002.
Worked example: HAND-CHECKED: 2⁶⁴ − 1 → 3. 4294967291 · 4294967279 = 18446743979220271189 → 4294967279
  (both factors prime and the product checked by hand-run script 2026-11-12; the test re-checks with `math/big`).
Best approach: trial division up to a small bound, then PrimeFactors. Plain trial division to √n is up to
  2³² divisions for a 64-bit semiprime and is rejected. Time and Space as NUM-INTEG-002.
Notes: Pollard finds A factor, not the smallest, so the full factorisation is needed in the worst case.

## Step 3: build order and dependencies

| Order | ID | Calls | S1 / other subjects |
|---|---|---|---|
| 1 | NUM-ALGEB-001 PolynomialEvaluate | `math.FMA` | none |
| 2 | NUM-ALGEB-002 PolynomialValueAndDerivative | `math.FMA` | none |
| 3 | NUM-ALGEB-003 PolynomialMultiply | `math.FMA` | none |
| 4 | NUM-ALGEB-004 LinearEvaluate | `math.FMA` | none |
| 5 | NUM-ALGEB-005 LinearSolve | — | none |
| 6 | NUM-ALGEB-006 QuadraticDiscriminant | `math.FMA`, `math.Frexp`, `math.Ldexp` | none |
| 7 | NUM-ALGEB-007 QuadraticRealRoots | 006 | none |
| 8 | NUM-ALGEB-008 QuadraticRoots | 006 | none |
| 9 | NUM-ALGEB-009 CubicRoots | 002, 008, `math.Cbrt` | none |
| 10 | NUM-ALGEB-010 PolynomialRoots | 001/002 (complex Horner, unexported), `slices.SortFunc` | none |

No S1 algorithm is needed (sorting ≤ n roots uses the standard library). Nothing here is BLOCKED.
Consumers in other subjects: num/geometry's tangent line needs 002 (Q6). The companion-matrix alternative for
010 would have been BLOCKED ON num/linalg; Aberth avoids that. The routed num/integer entries depend on
nothing but `math/bits`.
