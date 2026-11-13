# Progress: medicine

| Step | Note |
|---|---|
| orient | Read CLAUDE.md, rules, taxonomy (med/<specialty> unset), prob.md (Bayes = NUM-PROB-010/011). |
| research | 4 agents: reference-knowledgebase inventory; sources A (anthropometry/fluids/scales), B (cardio/resp), C (renal/labs/hepatic); then B2 full-text pass; lead read KDIGO 2024 PDF, Korteweg scan p. 530, Martin Fig. 2 image. Most publisher pages returned 403 to automated fetches. |
| decide | D-M2: ship only formulas read in the primary. 20 rows HELD-SOURCE, 14 HELD by policy. |
| build | Builder 1: renal, chemistry, hepatic (25). Builder 2: anthropometry, scales (15). Lead: cardio (3), respiratory (1). gofmt/vet/build only. |
| verify | go test ./med/...: 7 ok. Benchmarks: 44, all 0 B/op 0 allocs/op (out/bench.txt). Reviewer: 4 critical (suspected), 4 major, 5 minor. |
| fix | Roughton–Forster full title/pages and Zavorsky authors (confirmed by search); ErrDomain in cardio/resp; MELDNa/MELD3 scope documented (A-M9, A-M10); Dorwart and Fridericia DOIs; Figge test wording; Moens–Korteweg checks unrolled; full 180-cell Martin–Hopkins test against the lead's independent transcription. CDC labels, ADAG 2.594 and Tanaka example kept: verified in the research record. Re-ran chemistry, cardio, respiratory, hepatic: ok. |
| docs | catalog, PENDING_SHARED, PORT_CANDIDATES, sources/ written. |
