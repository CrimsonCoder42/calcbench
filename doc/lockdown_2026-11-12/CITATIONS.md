# Step 4: citations earlier runs flagged (2026-11-12)

Each item ends RESOLVED (with the citation and how it was read) or HELD (Deprecated line). Nothing ended HELD.
Notes appended to the matching doc/sources/ file.

## The named items

| Item | Where | Resolution |
|---|---|---|
| IE mortality Table 5 cells | MED-CARDIO-035 | RESOLVED. Park 2016, PMC4859286, Europe PMC fullTextXML, table-wrap jah31437-tbl-0005, every row label and weight read (park_lp.md). All 19 cells equal the code. The doc comment's unchecked "observed range 0–22" became the table's own span, −1 to 26. |
| 4AT labels | MED-NEURO-015 | RESOLVED. Bellelli 2014, PMC4066613, Methods "The 4AT": the three band sentences, quoted; they equal the constants verbatim (bellelli_g.md). Item points stay with the caller (supplement only). |
| Mifflin REE unit | LIF-NUTR-001/002 | RESOLVED on secondary sources. The primary's full text could not be reached (OUP 404, ScienceDirect 403) and its abstract prints no unit. Two ASN-journal papers reproduce the equation in kcal/d: Eklo 2026 AJCN Table 1 (PMC12975359) and Rather 2025 Curr Dev Nutr Table 1 (PMC12554929). The Units: line now cites Eklo (mifflin_md.md). |
| Traub–Kichen pages | MED-ANTHRO-006 | RESOLVED. Am J Hosp Pharm 1983;40(1):107-110, doi:10.1093/ajhp/40.1.107 (PubMed MedlinePgn 107-10; Crossref). The 121-122 / 40.1.121 record is Olson & Gallo, "Gluten in Pharmaceutical and Nutritional Products". DOI added to the Source line (traub_kichen_1983.md). |
| Epworth licence | MED-SCALE-001/002 | RESOLVED. epworthsleepinessscale.com/licenses: "A license is required to use the ESS … whether or not a license fee is payable". Only the scoring arithmetic and the band labels are shipped; the doc comment now says so. Band labels equal the site's verbatim (johns_1991_epworth.md). Devin: the consumer needs its own licence (LEFT FOR DEVIN). |
| %DV tie rule | LIF-NUTR-006 | RESOLVED, CODE CHANGED. FDA, Guidance for Industry: A Food Labeling Guide (January 2013), Appendix H, p. 129 (read in the FDA PDF): "exactly halfway between two whole numbers or higher (e.g., 2.5 to 2.99), the values round up". The request said half-even stays unless the FDA says otherwise, and the FDA does, so math.RoundToEven became math.Round (D-073, supersedes D-066). The parity case (2.5 of 20 → 13) now agrees with the reference. |
| Thomsen ellipsoid approximation | catalog/geometry.md | RESOLVED, not built: no primary source; no ID was ever given. NUM-GEOM-070 uses Carlson/DLMF. |

## Every "unconfirmed DOI", "not read" or "to confirm" line in a catalog or report

| Line | Resolution |
|---|---|
| catalog/geometry.md:571 (Braden 1986 DOI) | RESOLVED: Crossref 10.2307/2686282 = Braden, "The Surveyor's Area Formula", College Math J 17(4):326–337 (now aliased to 10.1080/07468342.1986.11972974). NUM-GEOM-032 second source. |
| catalog/geometry.md:734 (Adlaj 2012 DOI) | RESOLVED: Crossref 10.1090/noti879 = Adlaj, Notices AMS 59(8):1094. Not cited in code (alternative form only). |
| catalog/geometry.md:58, :282 ("Interface Angle", our reading) | RESOLVED: a catalog naming note; the function is NUM-GEOM-006, which is sourced. |
| catalog/geometry.md:1139 (helix verbatim) and :1155 | RESOLVED: NUM-GEOM-081 ships on its own Source lines (index test passes); the Gaussian row was never built. |
| catalog/combinatorics.md:258, :271 (Brualdi §2.2) | RESOLVED: NUM-COMBI-009/010 cite OEIS A111492/A144151 instead (A-013, math build); Brualdi is not cited in code. |
| catalog/combinatorics.md:11, :34 | RESOLVED: notes on requests and naming, not citations. |
| catalog/finance.md:684, :809 | RESOLVED: 684 says CONFIRMED; 809 notes one unconfirmed alternative, while FIN-INT-003/004 stand on QIA Eq. 4 and OpenStax Calculus §6.8 (both read). |
| catalog/medicine.md:31, 39, 59, 68, 111, 119 (BP table, conduction velocity, APLS, Bohr, Friedewald, Bazett) | RESOLVED as not built: HELD-SOURCE or SOURCE NOT FOUND rows with no ID; wave 2 of the port takes them. Bazett and Friedewald appear in code only as the name of an input (qtcBazettMs) or a comparison in a title. |
| catalog/medicine.md:397 (Kamath 2001 coefficient form) | RESOLVED: MED-HEPAT-002 computes the Kim 2021 form (first Source, read); Kamath is attribution only (D-052). |
| catalog/prob.md:366 (Altman & Bland 1994, "CONFIRM IN PRINT") | RESOLVED: the code cites Bland & Altman 2000, BMJ 320:1468 (Crossref confirmed; odds = P/(1 − P)); the 1994 note was moved to _delete by the math build and is not cited. |
| catalog/trig.md:28 (IEEE 754 §5.3.1) | RESOLVED: not cited by any num/trig function; the only IEEE 754 mention in code is algo/matrix/doc.go, which cites IEEE 754-2019 for float semantics in prose. |
| catalog/trig.md:457, :461 (Sinnott 1984, Vincenty 1975) | RESOLVED: NUM-TRIG-031 cites DLMF 4.42.8/4.21.28 and A&S 4.3.147 (read); the doc comment says Sinnott was not relied on; Vincenty not built. |
| catalog/stats.md:69 | RESOLVED: a dropped row (DROP), not a citation. |
| catalog/algorithms.md:306 / algorithms_core catalog/selection.md:25 (Frederickson–Johnson) | RESOLVED as not built: MedianOfTwoSorted is HELD in S10 with no ID (CHANGELOG S10 held list). |
| catalog/algorithms.md:511 / prefix.md:33 (Programming Pearls Column 8 page) | RESOLVED: ALG-PREFIX-009's first source is Laaksonen 2018 §2.4 Algorithm 3, pp. 21–23 (read first-hand); Bentley 1984 CACM 27(9):865–873 is Crossref-confirmed and Column 8 is confirmed from the table of contents; an oracle test checks every subarray. |
| catalog/algorithms.md:642 / twopointer.md:15 (Sedgewick–Wayne Ex. 1.4.15) | RESOLVED: ALG-TWOPTR-001 stands on Laaksonen 2018 §8.1 pp. 78–79 (read). ALG-TWOPTR-002 ThreeSum had only the unread exercise, so it gained a read Source: the algs4 booksite §1.4 Web Exercise "Sum of three", whose answer is the sort-then-two-pointer scan (read 2026-11-12). |
| catalog/algorithms.md:1669 / trie.md:11 (paper not read, Crossref confirmed) | RESOLVED: the trie functions cite Fredkin 1960 (Crossref) with Knuth Vol. 3 §6.3 and Manning–Raghavan–Schütze §3.2.1 for the permuterm functions; trie_test.go, autocomplete_top_k_test.go and prefix_suffix_index_test.go carry oracle tests. |
| catalog/algorithms.md:2258 / dp.md:171 (O(n log n) method not read) | RESOLVED as not built: HELD in S10 with no ID. |
| catalog/algorithms.md:3341 / ratelimit.md:102 | RESOLVED: the rate limiters cite ITU-T I.371 Annex A (GCRA, read) and RFC 2697 §3; Rathgeb 1991 (Crossref and abstract) is cited only for the window comparison. |
| catalog/algorithms.md:3425 / bloom.md:48 (Bloom 1970 full text) | RESOLVED: the Bloom functions cite Kirsch–Mitzenmacher 2008 §2 and Dillinger–Manolios 2004 §5.2 Algorithm 2 for the formulas and the hashing; Bloom 1970 is Crossref-confirmed attribution. |
| catalog/algorithms.md:3683 / stream.md:82 (Chan–Golub–LeVeque 1979 report) | RESOLVED: ALG-STREAM-004 also cites Chan, Golub, LeVeque 1983 (The American Statistician 37(3):242–247, read; chan_golub_leveque_1983.md). |
| port REPORT.md "Confirm IE mortality Table 5 / 4AT / Mifflin / Traub–Kichen / %DV / Epworth" | the named items above. |
