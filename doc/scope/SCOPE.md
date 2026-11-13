# Scope: slices and status

Status values: NOT STARTED · SPEC · BUILDING · VERIFYING · SHIPPED.

| Slice | Group | What | Spec | Status |
|---|---|---|---|---|
| S1 | ALG | Reference algorithms: search, selection, window, heap, merge, interval (13 functions) | [s1_reference_algorithms.md](../specs/s1_reference_algorithms.md) | SHIPPED 2026-11-12 |
| S2 | NUM | `num/logs`: the golden ratio constant (1) | [s2_logs.md](../specs/s2_logs.md) | SHIPPED 2026-11-12 |
| S3 | NUM | `num/integer`: primality, prime factors, smallest prime factor (3) | [s3_integer.md](../specs/s3_integer.md) | SHIPPED 2026-11-12 |
| S4 | NUM | `num/combinatorics`: factorial, binomial (exact and log), permutations, multisets, multinomial, circular, derangements, Catalan, Stirling, Bell, Stirling's approximation, big.Int twins (22) | [s4_combinatorics.md](../specs/s4_combinatorics.md) | SHIPPED 2026-11-12 |
| S5 | NUM | `num/algebra`: polynomial evaluate, derivative, multiply, roots (quadratic, cubic, Aberth–Ehrlich), linear (10) | [s5_algebra.md](../specs/s5_algebra.md) | SHIPPED 2026-11-12 |
| S6 | NUM | `num/trig`: conversions, normalisation, reciprocal and inverse functions, laws of sines and cosines, triangle solvers, central angle (31) | [s6_trig.md](../specs/s6_trig.md) | SHIPPED 2026-11-12 |
| S7 | NUM | `num/geometry`: points, lines, angles, triangles, polygons, circle, ellipse, solids, regular polyhedra, helix (80; 018/019 REMOVED → trig) | [s7_geometry.md](../specs/s7_geometry.md) | SHIPPED 2026-11-12 |
| S8 | NUM | `num/stats`: descriptive statistics, percentiles (Hyndman–Fan 1–9), dispersion, correlation, regression, grouped data, effect size (29; 007 held, 027 out) | [s8_stats.md](../specs/s8_stats.md) | SHIPPED 2026-11-12 |
| S9 | NUM | `num/prob`: probability rules, Bayes, odds, bounds, discrete and continuous distributions with quantiles (62) | [s9_prob.md](../specs/s9_prob.md) | SHIPPED 2026-11-12 |
| S10 | ALG | Algorithms core (requested as "S2"): sorting, search, selection, window, prefix, two-pointer, hashing, stack, queue, interval, heap, list, tree, trie, graph, union-find, dp, bits, intmath, text, enumerate, sample, matrix, geometry2d, ratelimit, bloom, seq, stream (330 functions; 17 held for want of a source, 6 routed to num/integer) | [../catalog/algorithms.md](../catalog/algorithms.md) | SHIPPED 2026-11-12 |
| S11 | FIN | Finance: `fin/interest`, `fin/loan`, `fin/invest`, `fin/business`, `fin/economics`, 122 functions (13 held back for sources, 2 removed as duplicates) | [catalog/finance.md](../catalog/finance.md) | SHIPPED 2026-11-12 |
| S12 | NUM | Linear algebra: Vec2/Vec3/n-D vectors, cylindrical and spherical coordinates, Mat2/Mat3/Mat4 and n×n matrices (LU, solve, det, inverse, rank, condition estimate, symmetric eigenvalues), quaternions (115 functions) | [catalog/linalg.md](../catalog/linalg.md) | SHIPPED 2026-11-12 |
| S13 | MED | Medicine: anthropometry, cardio, chemistry, hepatic, renal, respiratory, scales (44 functions, MED-ANTHRO-001…009, MED-CARDIO-001…003, MED-CHEM-001…010, MED-HEPAT-001…004, MED-RENAL-001…011, MED-RESP-001, MED-SCALE-001…006); dose rows and 20 source-held rows catalogued, not built | [catalog/medicine.md](../catalog/medicine.md) | SHIPPED 2026-11-12 |
| S14 | MED/LIF/UNT/NUM | Port of the reference knowledgebase calculators, wave 1: 209 functions in 23 packages (13 new med packages, life/fitness, life/nutrition, unit/temperature, num/series) | [port work folder README](../port_advomedix_2026-11-12/README.md) | SHIPPED 2026-11-12 |
| S15 | ALL | Lockdown: one index, the index test, doc-comment completeness, flagged citations, whole-tree test and benchmarks | [lockdown README](../lockdown_2026-11-12/README.md) | SHIPPED 2026-11-12 |
| S17 | NUM | Matrix gaps in `num/linalg`: 53 functions (NUM-LINALG-116…168): elementwise and Kronecker products, powers, minors, cofactors, adjugates, RREF and subspaces, checks, Cramer, characteristic polynomials, 2×2 and symmetric eigenproblems, power iteration, Cholesky, QR, MGS, least squares, triangular and tridiagonal solves, SVD, pseudo-inverse, condition number, norms, covariance ⇄ correlation | [catalog/matrix.md](../catalog/matrix.md) | SHIPPED 2026-11-12 |
| S18 | ALG/NUM/FIN/LIF/MED | numpy build: 208 functions from [numpy_gap GAPS.md](../numpy_gap_2026-11-13/GAPS.md) under DECIDED D1–D10. New packages: algo/reduce, num/signal, num/special, num/float, num/calculus. Additions to num/prob (87), num/linalg (23), num/algebra (18), num/stats (9) and others. Consolidation of private sums, NaN/Inf predicates and gcds onto algo/reduce and num/integer (D-081…D-085) | [numpy_build README](../numpy_build_2026-11-13/README.md) | SHIPPED 2026-11-13 |
| S19 | FIN/NUM/ALG | HFT and execution maths: 59 functions in fin/market (new), fin/invest, num/series, num/signal and algo/stream: order-book measures, returns and risk, time-series diagnostics, O(1)-per-tick streams (D-086) | [hft_math README](../hft_math_2026-11-13/README.md) | SHIPPED 2026-11-13 |
| S20 | ALG/NUM | toolbench requests R-1…R-3: algo/backoff (new, 4), num/integer CeilDiv, algo/cron (new, 2) (D-087) | [toolbench_requests README](../toolbench_requests_2026-11-13/README.md) | SHIPPED 2026-11-13 |
| — | NUM | QR-preconditioned SVD for m ≫ n; O(n³) adjugate of singular matrices via the SVD | [catalog/matrix.md](../catalog/matrix.md) | NOT STARTED (LATER) |
| — | NUM | Inline NaN/Inf expressions (81) not yet routed through reduce.AllFinite; von Mises moments (need BesselI1); weighted percentiles and histograms (need a source) | [numpy_build REPORT](../numpy_build_2026-11-13/REPORT.md) | NOT STARTED |
| — | NUM | `num/calculus` differential equations (tabulated integration, gradient and interpolation shipped in S18) | — | NOT STARTED |
| — | UNT | Unit conversions beyond temperature, kg⇄lb and cm⇄in (`unit/area`, `volume`, `duration`, `energy`, `pressure`, `speed`, `rf`, and the rest of `length` and `mass`) | — | NOT STARTED |
| — | UNT/ENG | Computer-science conversions and metrics (`unit/datasize`, `eng/computing`, `algo/checksum`, D-016) | — | NOT STARTED |
| — | SCI | Science formulas (`sci/*`) | — | NOT STARTED |
| — | ENG | Engineering formulas (`eng/*`) | — | NOT STARTED |
| — | FIN | Business and hospitality costing beyond S11 (`fin/business` hospitality and food-service rows), `fin/commodity`, `fin/personal` (D-021) | — | NOT STARTED |
| — | FIN | Trading maths beyond S19 (position sizing, pip and margin arithmetic, option Greeks beyond Black–Scholes price); S19's NEEDS SOURCE rows (DepthWeightedImbalance, TWAP, continuous Kelly) | [hft_math REPORT](../hft_math_2026-11-13/REPORT.md) | NOT STARTED |
| S16 | MED/LIF/NUM/UNT/ALG | Port wave 2: every NOT PORTED and REQUESTED row of S14 and its 17 candidates ended in a final state; 173 functions (18 + 2 withdrawn before release, D-075/D-076); new med/oncology, unit/mass, unit/length | [wave 2 README](../port_advomedix_wave2_2026-11-12/README.md) | SHIPPED 2026-11-12 |
| — | MED | The 276 NEEDS PDF entries of S16 (papers to drop into the wave 2 papers/ folder) | [NEEDS_PDF.md](../port_advomedix_wave2_2026-11-12/NEEDS_PDF.md) | NOT STARTED |
| — | LIF | Everyday calculators beyond fitness and nutrition (`life/sports`, `games`, `cooking`, `home`, `garden`, `travel`, `automotive`, `shopping`, `photography`, `music`, `brewing`) | — | NOT STARTED |

Packages are listed in [../taxonomy.md](../taxonomy.md). A slice gets a number when its spec or catalog is written.
