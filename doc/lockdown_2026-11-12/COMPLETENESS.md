# Step 3: doc comments and package docs (2026-11-12)

Checked by testing/index_test.go (TestIndex_docCommentShape, TestIndex_docGoAgreesWithCode) and by the
inventory script run in this pass.

| Check | Before Step 3 | After |
|---|---|---|
| Packages with doc.go, a "Package <name> …" sentence and an "ID short form:" line | 67/67 | 67/67 |
| doc.go lists every ID of its package, once, under the code's name | 1,071/1,071 | 1,071/1,071 |
| ID, Time, Space, Source, Tags present and in order | 1,071/1,071 | 1,071/1,071 |
| Calculators (unit, fin, med, life: 367) with Units: | 367/367 | 367/367 |
| Calculators with Example: | 367/367 | 367/367 |
| Medical functions (225) with Range: | 225/225 | 225/225 |
| Exported functions and methods with no ID and no covering ID (D-071) | 89 under the per-method reading; 2 under D-071 (Element.Next, Element.Prev) | 0 (D-071 counts a type returned by a covered method) |
| Functions with no catalog entry and no recoverable primary source (HELD, Deprecated) | 0 | 0 |

Source lines with no narrow locator (22 flagged by a pattern search, each read): 20 carry a locator the pattern
missed (Programming Pearls "Column 9.3", ITU-T I.371 "Annex A.2", RCP "Chart 1", Guy "problem E34", a patent
number, a technical-report number). Two are second sources beside a located first source: MED-ANTHRO-012's
'MUST' Report (the first Source cites the BAPEN PDF p. 3, Steps 1–4) and ALG-GEO2D-008's own proof (A-036,
beside Atallah 1985). None is HELD.
