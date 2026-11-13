# Matrix: catalog and build the gaps in num/linalg (2026-11-12)

One agent plus the reviewer at the end. Owns num/linalg/**, doc/matrix_2026-11-12/ and doc/catalog/matrix.md.
Started after the wave 2 port finished (its REPORT.md exists), so this run updates the shared files itself.
Catalog: ../catalog/matrix.md. Conventions: ../catalog/linalg.md (row-major []float64 with dimensions alongside,
dst reuse, fixed-size Mat2/Mat3/Mat4 values that never allocate, tolerances as parameters, radians).
Files: PROGRESS.md (one line per function) · REPORT.md · out/ (test and benchmark output).

Decisions (logged to doc/DECISIONS.md as D-078 onward):
- num/linalg comes after num/prob in D-009's build order, so it may import num/algebra (CubicRoots for 3×3
  eigenvalues). Lost: a second cubic solver inside linalg (two homes for one formula).
- "is square" is not a function: rows == cols is one comparison (D-008 by analogy). Lost: a wrapper.
- The covariance ⇄ correlation bridges live in num/linalg: they are diagonal scalings D^-1/2 Σ D^-1/2 of a matrix,
  and num/stats has no matrix type. Lost: num/stats (it works on samples, not on matrices).
- General non-symmetric eigenvalues for n > 2 are LATER: they need the implicitly shifted QR algorithm (Hessenberg
  reduction and Francis steps), a separate slice. Power iteration covers the dominant eigenpair.
