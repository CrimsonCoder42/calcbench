# Math build, 2026-11-12

Builds every cataloged math package in one run: `num/logs`, `num/integer`, `num/combinatorics`, `num/algebra`,
`num/trig`, `num/geometry`, `num/stats`, `num/prob`. Criteria: `doc/catalog/<subject>.md`, with every open
question decided by Devin's request of 2026-11-12 (recorded as D-ids in `doc/DECISIONS.md`).

- `PROGRESS.md`: one line per package.
- `notes_<package>.md`: each builder's hand-off (INDEX rows, sources, citations confirmed or held, findings).
- `REPORT.md`: the final report.
- `out/`: raw test and benchmark output (gitignored).

Build order (each may import those before it, never the reverse): logs, integer, combinatorics, algebra, trig,
geometry, stats, prob.
