# Work: linear algebra (vectors, matrices, quaternions, coordinates), 2026-11-12

Request: catalog the vCalc "Vectors" name list, define the standard set, build and test every KEEP in
`num/linalg`, verify once, and hand the shared-file rows to Devin. No stops; every decision is logged with the
alternative that lost (PENDING_SHARED.md §4). Owned paths: `num/linalg/**`, this folder, `doc/catalog/linalg.md`
(plus new files under `doc/sources/`, as the request asked).

| File | What |
|---|---|
| [../catalog/linalg.md](../catalog/linalg.md) | Cleaned list (136 rows), standard set, one research entry per function (115) |
| [PENDING_SHARED.md](PENDING_SHARED.md) | INDEX rows, SCOPE row, DECISIONS entries, source index rows, requests to other subjects |
| [PROGRESS.md](PROGRESS.md) | Step log |
| `out/` | Raw `go test` and benchmark output (gitignored) |

## Citations: how each was confirmed (2026-11-12)

| Source | Confirmed by |
|---|---|
| Strang, ILA 5th ed. | official table of contents (section titles); §10.6 sample PDF read (pp. 495–498, row-vector convention) |
| Kahan 2006, "Mindless" §12 | PDF text read, pp. 46–47 |
| Diebel 2006 | PDF text read: §2.4, §5.6, §5.8, §6.1–6.5, §6.12–6.13, §7.7 equation numbers |
| Shoemake 1985 | scanned PDF read by OCR (§2.3 p. 246, §3.2 pp. 247–248, §3.3 p. 248); DOI by search |
| Ogita, Rump, Oishi 2005 | PDF text read (Algorithms 3.1, 3.5, 4.4, 5.3; Proposition 5.5); DOI on Crossref |
| Jeannerod, Louvet, Muller 2013 | DOI (confirmed in the algebra pass) |
| Blue 1978; Hager 1984; Higham 1988; Smith 1961 | DOI, volume and pages on Crossref |
| Shepperd 1978 | journal, volume, pages and DOI by search (Crossref rate-limited) |
| Golub, Van Loan 4th ed. | section titles (§3.1, §3.2, §3.4, §8.5) from the table of contents by search |
| Higham 2002 §6.1 | table of contents by search |
| NIST DLMF §1.5(ii) | page read: eqs. 1.5.14, 1.5.16, angle ranges, the θ/φ swap note |
| ISO 80000-2:2019 | NOT read (HTTP 403); not cited anywhere |

No function was held back: every citation was confirmed. The Rodrigues 1840 paper could not be found on Crossref
and is not cited; the rotation formula cites Strang §10.6 eq. (1) and Diebel §7.7.

## Defects found during verification (all fixed, re-tested)

1. **FMA fusion**: `isFinite(x) = x − x == 0` became `FMA(a, b, −a·b)` after inlining `x = a·b`, so finite results
   were reported as ErrOverflow (Vec3Project, VecNormP, VecProject, MatConditionEstimate1). Now a comparison.
   The same hazard could break TwoProduct inside Dot2 (VecDotCompensated was 25% off on test 1); explicit
   float64 conversions now forbid fusion there and in det2.
2. **Badly scaled inverses**: the closed-form 2×2/3×3/4×4 inverses scaled by the largest element only, so
   diag(1, 1e−300, 1e−300) reported a false ErrSingular. Each row is now scaled by its own power of two.
3. **Subnormal scaling**: `math.Ldexp(1, −e)` overflows for subnormal inputs (e < −1022); the scale exponent is
   now clamped (`pow2Scale`).
4. **Repeated eigenvalues** (Mat3SymmetricEigenvalues): Smith's closed form carries ~√ε relative error for
   repeated eigenvalues. Not a code defect; now documented on the function and in the catalog; the test
   tolerance for that case states the reason.

## Reviewer findings (fresh-context review; raw output out/review_*.txt) and fixes

| # | Severity | Finding | Fix |
|---|---|---|---|
| 1 | HIGH | VecAngleBetween returned NaN, no error, for subnormal input (2^(−e)/n overflowed) | scale each element, then divide; test added |
| 2 | HIGH | Vec3TriangleArea returned +Inf, no error, when |a × b| > MaxFloat64 but half of it is finite | norm of (a × b)/2 on that path; test added |
| 3 | MEDIUM | closed-form inverses scaled rows only; column-badly-scaled matrices gave a false ErrSingular | row then column equilibration (m⁻¹ = C·w⁻¹·S); test added |
| 4 | MEDIUM | VecProject false ErrOverflow for a subnormal `onto` | as 1; test added |
| 5 | MEDIUM | det2 (and Vec2Cross/Dot, Vec3Cross, Mat2/3/4Det) false ErrOverflow when only a product overflowed | rescaled retry in det2; row-rescaled retry and Dot2 sums in Mat3Det/Mat4Det; tests added |
| 6 | LOW | Vec3Project/Reject/Reflect false ErrOverflow for |v| near MaxFloat64 | scale v by a power of two around the computation |
| 7 | LOW | QuatInverse returned zeros, no error, for |q| > MaxFloat64 | scaled norm, Ldexp back |
| 8 | LOW | CartesianToSpherical gave θ = π at the origin with z = −0 | zeros at R = 0 |
| 9 | LOW | Mat2SymmetricEigenvalues lost subnormal eigenvalues | power-of-two scaling |
| 10 | LOW | MatRank false overflow near MaxFloat64 | scale the work copy (ErrOverflow removed from its doc: now unreachable) |
| 11 | LOW | LUDecompose Time line said 2n³/3 multiply-adds | n³/3 multiply-adds (2n³/3 flops) |
| 12 | LOW | two Diebel §7.7 Source lines lacked "transposed" | added |
| — | note | 4×4 minor naming resembles Eberly's Laplace-expansion note | recorded in PENDING_SHARED §4; not cited (not an acceptable source) |
| — | tests | tolerances without reasons; grouped test names | reasons added; grouped names kept (subtests carry the function names), see report |

## Left for Devin

- Paste PENDING_SHARED.md §1–§5 into INDEX.md, SCOPE.md, DECISIONS.md and doc/sources/README.md.
- Q (optional): whether an iterative symmetric eigensolver (eigenvectors, full accuracy for repeated eigenvalues)
  should be the next linalg slice, before QR / least squares.
