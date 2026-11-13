# Catalog: matrices (`num/linalg`, SHORT `LINALG`): the gaps after the linear-algebra pass

Catalog and build pass, 2026-11-12. Input: the vCalc "Matrix" list, NAMES ONLY (no vCalc page was fetched).
Work folder: [../matrix_2026-11-12/](../matrix_2026-11-12/README.md). Conventions inherited unchanged from
[linalg.md](linalg.md): row-major `[]float64` with dimensions alongside; results in a `dst` reused when its
capacity suffices; fixed-size Mat2/Mat3/Mat4 values that never allocate; tolerances are parameters with a
documented default (a negative tol selects it); NaN/±Inf in → ErrNotFinite; never a silent NaN.

## Step 1: clean (42 rows)

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| 1 | Calculators: Cramer's Rule Calculator | DROP | — | a calculator; bundles the Cramer rows (#16–#20), built as NUM-LINALG-144/145 |
| 2 | Calculators: Matrix Multiplication | DROP | — | a calculator; bundles MatMul (NUM-LINALG-090) |
| 3 | Calculators: Understanding Eigenvalues and Eigenvectors of a 2x2 Matrix | DROP | — | a teaching page; its content is Mat2Eigen (NUM-LINALG-149) |
| 4 | 2X2 Determinant | REUSED | NUM-LINALG-060 Mat2Det | |
| 5 | 3 Equations 3 Unknowns | REUSED | NUM-LINALG-096 MatSolve (and NUM-LINALG-145 Mat3SolveCramer) | the general stable solver; Cramer is the tiny closed form |
| 6 | 3x3 Determinant | REUSED | NUM-LINALG-069 Mat3Det | |
| 7 | 3x3 Matrix Characteristics | KEEP | NUM-LINALG-135 Mat3Characteristics | a bundle of existing results in one struct |
| 8 | 3x4 times 4x5 | MERGE | → NUM-LINALG-090 MatMul | one general product |
| 9 | 4-Vector Norm | REUSED | NUM-LINALG-041 VecNorm | |
| 10 | 4x4 Determinant | REUSED | NUM-LINALG-081 Mat4Det | |
| 11 | 5-Vector Norm | MERGE | → NUM-LINALG-041 VecNorm | same function |
| 12 | Characteristic Polynomial of a 2x2 Matrix | KEEP | NUM-LINALG-146 Mat2CharPoly | |
| 13 | Characteristic Polynomial of a 3x3 Matrix | KEEP | NUM-LINALG-147 Mat3CharPoly | |
| 14 | Characteristics of a 2x2 Matrix | KEEP | NUM-LINALG-134 Mat2Characteristics | |
| 15 | Cramer's Rule (three equations) | KEEP | NUM-LINALG-145 Mat3SolveCramer | the five Cramer rows are two functions |
| 16 | Cramer's Rule (three equations, solved for z) | MERGE | → #15 | one component of the same solve |
| 17 | Cramer's Rule (two equations) | KEEP | NUM-LINALG-144 Mat2SolveCramer | |
| 18 | Cramer's Rule (two equations, solved for x) | MERGE | → #17 | |
| 19 | Cramer's Rule (two equations, solved for y) | MERGE | → #17 | |
| 20 | Determinant of a 2-by-2 Matrix | MERGE | → #4 (NUM-LINALG-060) | the four 2×2 determinant rows are one |
| 21 | Determinant of a 2x2 Matrix | MERGE | → #4 | |
| 22 | Determinant of a 3-by-3 Matrix | MERGE | → #6 (NUM-LINALG-069) | |
| 23 | Eigenvalues of a 2x2 Matrix | KEEP | NUM-LINALG-149 Mat2Eigen | general 2×2: real or complex pair, with eigenvectors |
| 24 | Inverse of a 2x2 Matrix | REUSED | NUM-LINALG-061 Mat2Inverse | |
| 25 | Inverse of a 3x3 Matrix | REUSED | NUM-LINALG-070 Mat3Inverse | |
| 26 | Matrix Test | DROP | — | a test page with no formula |
| 27 | Minor of a 3x3 Matrix | KEEP | NUM-LINALG-123 Mat3Minor | |
| 28 | Multiplication of a 2x2 Matrix by a 2x1 Matrix | REUSED | NUM-LINALG-058 Mat2MulVec | |
| 29 | Multiplication of a 2x2 Matrix by a Scalar | KEEP | NUM-LINALG-118 MatScale | one n-dimensional scale (fixed sizes: multiply the array) |
| 30 | Multiplication of a 2x3 Matrix and a 3x2 Matrix | MERGE | → NUM-LINALG-090 MatMul | |
| 31 | Multiplication of a 3x3 Matrix and a 3x1 Matrix | REUSED | NUM-LINALG-067 Mat3MulVec | |
| 32 | Multiplication of a 3x3 Matrix by a Scalar | MERGE | → #29 (NUM-LINALG-118) | |
| 33 | Multiply two 3x3 Matrices | REUSED | NUM-LINALG-066 Mat3Mul | |
| 34 | Row Echelon Form of a 3x3 Matrix | KEEP | NUM-LINALG-130 MatRowEchelon (and 131 MatRREF) | any size |
| 35 | Trace of 2X2 Matrix | REUSED | NUM-LINALG-062 Mat2Trace | the four trace rows are two IDs |
| 36 | Trace of 3X3 Matrix | REUSED | NUM-LINALG-071 Mat3Trace | |
| 37 | Trace of a 2x2 Matrix | MERGE | → #35 | |
| 38 | Trace of a 3x3 Matrix | MERGE | → #36 | |
| 39 | Transpose of a 3x3 Matrix | REUSED | NUM-LINALG-068 Mat3Transpose | |
| 40 | V3 - Vector Norm | REUSED | NUM-LINALG-017 Vec3Norm | |
| 41 | Vector Norm | MERGE | → NUM-LINALG-041 VecNorm | |
| 42 | WikiClip: Matrix Calculator List | DROP | — | a list of links |

Counts: KEEP 9 · REUSED 15 · MERGE 14 · DROP 4. (IDs of REUSED rows checked against INDEX.md before writing.)

## Step 2: the standard set (the real scope)

Every entry's full research (signature, formula, inputs, output, errors, source to the section, worked example,
variants, time and space) is its doc comment, the source of truth; Part 3 below lists them, generated from the
code. Sources used: Golub & Van Loan, Matrix Computations, 4th ed. (sections confirmed 2026-11-12: §2.3 Matrix
Norms, §3.1 Triangular Systems, §4.2 Positive Definite Systems, §5.2 The QR Factorization, §8.5 Jacobi Methods,
§8.6 Computing the SVD); Strang, Introduction to Linear Algebra, 5th ed. (whole table of contents confirmed:
§2.2, 2.3, 3.2, 3.3, 3.4, 3.5, 4.3, 4.4, 5.2, 5.3, 6.1, 6.4, 6.5, ch. 7, 11.2); Householder 1958 (JACM 5(4):339–342);
Hestenes 1958 (J SIAM 6(1):51–90); Rutishauser 1966 (Numer Math 9:1–10); Demmel & Veselić 1992 (SIMAX
13(4):1204–1245); Van Loan 2000 (J Comput Appl Math 123:85–100, Kronecker); Styan 1973 (LAA 6:217–240, Hadamard);
Knuth TAOCP Vol. 2 §4.6.3 (powers by squaring); Faddeev & Faddeeva 1963 (Leverrier–Faddeev); NIST e-Handbook
§6.5.4.1 and NIST Dataplot "CORRELATION" (covariance and correlation).

| Group | Function (ID) | Best approach (time, space) |
|---|---|---|
| Elementwise and products | MatAdd 116 · MatSub 117 · MatScale 118 · MatHadamard 119 | O(rc), 0 extra |
| | VecOuter 120 · MatKronecker 121 | O(mn), O(mnpq); result only |
| | MatPower 122 | O(n³ log k) by binary powering, two n² buffers |
| Structure | Mat3Minor 123 · Mat2/3/4Adjugate 124–126 (closed form, 0 alloc) | O(1) |
| | MatMinor 127 · MatCofactorMatrix 128 · MatAdjugate 129 | O(n³) via LU (det·A⁻¹) when nonsingular; O(n⁵) minors when singular |
| | MatRowEchelon 130 · MatRREF 131 (partial pivoting, tol) | O(r·c·min(r, c)) |
| | MatNullSpace 132 · MatColumnSpace 133 | from the RREF: O(r·c·min(r, c)) |
| | Mat2Characteristics 134 · Mat3Characteristics 135 | O(1); built from the existing functions |
| Checks (tolerance) | MatIsSymmetric 136 · MatIsDiagonal 137 · MatIsUpperTriangular 138 · MatIsLowerTriangular 139 · MatIsOrthogonal 140 · MatIsPositiveDefinite 141 · MatIsPositiveSemidefinite 142 · MatIsSingular 143 | O(n²) checks; O(n³) for orthogonal, definite, singular |
| Small solvers | Mat2SolveCramer 144 · Mat3SolveCramer 145 | O(1); MatSolve (096) is the stable general choice |
| Characteristic polynomial | Mat2CharPoly 146 · Mat3CharPoly 147 · MatCharPolyFaddeevLeVerrier 148 | O(1); O(n⁴) |
| Eigen | Mat2Eigen 149 · MatSymmetricEigenJacobi 150 · MatPowerIteration 151 | O(1); O(n³) per sweep; O(n²) per step |
| Factorisations | Cholesky 152 · CholeskySolve 153 | O(n³/3); O(n²) |
| | QRHouseholder 154 · GramSchmidtModified 155 · LeastSquaresQR 156 | O(mn²) |
| | ForwardSubstitution 157 · BackSubstitution 158 · TridiagonalSolve 159 | O(n²), O(n²), O(n) |
| | SVDJacobi 160 · MatPseudoInverse 161 · MatCondition2 162 | O(mn²) per sweep |
| Matrix norms | MatNormFrobenius 163 · MatNorm1 164 · MatNormInf 165 · MatNorm2 166 | O(rc); 166 by the SVD |
| Statistics bridges | CovarianceToCorrelation 167 · CorrelationToCovariance 168 | O(n²) |

DROPPED from the standard set: "is square" (rows == cols is one comparison; D-008 by analogy).
LATER: eigenvalues of a general non-symmetric matrix of order n > 2 (needs Hessenberg reduction and the shifted QR
algorithm; power iteration gives only the dominant pair).
HELD: none.

## Part 3: research entries (generated from the doc comments, the source of truth)

### MatAdd      (NUM-LINALG-116, built: num/linalg/mat_elementwise.go:41)

MatAdd returns the elementwise sum A + B of two rows×cols matrices in dst (dst may be a or b).
Errors: ErrDimension when a length disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
ErrOverflow when a sum is outside the float64 range.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (A + B entry by entry)
- Units: a and b in one unit; the result in that unit
- Example: [1 2; 3 4] + [5 6; 7 8] → [6 8; 10 12] (hand-checked)

### MatSub      (NUM-LINALG-117, built: num/linalg/mat_elementwise.go:63)

MatSub returns the elementwise difference A − B of two rows×cols matrices in dst (dst may be a or b).
Errors: ErrDimension when a length disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
ErrOverflow when a difference is outside the float64 range.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (A − B entry by entry)
- Units: a and b in one unit; the result in that unit
- Example: [5 6; 7 8] − [1 2; 3 4] → [4 4; 4 4] (hand-checked)

### MatScale      (NUM-LINALG-118, built: num/linalg/mat_elementwise.go:85)

MatScale returns k·A for a rows×cols matrix A in dst (dst may be a).
Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf in a or k;
ErrOverflow when a product is outside the float64 range.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (cA multiplies every entry by c)
- Units: a in unit A, k in unit K; the result in A·K
- Example: 3·[1 2; 3 4] → [3 6; 9 12] (hand-checked)

### MatHadamard      (NUM-LINALG-119, built: num/linalg/mat_elementwise.go:114)

MatHadamard returns the Hadamard (entrywise) product A ∘ B, (A ∘ B)ᵢⱼ = aᵢⱼ·bᵢⱼ, of two rows×cols matrices
in dst (dst may be a or b).
Errors: ErrDimension when a length disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
ErrOverflow when a product is outside the float64 range.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
- Source: Styan, "Hadamard products and multivariate statistical analysis", Linear Algebra and its Applications, 1973, 6:217–240, doi:10.1016/0024-3795(73)90023-2 (definition: the elementwise product A∘B)
- Units: a in unit A, b in unit B; the result in A·B
- Example: [1 2; 3 4] ∘ [5 6; 7 8] → [5 12; 21 32] (hand-checked)

### VecOuter      (NUM-LINALG-120, built: num/linalg/mat_products.go:20)

VecOuter returns the outer product u·vᵀ of an m-vector u and an n-vector v, the m×n matrix with element
(i, j) = uᵢ·vⱼ, in dst (which must not overlap u or v).
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a product is outside the float64 range;
ErrDimension when m·n overflows.

- Time: O(m·n)   m = len(u), n = len(v)
- Space: O(1) extra; O(m·n) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (a column times a row is a rank-one matrix)
- Units: u in unit U, v in unit V; the result in U·V
- Example: u = (1, 2), v = (3, 4, 5) → [3 4 5; 6 8 10] (hand-checked)

### MatKronecker      (NUM-LINALG-121, built: num/linalg/mat_products.go:53)

MatKronecker returns the Kronecker product A ⊗ B of the m×n matrix a and the p×q matrix b: the mp×nq
block matrix whose (i, j) block is aᵢⱼ·B, in dst (which must not overlap a or b).
Errors: ErrDimension when a length disagrees with its dimensions or mp·nq overflows; ErrNotFinite for a
NaN or ±Inf element; ErrOverflow when a product is outside the float64 range.

- Time: O(m·n·p·q)
- Space: O(1) extra; O(m·n·p·q) for the result, one allocation only if dst lacks capacity
- Source: Van Loan, "The ubiquitous Kronecker product", Journal of Computational and Applied Mathematics, 2000, 123(1–2):85–100, doi:10.1016/S0377-0427(00)00393-9, §1 (B ⊗ C as the block matrix of bᵢⱼC)
- Units: a in unit A, b in unit B; the result in A·B
- Example: [1 2; 3 4] ⊗ [0 5; 6 7] → [0 5 0 10; 6 7 12 14; 0 15 0 20; 18 21 24 28] (hand-checked)

### MatPower      (NUM-LINALG-122, built: num/linalg/mat_products.go:115)

MatPower returns Aᵏ for an n×n matrix and an integer k ≥ 0 (A⁰ = I) by binary powering (repeated
squaring), in dst (which must not overlap a).
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrOutOfRange when k < 0 or k > 2⁴⁰ (a policy cap, see
maxMatPower); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element leaves the float64 range.

- Time: O(n³·log k)   at most 2⌊log₂ k⌋ matrix products
- Space: O(n²) extra (two work buffers); O(n²) for the result, one allocation only if dst lacks capacity
- Source: Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.3 "Evaluation of Powers", Algorithm A (right-to-left binary method)
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖AB‖ ≤ ‖A‖‖B‖ bounds the growth of each product)
- Units: a in unit A; the result in Aᵏ
- Example: [1 1; 1 0]⁵ → [8 5; 5 3] (Fibonacci numbers; hand-checked)

### Mat3Minor      (NUM-LINALG-123, built: num/linalg/mat_cofactor_fixed.go:36)

Mat3Minor returns the minor M_ij of a 3×3 matrix: the determinant of the 2×2 matrix left after deleting
row i and column j (0-based). The cofactor is (−1)^(i+j)·M_ij.
Errors: ErrOutOfRange when i or j is outside 0–2; ErrNotFinite for a NaN or ±Inf element; ErrOverflow
when the minor is outside the float64 range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" (the cofactor C_ij = (−1)^(i+j) det M_ij, M_ij with row i and column j removed)
- Units: m in unit A; the result in A²
- Example: m = [1 2 3; 4 5 6; 7 8 10], i = 0, j = 0 → 5·10 − 6·8 = 2 (hand-checked)

### Mat2Adjugate      (NUM-LINALG-124, built: num/linalg/mat_cofactor_fixed.go:61)

Mat2Adjugate returns the adjugate (classical adjoint) adj m = Cᵀ = [d −b; −c a] of m = [a b; c d], the
matrix with m·adj m = adj m·m = det(m)·I. It exists for singular m too.
Errors: ErrNotFinite for a NaN or ±Inf element.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A; the 2×2 case [d −b; −c a])
- Units: m in unit A; the result in A
- Example: [1 2; 3 4] → [4 −2; −3 1] (hand-checked)

### Mat3Adjugate      (NUM-LINALG-125, built: num/linalg/mat_cofactor_fixed.go:79)

Mat3Adjugate returns the adjugate adj m = Cᵀ of a 3×3 matrix, each cofactor a 2×2 determinant within 2
ulps (Kahan's algorithm); m·adj m = det(m)·I, and adj m exists for singular m.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a cofactor is outside the float64 range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A)
- Units: m in unit A; the result in A²
- Example: [1 2 3; 0 1 4; 5 6 0] → [−24 18 5; 20 −15 −4; −5 4 1] (hand-checked; det = 1, so this is the inverse)

### Mat4Adjugate      (NUM-LINALG-126, built: num/linalg/mat_cofactor_fixed.go:110)

Mat4Adjugate returns the adjugate adj m = Cᵀ of a 4×4 matrix, each cofactor a 3×3 determinant (a scalar
triple product summed by Dot2); m·adj m = det(m)·I, and adj m exists for singular m.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a cofactor is outside the float64 range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A)
- Units: m in unit A; the result in A³
- Example: diag(1, 2, 3, 4) → diag(24, 12, 8, 6) (hand-checked)

### MatMinor      (NUM-LINALG-127, built: num/linalg/mat_cofactor.go:51)

MatMinor returns the minor M_ij of an n×n matrix (n ≥ 2): the determinant of the (n−1)×(n−1) matrix left
after deleting row i and column j (0-based), by MatDet (LU with partial pivoting).
Errors: ErrDimension when n < 2 or len(a) ≠ n²; ErrOutOfRange when i or j is outside 0…n−1; ErrNotFinite
for a NaN or ±Inf entry; ErrOverflow when the minor is outside the float64 range.

- Time: O(n³)   one LU of order n−1
- Space: O((n−1)²) extra (the submatrix and its LU)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" (M_ij: remove row i and column j; C_ij = (−1)^(i+j) det M_ij)
- Units: a in unit A; the result in A^(n−1)
- Example: [1 2 3; 4 5 6; 7 8 10], i = 1, j = 2 → det [1 2; 7 8] = −6 (hand-checked)

### MatCofactorMatrix      (NUM-LINALG-128, built: num/linalg/mat_cofactor.go:80)

MatCofactorMatrix returns the matrix of cofactors C, C_ij = (−1)^(i+j)·M_ij, of an n×n matrix in dst. When
A is nonsingular it uses C = det(A)·A⁻ᵀ from one LU (O(n³)); when A is singular (det A = 0) it computes
every minor by its own LU (O(n⁵)), so singular matrices, whose cofactors may still be nonzero, get the
exact definition. n = 1 gives C = [1]. dst must not overlap a.
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOverflow when a
cofactor is outside the float64 range.

- Time: O(n³) when A is nonsingular; O(n⁵) when it is singular
- Space: O(n²) extra; O(n²) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A, so C = det A·(A⁻¹)ᵀ)
- Units: a in unit A; the result in A^(n−1)
- Example: [1 2; 3 4] → [4 −3; −2 1] (hand-checked)

### MatAdjugate      (NUM-LINALG-129, built: num/linalg/mat_cofactor.go:101)

MatAdjugate returns the adjugate adj A = Cᵀ (the transposed cofactor matrix) of an n×n matrix in dst, so
A·adj A = adj A·A = det(A)·I and det(adj A) = det(A)^(n−1). Nonsingular A: adj A = det(A)·A⁻¹ from one LU;
singular A: every cofactor by its own LU (O(n⁵)). n = 1 gives [1]. dst must not overlap a.
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOverflow when an
entry is outside the float64 range.

- Time: O(n³) when A is nonsingular; O(n⁵) when it is singular (n minors of order n−1 each; an O(n³) route through the SVD exists and is LATER)
- Space: O(n²) extra; O(n²) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (A·Cᵀ = (det A)·I)
- Units: a in unit A; the result in A^(n−1)
- Example: [1 2; 2 4] (singular) → [4 −2; −2 1] (hand-checked: A·adj A = 0 = det A·I)

### MatRowEchelon      (NUM-LINALG-130, built: num/linalg/echelon.go:93)

MatRowEchelon returns a row echelon form of a rows×cols matrix by Gaussian elimination with partial
pivoting (the largest |entry| in the column becomes the pivot), in dst, with the pivot columns. An entry
with |x| ≤ tol·max|aᵢⱼ| is treated as zero (tol < 0 selects max(rows, cols)·2⁻⁵²) and set to exactly 0, so
the staircase shape is exact. The echelon form is not unique (it depends on the pivoting); the pivot
columns and their number (the rank) are. dst must not overlap a.
Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(r·c·min(r, c))   r = rows, c = cols
- Space: O(min(r, c)) extra for the pivot list; O(r·c) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.2 "The Idea of Elimination" and §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (echelon form U, pivot columns)
- Units: a in unit A; the result in A
- Example: [1 2 1; 2 4 0; 3 6 3] → pivots in columns 0 and 2 (rank 2) (hand-checked)

### MatRREF      (NUM-LINALG-131, built: num/linalg/echelon.go:117)

MatRREF returns the reduced row echelon form R of a rows×cols matrix (each pivot 1, zeros above and below
it) by Gauss–Jordan elimination with partial pivoting, in dst, with the pivot columns. |x| ≤ tol·max|aᵢⱼ|
is treated as zero (tol < 0 selects max(rows, cols)·2⁻⁵²). R is unique in exact arithmetic; in floating
point the decision "zero or pivot" is the tolerance's. The RREF of an invertible matrix is I. dst must not
overlap a.
Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(r·c·min(r, c))   r = rows, c = cols
- Space: O(min(r, c)) extra; O(r·c) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (the reduced row echelon form R = rref(A)) and §3.3 "The Complete Solution to Ax = b"
- Units: a in unit A; the result dimensionless (pivot rows are scaled to 1)
- Example: [1 2 2 4; 3 8 6 16] → R = [1 0 2 0; 0 1 0 2], pivots 0 and 1 (Strang §3.2 example; hand-checked)

### MatNullSpace      (NUM-LINALG-132, built: num/linalg/echelon.go:141)

MatNullSpace returns a basis of the null space {x : A·x = 0} of a rows×cols matrix as the columns of a
cols×k matrix in dst (k = cols − rank), with k: from R = rref(A), one special solution per free column f
(x_f = 1, x at each pivot column p = −R[row of p, f], other free entries 0). The basis is not orthonormal;
orthonormalise it with QRHouseholder if needed. tol as in MatRREF. dst must not overlap a; k = 0 returns
a zero-length dst.
Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(r·c·min(r, c))   the RREF dominates
- Space: O(r·c) extra (the RREF); O(c·k) for the result
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (special solutions: one per free variable; the nullspace matrix N)
- Units: a in unit A; the basis dimensionless
- Example: [1 2 2 4; 3 8 6 16] → basis columns (−2, 0, 1, 0) and (0, −2, 0, 1) (Strang §3.2 example; hand-checked)

### MatColumnSpace      (NUM-LINALG-133, built: num/linalg/echelon.go:179)

MatColumnSpace returns a basis of the column space of a rows×cols matrix: the pivot columns of A itself
(found by MatRREF's pivoting), as the columns of a rows×rank matrix in dst, with the rank. tol as in
MatRREF. dst must not overlap a.
Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(r·c·min(r, c))   the RREF dominates
- Space: O(r·c) extra (the RREF); O(r·rank) for the result
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.5 "Dimensions of the Four Subspaces" (the pivot columns of A are a basis for its column space C(A)) and §3.4 "Independence, Basis and Dimension"
- Units: a in unit A; the basis in A
- Example: [1 2 2 4; 3 8 6 16] → columns (1, 3) and (2, 8), rank 2 (Strang §3.2 example; hand-checked)

### Mat2Characteristics      (NUM-LINALG-134, built: num/linalg/mat_characteristics.go:45)

Mat2Characteristics returns the determinant (Mat2Det), trace (Mat2Trace), rank (as MatRank: complete pivoting, tol < 0
selects max(n)·ε·max|aᵢⱼ|), inverse when it exists (Mat2Inverse) and eigenvalues (Mat2Eigen) of a 2×2 matrix in
one result, built from the existing functions.
Errors: ErrNotFinite for a NaN or ±Inf element or tol; ErrOverflow when a quantity is outside the float64
range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (determinant and inverse) and §6.1 "Introduction to Eigenvalues" (trace, determinant and eigenvalues of a 2×2)
- Units: m in unit A; Det in A², Trace and eigenvalues in A, Inverse in 1/A
- Example: [4 1; 2 3] → det 10, trace 7, rank 2, inverse [0.3 −0.1; −0.2 0.4], eigenvalues 2, 5 (hand-checked)

### Mat3Characteristics      (NUM-LINALG-135, built: num/linalg/mat_characteristics.go:88)

Mat3Characteristics returns the determinant (Mat3Det), trace (Mat3Trace), rank (as MatRank: complete pivoting, tol < 0
selects max(n)·ε·max|aᵢⱼ|), inverse when it exists (Mat3Inverse), characteristic polynomial (Mat3CharPoly) and
eigenvalues (the roots of that polynomial by CubicRoots, NUM-ALGEB-009, each accurate to its own
conditioning) of a 3×3 matrix in one result. For a symmetric matrix Mat3SymmetricEigenvalues is the more
accurate eigenvalue source.
Errors: ErrNotFinite for a NaN or ±Inf element or tol; ErrOverflow when a quantity is outside the float64
range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" and §6.1 "Introduction to Eigenvalues" (eigenvalues are the roots of det(A − λI) = 0; Mat3CharPoly returns the monic det(λI − A) = −det(A − λI), same roots)
- Units: m in unit A; Det in A³, Trace and eigenvalues in A, Inverse in 1/A
- Example: diag(1, 2, 3) → det 6, trace 6, rank 3, inverse diag(1, 1/2, 1/3), eigenvalues 1, 2, 3 (hand-checked)

### MatIsSymmetric      (NUM-LINALG-136, built: num/linalg/mat_checks.go:27)

MatIsSymmetric reports whether an n×n matrix is symmetric: |aᵢⱼ − aⱼᵢ| ≤ tol·max|a| for every pair
(tol < 0 selects n·2⁻⁵²; tol = 0 asks for exact symmetry).
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(n²)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.4 "Symmetric Matrices" (Aᵀ = A)
- Units: a in any unit; the result is a truth value
- Example: [1 2; 2 3] → true; [1 2; 3 4] → false (hand-checked)

### MatIsDiagonal      (NUM-LINALG-137, built: num/linalg/mat_checks.go:78)

MatIsDiagonal reports whether every off-diagonal entry of a rows×cols matrix has |aᵢⱼ| ≤ tol·max|a|
(tol < 0 selects max(rows, cols)·2⁻⁵²; tol = 0 asks for exact zeros).
Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.4 "Symmetric Matrices" (the diagonal matrix Λ) and §2.6 "Elimination = Factorization: A = LU" (diagonal D in LDU)
- Units: a in any unit; the result is a truth value
- Example: [2 0; 0 3] → true; [2 1e−20; 0 3] with tol 0 → false (hand-checked)

### MatIsUpperTriangular      (NUM-LINALG-138, built: num/linalg/mat_checks.go:94)

MatIsUpperTriangular reports whether every entry below the diagonal of a rows×cols matrix has
|aᵢⱼ| ≤ tol·max|a| (tol < 0 selects max(rows, cols)·2⁻⁵²).
Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1), no allocation
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (upper triangular: uᵢⱼ = 0 for i > j)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (upper triangular U)
- Units: a in any unit; the result is a truth value
- Example: [1 2; 0 3] → true; [1 2; 4 3] → false (hand-checked)

### MatIsLowerTriangular      (NUM-LINALG-139, built: num/linalg/mat_checks.go:110)

MatIsLowerTriangular reports whether every entry above the diagonal of a rows×cols matrix has
|aᵢⱼ| ≤ tol·max|a| (tol < 0 selects max(rows, cols)·2⁻⁵²).
Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1), no allocation
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (lower triangular: lᵢⱼ = 0 for i < j)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (lower triangular L)
- Units: a in any unit; the result is a truth value
- Example: [1 0; 2 3] → true; [1 5; 2 3] → false (hand-checked)

### MatIsOrthogonal      (NUM-LINALG-140, built: num/linalg/mat_checks.go:125)

MatIsOrthogonal reports whether an n×n matrix Q is orthogonal: ‖QᵀQ − I‖_max ≤ tol (absolute, since QᵀQ
is dimensionless; tol < 0 selects n·2⁻⁵²·8 to allow a few roundings per entry).
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(n³)   the product QᵀQ
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.4 "Orthonormal Bases and Gram-Schmidt" (QᵀQ = I; a square Q is orthogonal)
- Units: a dimensionless; the result is a truth value
- Example: the rotation [0 −1; 1 0] → true; [1 1; 0 1] → false (hand-checked)

### MatIsPositiveDefinite      (NUM-LINALG-141, built: num/linalg/mat_checks.go:163)

MatIsPositiveDefinite reports whether a symmetric n×n matrix is positive definite, by attempting its
Cholesky factorisation (it succeeds exactly when every pivot is > 0); a matrix that is not symmetric to
tol (MatIsSymmetric) is reported false.
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(n³)   n³/3 flops (Cholesky)
- Space: O(n²) extra (the factor)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (positive pivots ⇔ positive eigenvalues ⇔ xᵀAx > 0)
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §4.2 "Positive Definite Systems" (Cholesky exists exactly for symmetric positive definite A)
- Units: a in any unit; the result is a truth value
- Example: [2 −1; −1 2] → true; [1 2; 2 1] → false (eigenvalues 3, −1) (hand-checked)

### MatIsPositiveSemidefinite      (NUM-LINALG-142, built: num/linalg/mat_checks.go:193)

MatIsPositiveSemidefinite reports whether a symmetric n×n matrix is positive semidefinite: its smallest
eigenvalue (MatSymmetricEigenJacobi, run with the same tol for its symmetry test and its stopping test) is
≥ −tol·max|λ| (tol < 0 selects n·2⁻⁵²·8 here and Jacobi's own default there). Cholesky cannot decide
this, since a semidefinite matrix has a zero pivot. A matrix not symmetric to tol is reported false.
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrNoConvergence
from the eigenvalue sweep.

- Time: O(n³) per Jacobi sweep
- Space: O(n²) extra (the working copy)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (semidefinite: all eigenvalues ≥ 0)
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.5 "Jacobi Methods"
- Units: a in any unit; the result is a truth value
- Example: [1 1; 1 1] → true (eigenvalues 0, 2); [1 2; 2 1] → false (hand-checked)

### MatIsSingular      (NUM-LINALG-143, built: num/linalg/mat_checks.go:220)

MatIsSingular reports whether an n×n matrix is singular to working precision: its numerical rank
(MatRank with tol; tol < 0 selects MatRank's default) is below n.
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.

- Time: O(n³)   complete-pivoting elimination
- Space: O(n²) extra (MatRank's working copy)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.4 "Independence, Basis and Dimension" (singular: dependent columns, rank < n)
- Units: a in any unit; the result is a truth value
- Example: [1 2; 2 4] → true; [1 2; 3 4] → false (hand-checked)

### Mat2SolveCramer      (NUM-LINALG-144, built: num/linalg/mat_cofactor_fixed.go:158)

Mat2SolveCramer solves the 2×2 system m·x = b by Cramer's rule, xⱼ = det(Mⱼ)/det(m) with Mⱼ = m whose
column j is replaced by b, each determinant within 2 ulps (Kahan's algorithm). It is the closed form for
a tiny system and needs no allocation; for an ill-conditioned or larger system MatSolve (LU with partial
pivoting, NUM-LINALG-096) is the stable general choice.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrSingular when det m = 0; ErrOverflow when a component is
outside the float64 range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (xⱼ = det Bⱼ / det A)
- Units: m in unit A, b in unit B; x in B/A
- Example: 2x + y = 5, x + 3y = 10 → x = (1, 3) (hand-checked)

### Mat3SolveCramer      (NUM-LINALG-145, built: num/linalg/mat_cofactor_fixed.go:186)

Mat3SolveCramer solves the 3×3 system m·x = b by Cramer's rule, each determinant a scalar triple product
summed by Dot2. It is the closed form for a tiny system and needs no allocation; MatSolve (NUM-LINALG-096)
is the stable general choice when m is ill-conditioned.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrSingular when det m = 0; ErrOverflow when a component is
outside the float64 range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (xⱼ = det Bⱼ / det A)
- Units: m in unit A, b in unit B; x in B/A
- Example: x + y + z = 6, 2y + 5z = −4, 2x + 5y − z = 27 → x = (5, 3, −2) (hand-checked)

### Mat2CharPoly      (NUM-LINALG-146, built: num/linalg/mat_cofactor_fixed.go:220)

Mat2CharPoly returns the characteristic polynomial det(λI − m) = λ² − tr(m)·λ + det(m) of a 2×2 matrix
as ascending coefficients (D-006): [det m, −tr m, 1].
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a coefficient is outside the float64
range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (det(A − λI) = λ² − (trace)λ + det for 2×2)
- Units: m in unit A; coefficient k in A^(2−k)
- Example: [4 1; 2 3] → λ² − 7λ + 10, i.e. [10, −7, 1] (hand-checked)

### Mat3CharPoly      (NUM-LINALG-147, built: num/linalg/mat_cofactor_fixed.go:244)

Mat3CharPoly returns the characteristic polynomial det(λI − m) = λ³ − tr(m)·λ² + E₂·λ − det(m) of a 3×3
matrix as ascending coefficients (D-006): [−det m, E₂, −tr m, 1], where E₂ is the sum of the three
principal 2×2 minors.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a coefficient is outside the float64
range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (det(A − λI) = 0; returned here as the monic det(λI − A) = −det(A − λI), same roots: the λ² coefficient is −trace, the λ coefficient the sum of the principal 2×2 minors, the constant −det)
- Units: m in unit A; coefficient k in A^(3−k)
- Example: diag(1, 2, 3) → (λ − 1)(λ − 2)(λ − 3) = λ³ − 6λ² + 11λ − 6, i.e. [−6, 11, −6, 1] (hand-checked)

### MatCharPolyFaddeevLeVerrier      (NUM-LINALG-148, built: num/linalg/mat_cofactor.go:165)

MatCharPolyFaddeevLeVerrier returns the characteristic polynomial det(λI − A) of an n×n matrix as n + 1
ascending coefficients (D-006; the λⁿ coefficient is 1) by the Faddeev–LeVerrier recursion
M₀ = 0, c_n = 1; Mₖ = A·Mₖ₋₁ + c_{n−k+1}·I, c_{n−k} = −tr(A·Mₖ)/k. It needs only matrix products and traces
and is exact in exact arithmetic, but in floating point the coefficients can lose all accuracy for large or
ill-scaled n (the recursion amplifies rounding error roughly like ‖A‖ⁿ/n!·u), so it is for small, well-
scaled matrices; roots of the result are poorly conditioned functions of the coefficients, so eigenvalues
should come from MatSymmetricEigenJacobi or Mat2Eigen, not from these roots. dst must not overlap a.
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOverflow when a
coefficient leaves the float64 range.

- Time: O(n⁴)   n matrix products of order n
- Space: O(n²) extra (two work matrices); O(n) for the result, one allocation only if dst lacks capacity
- Source: Hou, "Classroom Note: A Simple Proof of the Leverrier–Faddeev Characteristic Polynomial Algorithm", SIAM Review, 1998, 40(3):706–709, doi:10.1137/S003614459732076X (the recursive relations of the Leverrier–Faddeev method and their proof)
- Source: Faddeev, Faddeeva, Computational Methods of Linear Algebra, W. H. Freeman, 1963 (Leverrier's method in Faddeev's modification; attribution, no page cited)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (det(A − λI); the sum of the eigenvalues is the trace)
- Units: a in unit A; coefficient k in A^(n−k)
- Example: [2 1; 1 2] → λ² − 4λ + 3, i.e. [3, −4, 1] (hand-checked)

### Mat2Eigen      (NUM-LINALG-149, built: num/linalg/eigen.go:40)

Mat2Eigen returns the eigenvalues and eigenvectors of a general real 2×2 matrix m = [a b; c d]: the roots
λ = (a + d)/2 ± √q of λ² − (a + d)λ + (ad − bc) = 0, with the discriminant q = ((a − d)/2)² + bc formed
directly (not as trace² − 4·det, which loses √u of accuracy for nearly equal eigenvalues), a real pair or a
complex-conjugate pair; the root of smaller modulus comes from Vieta's det/λ, so neither cancels. For each λ
the eigenvector is (b, λ − a) or (λ − d, c), whichever is larger, normalised to unit length. m is scaled by an
exact power of two first, so no intermediate over- or underflows.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an eigenvalue is outside the float64 range.

- Time: O(1)
- Space: O(1), no allocation
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (λ² − (trace)λ + det = 0; eigenvectors from the null space of A − λI; complex eigenvalues of a real matrix come in conjugate pairs)
- Units: m in unit A; eigenvalues in A; eigenvectors dimensionless
- Example: [4 1; 2 3] → λ = 2, 5 with vectors ∝ (1, −2), (1, 1); [0 −1; 1 0] → λ = ±i (hand-checked)

### MatSymmetricEigenJacobi      (NUM-LINALG-150, built: num/linalg/eigen.go:119)

MatSymmetricEigenJacobi returns the eigenvalues of a symmetric n×n matrix in ascending order and, when
wantVectors is true, an n×n matrix whose column k is a unit eigenvector for value k, by the cyclic Jacobi
method: sweeps of 2×2 symmetric Schur rotations over every off-diagonal pair (p, q) in row order until the
off-diagonal Frobenius norm is at most tol·‖A‖_F (tol < 0 selects n·2⁻⁵²). Jacobi computes small
eigenvalues to high relative accuracy (Demmel–Veselić 1992). Results go to valuesDst (length n) and
vectorsDst (length n², untouched when wantVectors is false); neither may overlap a.
Errors: ErrDimension for n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf element; ErrNotSymmetric when
|aᵢⱼ − aⱼᵢ| > tol'·max|a| for some pair (tol' = max(tol, default)); ErrNoConvergence after 60 sweeps.

- Time: O(n³) per sweep   n = order; typically 5–10 sweeps
- Space: O(n²) extra (the working copy and, with vectors, one n² buffer to reorder them; O(n) for the sort index); O(n) and O(n²) for the results, one allocation each only if a dst lacks capacity
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.5 "Jacobi Methods" (the 2-by-2 symmetric Schur decomposition, the cyclic-by-row Jacobi algorithm, off(A))
- Source: Rutishauser, "The Jacobi method for real symmetric matrices", Numerische Mathematik, 1966, 9(1):1–10, doi:10.1007/BF02165223
- Source: Demmel, Veselić, "Jacobi's Method is More Accurate than QR", SIAM Journal on Matrix Analysis and Applications, 1992, 13(4):1204–1245, doi:10.1137/0613074
- Units: a in unit A; eigenvalues in A; eigenvectors dimensionless
- Example: [2 1; 1 2] → (1, 3) with vectors ∝ (1, −1), (1, 1); [2 −1 0; −1 2 −1; 0 −1 2] → 2 − √2, 2, 2 + √2 (hand-checked)

### MatPowerIteration      (NUM-LINALG-151, built: num/linalg/eigen.go:262)

MatPowerIteration returns the dominant eigenpair of an n×n matrix by the power method: x ← A·x/‖A·x‖,
λ = xᵀA·x (Rayleigh quotient), stopping when ‖A·x − λx‖₂ ≤ tol·|λ| (tol < 0 selects 1e−12) or after
maxIter steps. It converges when A has a single eigenvalue of largest modulus, real, with |λ₁| > |λ₂|, and
x0 has a component along its eigenvector; the error then falls like |λ₂/λ₁|^k. vecDst receives the
vector (length n; it must not overlap a or x0).
Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf element; ErrZeroVector when x0
or an iterate A·x is zero; ErrOutOfRange when maxIter < 1; ErrNoConvergence after maxIter steps (the
condition above fails: equal-modulus or complex dominant eigenvalues).

- Time: O(k·n²)   k = iterations used ≤ maxIter
- Space: O(n) extra (one work vector); O(n) for the result, one allocation only if vecDst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (uₖ = Aᵏu₀ = Σ cᵢλᵢᵏxᵢ is dominated by the eigenvalue of largest |λ|) and §6.4 "Symmetric Matrices" (the Rayleigh quotient xᵀAx/xᵀx)
- Units: a in unit A; the eigenvalue in A; the vector dimensionless
- Example: [2 1; 1 2], x0 = (1, 0) → λ = 3, x ∝ (1, 1) (hand-checked)

### Cholesky      (NUM-LINALG-152, built: num/linalg/cholesky.go:27)

Cholesky returns the lower-triangular L with positive diagonal such that A = L·Lᵀ, for a symmetric
positive-definite n×n matrix A, in dst (which must not overlap a). Only the lower triangle of a is read;
symmetry is not checked (MatIsSymmetric does that). Entries of dst above the diagonal are set to 0.
Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry;
ErrNotPositiveDefinite when a pivot is ≤ 0 (A is not positive definite, to working precision).

- Time: O(n³)   n³/3 flops
- Space: O(1) extra; O(n²) for the result, one allocation only if dst lacks capacity
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §4.2 "Positive Definite Systems" (the Cholesky factorization A = GGᵀ, existence for symmetric positive definite A, the gaxpy Cholesky algorithm)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (A = LDLᵀ with positive pivots; A = (L√D)(L√D)ᵀ is Cholesky)
- Units: a in unit A; L in √A
- Example: A = [4 2; 2 3] → L = [2 0; 1 √2] (hand-checked: L·Lᵀ = [4 2; 2 3])

### CholeskySolve      (NUM-LINALG-153, built: num/linalg/cholesky.go:76)

CholeskySolve solves Ax = b given the Cholesky factor L of A (A = L·Lᵀ, as returned by Cholesky) by a
forward solve Ly = b and a back solve Lᵀx = y, in dst (dst may be b).
Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrSingular when a
diagonal entry of L is zero; ErrOverflow when a component leaves the float64 range.

- Time: O(n²)   2n² flops
- Space: O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §4.2 "Positive Definite Systems" (solving Ax = b via Gy = b, Gᵀx = y) and §3.1 "Triangular Systems"
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (Cholesky) and §2.6 "Elimination = Factorization: A = LU" (two triangular solves)
- Units: L in √A, b in unit B; x in B/A
- Example: L = [2 0; 1 √2] (A = [4 2; 2 3]), b = (8, 7) → x = (1.25, 1.5) (hand-checked: 4·1.25 + 2·1.5 = 8, 2·1.25 + 3·1.5 = 7)

### QRHouseholder      (NUM-LINALG-154, built: num/linalg/qr.go:96)

QRHouseholder returns the thin QR factorisation A = Q·R of an m×n matrix (m ≥ n) by Householder
reflections: Q is m×n with orthonormal columns (in qDst) and R is n×n upper triangular (in rDst). Each
reflector is I − βvvᵀ with v chosen so that no cancellation occurs (Householder 1958); the factorisation is
backward stable. Neither dst may overlap a.
Errors: ErrDimension unless m ≥ n ≥ 1 and len(a) = m·n; ErrNotFinite for a NaN or ±Inf element;
ErrOverflow when an element of Q or R leaves the float64 range.

- Time: O(m·n²)   about 4mn² − 4n³/3 flops including forming Q
- Space: O(m·n) extra (the reflectors and a working copy); O(m·n + n²) for the results
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §5.2 "The QR Factorization" (Householder QR; the thin QR factorization) and the Householder vector computation it uses
- Source: Householder, "Unitary Triangularization of a Nonsymmetric Matrix", Journal of the ACM, 1958, 5(4):339–342, doi:10.1145/320941.320947
- Units: a in unit A; Q dimensionless; R in A
- Example: [3 0; 4 5] → |R| = [5 4; 0 3], Q·R = A and QᵀQ = I (hand-checked; the signs follow the reflectors: here R = [−5 −4; 0 −3])

### GramSchmidtModified      (NUM-LINALG-155, built: num/linalg/qr.go:163)

GramSchmidtModified returns the thin QR factorisation A = Q·R of an m×n matrix (m ≥ n) by the modified
Gram–Schmidt process, which orthogonalises each remaining column against qₖ as soon as qₖ is formed. Its
loss of orthogonality is ‖QᵀQ − I‖ ≈ u·κ(A), far better than classical Gram–Schmidt but worse than
Householder (use QRHouseholder when κ(A) is large). Q is m×n in qDst, R n×n in rDst; neither may overlap a.
Errors: ErrDimension unless m ≥ n ≥ 1 and len(a) = m·n; ErrNotFinite for a NaN or ±Inf element;
ErrRankDeficient when a column's remaining norm is ≤ tol times its original norm (tol < 0 selects
n·2⁻⁵²), i.e. it depends on earlier columns.

- Time: O(m·n²)   2mn² flops
- Space: O(n) extra (the original column norms); O(m·n + n²) for the results, one allocation each only if a dst lacks capacity
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §5.2 "The QR Factorization" (classical and modified Gram–Schmidt; the MGS orthogonality bound)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.4 "Orthonormal Bases and Gram-Schmidt" (subtract projections one at a time; A = QR)
- Units: a in unit A; Q dimensionless; R in A
- Example: [1 1; 1 0; 0 1] → q₁ = (1, 1, 0)/√2, q₂ = (1, −1, 2)/√6, R = [√2 1/√2; 0 √(3/2)] (hand-checked)

### LeastSquaresQR      (NUM-LINALG-156, built: num/linalg/qr.go:227)

LeastSquaresQR solves the full-rank least-squares problem min ‖A·x − b‖₂ for an m×n matrix (m ≥ n) by
Householder QR: Qᵀb is formed by applying the reflectors, R·x = (Qᵀb)₀…ₙ₋₁ is solved by back substitution,
and the residual norm is ‖(Qᵀb)ₙ…ₘ₋₁‖₂. Normal equations are never formed (they square κ(A)). xDst
receives x (length n; it must not overlap a or b).
Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrRankDeficient when
|rₖₖ| ≤ tol·max|rᵢᵢ| for some k (tol < 0 selects max(m, n)·2⁻⁵²); ErrOverflow when a coefficient leaves
the float64 range.

- Time: O(m·n²)   2mn² − 2n³/3 flops
- Space: O(m·n) extra (a working copy and the reflectors); O(n) for the result
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.3 "Least Squares Approximations" and §4.4 "Orthonormal Bases and Gram-Schmidt" (A = QR turns AᵀAx̂ = Aᵀb into Rx̂ = Qᵀb)
- Units: a in unit A, b in unit B; x in B/A; the residual norm in B
- Example: fit y = c + d·t to (0, 6), (1, 0), (2, 0): a = [1 0; 1 1; 1 2], b = (6, 0, 0) → x = (5, −3), residual √6 (Strang §4.3 worked example)

### ForwardSubstitution      (NUM-LINALG-157, built: num/linalg/triangular.go:38)

ForwardSubstitution solves Lx = b for an n×n lower-triangular L (entries above the diagonal are ignored)
by row-oriented forward substitution, in dst (dst may be b).
Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrSingular when a
diagonal entry is zero; ErrOverflow when a component leaves the float64 range.

- Time: O(n²)   n² flops
- Space: O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (row-oriented forward substitution)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (forward elimination Lc = b)
- Units: l in unit L, b in unit B; x in B/L
- Example: L = [2 0; 3 4], b = (2, 11) → x = (1, 2) (hand-checked)

### BackSubstitution      (NUM-LINALG-158, built: num/linalg/triangular.go:74)

BackSubstitution solves Ux = b for an n×n upper-triangular U (entries below the diagonal are ignored)
by row-oriented back substitution, in dst (dst may be b).
Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrSingular when a
diagonal entry is zero; ErrOverflow when a component leaves the float64 range.

- Time: O(n²)   n² flops
- Space: O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (row-oriented back substitution)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.2 "The Idea of Elimination" (back substitution on the upper triangular system)
- Units: u in unit U, b in unit B; x in B/U
- Example: U = [2 1; 0 4], b = (4, 8) → x = (1, 2) (hand-checked)

### TridiagonalSolve      (NUM-LINALG-159, built: num/linalg/triangular.go:111)

TridiagonalSolve solves Tx = d for an n×n tridiagonal T given by its sub-diagonal sub (len n−1), diagonal
diag (len n) and super-diagonal sup (len n−1), by the Thomas algorithm (Gaussian elimination without
pivoting specialised to three bands), in dst (which must not overlap an input). Without pivoting it is
stable for diagonally dominant or symmetric positive-definite T; for other T use MatSolve.
Errors: ErrDimension for mismatched lengths or n < 1; ErrNotFinite for a NaN or ±Inf entry; ErrSingular
when a pivot is zero; ErrOverflow when a component leaves the float64 range.

- Time: O(n)   8n flops
- Space: O(n) extra (the modified super-diagonal, one allocation); O(n) for the result
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (a tridiagonal matrix has bidiagonal L and U, O(n) work)
- Units: T in unit T, d in unit D; x in D/T
- Example: sub = (1, 1), diag = (4, 4, 4), sup = (1, 1), d = (5, 6, 5) → x = (1, 1, 1) (hand-checked)

### SVDJacobi      (NUM-LINALG-160, built: num/linalg/svd.go:38)

SVDJacobi returns the thin SVD of an m×n matrix by one-sided Jacobi (Hestenes): plane rotations applied
to pairs of columns of a working copy of A (of Aᵀ when m < n) until every pair is orthogonal to within
tol (|uᵢ·uⱼ| ≤ tol·‖uᵢ‖‖uⱼ‖; tol < 0 selects max(m, n)·2⁻⁵²); the column norms are the singular values.
One-sided Jacobi computes the singular values to high relative accuracy (Demmel–Veselić 1992). U and V are
formed only when wantVectors is true.
Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element;
ErrNoConvergence after 60 sweeps.

- Time: O(m·n·min(m, n)) per sweep   typically 6–10 sweeps
- Space: O(m·n + min(m, n)²) for the working copy, the rotation accumulator and the result; for m ≫ n a QR preconditioning step would cut the per-sweep cost to O(n³) (LATER)
- Source: Hestenes, "Inversion of Matrices by Biorthogonalization and Related Results", Journal of the Society for Industrial and Applied Mathematics, 1958, 6(1):51–90, doi:10.1137/0106005
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.6 "Computing the SVD" (Jacobi SVD procedures)
- Source: Demmel, Veselić, "Jacobi's Method is More Accurate than QR", SIAM Journal on Matrix Analysis and Applications, 1992, 13(4):1204–1245, doi:10.1137/0613074
- Units: a in unit A; singular values in A; U and V dimensionless
- Example: [3 0; 4 5] → S = (3√5, √5) ≈ (6.708204, 2.236068) (hand-checked: σ₁σ₂ = |det| = 15, σ₁² + σ₂² = ‖A‖_F² = 50)

### MatPseudoInverse      (NUM-LINALG-161, built: num/linalg/svd.go:165)

MatPseudoInverse returns the Moore–Penrose pseudo-inverse A⁺ = V·diag(1/σᵢ)·Uᵀ of an m×n matrix, an n×m
matrix, from SVDJacobi; singular values σᵢ ≤ tol·σ₁ are treated as zero (tol < 0 selects
max(m, n)·2⁻⁵²), which gives the minimum-norm least-squares solution x = A⁺b and A·A⁺·A = A. dst must not
overlap a.
Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element;
ErrNoConvergence from the SVD; ErrOverflow when an element leaves the float64 range.

- Time: O(m·n·min(m, n)) per SVD sweep, plus O(m·n·min(m, n)) to form A⁺
- Space: O(m·n) for the SVD; O(m·n) for the result, one allocation only if dst lacks capacity
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, ch. 7 "The Singular Value Decomposition (SVD)" (the pseudoinverse A⁺ = VΣ⁺Uᵀ)
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.6 "Computing the SVD" (the SVD as the tool for rank-deficient problems and the pseudo-inverse)
- Units: a in unit A; the result in 1/A
- Example: [1 2; 2 4] (rank 1) → A⁺ = A/25 = [0.04 0.08; 0.08 0.16] (hand-checked: A·A⁺·A = A)

### MatCondition2      (NUM-LINALG-162, built: num/linalg/svd.go:208)

MatCondition2 returns the 2-norm condition number κ₂(A) = σ_max/σ_min of an m×n matrix from its singular
values (SVDJacobi), k = min(m, n) of them.
Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element; ErrSingular
when σ_min = 0; ErrNoConvergence from the SVD; ErrOverflow when the ratio is outside the float64 range.

- Time: O(m·n·min(m, n)) per SVD sweep
- Space: O(m·n), the SVD's working copy
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (‖A‖ = σ_max; c = σ_max/σ_min)
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖₂) and §8.6 "Computing the SVD"
- Units: a in unit A; the result dimensionless
- Example: [3 0; 4 5] → 3√5/√5 = 3 (hand-checked)

### MatNormFrobenius      (NUM-LINALG-163, built: num/linalg/mat_norms.go:16)

MatNormFrobenius returns the Frobenius norm ‖A‖_F = √(Σᵢⱼ aᵢⱼ²) of a rows×cols matrix, scaled by a power
of two as VecNorm is, so no intermediate square overflows or underflows.
Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
ErrOverflow when the norm is outside the float64 range.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1), no allocation
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (the Frobenius norm ‖A‖_F)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (‖A‖_F² is the sum of all aᵢⱼ²)
- Units: a in unit A; the result in A
- Example: [1 2; 3 4] → √30 ≈ 5.477225575 (hand-checked)

### MatNorm1      (NUM-LINALG-164, built: num/linalg/mat_norms.go:36)

MatNorm1 returns the 1-norm ‖A‖₁ = maxⱼ Σᵢ |aᵢⱼ|, the largest absolute column sum, of a rows×cols matrix.
Each column sum is compensated (Neumaier) so the result does not depend on summation order.
Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
ErrOverflow when a column sum is outside the float64 range.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1), no allocation
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖₁ = max over j of Σᵢ |aᵢⱼ|)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (the l¹ norm of A: largest column sum)
- Units: a in unit A; the result in A
- Example: [1 −2; 3 4] → max(1 + 3, 2 + 4) = 6 (hand-checked)

### MatNormInf      (NUM-LINALG-165, built: num/linalg/mat_norms.go:70)

MatNormInf returns the infinity norm ‖A‖∞ = maxᵢ Σⱼ |aᵢⱼ|, the largest absolute row sum, of a rows×cols
matrix, each row sum compensated (Neumaier).
Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
ErrOverflow when a row sum is outside the float64 range.

- Time: O(r·c)   r = rows, c = cols
- Space: O(1), no allocation
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖∞ = max over i of Σⱼ |aᵢⱼ|)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (the l∞ norm of A: largest row sum)
- Units: a in unit A; the result in A
- Example: [1 −2; 3 4] → max(1 + 2, 3 + 4) = 7 (hand-checked)

### MatNorm2      (NUM-LINALG-166, built: num/linalg/svd.go:237)

MatNorm2 returns the spectral norm ‖A‖₂ = σ_max, the largest singular value of an m×n matrix
(SVDJacobi).
Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element;
ErrNoConvergence from the SVD.

- Time: O(m·n·min(m, n)) per SVD sweep
- Space: O(m·n), the SVD's working copy
- Source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖₂ is the square root of the largest eigenvalue of AᵀA)
- Source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (‖A‖ = σ_max)
- Units: a in unit A; the result in A
- Example: [3 0; 4 5] → 3√5 ≈ 6.708204 (hand-checked)

### CovarianceToCorrelation      (NUM-LINALG-167, built: num/linalg/covariance.go:24)

CovarianceToCorrelation returns the correlation matrix R = D^(−1/2)·Σ·D^(−1/2), rᵢⱼ = σᵢⱼ/(σᵢσⱼ) with
σᵢ = √σᵢᵢ, of an n×n covariance matrix Σ, in corrDst (the diagonal is exactly 1); stdDst receives the
standard deviations σᵢ (length n), which CorrelationToCovariance needs to go back. Symmetry is not
checked (MatIsSymmetric does that); entries are computed as σᵢⱼ/σᵢ/σⱼ so no product overflows.
Errors: ErrDimension when n < 1 or len(cov) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOutOfRange when
a variance σᵢᵢ is ≤ 0; ErrOverflow when an entry σᵢⱼ/(σᵢσⱼ) leaves the float64 range (cov is then not a
covariance matrix).

- Time: O(n²)
- Space: O(1) extra; O(n²) and O(n) for the results, one allocation each only if a dst lacks capacity
- Source: NIST/SEMATECH e-Handbook of Statistical Methods, §6.5.4.1 "Mean Vector and Covariance Matrix", https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc541.htm (read 2026-11-12) (the variance-covariance matrix S)
- Source: NIST, Dataplot Reference Manual, "CORRELATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm (read 2026-11-12) (correlation = covariance/(s_x·s_y))
- Units: cov in X·Y units; the correlation dimensionless; std in X
- Example: Σ = [4 2; 2 9] → σ = (2, 3), R = [1 1/3; 1/3 1] (hand-checked)

### CorrelationToCovariance      (NUM-LINALG-168, built: num/linalg/covariance.go:65)

CorrelationToCovariance returns the covariance matrix Σ = D^(1/2)·R·D^(1/2), σᵢⱼ = rᵢⱼ·σᵢ·σⱼ, from an n×n
correlation matrix R and the standard deviations σᵢ (each > 0), in dst.
Errors: ErrDimension when n < 1, len(corr) ≠ n² or len(std) ≠ n; ErrNotFinite for a NaN or ±Inf entry;
ErrOutOfRange when a standard deviation is ≤ 0; ErrOverflow when an entry leaves the float64 range.

- Time: O(n²)
- Space: O(1) extra; O(n²) for the result, one allocation only if dst lacks capacity
- Source: NIST, Dataplot Reference Manual, "CORRELATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm (read 2026-11-12) (covariance = correlation·s_x·s_y)
- Source: NIST/SEMATECH e-Handbook of Statistical Methods, §6.5.4.1 "Mean Vector and Covariance Matrix", https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc541.htm (read 2026-11-12)
- Units: corr dimensionless, std in X; the result in X²
- Example: R = [1 1/3; 1/3 1], σ = (2, 3) → Σ = [4 2; 2 9] (hand-checked)

