# Brief: numpy build (2026-11-13), S18

EDITS CODE. Risk NORMAL. Builders read this file and the files it names; they do not re-orient.
Rules that apply to every line: .claude/CLAUDE.md, .claude/rules/{go-functions,go-testing,go-complexity,numeric,
sources}.md (read them once, they are short). The work list: doc/numpy_gap_2026-11-13/GAPS.md and the area file it
cites in brackets (math.md, stats.md, linalg.md, fft_special.md). Read the rows you own whole. Before building a
row, re-verify its "none" by grep of `doc/numpy_gap_2026-11-13/calcbench_functions.md` and the code: a HAVE found late
is a REUSED row (note it in PROGRESS.md), never a second function.

## DECIDED (verbatim from the request; do not ask)

D1 New packages, SHORTs per doc/taxonomy.md (add rows there if absent): num/signal (SIG) -- fft, ifft, rfft, irfft,
   fftfreq, rfftfreq, fftshift, ifftshift on []float64 / []complex128, radix-2 plus Bluestein for other lengths,
   hanning, hamming, blackman, bartlett, kaiser, convolve, correlate (full/same/valid) · num/special (SPEC) -- gamma,
   lgamma wrappers only where Go's math lacks them, beta, incomplete gamma/beta (regularized), Bessel i0, sinc,
   Legendre, Chebyshev, Hermite, Laguerre evaluation and roots · num/float (FLT) -- isclose, allclose, nextafter,
   spacing, nan_to_num, copysign helpers Go lacks. Polynomial arithmetic (polyadd, polysub, polymul, polydiv, polyder,
   polyint, polyfromroots, polycompanion, roots via companion eigenvalues, polyfit any degree) goes into num/algebra.
D2 No NaN-skipping variants. Drop every nan*/fmax/fmin row; log them in GAPS.md as DROPPED (D2).
D3 isclose is symmetric: |a-b| <= max(rel*max(|a|,|b|), abs). Source PEP 485; Knuth TAOCP vol 2 §4.2.2.
D4 Percentile/quantile keep the nine Hyndman & Fan (1996) methods. numpy-only methods are DROPPED (D4).
D5 New package algo/reduce (RED): Sum (Neumaier compensated), Prod, CumSum, CumProd, Diff(n), ArgMin, ArgMax,
   AllFinite. Sources: Neumaier 1974; Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., §4.3.
   Every private compensated sum, summation helper and NaN/Inf check in the tree is replaced by a call to algo/reduce.
D6 gcd, lcm, extended gcd go in num/integer (Knuth TAOCP vol 2 §4.5.2). The three private gcds are replaced.
D7 fin/invest/black_scholes_call.go uses NUM-PROB-048 for the normal CDF. NPV, IRR, MIRR, profitability index and
   Macaulay duration use NUM-ALGEB-001/002 for polynomial evaluation. Same numbers to 1e-12 or the agent reports why.
D8 Matrix gaps in num/linalg: general eigenvalues (Hessenberg + shifted QR, Golub & Van Loan 4th ed. ch. 7), rank-
   deficient lstsq via SVD (ch. 5), slogdet, matrix_rank, pinv, kron, and whatever else GAPS.md lists there.
D9 Distributions in num/prob: every row GAPS.md lists (gamma, beta, Student t, F, Weibull, Cauchy, multivariate
   normal pdf/logpdf, binomial and Poisson quantiles). Primary sources: NIST DLMF for the special functions,
   Johnson, Kotz & Balakrishnan for the distributions. pdf, cdf, quantile, mean, variance each where defined.
D10 Every function follows the doc-comment template (ID, Time, Space, Source to the chapter/equation, Tags). IDs
   continue each prefix's sequence; never reused. "No source, no ship": a row with no primary source is left in
   GAPS.md as NEEDS SOURCE, not built.

## Lead's rulings on what DECIDED leaves open (logged as D-081…D-084 at the end)

- **R1 Import order and the polynomial split (D-082).** num/linalg already imports num/algebra (D-078), so
  num/algebra can never import num/linalg. Coefficient arithmetic (add, sub, mul = existing NUM-ALGEB-003,
  div, der, int, fromroots, companion matrix as a []float64, trim, val2d, valfromroots) goes in num/algebra.
  **Roots via companion eigenvalues and polyfit (Vandermonde + QR) go in num/linalg**, because they need the
  eigen-solver and LeastSquaresQR. Lost: both in num/algebra (an import cycle).
  Inside num/: special imports nothing from num/; prob may import special; signal may import special; linalg may
  import algebra; stats may import prob, linalg and special. Everything may import algo/reduce. algo/** imports
  nothing from num/.
- **R2 algo/reduce API (D-081), fixed now so the other builders can code against it:**
  ```go
  func Sum(values []float64) (float64, error)                // Neumaier; ErrNotFinite, ErrOverflow; empty → 0, nil
  func NeumaierAdd(sum, comp, x float64) (newSum, newComp float64) // one step, for running/streaming sums
  func TwoSum(a, b float64) (s, e float64)                    // Knuth's error-free sum (s + e = a + b exactly)
  func Prod(values []float64) (float64, error)               // empty → 1, nil; ErrNotFinite, ErrOverflow
  func CumSum(dst, values []float64) ([]float64, error)      // compensated running sum, len n
  func CumProd(dst, values []float64) ([]float64, error)
  func Diff(dst, values []float64, order int) ([]float64, error) // n-th forward difference, len n-order
  func ArgMin[T cmp.Ordered](values []T) (int, error)        // first index on ties; ErrEmpty; NaN → ErrNotFinite (float types)
  func ArgMax[T cmp.Ordered](values []T) (int, error)
  func AllFinite(values ...float64) bool                     // true when no NaN or ±Inf; AllFinite() == true
  ```
  `AllFinite` is variadic so it replaces both scalar checks (`AllFinite(a, b)`) and slice checks
  (`AllFinite(xs...)`) with no allocation. NeumaierAdd and TwoSum are the "summation helpers" D5 consolidates;
  they are additions to D5's list (D-081). Sentinels: ErrEmpty, ErrNotFinite, ErrOverflow, ErrOutOfRange.
- **R3 num/integer (D6):** `GCD(a, b uint64) uint64` (binary, Knuth §4.5.2 Algorithm B), `LCM(a, b uint64)
  (uint64, error)` (ErrOverflow), `ExtendedGCD(a, b int64) (g, x, y int64, err error)` (Algorithm X; overflow only
  at MinInt64). uint64 because all three private gcds are uint64.
- **R4 num/special API used by others:** `RegularizedGammaP(a, x float64) (float64, error)`,
  `RegularizedGammaQ(a, x float64) (float64, error)`, `RegularizedBeta(x, a, b float64) (float64, error)` (I_x(a,b)),
  `Beta(a, b float64) (float64, error)`, `LogBeta(a, b float64) (float64, error)`, `BesselI0(x float64) (float64, error)`,
  `Sinc(x float64) float64` (normalised, sin(πx)/(πx)), `RiemannZeta(s float64) (float64, error)` for s > 1.
  Go already has math.Gamma, math.Lgamma, math.Erf, math.Nextafter, math.Copysign: never wrapped (D-012).
- **R5 NaN/Inf consolidation scope.** "NaN/Inf check" = a named private helper (`isFinite`, `checkFinite`,
  `firstNotFinite`, ...). Its predicate becomes `reduce.AllFinite`; a wrapper that only formats the package's own
  sentinel error may stay, its body calling reduce.AllFinite. Inline `math.IsNaN(x) || math.IsInf(x, 0)`
  expressions (81 today) are not helpers and are left (listed as LATER).
- **R6 Moving replaced code.** A helper inside a file that stays cannot be `mv`-ed: before the first edit, copy the
  original file with `cp -n <path> _delete/<same path>` (if that path already exists in _delete from an earlier run,
  use `_delete/<same path>.2026-11-13`), then edit, then add one line to `_delete/LEDGER.md`:
  `<repo path> | replaced by <ID or function> | 2026-11-13`. A whole file that is no longer needed is `mv -n`-ed.
  `mkdir -p` the target first. Never rm.
- **R7 Already built, REUSED (not rebuilt):** matrix_rank NUM-LINALG-099 MatRank, pinv NUM-LINALG-161, kron -121,
  nextafter/copysign (Go math), lgamma/gamma (Go math), polymul = NUM-ALGEB-003.

## Next free ID per prefix (INDEX.md, 2026-11-13, REMOVED rows counted)

ALG-RED 001 (new) · NUM-SIG 001 (new) · NUM-SPEC 001 (new) · NUM-FLT 001 (new) · NUM-INTEG 004 · NUM-PROB 069 ·
NUM-LINALG 169 · NUM-ALGEB 011 · NUM-STATS 038. Other prefixes: highest in INDEX.md plus one.
IDs are numbered in the order you build; never reuse one. A new package's doc.go states its SHORT.

## Owners (exclusive paths; nothing else is touched)

| Builder | Owns | Rows (GAPS.md group) |
|---|---|---|
| B1 | algo/reduce (new), num/integer (new files only), algo/selection, algo/sorting, algo/hashing | D5, D6; §7 argmax/argmin/argpartition/cumprod/diff/sum/prod; §9 bincount, lexsort, unique_*, isin, setxor1d |
| B2 | num/prob | D9, §1 (all 29); exports nothing special needs; uses num/special (R4) |
| B3 | num/signal (new) | D1 signal, §3 (17); kaiser uses special.BesselI0 (R4) |
| B4 | num/linalg, num/algebra | D8 §4, D1 polynomial §5, §7 multiply/divide (VecHadamard, elementwise divide), §10 scalar algebra |
| B5 | num/special (new), num/float (new) | D1 special and float, D3, §11, §13, orthopoly evaluation and roots (§6, reduced by D1 to evaluation and roots) |
| B6 | num/stats, num/calculus (new, planned in taxonomy as CALC), num/logs, num/trig | §8 (D4 applies), §12, §15, §16 |
| C1…C4 | W2, after B1: algo/** · num/** · fin/**, life/**, unit/** · med/** | R5, R6, D5–D7 consolidation |

## Every builder

- Build, write doc.go entries (a new package: doc.go with the package comment, SHORT and the list) and _test.go
  files with oracle values from the cited source (small fixed inputs; each test under 60 s; no fuzz without
  -fuzztime). Benchmarks with b.ReportAllocs().
- Run ONLY: `gofmt -l <your dirs>`, `go build ./...`, `go vet ./<your pkg>/...`. **Never go test.** No `timeout`
  binary exists on this Mac: do not use it; these three commands finish on their own.
- Do not edit INDEX.md, doc/taxonomy.md, DECISIONS, CHANGELOG, SCOPE or anything shared: the lead does.
- Append your rows to doc/numpy_build_2026-11-13/PROGRESS.md: `ID | Function | file | GAPS row | BUILT` or
  `GAPS row | REUSED <ID>` or `GAPS row | NEEDS SOURCE (why)` or `GAPS row | DROPPED (D2/D4)`.
- 25 minutes of wall time. If not finished, stop cleanly (code builds), and hand back the list of what remains.
- Hand back at most 15 lines: done / partial, ID range, what remains, anything for the lead.

## Baseline (2026-11-13, before any edit)

`go build ./...` ok · `go vet ./...` ok (no output) · `go test ./... -count=1 -run XXX -timeout 120s`: 71 ok ·
`go test -list` totals: 4,583 Test functions; 5,849 Test+Benchmark+Example+Fuzz. INDEX.md: 1,325 rows = 1,297 live
+ 28 REMOVED. INDEX is regenerated by the lead with the D-070 generator.
