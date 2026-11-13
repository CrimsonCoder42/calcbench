# Report: port of the reference knowledgebase into calcbench (2026-11-12)

**Counts (LEDGER.md, 902 rows).** Entries read: 692 calculator JSON, 136 range JSON, 37 pipeline files, 32 series functions, plus the blood-sugar slice (23 rows). Outcomes: PORTED 171 · REUSED 44 · NOT PORTED 682 (15 of them are the blood-sugar rows deferred by B10 and re-done by B12) · REQUESTED 5 (2 were resolved by B5; 3 go to other owners). Pending triage: 118 rows over 49 files, nothing built (PENDING_TRIAGE.md). The medicine agent's packages were reused, never edited (D-P1): MED-SCALE-003/005, MED-RENAL-005…010, MED-HEPAT-001…004, MED-CHEM-005/009, and others.

**Per package (functions, final `go test -count=1 ./...`, out/test_final.txt).** 209 new functions. All 24 packages pass:
- New packages: cardio +42 (45 total) ok · respiratory +25 (26) ok · endocrine 23 ok · neuro 15 ok · trauma 14 ok · critical 11 ok · hematology 10 ok · fitness 9 ok · series 8 ok · nutrition 7 ok · gastro 7 ok · psychiatry 7 ok · rheumatology 5 ok · renal +5 (16) ok · anthropometry +4 (13) ok · pediatrics 4 ok · temperature 4 ok · scales +2 (8) ok · infectious 2 ok · thermoregulation 2 ok · emergency 1 ok · hepatic +1 (5) ok · obstetric 1 ok.
- Benchmarks: 271 (out/bench.txt). Allocations: 0 for every scalar function; 1 (the result) for Splits and the bucket/window functions; 21–38 for Theil–Sen.

**Citation corrections (146 ledger rows carry one). The 5 most serious:**
1. Harris–Benedict uses 6.775, a transposition of 6.7550 (Harris & Benedict 1918, p. 373).
2. The US PEF coefficients are not Hankinson 1999's.
3. NTS: the cited PMID is an epilepsy paper, and the GCS weight 0.4006 is missing.
4. Seven general/ entries cite a different instrument from the one they compute (cdiva, smast, smast_g, dast_10, rts, nyha, ace). Centor computes McIsaac, and the Rochester criteria are Jaskiewicz 1994.
5. PIAMA, PARS, Newsom, pediatric ETT, gmaws, hii_impairment, chalice_rule, cpes_score and modified_shock_index cite unrelated papers.

**Parity disagreements (tests assert the primary; the ledger records the reference's value):**
- NEWS2 SpO₂ Scale 2 at ≥ 93% on air: RCP Chart 1 scores 0 (two places).
- Centor without the McIsaac age term.
- Pediatric sleep < 4 months: AASM gives no value.
- CRAFFT case 2: the 1999 paper has no Part-A gate.
- GAD-7: the manual's labels, not the reference's.
- NTS ×0.4006 GCS weight.
- TASH: one fracture item scored by severity.
- STONE: race item and haematuria points.
- PECARN febrile infant: ANC 4090, PCT 1.71.
- Shapiro rule.
- Ottawa ankle: the 1992 rule.
- IPF GL: men's minimum 40 kg.
- Peak-flow 80% is GREEN (NHLBI).
- Hypothermia °F edges 82.4/89.6.
- OGTT 200 mg/dL is "Diabetes" (≥ 200).
- Time in range of exactly 70% is "not met" (> 70%).
- Harris–Benedict cases 1/3/5.
- %DV tie: 13 → 12, half-even (D-P12).

**Complexity improvements:** Theil–Sen O(n²) → expected O(n log n) (randomized slope selection, Dillencourt–Mount–Netanyahu 1992 / Matoušek 1991):
- 1e3: 44.8 ms → 1.40 ms.
- 1e4: 6.29 s → 16.1 ms.
- 1e5: ~40 GB → 194 ms and 12.8 MB.
The scalar ports drop the reference's per-call unit conversion and re-checks; they are O(1) with 0 allocations.

**Reviewer (one pass):** all 24 packages ok; no formula, points or label error in a sample of 50+ rows, including the 11 cardio tables read through a summarising fetch. Its own network was blocked, so those checks were against its knowledge. Fixed:
- C→F overflow (ErrOverflow); ROX overflow and AKIN at huge inputs (ErrDomain); a test for each.
- num/series cited library documentation → Reingold & Dershowitz plus ISO 8601-1.
- The duplicated NEWS2 SpO₂ table is guarded by a cross-check test (D-P14).
- Shanghai range 2–9; pack-years DOI; "patient" wording (D-P11, D-P15); Theil–Sen Space O(max(n, 2²²)).
- IE score: removed the invented upper age of 120; REMOVED INDEX rows added; ledger headers.
- TRiP(cast) re-checked verbatim (EClinicalMedicine 2020 Table 1): it matches.
- Only the affected packages were re-run; the final whole-tree run is as above.

**Consumer change list:** CONSUMER_CHANGES.md, one line per ported or reused entry (215 lines). The common changes:
- Pass inputs in the unit named by the function, converting at the boundary.
- Map ErrOutOfRange, ErrNotFinite, ErrDomain and ErrNotTestable.
- Use the primary's labels.
- Drop the reference's classifier plumbing (coordinator, all_readings); it stays in the consumer.

**FOUND, NOT FIXED (another subject's live packages):**
- algo/prefix TestRectangleSum_oracle: an index error on an empty table.
- algo/selection TestKthLargest_knuthExample: 512 vs 509.

**LEFT FOR DEVIN:**
- Paste PENDING_SHARED.md (212 INDEX rows, taxonomy, SCOPE, D-P1…15, CHANGELOG, requests) and move sources/ (161 notes).
- Confirm IE mortality Table 5: its cells were seen only in a fetch summary.
- Confirm the 4AT labels.
- Confirm the Mifflin REE unit (kcal/d).
- Settle the Traub–Kichen page dispute.
- Decide the %DV tie rule.
- Decide the Epworth licence.
- Point the next wave at PENDING_TRIAGE.md's 17 candidates and the NOT PORTED rows that only need a full text read in a browser (most publishers refused automated fetches).
