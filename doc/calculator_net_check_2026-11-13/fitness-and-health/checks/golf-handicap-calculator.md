# golf-handicap-calculator (HANDOFF from A-OTHER; notes ../other/checks/golf-handicap-calculator.md)

## FIELDS (2026-11-13)
calform: per round N = 1..20: rN (course rating), lN (slope), cN (18-hole score), nN (9-hole score), aN; defaults rounds 1-6 (round 2 a 9-hole score).
calform2 (course handicap): hindex 10 · ccourse 70 · cslope 121 · cpar 72.
Result: Handicap Index (0.1); the differentials sit behind a link the driver does not capture.

## RESEARCH
The R&A / USGA Rules of Handicapping (effective 1 January 2024), randa.org/roh/the-rules-of-handicapping/rule-5 read 2026-11-13:
Rule 5.2a table (read from the page's table image, first hand), Rule 5.2b "lowest 8 of the most recent 20 ... round to the nearest tenth",
Rule 5.3 maximum 54.0. Rule 5.1 (the Score Differential formula, 18- and 9-hole) is not in the page text and usga.org returned 403:
NOT built (NEEDS SOURCE). Course Handicap (Rule 6.1, calform2) not built for the same reason.

## BUILD
LIF-SPORT-001 HandicapIndexWHS — new package life/sports (taxonomy SHORT SPORT): life/sports/doc.go, handicap_index_whs.go, _test.go.
Takes Score Differentials (caller's), exact arithmetic in tenths; ties half-even (Rules give no tie rule — LEAD.md).

## COMPARE (scratch: checks/scratch/golf-handicap-calculator.go; differentials prepared with the Rule 5.1 formula as quoted by A-OTHER, PCC 0)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| their default 18-hole rounds 1,3,4,5,6 (diffs 13.2, 14.3, 18.3, 11.3, 17.0); 9-hole round 2 present | 5 diffs → lowest 1 = 11.3 | 11.3 | MATCH (their 9-hole round did not enter the count) |
| rounds 1, 3, 5 only | 3 diffs → 11.3 − 2.0 = 9.3 | 9.3 | MATCH |
DRIVER: the differentials list ("Click here to view the differentials") is not captured.
Submissions: 2 fields + 2 runs.
