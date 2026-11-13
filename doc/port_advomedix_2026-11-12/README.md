# Work: port of the reference knowledgebase calculators into calcbench, 2026-11-12

Request: port every calculator of the read-only reference knowledgebase (source paths appear only in LEDGER.md and
ledger/) one by one, verified against primary sources; lead + builders + one reviewer; no stops.

| File | What |
|---|---|
| [PORT_CONVENTIONS.md](PORT_CONVENTIONS.md) | Binding rules for builders (identity, the unit of work, policy outcomes, IDs, packages) |
| [DECISIONS.md](DECISIONS.md) | D-P1… decisions of this run, each with the alternative that lost |
| [briefs/](briefs/) | Per-builder entry lists (routing by home package) |
| [ledger/](ledger/) | Builders' ledger sections; merged into LEDGER.md by the lead |
| [LEDGER.md](LEDGER.md) | One row per source entry |
| [PROGRESS.md](PROGRESS.md) | One line per wave event / entry batch as it lands |
| [PENDING_SHARED.md](PENDING_SHARED.md) | INDEX, SCOPE, taxonomy, DECISIONS rows for the merge pass |
| [REPORT.md](REPORT.md) | Final report |
| `out/` | Raw test and benchmark output |
