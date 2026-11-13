# HFT and execution maths (2026-11-13), slice S19

This slice covers order-book measures, returns and risk, time-series diagnostics, and O(1)-per-tick streaming
primitives. It holds formulas only.

- [BRIEF.md](BRIEF.md): what was DECIDED, notes on reuse and imports, who owns what, the next IDs, and the baseline.
- [PROGRESS.md](PROGRESS.md): the stages.
- [B1.md](B1.md) fin/market · [B2.md](B2.md) fin/invest · [B3.md](B3.md) num/series and num/signal ·
  [B4.md](B4.md) algo/stream: one row per function (BUILT, REUSED or NEEDS SOURCE).
- [REPORT.md](REPORT.md): the final report.
- out/: raw test output.

## Summary

| Builder | Package | IDs | BUILT | REUSED | NEEDS SOURCE |
|---|---|---|---|---|---|
| B1 | fin/market (new) | FIN-MKT-001…017 | 17 | 1 (RelativeSpread → FIN-INVEST-030) | 2 (DepthWeightedImbalance, TWAP) |
| B2 | fin/invest | FIN-INVEST-032…052 | 21 | 2 (SharpeRatio → 025; one-pair return → 010/012) | 1 (continuous Kelly) |
| B3 | num/series, num/signal | NUM-SERIES-012…019, NUM-SIG-016 | 9 | 1 (SpreadZScore → NUM-STATS-028) | 0 |
| B4 | algo/stream | ALG-STREAM-005…016 | 12 | 0 | 0 |
| **Total** | | | **59** | **4** | **3** |

Functions: 1,505 → 1,564. INDEX == code == doc.go: yes. Tests: all ok, plain and -race, after one round of fixes.
Decision: D-086.
