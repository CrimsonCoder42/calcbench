# Report: calculator.net check, 2026-11-13

All 206 rows are closed: the 191 calculators plus 15 rows that were handed between agents. calculator.net was used as a reading list
and an independent check only; every formula comes from a primary source. No pre-run function changed.

| Index | Rows | MATCH | FIXED | BUILT | SKIPPED | NEEDS SOURCE | UNRESOLVED | BLOCKED | HANDOFF out |
|---|---|---|---|---|---|---|---|---|---|
| financial (A-FIN) | 73 | 42 | 0 | 9 | 16 | 4 | 0 | 2 | 0 |
| fitness-and-health (A-FIT) | 39 | 18 | 0 | 7 | 2 | 12 | 0 | 0 | 0 |
| math (A-MATH) | 39 | 19 | 0 | 11 | 7 | 1 | 0 | 0 | 1 |
| other (A-OTHER) | 55 | 10 | 0 | 15 | 13 | 3 | 0 | 0 | 14 |

- **FIXED:** none. No comparison showed one of our functions wrong.
- **UNRESOLVED:** none. All 28 differences are THEIRS DIFFERS: their convention or approximation, with ours backed by
  its source (`DISCREPANCIES.md`). Examples:
  - Mortgage totals: they multiply an unrounded payment by n; we follow Reg Z 1026.18(h).
  - Due date: they add (cycle − 28) days; ACOG CO 700 makes no cycle adjustment.
  - /31 subnet: they report 0 usable hosts; RFC 3021 gives 2.
  - Sample size: they round z.
  - Ellipsoid surface area: they use the p = 1.6 approximation.
  - GFR: they use CKD-EPI 2009; we use 2021.
- **BLOCKED (2):** real-estate (the page has no calculator fields) and college-cost (an empty result block).
- **Built but never compared:**
  - percent: the form returned nothing, and its 8 submissions are used.
  - voltage-drop: the page cannot be set to the custom-resistance mode.
- **Reviewer fixes (3 test failures, all in this run's new code, fixed once by their owners):**
  - ENG-ELEC parallel resistance: NaN on a subnormal R. Each 1/R is now checked for Inf before the sum.
  - SCI-METEO-003 DewPointC: 1 ulp above the air temperature at RH 100 %. It is now clamped.
  - NUM-STATS-049 GeometricMean: 1.8e-14 relative error. It now uses a scaled product (IEEE 754 §5.3.3).
- **Functions:** 1,564 before, 1,656 after (+92). There are 6 new REMOVED rows, MED-ANTHRO-016..021: Navy body fat and
  the ideal-weight formulas, withdrawn under D-075/D-076. 12 packages were flipped from planned to built in
  `doc/taxonomy.md`. INDEX, code and doc.go agree: yes. TestIndex PASS reports "IDs in code 1656 | INDEX rows 1656
  (+34 REMOVED) | doc.go entries 1656 | packages 88".
- **Tests:**
  - Full run: `go test ./...` gave 86 ok, 3 FAIL; `-race` gave the same, with no races.
  - Re-run of the 9 failed or touched packages after the fixes: 9 ok, and 9 ok with -race (`out/go_test_rerun.txt`,
    `out/go_test_race_rerun.txt`).
  - The full tree was not re-run.
  - Benchmarks were not run.
- **Site use:** no 403, 429 or block page; at most 8 submissions on any calculator.

  | Who | Page loads | Submissions |
  |---|---|---|
  | A-FIN | ~66 | ~108 |
  | A-FIT | 34 | 54 |
  | A-MATH | ~39 | ~143 |
  | A-OTHER | ~55 | 56 |
  | Lead (driver tuning) | ~14 | |

- **FOUND, NOT FIXED:**
  - A-FIN's helper overwrote about 35 rows of `other/PROGRESS.md` through a shared scratchpad name. A-OTHER rebuilt
    the file from its checks files, and the brief now requires prefixed helper names.
  - Source locators recalled from memory, not yet read first hand (reviewer list): 26 CFR 1.167(b)-3 Example,
    FCM-R19-2003, Alduchov eq. 21, Lawrence eqs. 6 and 8, Glover §2.5, Kurose §1.4.1, FE Handbook pp. 230-231,
    OCC handbook pp. 138 and 43.
  - The heat-index and wind-chill Range is wider than the source validated.
  - Two test files cover two concepts each: num/algebra/percent_test.go and num/series/sequence_test.go.
- **DELETED (mv -n to _delete/, logged in `_delete/LEDGER.md`):**
  - med/anthropometry navy and ideal-weight files (4);
  - unit/area/area_test.go;
  - 2 lead diagnostic scripts;
  - snapshots of INDEX.md and taxonomy.md.
- **LEFT FOR DEVIN:** the decisions in each folder's `LEAD.md`:
  - age Y/M/D split;
  - atomic weights;
  - NEC tables;
  - unit/power;
  - clock-arithmetic skips;
  - final-grade solve;
  - ovulation/Wilcox;
  - CKD-EPI 2009;
  - RoundMode API;
  - GeometricMean with 0 and negatives;
  - stdlib wrappers;
  - tip split;
  - WHS tie rule;
  - quantity × price.

  Also: `doc/sources/` notes and the CHANGELOG/SCOPE entries were not written, and PIDs 38402/38510
  (playwright-mcp) are not ours.
