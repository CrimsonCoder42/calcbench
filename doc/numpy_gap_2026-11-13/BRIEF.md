# Brief: numpy gap analysis (2026-11-13)

READ-ONLY run, risk LOW. No code, no tests, no INDEX.md or doc.go changes, no go build, nothing under _delete/.
Git is read-only. No background processes; every command bounded. Never curl or wget. Never open .env files.

## The question

calcbench is meant to play numpy's role in Go: the pure maths layer everything else imports. It is formula-first
(named, cited calculators); numpy is array-first. For each numpy routine in your area: does calcbench have the
same mathematics, and if not, does it belong in calcbench?

## calcbench inventory

- `doc/numpy_gap_2026-11-13/calcbench_functions.md`: every live exported function (1,297; generated from INDEX.md,
  which also holds 28 REMOVED rows, not listed). Columns: ID | Function | Package | Does. **Grep this file**; open a
  package only to confirm a signature (`go doc ./num/stats Variance`, bounded).
- 70 packages (`go list ./...`, plus `testing`). The maths lives mostly in: num/linalg (168), num/geometry (80),
  num/prob (68), num/stats (35), num/trig (31), num/combinatorics (22), num/series (11), num/algebra (10),
  num/integer (3), num/logs (1); array-ish algorithms in algo/sorting (10), algo/search (10), algo/selection (4),
  algo/prefix (14), algo/window (16), algo/hashing (13), algo/stream (4), algo/bits (15), algo/intmath (16).
- Homes: doc/taxonomy.md is the map. NUM packages: algebra ALGEB, calculus CALC (planned: differentiation,
  integration, ODEs), combinatorics COMBI, geometry GEOM, integer INTEG (gcd, lcm, primes, modular), linalg LINALG,
  logs LOGS, prob PROB, series SERIES, stats STATS, trig TRIG. A home not in the taxonomy is a NEW package: propose
  `num/<name>` with a SHORT of at most 6 capitals, and mark it NEW. Import order: algo < num < unit < domains.
- Rules worth knowing: .claude/rules/sources.md (primary sources only; numpy's docs are a catalogue, never the
  citation), .claude/rules/numeric.md (float policy). You need not read the rest.

## Scope rule (verbatim from the request; apply per routine, do not ask)

calcbench takes pure functions on scalars and []float64 / [][]float64, stdlib only, no new types with methods
beyond what exists today. So:
- IN SCOPE: elementwise math, reductions (sum/mean/var/std/min/max/argmin/argmax/cumsum/cumprod/diff/prod),
  rounding, clip, interpolation, histogram/bincount/digitize/percentile/quantile, correlate/convolve/cov/corrcoef,
  sorting/searching helpers numpy exposes as routines (searchsorted, unique, partition, lexsort), set routines,
  linalg routines (norm, det, inv, solve, lstsq, qr, cholesky, svd, eig, pinv, matrix_rank, kron, trace), polynomial
  (polyfit, polyval, roots, polyder, polyint, companion), window functions (hanning, hamming, blackman, bartlett,
  kaiser), FFT (fft, ifft, rfft, fftfreq -- real-valued, radix-2 and Bluestein at most), special functions numpy
  ships (sinc, i0, gcd, lcm, logaddexp, expm1, log1p, hypot, deg2rad...), floating helpers (isclose, allclose,
  nextafter, spacing, nan_to_num, nanmean family), random only as the named distributions' math, not generators.
- OUT OF SCOPE (list once, do not inventory): ndarray, dtypes, broadcasting, strides, views, ufunc machinery, masked
  arrays, structured arrays, I/O (load/save/genfromtxt), string routines, datetime64, random bit generators,
  testing/distutils/f2py/ctypeslib, matrix class, einsum's general engine (a fixed-case dot/outer/trace is fine).

Interpretation for agents: a routine that Go's `math` package already provides (sin, exp, floor, hypot, expm1,
log1p, nextafter, ...) on a scalar is **HAVE (Go stdlib: math.X)**, not MISSING: calcbench never wraps math (D-012).
Its slice form (np.sin over an array) is the question for the vec-layer verdict, so note "stdlib scalar" in the
calcbench-match column. Array creation/manipulation (zeros, reshape, concatenate, stack, tile...) is ndarray
machinery: OUT OF SCOPE, listed once, unless it is real maths (e.g. linspace, logspace, meshgrid are borderline:
judge and say why).

## Agent split (each writes doc/numpy_gap_2026-11-13/<area>.md, nothing else)

| Agent | File | numpy pages |
|---|---|---|
| A1 | `math.md` | routines.math, routines.logic, the isclose/allclose family, floating-point helpers |
| A2 | `stats.md` | routines.statistics, histograms, routines.sort (sorting, searching, counting), routines.set |
| A3 | `linalg.md` | numpy.linalg (routines.linalg), numpy.polynomial (polynomial package + legacy poly1d routines) |
| A4 | `fft_special.md` | numpy.fft, window functions, numpy.random distributions (math only), special functions (i0, sinc...) |

Authority: 1. calcbench code (calcbench_functions.md, go doc, doc.go) · 2. https://numpy.org/doc/stable/reference/routines.html
and each module page under it · 3. https://github.com/numpy/numpy/tree/main/numpy only to confirm module names. If a
page cannot be fetched, say so in your file and use the next page.

## Row format (one row per numpy routine on your pages)

`numpy name | what it computes (one line) | calcbench match (ID, "match?: <ID>" for a guess, "stdlib: math.X", or none) |
verdict | proposed home + primary source (MISSING/PARTIAL only)`

Verdicts: **HAVE** · **PARTIAL** (say what is missing: no axis, no weights, no nan variant, only one dimension...) ·
**MISSING, IN SCOPE** · **OUT OF SCOPE**. A match is the same mathematics, not the same name: grep by meaning.
Primary sources: NIST DLMF (with section), Golub & Van Loan 4th ed., Press et al. Numerical Recipes 3rd ed., Knuth
TAOCP, Higham, Abramowitz & Stegun, Welford 1962, Hyndman & Fan 1996, Harris 1978 (windows), Cooley & Tukey 1965,
Bluestein 1970, Kaiser 1974, Oppenheim & Schafer... cite to the section; numpy's docs are never the citation.

File layout: a header (pages read, pages that failed), a counts line (HAVE / PARTIAL / MISSING IN SCOPE / OUT OF
SCOPE), the table, then an "Out of scope, listed once" line, then "Hand loops in calcbench" (IDs you noticed that
loop over slices by hand where numpy has one routine, e.g. a calculator that sums or means its own slice).
Finish within 15 minutes; if you run out of time, write what you have and mark the file PARTIAL at the top.
