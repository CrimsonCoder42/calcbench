# calculator.net check, 2026-11-13

This run takes each calculator in `doc/calculator_net_map_2026-11-13/map.json` (191 rows) and does three things:

1. Researches it from primary sources.
2. Builds or corrects our own function.
3. Compares our output with calculator.net's on the same inputs, as an independent check.

calculator.net is a reading list, never a source (BRIEF.md, The rule).

| File | What |
|---|---|
| `BRIEF.md` | the procedure, conduct on their site, owned paths, next free IDs, baseline |
| `DISCREPANCIES.md` | every difference from theirs and its verdict, merged by the lead from each folder |
| `REPORT.md` | the final report |
| `financial/` (A-FIN), `fitness-and-health/` (A-FIT), `math/` (A-MATH), `other/` (A-OTHER) | `PROGRESS.md` (one row per calculator), `checks/<slug>.md`, `checks/scratch/*.go` (`//go:build ignore`), `DISCREPANCIES.md`, `LEAD.md` |
| `out/` (gitignored) | driver proof runs on the Mortgage Calculator, vet baseline, test output |

Drivers: `testing/web/scripts/calculator_net_fields.mjs` and `calculator_net_run.mjs` (Playwright, headless
Chromium, run from `testing/web/`).
