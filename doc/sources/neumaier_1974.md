# Neumaier 1974

**Citation**: Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", *ZAMM*, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (Crossref 2026-11-12).

Taken: compensated summation whose correction term is exact whichever operand is larger. Used by ALG-WINDOW-001, -002. Checked: oracle against exact rational sums (math/big), and the case where 1e16 leaves the window.

## num (math build, 2026-11-12)

Also used by NUM-GEOM-032, -033 (shoelace area and polygon perimeter; and the AGM sum inside NUM-GEOM-046), NUM-PROB-007, -018 (and the compensated sums inside NUM-PROB-009, -016, -017, -027, -033, -036, -042), NUM-STATS-001 (and every compensated sum in num/stats, through the unexported core in compensated_sum.go). Checked: exact big.Float shoelace oracle (the unit square at 1e15 gives exactly 1); StRD NumAcc fixtures for the mean; exact big.Rat oracles in num/prob.


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-STREAM-001.

Citations as written in the doc comments:
- Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).
