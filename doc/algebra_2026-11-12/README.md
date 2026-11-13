# Work: algebra catalog pass, 2026-11-12

Request: clean the vCalc "Algebra" name list, research every KEEP from primary sources, order the build, and
write the shared-file rows for Devin. DOCS ONLY: no Go code, nothing under `num/` or `algo/`.
Owned paths: `doc/catalog/algebra.md` and this folder. Other subjects ran in the same tree at the same time.

| File | What |
|---|---|
| [../catalog/algebra.md](../catalog/algebra.md) | The deliverable: Step 1 table, Step 2 entries, Step 3 order |
| [PENDING_SHARED.md](PENDING_SHARED.md) | INDEX, SCOPE, taxonomy rows and requests to other subjects, ready to paste |
| [PROGRESS.md](PROGRESS.md) | Step log |

## Questions for Devin (each: recommendation · alternative · why it lost)

- **Q1 number theory home.** Recommend a new package `num/integer` (SHORT `INTEG`) for primes,
  factorisation, gcd, lcm, modular arithmetic. Alternative: keep Smallest Prime Factor in `num/algebra`. Lost
  because the taxonomy line for algebra is roots, polynomials, basic operations; number theory has its own
  sources (Knuth §4.5) and will grow.
- **Q2 API shapes.** Recommend ascending coefficients (`coeffs[i]` multiplies x^i); a zero leading
  coefficient is `ErrDegenerate` (no silent drop to a lower degree); "no real roots" is a result (n = 0), not
  an error. Alternative: descending order, as written on paper and in NumPy's `polyval`. Lost because index ≠
  power invites off-by-one errors and appending a higher term shifts every index.
- **Q3 general root finder.** Recommend Aberth–Ehrlich with Bini's start points: O(n²) a sweep, no
  deflation error, no dependency. Alternatives: Laguerre with deflation (Numerical Recipes §9.5) lost on
  accumulated deflation error and the polishing it needs; companion-matrix eigenvalues lost on O(n³) and being
  BLOCKED ON num/linalg; Jenkins–Traub lost on size (three stages, hard to verify) for no gain at our degrees.
- **Q4 stdlib rows.** 10^x and e^x were listed as both MOVE num/logs and DROP; I chose DROP (`math.Pow10`,
  `math.Pow`, `math.Exp`). e^x−1 and Scalb are MOVED as asked, but `math.Expm1` and `math.Ldexp` cover them;
  recommend num/logs drops them too. Alternative: thin wrappers for discoverability. Lost: two names for
  one thing.
- **Q5 percent change and ratio.** Recommend leaving them for the Basic Operations pass (same package), since
  they are not on this list. Alternative: add them here as 011–012. Lost: two passes numbering into one
  package invites duplicates.
- **Q6 imports inside `num`.** The tangent line (geometry) needs NUM-ALGEB-002. Recommend allowing
  `num/geometry` → `num/algebra` (acyclic, algebra imports no other num package). Alternative: build the
  tangent line in `num/algebra`. Lost: it is geometry by the request's own routing.
