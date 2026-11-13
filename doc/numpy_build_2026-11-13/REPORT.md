# Report: numpy build, S18 (2026-11-13)
**Builders** (all DONE; leftovers relaunched on their remaining list): · B1: ALG-RED-001…010 · NUM-INTEG-004…006 ·
ALG-SELECT-005 · ALG-SORT-011…014 · ALG-HASH-014…017. · B2 (4 passes): NUM-PROB-069…155, 87 distribution functions. ·
B3: NUM-SIG-001…015: FFT family, 5 windows, Convolve, Correlate. · B4 (2 passes): NUM-LINALG-169…191,
NUM-ALGEB-011…028. · B5 (2 passes): NUM-SPEC-001…020, NUM-FLT-001…006. · B6: NUM-STATS-038…046 · NUM-CALC-001…004 ·
NUM-LOGS-002…004 · NUM-TRIG-032.
**Consolidation** (consolidation.md: 231 rows, 226 replaced, 5 not done with a reason, D-084): · No private
compensated sum, twoSum, neumaierAdd or isFinite remains outside algo/reduce. · Chi-square's private gamma now calls
NUM-SPEC-001/002. · The num/integer and num/combinatorics gcds now call NUM-INTEG-004. · Not done: algo/geometry2d
binaryGCD (algo may not import num) and algo/tree's generic checkFinite[T]. · Each edited file was copied to
`_delete/` first (237 LEDGER lines). Whole files moved: 3 empty strays.
**GAPS.md (142 rows):** BUILT 112 · PARTLY BUILT 10 · DROPPED 20. · Dropped: 17 under D2, 3 under D1. · Also dropped
inside PARTLY rows: orthogonal-polynomial coefficient arithmetic (D1) and numpy's 4 extra percentile methods (D4). ·
NEEDS SOURCE: percentile and histogram weights.
**Count:** 1,297 → 1,505 functions (+208); 70 → 75 packages. INDEX == code == doc.go: yes (1505 | 1505 (+28 REMOVED) |
1505).
**Tests** (bounded by go's own `-timeout`; macOS has no `timeout` binary): · Full run: 71 ok, 5 FAIL. Of the 15
findings, 4 were code bugs (Sinc ×2, Dirichlet allocating, the grid-2D doc example) and 11 were wrong test
expectations. · `-race`: the same 5 failed; no race. · All 15 were fixed by their owners. Rerun of the 5 packages: ok,
both plain and -race. gofmt, vet and build clean.
**D7:** · Black–Scholes matches the old code bit for bit. · NPV is within 1e-12 relative where the flows don't cancel;
where they cancel, it stays within the Horner bound (fin/interest/d7_equivalence_test.go:83). · IRR, MIRR,
profitability index and Macaulay are not converted (D-084): they need an overflow sign, a second derivative or an O(n)
buffer.
**Decisions** D-081…D-085: · reduce API, with NeumaierAdd and TwoSum added; · polynomial roots and fit live in
num/linalg; · the multivariate normal takes L; · the scope of the consolidation; · error-returning choices where numpy
clamps.
**FOUND, NOT FIXED:** · 81 inline IsNaN/IsInf expressions remain. · Von Mises moments need BesselI1. · Builders
flagged citations they wrote from memory (GvL §7.4.6, H&J §5.6, DLMF §1.4, A&S 25.4.38/40, JKB/JKK sections, Proakis
§2.6.1). · The doc/sources notes are unwritten.
**LEFT FOR DEVIN:** `/tmp/x` (outside the repo) · empty `_delete/` · re-send the parked HFT request; its start check
now passes.
