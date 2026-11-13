# Work: medicine catalog and build, 2026-11-12

Request: MEDICINE: CATALOG AND BUILD, NO STOPS. Owned paths: `med/**`, this folder, `doc/catalog/medicine.md`.
Shared files (INDEX, SCOPE, taxonomy, DECISIONS, ASSUMPTIONS, CHANGELOG, TODO, doc/sources) are not edited:
their rows are in PENDING_SHARED.md.

| File | What |
|---|---|
| [../catalog/medicine.md](../catalog/medicine.md) | The catalog: cleaning table, standard set, 44 built entries, HELD / HELD-SOURCE / not found |
| [CONVENTIONS.md](CONVENTIONS.md) | Binding conventions for every med package (errors, units, Range, categories, tests) |
| [spec_renal_chemistry_hepatic.md](spec_renal_chemistry_hepatic.md), [spec_anthropometry_scales.md](spec_anthropometry_scales.md) | Build specs given to the builders |
| [DECISIONS_DRAFT.md](DECISIONS_DRAFT.md) | D-M1…D-M9, A-M1…A-M8 (A-M9, A-M10 in PENDING_SHARED) |
| [PENDING_SHARED.md](PENDING_SHARED.md) | INDEX rows (44), SCOPE, taxonomy, ledgers, TODO, moves |
| PORT_CANDIDATES (moved 2026-11-12 to the S14 port work folder as PORT_CANDIDATES_medicine_run.md: it names the consumer's files) | reference-knowledgebase inventory with file:line → calcbench ID / NEW / HELD |
| [sources/](sources/) | One note per source, to move into doc/sources/ |
| `out/bench.txt` | Raw benchmark output (gitignored) |
