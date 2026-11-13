# S5: Algebra (`num/algebra`)

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/` (hand-off: `notes_algebra.md`)
Decisions: Devin's request of 2026-11-12 (catalog Q2, Q3, Q4, Q5, Q6; D-ids assigned at merge) · Assumptions: none

Detail for every entry (formula, inputs, alternatives) is in [../catalog/algebra.md](../catalog/algebra.md), Step 2.
Conventions: coefficients ascending (`coeffs[i]` multiplies xⁱ); every value dimensionless; "no real roots" is a
result, not an error; a zero leading coefficient is `ErrDegenerate`. Sentinels: `ErrEmpty`, `ErrNotFinite`,
`ErrOverflow`, `ErrDegenerate`, `ErrNoConvergence`, each wrapped with the parameter name.
Imports: standard library only (`math`, `math/cmplx`, `slices`, `cmp`, `errors`, `fmt`). `num/geometry` may import
this package (tangent line uses NUM-ALGEB-002).

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| `NUM-ALGEB-001` | `PolynomialEvaluate` | `func PolynomialEvaluate(coeffs []float64, x float64) (float64, error)` | dimensionless |
| `NUM-ALGEB-002` | `PolynomialValueAndDerivative` | `func PolynomialValueAndDerivative(coeffs []float64, x float64) (value, derivative float64, err error)` | derivative per unit of x |
| `NUM-ALGEB-003` | `PolynomialMultiply` | `func PolynomialMultiply(a, b []float64) ([]float64, error)` | dimensionless |
| `NUM-ALGEB-004` | `LinearEvaluate` | `func LinearEvaluate(slope, intercept, x float64) (float64, error)` | slope y-unit per x-unit |
| `NUM-ALGEB-005` | `LinearSolve` | `func LinearSolve(slope, intercept, y float64) (float64, error)` | result in x-unit |
| `NUM-ALGEB-006` | `QuadraticDiscriminant` | `func QuadraticDiscriminant(a, b, c float64) (float64, error)` | dimensionless |
| `NUM-ALGEB-007` | `QuadraticRealRoots` | `func QuadraticRealRoots(a, b, c float64) (lo, hi float64, n int, err error)` | dimensionless |
| `NUM-ALGEB-008` | `QuadraticRoots` | `func QuadraticRoots(a, b, c float64) (r1, r2 complex128, err error)` | dimensionless |
| `NUM-ALGEB-009` | `CubicRoots` | `func CubicRoots(a, b, c, d float64) ([3]complex128, error)` | dimensionless |
| `NUM-ALGEB-010` | `PolynomialRoots` | `func PolynomialRoots(coeffs []float64) ([]complex128, error)` | dimensionless |

## Per function

### `NUM-ALGEB-001` `PolynomialEvaluate`
- **Source**: Knuth TAOCP Vol. 2, 3rd ed., 1997, §4.6.4, pp. 486–488; Higham 2002 §5.1, p. 94; DLMF §1.11(i).
- **Best known**: O(n) time, n multiply-adds (optimal for a general polynomial, Knuth §4.6.4); O(1) space.
- **Approach**: Horner with `math.FMA`. **Lost**: power-by-power with `math.Pow` (n log n work, worse rounding).
- **Input modified**: no. **Worked example**: [−1, 2, −6, 2] at 3 → 5 (hand-checked; no numeric example in the sources).
- **Errors**: `ErrEmpty` (empty), `ErrNotFinite` (NaN/Inf x or coefficient), `ErrOverflow` (non-finite value).

### `NUM-ALGEB-002` `PolynomialValueAndDerivative`
- **Source**: DLMF §1.11(i) extended Horner scheme; Higham 2002 §5.2.
- **Best known**: O(n), 2n FMAs, O(1) space. **Approach**: one pass, two recurrences. **Lost**: forming derivative
  coefficients (O(n) allocation, second pass).
- **Input modified**: no. **Worked example**: [−1, 2, −6, 2] at 3 → (5, 20), hand-checked.
- **Errors**: as 001; `ErrOverflow` when either result is not finite.

### `NUM-ALGEB-003` `PolynomialMultiply`
- **Source**: Cormen et al., 3rd ed., 2009, §30.1, eqs. (30.1)–(30.2).
- **Best known**: O(n log n) by FFT; schoolbook O(n·m) chosen. **Lost**: FFT and Karatsuba: faster only for degree
  in the hundreds, and FFT loses accuracy on coefficients of mixed magnitude (catalog, 003 Variants).
- **Input modified**: no. **Worked example**: (2x + 3)(x − 4) → [−12, −5, 2], hand-checked.
- **Errors**: `ErrEmpty`, `ErrNotFinite`, `ErrOverflow`.

### `NUM-ALGEB-004` `LinearEvaluate` · `NUM-ALGEB-005` `LinearSolve`
- **Source**: Abramson, College Algebra 2e, OpenStax, 2021, §4.1 (slope-intercept form); IEEE 754-2019 cl. 5.4.1 (FMA).
- **Best known**: O(1). **Approach**: one FMA; (y − b)/m. **Lost**: m*x + b in two roundings.
- **Worked example**: m = 2, b = 3, x = 4 → 11; y = 11 → 4 (hand-checked).
- **Errors**: `ErrNotFinite`, `ErrOverflow`; 005 also `ErrDegenerate` (slope == 0 exactly, no tolerance).

### `NUM-ALGEB-006` `QuadraticDiscriminant`
- **Source**: Jeannerod, Louvet, Muller 2013 (≤ 2u with FMA); Boldo 2009; Goldberg 1991 §1.4, p. 10 (example).
- **Best known**: O(1), 3 FMAs plus exact power-of-two scaling (Frexp/Ldexp). **Lost**: naive b*b − 4*a*c
  (loses every digit when b² ≈ 4ac); a single common scale factor (underflows the smaller term needlessly).
- **Worked example**: Goldberg's a = 1.22, b = 3.34, c = 2.28 → 0.0292, within 2u of the exact value of the rounded inputs.
- **Errors**: `ErrNotFinite`; `ErrOverflow` when D itself is out of range.

### `NUM-ALGEB-007` `QuadraticRealRoots` · `NUM-ALGEB-008` `QuadraticRoots`
- **Source**: Numerical Recipes in C, 2nd ed., 1992, §5.6, eqs. (5.6.4)–(5.6.5); Higham 2002 §1.8; DLMF §1.11(iii).
- **Best known**: O(1). **Approach**: q = −½(b + sgn(b)√D) on scaled values; complex branch re and im separately.
  **Lost**: textbook (−b ± √D)/(2a) (cancellation).
- **Worked example**: x² − 3x + 2 → 1, 2; x² − 10⁸x + 1 → 10⁻⁸, 10⁸ within 2 ulps; x² + 2x + 5 → −1 ± 2i.
- **Errors**: `ErrDegenerate` (a = 0), `ErrNotFinite`, `ErrOverflow`. D < 0 in 007 is n = 0, not an error.

### `NUM-ALGEB-009` `CubicRoots`
- **Source**: Numerical Recipes in C, 2nd ed., §5.6, eqs. (5.6.10)–(5.6.12), (5.6.15)–(5.6.17); DLMF §1.11(iii); Kahan 1986.
- **Best known**: O(1). **Approach**: exact power-of-two scaling of x; Viète (three real: keep only the root of
  largest modulus) or Cardano (one real) for one real root x₁; ≤ 3 Newton steps on the scaled cubic and on the
  original coefficients (power-of-two normalised); deflation to a·x² + f1·x + f0 backward (from d) when
  |a|·|x₁|³ ≥ |d|, forward otherwise, so the larger-root division never cancels; 008 for the quadratic (small root
  by Vieta); ≤ 3 Newton steps (complex for a pair) on each quadratic root. Fix 2026-11-12: the earlier version
  deflated forward by the largest root and took Viète's small root directly, both cancelling; see
  `doc/math_build_2026-11-12/fix_algebra.md`. **Lost**: Cardano for all three roots (unstable under roundoff,
  Kahan 1986); PolynomialRoots (iterative, allocates).
- **Worked example**: DLMF Example z³ − 6z² + 6z − 2. Widely spread cases checked against 8192-bit references.
- **Errors**: `ErrDegenerate` (a = 0), `ErrNotFinite`, `ErrOverflow`.

### `NUM-ALGEB-010` `PolynomialRoots`
- **Source**: Aberth 1973; Ehrlich 1967; Bini 1996 (start points equispaced on circles with radii from the Newton
  polygon); Graillat & Ménissier-Morain 2012 (compensated Horner in complex arithmetic, Algorithms 2.1, 3.1, 3.4,
  5.2, 5.4, Theorem 5.5).
- **Best known**: O(n²) per sweep; k sweeps, k ≤ 100 + n. **Approach**: Aberth–Ehrlich, Gauss–Seidel updates,
  p and p′ by compensated Horner (twice-working-precision accuracy), reversed-polynomial evaluation for |z| > 1.
  Convergence test, our own: |Δ| ≤ 4u·|z|, or the correction stops decreasing while |p(z)| ≤ (6n·u)²·Σ|aᵢ||z|ⁱ.
  Angle offset σ = 0.7: our own choice (any value off the conjugation-symmetric set works). Real snap when
  |p(Re z)| ≤ |p(z)|; conjugate pairing after. A root whose start radius underflows is returned as 0 (not an
  error). Fix 2026-11-12: Bini's residual stopping rule (weights 1 + 3.8i) and σ = 0.7 were sourced only from
  Bini's code comments; the residual rule accepted points 0.5 from Wilkinson-20's roots. **Lost**: Laguerre with
  deflation (deflation error), companion-matrix eigenvalues (O(n³), needs num/linalg), Jenkins–Traub (size);
  plain Horner (Wilkinson-20 roots only to ~1e-2, measured).
- **Worked example**: DLMF quartic Example z⁴ − 4z³ + 5z + 2. Wilkinson-20 against 8192-bit Newton references.
- **Errors**: `ErrEmpty` (len < 2), `ErrDegenerate` (leading 0), `ErrNotFinite`, `ErrOverflow` (n·Σ|aᵢ| or a root
  above the float64 range), `ErrNoConvergence` (approximations still returned).

## Out of scope

Complex coefficients; k-th derivatives (extended Horner, DLMF §1.11(i)); running error bound variant
(Higham §5.1); FFT multiplication; real-roots-only cubic; percent change and ratio (Basic Operations pass,
IDs from NUM-ALGEB-011). e^x−1 and scalb are dropped everywhere (`math.Expm1`, `math.Ldexp`).
