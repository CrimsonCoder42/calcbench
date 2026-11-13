# Progress: geometry catalog pass

| Time | Step | Note |
|---|---|---|
| 2026-11-12 | start | Read CLAUDE.md, taxonomy.md, sources.md, numeric.md, SCOPE.md, DECISIONS.md. `doc/catalog/` did not exist; created it (owned path). |
| 2026-11-12 | research | Two read-only research agents: (1) DLMF ch. 19, Carlson 1995, Kahan 2014, Ramanujan 1914, Almkvist-Berndt 1988, MathWorld pointers; (2) CRC 32nd ed. ch. 4, Bronshtein 6th ed. ch. 3, Coxeter section locators. |
| 2026-11-12 | research back | DLMF 19.8/19.9/19.33, Kahan §2/§3/Table 1, Ramanujan §15, Carlson 1995 confirmed. Thomsen ellipsoid approximation: no primary source. CRC32 §4.6.2 and §4.10.3 confirmed by errata; other CRC32 subsections inferred (†); BSMM subsections from a snippet (‡). |
| 2026-11-12 | checks | Scratch computations (session scratchpad, not the repo): every HAND-CHECKED value; AGM/MAGM perimeter vs quadrature (needs scaling to a = 1; unscaled failed at 1e6:1); Carlson R_G ellipsoid area vs prolate closed form (exact) and surface quadrature (1e-6). |
| 2026-11-12 | STEP 1a/1b/2/3 | `doc/catalog/geometry.md`: 51 input rows cleaned; 82 IDs reserved, NUM-GEOM-001…082, in build order. |
| 2026-11-12 | STEP 4 | `PENDING_SHARED.md`: 82 INDEX rows, SCOPE row, taxonomy proposal, MOVE requests, sources to create. |
| 2026-11-12 | stop | Awaiting Devin's approval. Nothing built. |
