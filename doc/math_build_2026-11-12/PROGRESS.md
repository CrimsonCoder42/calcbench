# Progress

| Package | Status |
|---|---|
| num/logs | BUILT (lead): 1 constant NUM-LOGS-001 GoldenRatio |
| num/integer | BUILT: 3 functions NUM-INTEG-001…003, none held |
| num/combinatorics | BUILT: 22 functions NUM-COMBI-001…022 (incl. 6 big twins), none held; Brualdi §2.2 dropped, 009/010 cite OEIS A111492/A144151; Derangements limited to n ≤ 20 (catalog's 21 overflows) |
| num/algebra | BUILT: 10 functions NUM-ALGEB-001…010, none held; Stewart→OpenStax College Algebra 2e §4.1, Knuth §4.6→CLRS §30.1, NR 2nd ed. §5.6; Bini constants confirmed only from Bini's reference-code comments |
| num/trig | BUILT: 31 functions NUM-TRIG-001…031, none held; A&S 4.4.13 and Kahan 1987 dropped (covered by DLMF, A&S 4.4.4, POSIX); Sinnott 1984 via two secondary sources |
| num/stats | BUILT: 29 functions (001–006, 008–026, 028–031); HELD 007 GroupedMode (Spiegel ch. 3 formula not seen; code in _delete/num/stats/), 027 WeightedGeometricMean OUT (no source); Hogg-McKean-Craig cited as 8th ed. §1.9 Def. 1.9.2 (7th ed. numbering unconfirmed) |
| num/prob | BUILT: 62 functions NUM-PROB-001…062, none held; NegativeBinomialTrials* naming; Loader route kept; Ross cited at section level; odds formats and Cauchy held (no ID) |
| num/geometry | BUILT: 80 functions NUM-GEOM-001…082 except 018/019 (REMOVED → NUM-TRIG-020/021), none held; CRC cited as 31st ed. (2003, read); Bronshtein dropped |
| VERIFY | go test ./num/... all ok (out/go_test_final.txt); reviewer: 14 findings, all HIGH/MEDIUM fixed; INDEX regenerated (253 rows) |
| DONE | REPORT.md written; /cleanup counts agree |
