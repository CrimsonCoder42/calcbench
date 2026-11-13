# Report: lockdown (2026-11-12)

**Pending files: 9 of 9 applied, 1,139 lines** (APPLY_LOG.md; each ends "APPLIED 2026-11-12"): finance 286, linalg 210,
geometry 151, prob 111, medicine 108, algebra 74, trig 74, port 63, combinatorics 62. D-018…D-069 renumbered with a
cross-reference; A-038…A-057; SCOPE S11–S15; 4 CHANGELOG entries with file:line; requests → 15 TODO, 5 SATISFIED;
188 source notes moved (mv -n), 232 README rows. INDEX.md generated from code (D-070); taxonomy one row per package
with status, 21 packages added (D-072).

**Counts | index test:** functions 1,071 | INDEX rows 1,071 (+8 REMOVED) | doc.go 1,071; 67 packages.
`testing/index_test.go` (shape with D-071 coverage, INDEX word for word, doc.go, duplicates, taxonomy): `ok`; a
mutation check caught a dropped row and a renamed cell.

**Doc-comment completeness:** 1,071/1,071 before and after (Units 367/367, Example 367/367, Range 225/225). HELD: none.

**Citations (CITATIONS.md), all resolved, none HELD:** IE Table 5, 19/19 cells equal (PMC XML) · 4AT labels verbatim ·
Traub–Kichen pp. 107-110, doi:10.1093/ajhp/40.1.107 · Epworth: licence required, arithmetic only · Mifflin kcal/d on
secondary sources (Eklo 2026, Rather 2025) · %DV: FDA Food Labeling Guide App. H p. 129 rounds an exact half UP, code
changed `life/nutrition/percent_daily_value.go:24` (D-073) · ThreeSum gained a read source
`algo/twopointer/three_sum.go:16` · 22 catalog/report lines resolved.

**Whole tree:** gofmt -l . clean · go vet ./... clean · go build ./... ok ·
`go test ./... -timeout 30m -count=1`: 68 ok, 0 FAIL (out/go_test_final.txt) · `-race -timeout 45m`: 68 ok, 0 FAIL
(out/go_test_race.txt). Reds: none. algo/prefix and algo/selection were already green (fixed by S10).

**Benchmarks:** 1,879 lines, out/bench_all.txt (200ms, count 1, three bounded calls). Allocating scalars 0; 46 growth
candidates, all explained (out/bench_flags.txt). Flags: none.

**Product-name grep:** 2 consumer files moved into the S14 folder; product and repository names rewritten in 161
source notes, 10 algorithm docs, 1 test comment and 2 decision texts; population prose rewritten. Kept by D-074: the
S14 folder name in links (CHANGELOG, doc/README, SCOPE) and the request file; the identity word in verbatim titles
and quotations (84 files).
**By group (packages / IDs):** ALG 29/343 · NUM 10/361 · UNT 1/4 · FIN 5/122 · MED 20/225 · LIF 2/16 · 67/1,071.
**Processes:** none of mine alive; gopls 8537, 8543, 61562, 61563 (editor servers) left alone.

**FOUND, NOT FIXED:** the Mifflin unit needs the primary's full text to confirm. LIF-NUTR-006 changes at exact
ties (12.5 → 13): a consumer that pinned the old value must update.
**MOVED TO _delete/:** the empty doc/medicine_2026-11-12/sources/ and doc/port_advomedix_2026-11-12/sources/.

**LEFT FOR DEVIN:** git init · empty _delete/ · HELD list: none · Epworth licence for the consumer's use · rule-sync
suggestions from linalg §8 (.claude/rules/numeric.md FMA note; grouped test names), since .claude/ is yours.
