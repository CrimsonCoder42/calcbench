# Blum, Floyd, Pratt, Rivest, Tarjan 1973

**Citation**: Blum, Floyd, Pratt, Rivest, Tarjan, "Time bounds for selection", *Journal of Computer and System Sciences*, 1973, 7(4):448–461, doi:10.1016/S0022-0000(73)80033-9 (Crossref 2026-11-12).

Taken: median of medians of groups of 5 as a pivot; linear worst case. Used by ALG-SELECT-001, -002 (fallback path). Checked: `TestSelectRange_medianOfMediansOracle` forces the path on 400 inputs; `BenchmarkSelectRangeMedianOfMedians` shows linear growth.


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-SELECT-004.

Citations as written in the doc comments:
- Blum, Floyd, Pratt, Rivest, Tarjan, "Time bounds for selection", Journal of Computer and System Sciences, 1973, 7(4):448–461, doi:10.1016/S0022-0000(73)80033-9

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).
