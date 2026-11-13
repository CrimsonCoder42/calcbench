# Report: port wave 2 (2026-11-12)
**Rows taken in:** 715 wave 1 leftovers + 17 triage candidates → 734 ledger rows (some candidates split). Final states
(LEDGER.md): PORTED 87 · PORTED-SECONDARY 57 · REUSED 25 · NEEDS PDF 292 · HELD BY POLICY 63 · NOT A CALCULATION 110 ·
SOURCE DOES NOT SUPPORT 100. Every row has exactly one.
**Both waves:** entries examined 902 + 17; now in calcbench 384 (wave 1: 171 PORTED, 44 REUSED; wave 2: 144 ported,
25 REUSED); held by policy 63; not a calculation 110; functions added 209 + 173 = 382.
**Per package** (functions added; `go test ./... -timeout 30m -count=1`: 71 ok, 0 FAIL, out/go_test_final.txt):
med/respiratory 27 · cardio 18 · neuro 18 · psychiatry 11 · renal 11 · critical 10 · trauma 9 · endocrine 8 ·
hematology 7 · gastro 6 · scales 5 · pediatrics 4 · chemistry 3 · hepatic 3 · rheumatology 3 · anthropometry 2 ·
emergency 2 · infectious 2 · obstetric 1 · oncology 1 (new) · num/prob 6 · num/stats 6 · num/series 3 · algo/window 1 ·
life/fitness 2 · unit/mass 2 (new) · unit/length 2 (new); all ok. med/pediatrics/data: WHO weight velocity, SHA-256 tested.
**CONFIRM.md:** 57 rows, each with two sources verified to cite the primary (out/verify_A.md, out/verify_B.md); BODE's
replacement cites Celli 2004 with wrong page numbers.
**NEEDS_PDF.md: 292 rows.** Ten papers unlock the most:
Fleming 2011 Web Appendix (PMC3789232; 14 rows) · Flynn 2017 Tables 4–5 (PMID 28827377; 18 paediatric BP rows) ·
Whelton 2017 (PMID 29133356; 4) · Hankinson 1999 (PMID 9872837) · D'Agostino 2008 (PMID 18212285) · Mena 2005
(PMID 15716690) · SCORE2/SCORE2-OP/PREVENT supplements · Bergstrom 1987 (PMID 3299278) · Mentzer 1973 (PMID 4123424) ·
Service 1970 (PMID 5469118).
**CONSUMER_DEFECTS.md:** about 129 rows, 20 high. Five most serious: aap_peds_htn (home-made height adjustment) ·
rescue_ihca (invented mortality %) · hiv_ckd (wrong point weights) · SAPS 3 (points disagree with Table 1 in most rows) ·
canadian_tia (item set matches no table).
**Parity disagreements** (the tests assert the primary): HEART age 65 · POMPE-C 0.799/−0.0636 vs 0.800/−0.063 · Utah
COVID items · SIRS missing PaCO2 and band criteria. %DV ties now agree (D-073).
**Reviewer:** 68 ok, 3 FAIL. Fixed: Oakland 35; RIFLE edge → None; wording (NEXUS, NEDOCS EDCensus, LAMS); SAPS 3 doc
gives 217 and 251; ALG-WINDOW-016 Example; OSTA/RIFLE Source lines. Two-source bar: 14 entries failed → 18 functions
withdrawn (D-076), plus Braden and Mentzer (D-075): 20 IDs retired, REMOVED rows; 8 entries re-sourced. Re-run: touched
packages, then the whole tree.
**Index test:** `ok`. 1,244 functions | 1,244 INDEX rows (+28 REMOVED) | 1,244 doc.go entries; 70 packages.
**Shared files updated:** INDEX regenerated · doc.go entries placed · taxonomy (oncology, mass, length built) · SCOPE S16 ·
D-075…D-077, A-058…A-067 · CHANGELOG · 142 notes → doc/sources/ · .gitignore `doc/**/papers/` · CONSUMER_CHANGES +269.
**FOUND, NOT FIXED:** none.
**LEFT FOR DEVIN:** drop PDFs into papers/ · MDRD race coefficient (kept, A-067) · SAPS 3 maximum 251 (A-063) · ARV
unweighted (A-066) · peak-flow ages (A-061/062) · sync: .claude/CLAUDE.md:3 names the product; the cleanup skill's
scope line is stale; linalg rule-sync notes. Processes: none of mine alive; gopls 8537, 8543, 61562, 61563 untouched.
