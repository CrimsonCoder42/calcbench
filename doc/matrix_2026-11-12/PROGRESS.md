# Progress (one line per function)

- [x] Start: wave 2 port REPORT.md exists (finished); this run owns num/linalg and updates shared files at the end.
- [x] Re-verified the reuse list against the code: Mat2/3/4 Det, Inverse, Mul, MulVec, Trace, Transpose, Identity; Mat2/Mat3SymmetricEigenvalues; MatDet, MatInverse, MatMul, MatVecMul, MatRank, MatSolve, MatTrace, MatTranspose, MatIdentity, LUDecompose, LUSolve, MatConditionEstimate1; VecNorm, VecNorm1, VecNormInf, VecNormP: all present (NUM-LINALG-061…100, 041…044).
- [x] Built (no tests run yet): NUM-LINALG-116…122 (elementwise, outer, Kronecker, power), 123…126 (minor, adjugates), 144…147 (Cramer, char polys), 152…153 (Cholesky), 157…159 (triangular, tridiagonal), 163…165 (norms).
- [x] Built: all 53 (NUM-LINALG-116…168), tests in matrix_ops_test.go and matrix_factor_test.go, doc.go entries; gofmt/vet/build clean.
- [x] Verify: go test ./num/linalg/... 1 FAIL (MatPowerIteration A·x = 0 case: code fixed to the documented ErrZeroVector) → ok; 92 new benchmarks in out/bench_new.txt, fixed-size paths 0 allocs; reviewer started.
- [x] Reviewer: 2 BLOCKER, 6 MAJOR, 8 MINOR; fixed: QR/LS power-of-two scaling and normalised Householder vectors, Mat2Eigen direct discriminant + scaling, PSD tol pass-through, covariance overflow, Mat3Characteristics 0 allocs (rankFixed), Faddeev cited to Hou 1998, citation and Space fixes, SVD doc, MatPower cap documented, new regression tests; final go test ./num/linalg/... ok.
- [x] Shared files updated; index test ok 1,297|1,297|1,297; whole tree 71 ok; REPORT.md written.
