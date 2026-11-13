# Report: math build, 2026-11-12 (S2–S9)

**Built** (catalog order; all functions have code, doc comments, tests and benchmarks; every package `go test -count=1` ok in out/go_test_final.txt):
| Package | Functions | Benchmarks (out/bench_<pkg>.txt, + _fix) |
|---|---|---|
| num/logs | 1 (constant NUM-LOGS-001) | 2 ns/op, 0 allocs |
| num/integer | 3 | IsPrime/PrimeFactors 32/64-bit; PrimeFactors 64-bit worst 2.0 ms; the result is the only allocation |
| num/combinatorics | 22 (16 + 6 big.Int twins) | uint64 forms 0 allocs; BellBig n=1000 25 ms |
| num/algebra | 10 | Horner 0 allocs; PolynomialRoots deg 1e3 126 ms, 4 allocs (2.6× slower after the accuracy fix) |
| num/trig | 31 | all 0 allocs; slowest SSA 120 ns |
| num/geometry | 80 (018/019 REMOVED → TRIG-020/021) | all 0 allocs; polygon 1e5 vertices 2.1 ms (linear) |
| num/stats | 29 (007 held, 027 out) | linear growth; MAD 1 alloc; Mode 1e5 8.3 ms (sort) |
| num/prob | 62 | all 0 allocs; CDF growth ~√n as documented |
INDEX.md: 253 rows (14 ALG + 237 NUM + 2 REMOVED); cleanup counts: 238 declarations = 238 rows = 238 doc.go entries.

**Reviewer** (14 findings): fixed CubicRoots wrong roots for spread coefficients (deflation direction, Vieta small root, Newton
polish); PolynomialRoots Wilkinson-20 forward error (compensated Horner, Graillat & Ménissier-Morain 2012); stats silent
underflow at 1e-200 (power-of-two rescaling); odds functions 014–017 re-sourced (Štrumbelj unreadable → Wheatcroft 2021,
Egidi et al. 2018, Wunderlich & Memmert 2018); SAS/ASA/AAS/SSA near π; Bini constants (3.8 removed, 0.7 now our own);
DiscreteVariance correction; Uniform wide ranges; MAD allocs; doc field order; DLMF release; Sinnott and Altman–Bland
lines dropped. Fixed before review: 11 test failures (mostly test data; AngleBetweenLinesRad made exact for parallel lines).

**Citations not confirmed → held**: NUM-STATS-007 GroupedMode (Spiegel ch. 3 formula not seen; code in _delete/num/stats/);
NUM-STATS-027 WeightedGeometricMean (no source). Dropped lines (function still covered): Brualdi §2.2, Bronshtein, CRC 32nd
(→ 31st ed., read), Stewart (→ OpenStax), Knuth §4.6 (→ CLRS §30.1), Kahan 1987, A&S 4.4.13, Hampel 1974, Štrumbelj 2014,
Altman–Bland 1994, Sinnott 1984. Aberth 1973 formula text not re-read first-hand (record confirmed; catalog CHECKED).
Cited URLs were not re-fetched in verify (curl denied); builders opened them during the build.

**FOUND, NOT FIXED** (other live subjects): algo/selection new KthLargest test fails (509 vs 512, files of 15:07);
algo/prefix TestRectangleSum_oracle (0×0 table); algo/stack BalancedBrackets case. `.claude/CLAUDE.md:17` lists num topics
without `integer`. Catalogs still word Q's as open and carry the values the builders corrected (doc/TODO.md).

**LEFT FOR DEVIN**: overround definition (Wheatcroft Σ−1 vs Wunderlich Σ; 016 follows Wheatcroft); PolynomialRoots start-angle
0.7 unsourced own choice; Spiegel check to release 007; Hogg 8th vs 7th ed. (A-026); StirlingApprox/big twins have no n cap;
medicine PENDING_SHARED.md not folded (not this run); second source for GEOM-057 (A-032).

**MOVED TO _delete/**: num/stats/grouped_mode.go (held) · num/algebra/zz_debug_test.go, zz_stress_test.go (fixer scratch) ·
doc/sources/strumbelj_2014.md, altman_bland_1994.md, sinnott_1984.md (no longer cited).
