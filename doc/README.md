# calcbench docs

The function index is [`../INDEX.md`](../INDEX.md); the entry point is [`../README.md`](../README.md); the rules
are in `../.claude/` (CLAUDE.md and rules/).

## Shared files

| File | What it holds |
|---|---|
| [taxonomy.md](taxonomy.md) | Every package, built or planned, its SHORT, and where outside catalog categories land |
| [scope/SCOPE.md](scope/SCOPE.md) | The slices and their status, including what is NOT STARTED |
| [specs/](specs/) | One spec per slice, written before code ([TEMPLATE.md](specs/TEMPLATE.md); S1–S9) |
| [DECISIONS.md](DECISIONS.md) | Decisions D-001…D-080, each with the alternative that lost; a cross-reference maps run-local numbers |
| [ASSUMPTIONS.md](ASSUMPTIONS.md) | What was proceeded on without a decision (A-001…A-067) |
| [CHANGELOG.md](CHANGELOG.md) | What shipped, by date, with IDs and file:line |
| [TODO.md](TODO.md) | Known open work, including requests between subjects |
| [sources/](sources/) | One note per outside source: what was taken, the IDs that use it, how it was checked ([README](sources/README.md) indexes 413) |
| [requests/](requests/) | Requests Devin wrote as files ([lockdown_2026-11-12.md](requests/lockdown_2026-11-12.md)) |

## Catalogs

Catalog passes clean an outside list of names into KEEP / MOVE / MERGE / DROP / LATER and research each KEEP:
[algebra](catalog/algebra.md) · [algorithms](catalog/algorithms.md) · [combinatorics](catalog/combinatorics.md) ·
[finance](catalog/finance.md) · [geometry](catalog/geometry.md) · [linalg](catalog/linalg.md) ·
[matrix](catalog/matrix.md) · [medicine](catalog/medicine.md) · [prob](catalog/prob.md) · [stats](catalog/stats.md) · [trig](catalog/trig.md).

## Work folders

One folder per piece of work (`<topic>_<YYYY-MM-DD>/`): README, PROGRESS, REPORT, `out/` (raw output, gitignored).
Each applied PENDING_SHARED.md ends "APPLIED 2026-11-12".

| Folder | Work |
|---|---|
| [reference_algorithms_2026-11-12/](reference_algorithms_2026-11-12/) | S1 reference algorithms (built) |
| [catalog_stats_2026-11-12/](catalog_stats_2026-11-12/) | statistics catalog pass (docs only) |
| [algebra_2026-11-12/](algebra_2026-11-12/) · [combinatorics_2026-11-12/](combinatorics_2026-11-12/) · [geometry_2026-11-12/](geometry_2026-11-12/) · [prob_2026-11-12/](prob_2026-11-12/) · [trig_2026-11-12/](trig_2026-11-12/) | catalog passes (docs only), built by the math build |
| [math_build_2026-11-12/](math_build_2026-11-12/) | S2–S9: num/logs, integer, combinatorics, algebra, trig, geometry, stats, prob |
| [algorithms_core_2026-11-12/](algorithms_core_2026-11-12/) | S10 algorithms core (28 packages; per-package catalogs in catalog/) |
| [finance_2026-11-12/](finance_2026-11-12/) | S11 finance (research/ holds the per-package research) |
| [linalg_2026-11-12/](linalg_2026-11-12/) | S12 num/linalg |
| [medicine_2026-11-12/](medicine_2026-11-12/) | S13 medicine (conventions, specs, decision draft) |
| [port_advomedix_2026-11-12/](port_advomedix_2026-11-12/) | S14 port of the reference knowledgebase, wave 1 (ledger, conventions, consumer changes) |
| [port_advomedix_wave2_2026-11-12/](port_advomedix_wave2_2026-11-12/) | S16 port wave 2 (every leftover row in a final state; CONFIRM, NEEDS_PDF, CONSUMER_DEFECTS) |
| [matrix_2026-11-12/](matrix_2026-11-12/) | S17 matrix gaps in num/linalg |
| [lockdown_2026-11-12/](lockdown_2026-11-12/) | S15 lockdown: merge of every PENDING_SHARED, the index test, citations, whole-tree test and benchmarks |
