# Report: matrix gaps in num/linalg (2026-11-12)
**Input list (42 vCalc names, catalog/matrix.md Step 1):** KEEP 9 · REUSED 15 · MERGE 14 · DROP 4 (three calculator
pages, Matrix Test, the WikiClip; each note says what it bundles).
**Reused, re-verified against the code, never rewritten:** Mat2/3/4 Det, Inverse, Mul, MulVec, Trace, Transpose,
Identity; Mat2/Mat3SymmetricEigenvalues; MatDet, MatInverse, MatMul, MatVecMul, MatRank, MatSolve, MatTrace,
MatTranspose, MatIdentity, LUDecompose, LUSolve, MatConditionEstimate1; VecNorm, VecNorm1, VecNormInf, VecNormP.
**Built: 53 functions, NUM-LINALG-116…168** (num/linalg; file:line for each in doc/CHANGELOG.md S17):
- Elementwise/products: 116 MatAdd · 117 MatSub · 118 MatScale · 119 MatHadamard · 120 VecOuter · 121 MatKronecker
  · 122 MatPower (mat_elementwise.go, mat_products.go:115).
- Structure: 123 Mat3Minor · 124–126 Mat2/3/4Adjugate · 127 MatMinor (mat_cofactor.go:51) · 128 MatCofactorMatrix
  · 129 MatAdjugate · 130 MatRowEchelon · 131 MatRREF · 132 MatNullSpace · 133 MatColumnSpace · 134/135
  Mat2/Mat3Characteristics. Checks: 136–143 (symmetric, diagonal, upper, lower, orthogonal, PD, PSD, singular).
- Solvers/polynomials: 144/145 Mat2/Mat3SolveCramer · 146/147 Mat2/Mat3CharPoly · 148 MatCharPolyFaddeevLeVerrier
  (mat_cofactor.go:165). Eigen: 149 Mat2Eigen · 150 MatSymmetricEigenJacobi · 151 MatPowerIteration (eigen.go).
- Factorisations: 152 Cholesky · 153 CholeskySolve · 154 QRHouseholder · 155 GramSchmidtModified · 156
  LeastSquaresQR · 157/158 Forward/BackSubstitution · 159 TridiagonalSolve · 160 SVDJacobi · 161 MatPseudoInverse
  · 162 MatCondition2. Norms/statistics: 163–166 MatNorm{Frobenius,1,Inf,2} · 167/168 Covariance↔Correlation.
**HELD:** none. **DROPPED from the standard set:** "is square" (one comparison, D-078).
**LATER:** non-symmetric eigenvalues for n > 2 (needs Hessenberg + shifted QR) · QR-preconditioned SVD for m ≫ n ·
O(n³) adjugate of a singular matrix via the SVD (now O(n⁵), documented at mat_cofactor.go:95).
**Tests:** `go test ./num/linalg/... -timeout 10m -count=1` → `ok num/linalg` (out/go_test_final.txt). Oracles,
fixed seed, capped sizes: A = QR, QᵀQ = I · A = LLᵀ · A·A⁺·A = A · A·v = λv · RREF(invertible) = I · char-poly roots
= eigenvalues · Cramer = MatSolve · det(adj A) = det(A)^(n−1); Hilbert matrices (Cholesky to n = 10, κ₂(H₈) ≈
1.5258e10), extreme scales (matrix_review_test.go:11), every documented error. Whole tree: 71 ok, 0 FAIL.
**Benchmarks** (out/bench_new.txt, 92 lines, benchtime 100ms): every fixed-size path 0 allocs/op · MatAdd 1e5
176 µs, 0 allocs (dst reused) · TridiagonalSolve 1e5 1.31 ms · Cholesky n=100 142 µs · Householder QR n=100
1.54 ms · Jacobi eigen n=30 389 µs · SVD n=30 364 µs · RREF n=100 309 µs. Growth matches documented Time.
**Reviewer:** 2 BLOCKER, 6 MAJOR, 8 MINOR, all fixed, package re-run green:
- QR and least squares silently wrong at 1e±170 → power-of-two scaling, normalised Householder vectors (D-080).
- Mat2Eigen √u-inaccurate near defective → direct discriminant + scaling (matrix_review_test.go:44).
- PSD check raised an undocumented error → tolerance passed through · covariance silent Inf → ErrOverflow.
- Mat3Characteristics allocated → rankFixed (D-079, matrix_review_test.go:140) · Faddeev locator → Hou 1998.
- Minors: GVL cited only by confirmed sections, Strang titles, Space lines, SVD doc, MatPower cap, staircase test.
**Shared files updated** (wave 2 had finished): INDEX regenerated, index test ok 1,297 | 1,297 | 1,297 (+28
REMOVED); doc.go +53; taxonomy; SCOPE S17 + LATER row; D-078…D-080; CHANGELOG S17; doc/sources: 7 new notes plus
additions to GVL, Strang, NIST, Dataplot, Knuth; README counts. No PENDING_SHARED.md needed.
**FOUND, NOT FIXED:** GVL table of contents not fetchable; only 6 section titles confirmed, others not cited ·
Faddeev & Faddeeva 1963 is attribution only, no page.
**MOVED TO _delete/:** none. **LEFT FOR DEVIN:** nothing to decide. Processes: none of mine alive.
